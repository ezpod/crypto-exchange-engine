package main

import (
	"github.com/Shopify/sarama"
	"log"
	"time"
	"github.com/ezpod/crypto-exchange-engine/engine"
	"fmt"
	"math/rand"
	"github.com/panjf2000/ants"
	"sync"
)

func main() {
	producer := createProducer()
	defer producer.AsyncClose()
	pool, _ := ants.NewPool(2000)
	defer pool.Release()
	var wg sync.WaitGroup

	go func(p sarama.AsyncProducer) {
		for {
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
		}
	}(producer)

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		pool.Submit(sendOrder(producer))
		time.Sleep(100 * time.Nanosecond)
	}
	wg.Wait()

}

func sendOrder(producer sarama.AsyncProducer)  func() {
	return func(){

		strKey :=  GetRandomString(10)
		srcValue := engine.Order{
			Amount: uint64(RandInt64(int64(1), int64(10000))),
			Price:  uint64(RandInt64(int64(1), int64(100))),
			ID:     strKey,
			Side:   int8(RandInt64(int64(1), int64(3))-int64(1)),
		}

		value := fmt.Sprintf(string(srcValue.ToJSON()))
		msg := &sarama.ProducerMessage{
			Topic: "lzpOrders",
		}

		msg.Key = sarama.StringEncoder(strKey)
		msg.Value = sarama.ByteEncoder(value)
		//fmt.Println(value)

		producer.Input() <- msg
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

// Interval Random Number
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	time.Sleep(100 * time.Nanosecond)
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

// Generate the specified digit string
func  GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	time.Sleep(100 * time.Nanosecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}