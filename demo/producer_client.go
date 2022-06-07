package main

/*
	生产者客户端
*/

import (
	"MQ/demo/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	//监听
	conn, err := grpc.Dial("127.0.0.1:12345", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewBrokerClient(conn)

	//生产者:每隔1秒，生成10条消息
	t := time.NewTicker(time.Second)
	var id = int32(1)
	for {
		select {
		case <-t.C:
			fmt.Println("生产者生成消息数:", id)
			_, err := c.Process(context.Background(), &proto.Msg{
				Id:         id,
				Topic:      "winter",
				MsgType:    1,
				MsgContent: "第" + string(id) + "消息:下雪了",
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			id++
		}
	}
}
