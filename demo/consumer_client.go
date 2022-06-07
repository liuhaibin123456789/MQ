package main

import (
	"MQ/demo/proto"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

/*
	消费者客户端
*/

func main() {
	//监听
	conn, err := grpc.Dial("127.0.0.1:12345", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := __.NewBrokerClient(conn)

	//延时5秒，之后再每隔2秒消费消息
	time.Sleep(time.Second * 5)

	t := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-t.C:
			response, err := c.Process(context.Background(), &__.Msg{
				MsgType: 2, //消费者从远程rpc服务端拉取消息
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			msg := new(__.Msg)
			//还原出结构体,找出id
			err = json.Unmarshal(response.O.(*__.Response_Data).Data, msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			//消费
			fmt.Println("消费者消费信息：", msg)

			_, err = c.Process(context.Background(), &__.Msg{
				Id:      msg.Id,
				MsgType: 3, //ack
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		default:
			time.Sleep(time.Second)
		}
	}
}
