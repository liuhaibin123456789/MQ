package queue

import "errors"

/*
	封装队列操作
*/

type Queue struct {
	//队列初始化长度
	MaxSize int
	//队列头指针
	Front int
	//队列尾指针
	Rear int
	//队列元素
	Array []interface{}
}

// NewQueue 初始化队列模型
func NewQueue(maxSize int) *Queue {
	return &Queue{
		MaxSize: maxSize,
		Front:   -1,
		Rear:    -1,
		Array:   make([]interface{}, 0, maxSize),
	}
}

// Add 队尾添加队列元素
func (q *Queue) Add(value interface{}) (err error) {
	//队列已满
	if q.Rear == q.MaxSize-1 {
		err = errors.New("the queue is full")
		return
	}
	//队列未满
	q.Rear++
	q.Array[q.Rear] = value
	return
}

// Get 队头取出队列元素
func (q *Queue) Get() (value interface{}, err error) {
	//队列元素为空
	if q.Front == q.Rear {
		err = errors.New("the queue is empty")
		return
	}
	//取队列元素
	q.Front++
	value = q.Array[q.Front]
	return
}

//// GetCopy 取出队列元素的副本，并未改变队列结构
//func (q *Queue) GetCopy() (value interface{}, err error) {
//	//队列元素为空
//	if q.Front == q.Rear {
//		err = errors.New("the queue is empty")
//		return
//	}
//	//取队列元素
//	//q.Front++
//	value = q.Array[q.Front]
//	return
//}

//// GetByIndex 取出某个位置的队列元素
//func (q *Queue) GetByIndex() {
//
//}
