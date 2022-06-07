package main

/*
使用http请求请求某个重要功能过后放到消息队列里，针对这个功能单独定制一个topic，
在这个topic下实现：经过这个功能过后，然后从消息队列取出请求，再定时定量的去执行后续一些附带的不重要的功能
消息队列里存的是：对应的请求身份标识（token）（可能还有请求所需数据一起打包放入消息队列的消息里）json
*/

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var Address = []string{"120.26.125.147:9092"}

func main() {
	syncProducer(Address)

	//ASyncProducer异步生产者
	//ASyncProducer(Address)

}

//同步消息模式
func syncProducer(address []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 5 * time.Second

	//并发量小时，使用同步生产者；大时使用异步生产者
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

func ASyncProducer(address []string, limit int64) {
	config := sarama.NewConfig()
	// 异步生产者不建议把 Errors 和 Successes 都开启，一般开启 Errors 就行
	// 同步生产者就必须都开启，因为会同步返回发送成功或者失败
	config.Producer.Return.Errors = false   // 设定需要返回错误信息
	config.Producer.Return.Successes = true // 设定需要返回成功信息
	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		log.Fatal("NewSyncProducer err:", err)
	}
	defer producer.AsyncClose()

	topic := "test-hello-world"

	go func() {
		// [!important] 异步生产者发送后必须把返回值从 Errors 或者 Successes 中读出来 不然会阻塞 sarama 内部处理逻辑 导致只能发出去一条消息
		for {
			select {
			case s := <-producer.Successes():
				log.Printf("[Producer] key:%v msg:%+v \n", s.Key, s.Value)
			case e := <-producer.Errors():
				if e != nil {
					log.Printf("[Producer] err:%v msg:%+v \n", e.Msg, e.Err)
				}
			}
		}
	}()
	// 异步发送
	for i := int64(0); i < limit; i++ {
		str := strconv.Itoa(int(time.Now().UnixNano()))
		msg := &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(str)}
		// 异步发送只是写入内存了就返回了，并没有真正发送出去
		// sarama 库中用的是一个 channel 来接收，后台 goroutine 异步从该 channel 中取出消息并真正发送
		producer.Input() <- msg

		atomic.AddInt64(&i, 1)
		if atomic.LoadInt64(&i)%1000 == 0 {
			log.Printf("已发送消息数:%v\n", i)
		}

	}
	log.Printf("发送完毕 总发送消息数:%v\n", limit)
}
