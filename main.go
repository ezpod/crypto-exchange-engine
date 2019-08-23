package main

import (
	"github.com/ezpod/crypto-exchange-engine/engine"
	"log"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"time"
	"fmt"
	"sync"
	"github.com/panjf2000/ants"
)

func main() {

	// create the consumer and listen for new order messages
	consumer := createConsumer()

	// create the producer of trade messages
	producer := createProducer()

	// create the order book
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 100),
		SellOrders: make([]engine.Order, 0, 100),
	}

	// create a signal channel to know when we are done
	done := make(chan bool)

	// start processing orders
	go func() {
		for msg := range consumer.Messages() {

			var order engine.Order
			// decode the message
			order.FromJSON(msg.Value)
			// process the order
			//trades :=
			trades := book.Process(order)

			// send trades to message queue
			pool, _ := ants.NewPool(2000)
			defer pool.Release()

			var wg sync.WaitGroup

			for _, trade := range trades {
				pool.Submit(sendTrade(producer, trade, &wg))
			}

			// mark the message as processed
			consumer.MarkOffset(msg, "")
			wg.Wait()
		}
		done <- true
	}()

	// wait until we are done
	<-done
}

//
// Create the consumer
//

func sendTrade(producer sarama.AsyncProducer, trade engine.Trade, wg *sync.WaitGroup)  func() {
	return func(){
		go func(p sarama.AsyncProducer) {
			select {
			case suc := <-p.Successes():
				if suc != nil {
					fmt.Printf("succeed, offset=%d, timestamp=%s, partitions=%d\n", suc.Offset, suc.Timestamp.String(), suc.Partition)
					wg.Done()
					//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
				}
			case fail := <-p.Errors():
				if fail != nil {
					fmt.Printf("error= %v\n", fail.Err)
					wg.Done()
				}
			}
		}(producer)

		wg.Add(1)
		value := fmt.Sprintf(string(trade.ToJSON()))
		msg := &sarama.ProducerMessage{
			Topic: "lzpTrades",
		}
		msg.Value = sarama.ByteEncoder(value)
		//fmt.Println(value)

		producer.Input() <- msg
	}
}

func createConsumer() *cluster.Consumer {
	// define our configuration to the cluster
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = false
	config.Group.Return.Notifications = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// create the consumer
	consumer, err := cluster.NewConsumer([]string{"127.0.0.1:9092"}, "lzpConsumer", []string{"lzpOrders"}, config)
	if err != nil {
		log.Fatal("Unable to connect consumer to kafka cluster")
	}
	go handleErrors(consumer)
	go handleNotifications(consumer)
	return consumer
}

func handleErrors(consumer *cluster.Consumer) {
	for err := range consumer.Errors() {
		log.Printf("Error: %s\n", err.Error())
	}
}

func handleNotifications(consumer *cluster.Consumer) {
	for ntf := range consumer.Notifications() {
		log.Printf("Rebalanced: %+v\n", ntf)
	}
}

//
// Create the producer
//


func createProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 2 * time.Second
	config.Producer.Retry.Max = 3
	//config.Version = sarama.V0_10_0_1
	producer, err := sarama.NewAsyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatal("Unable to connect producer to kafka server")
	}
	return producer
}
