package npubo

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
)

var (
	ErrTopicNotFound     = errors.New("topic not found")
	ErrSubscriberTimeout = errors.New("subscriber timeout")
	ErrInvaildTopic      = errors.New("invaild topic")
	ErrNilNode           = errors.New("nil node ")
)

type (
	// 订阅消息回调
	Call func(sub *Subscriber, val interface{}) error

	// 推送错误返回回调
	ErrCall func(sub *Subscriber, e error)

	ChanCall struct {
		Subscriber *Subscriber
		Content    interface{}
	}

	// 前缀树节点
	Node struct {
		Calls    map[string]*Call
		NextNode map[string]*Node
	}

	// 推送结构
	Publisher struct {
		Node       *Node
		Root       *Node
		timeout    int
		openChan   bool
		rwLock     *sync.RWMutex
		workerLock *sync.Mutex
	}

	// 订阅结构
	Subscriber struct {
		Node      *Node           // 所在节点
		Publisher *Publisher      // 所在推送
		Topic     string          // 订阅路径
		CallTopic string          // 推送订阅路径
		CId       string          // 客户端Id
		C         chan (ChanCall) // 通道订阅
		isWorker  bool
	}
)

// 初始化节点
func newNode() *Node {
	return &Node{
		NextNode: make(map[string]*Node),
		Calls:    make(map[string]*Call),
	}
}

func NewPublisher(timeout int, openChan bool) *Publisher {
	return &Publisher{
		Node: newNode(),
		Root: &Node{
			Calls:    make(map[string]*Call),
			NextNode: nil,
		},
		openChan:   openChan,
		timeout:    timeout,
		rwLock:     &sync.RWMutex{},
		workerLock: &sync.Mutex{},
	}
}

func (that *Subscriber) RootEvict(c_id string, call Call) error {
	if that.isWorker {
		that.Publisher.workerLock.Lock()
		defer that.Publisher.workerLock.Unlock()
	} else {
		that.Publisher.rwLock.Lock()
		defer that.Publisher.rwLock.Unlock()
	}
	if that.Node == nil {
		return nil
	}

	delete(that.Node.Calls, that.CId)
	that.Node = nil
	return nil
}

// 取消订阅
func (that *Subscriber) Evict() error {
	if that.isWorker {
		that.Publisher.workerLock.Lock()
		defer that.Publisher.workerLock.Unlock()
	} else {
		that.Publisher.rwLock.Lock()
		defer that.Publisher.rwLock.Unlock()
	}
	if that.Node == nil {
		return nil
	}

	delete(that.Node.Calls, that.CId)
	defer func() { recover() }()
	close(that.C)
	that.Node = nil
	return nil
}

// 重写订阅函数
func (that *Subscriber) RewriteCall(call Call) error {
	if that.isWorker {
		that.Publisher.workerLock.Lock()
		defer that.Publisher.workerLock.Unlock()
	} else {
		that.Publisher.rwLock.Lock()
		defer that.Publisher.rwLock.Unlock()
	}
	if that.Node == nil {
		return ErrNilNode
	}

	that.Node.Calls[that.CId] = &call
	return nil
}

// 关闭实例
func (that *Publisher) Close() {
	that.rwLock.RLock()
	defer that.rwLock.RUnlock()
	that.Node, that.Root = nil, nil

	runtime.GC()

	that.Node, that.Root = newNode(), newNode()
	fmt.Println(that)
}
