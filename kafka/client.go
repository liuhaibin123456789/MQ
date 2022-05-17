package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
	//	"github.com/bsm/sarama-cluster"
)

var Address = []string{"120.26.125.147:9092"}

func main() {
	syncProducer(Address)
}

//同步消息模式
func syncProducer(address []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()

	topic := "test-hello-world"
	srcValue := "sync: this is a message. index=%d "
	srcKey := "msgKey: %d"

	for i := 0; i < 100; i++ {
		value := fmt.Sprintf(srcValue, i)
		key := fmt.Sprintf(srcKey, i)

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder(key),
			Value: sarama.ByteEncoder(value),
		}
		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		} else {
			fmt.Fprintf(os.Stdout, value+"发送成功，partition=%d, offset=%d \n", part, offset)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
