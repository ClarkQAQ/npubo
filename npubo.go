package npubo

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type (
	// 订阅消息回调
	Call func(sub *Subscriber, val interface{}) error

	// 推送错误返回回调
	ErrCall func(sub *Subscriber, e error)

	// 前缀输节点
	Node struct {
		Calls    map[string]*Call
		NextNode map[string]*Node
	}

	// 推送结构
	Publisher struct {
		Node       *Node
		timeout    int
		rwLock     *sync.RWMutex
		workerLock *sync.Mutex
	}

	// 订阅结构
	Subscriber struct {
		Node      *Node      // 所在节点
		Publisher *Publisher // 所在推送
		Topic     string     // 订阅路径
		CId       string     // 客户端Id
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

func NewPublisher(timeout int) *Publisher {
	return &Publisher{
		Node:       newNode(),
		timeout:    timeout,
		rwLock:     &sync.RWMutex{},
		workerLock: &sync.Mutex{},
	}
}

// 订阅
func (that *Publisher) Subscribe(topic string, c_id string, call Call) (*Subscriber, error) {
	that.rwLock.Lock()
	defer that.rwLock.Unlock()

	nowNode := that.Node.NextNode
	cals := strings.Split(topic, "/")
	sub := &Subscriber{
		Topic: topic,
		CId:   c_id,
		Node:  nil,
	}

	if strings.Contains(topic, "*") || cals[0] == "" {
		return sub, errors.New("invaild topic")
	}

	for i, v := range cals {
		if _, ok := nowNode[v]; !ok {
			nowNode[v] = newNode()
		}
		if i == len(cals)-1 {
			nowNode[v].Calls[c_id] = &call
			sub.Node = nowNode[v]
		}
		nowNode = nowNode[v].NextNode
	}
	return sub, nil
}

// 发布消息
func (that *Publisher) Publish(topic string, val interface{}, errBack ErrCall) {
	that.rwLock.RLock()
	defer that.rwLock.RUnlock()

	nowNode := that.Node.NextNode
	cals := strings.Split(topic, "/")

	topicRecode := []string{}
	for i, v := range cals {
		if _, ok := nowNode[v]; ok { // 正常匹配路径
			topicRecode = append(topicRecode, v)
			if i == len(cals)-1 {
				for c_id, call := range nowNode[v].Calls {
					sub := &Subscriber{
						Topic:     strings.Join(topicRecode, "/"),
						CId:       c_id,
						Node:      nowNode[v],
						Publisher: that,
						isWorker:  true,
					}
					e := that.callSubscriber(call, sub, val)
					if e != nil && errBack != nil {
						that.callErrBack(errBack, sub, e)
					}
				}
			}
			nowNode = nowNode[v].NextNode
		} else if v == "*" { // 匹配通配符
			that.callAllNode(nowNode, errBack, topicRecode, val)
			break
		} else {
			break
		}
	}
}

func (that *Publisher) callAllNode(nowNode map[string]*Node, errBack ErrCall, topicInit []string, val interface{}) {
	var wg sync.WaitGroup
	for t, n := range nowNode {
		topic := topicInit
		topic = append(topic, t)
		for c_id, call := range n.Calls {
			sub := &Subscriber{
				Topic:     strings.Join(topic, "/"),
				CId:       c_id,
				Node:      n,
				Publisher: that,
				isWorker:  true,
			}
			e := that.callSubscriber(call, sub, val)
			if e != nil && errBack != nil {
				that.callErrBack(errBack, sub, e)
			}
		}

		wg.Add(1)
		go func(nextNode map[string]*Node, topic []string, val interface{}) {
			defer wg.Done()
			that.callAllNode(nextNode, errBack, topic, val)
		}(n.NextNode, topic, val)
	}
	wg.Wait()
}

func (that *Publisher) callSubscriber(call *Call, sub *Subscriber, val interface{}) (e error) {
	done := make(chan byte, 0)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				e = errors.New(fmt.Sprint(r))
			}
		}()
		e = (*call)(sub, val)
		done <- 0
	}()
	select {
	case <-done:
		return e
	case <-time.After(time.Microsecond * time.Duration(that.timeout)):
		close(done)
		return errors.New("subscriber timeout")
	}
}

func (that *Publisher) callErrBack(errBack ErrCall, sub *Subscriber, e error) {
	done := make(chan byte, 0)
	go func() {
		defer func() { recover() }()
		errBack(sub, e)
		done <- 0
	}()
	select {
	case <-done:
		return
	case <-time.After(time.Microsecond * time.Duration(that.timeout)):
		close(done)
		return
	}
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
		return errors.New("node is nil")
	}

	that.Node.Calls[that.CId] = &call
	return nil
}

// 关闭实例
func (that *Publisher) Close() {
	that.rwLock.RLock()
	defer that.rwLock.RUnlock()
	that.Node = newNode()
}
