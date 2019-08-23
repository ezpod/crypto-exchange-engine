package main

import (
	"github.com/bsm/sarama-cluster"
	"github.com/Shopify/sarama"
	"log"
	"fmt"
	"time"
)

func main() {
	consumer := createConsumer()

	done := make(chan bool)
	go func() {
		for{
			for msg := range consumer.Messages() {
				fmt.Println(time.Now().Unix(), string(msg.Value))

				// Monitor transaction processing, record testing, etc.

				// mark the message as processed
				consumer.MarkOffset(msg, "")
			}
		}
		done <- true
	}()
	<-done
}

//
// Create the consumer
//

func createConsumer() *cluster.Consumer {
	// define our configuration to the cluster
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = false
	config.Group.Return.Notifications = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// create the consumer
	consumer, err := cluster.NewConsumer([]string{"127.0.0.1:9092"}, "lzpConsumer1", []string{"lzpTrades"}, config)
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
