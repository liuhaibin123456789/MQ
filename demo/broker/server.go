package broker

import (
	"MQ/demo/broker/proto"
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type brokerService struct {
	//内嵌接口，保证方法实现
	__.__

	//内嵌锁，实现并发控制
	brokerLock sync.RWMutex

	//消息队列
	msgList *list.List

	//队列长度
	len int
}

//维护一个全局标记
var en *list.Element

// Process broker处理消息收发服务：将收到的消息暂存至队列里（消息堆积只提供一个50容量队列）生产者与消费者何时
func (s *brokerService) Process(context context.Context, msg *__.Msg) (res *__.Response, err error) {
	res = new(__.Response)

	//初始化队列
	s.len = s.msgList.Len()

	switch msg.MsgType {

	//生产者发送消息
	case 1:
		if e := s.msgList.PushBack(msg); e == nil {
			res.Ok = false
			err = errors.New("insert element into list wrong")
			return
		}
		s.len = s.msgList.Len()
		fmt.Println(msg)
		res.Ok = true
		res.O = &__.Response_Data{Data: []byte("消息成功放入消息队列")}

	//消费者消费消息
	case 2:
		if en == nil {
			en = s.msgList.Front()
		}
		e := en
		//将下一次的队列元素存起
		en = e.Next()
		//将消息从消息队列里取出来，但是没有删除
		if e == nil {
			err = errors.New("the list is empty")
			return
		}

		//将消息写入返回结果
		bytes, err1 := toBytes(e.Value)
		if err1 != nil {
			err = err1
			return
		}

		res.O = &__.Response_Data{Data: bytes}
		return

	//ack删除
	case 3:
		s.brokerLock.Lock()
		//消费者确认消费到消息后，再将消息的某个id发送过来，同时消息的类型标记为3
		for e := s.msgList.Front(); e != nil; e = e.Next() {
			if e.Value.(*__.Msg).Id == msg.Id {
				s.msgList.Remove(e)
				s.len = s.msgList.Len()
				break
			}
		}
		res.Ok = true
		res.O = &__.Response_Data{Data: []byte("ack删除成功")}
		s.brokerLock.Unlock()
	default:
		err = errors.New("the message type is wrong")
	}
	return
}

func toBytes(msg interface{}) (msgString []byte, err error) {
	m := msg.(*__.Msg)
	msgString, err = json.Marshal(m)
	return
}

func main() {
	conn, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println(err)
		return
	}

	s := grpc.NewServer()

	__.RegisterBrokerServer(s, &brokerService{msgList: list.New()})

	if err := s.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
