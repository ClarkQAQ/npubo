package npubo

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// 发布消息
func (that *Publisher) Publish(topic string, val interface{}, errBack ErrCall) {
	that.rwLock.RLock()
	defer that.rwLock.RUnlock()

	// 觉得慢可以加协程
	that.callNode(that.Root, "*", topic, errBack, val)

	nowNode := that.Node.NextNode
	cals := strings.Split(topic, "/")

	topicRecode := []string{}
	for i, v := range cals {
		if _, ok := nowNode[v]; ok { // 正常匹配路径
			topicRecode = append(topicRecode, v)
			if i == len(cals)-1 {
				that.callNode(nowNode[v], strings.Join(topicRecode, "/"), topic, errBack, val)
			}
			nowNode = nowNode[v].NextNode
		} else if v == "*" { // 匹配通配符
			that.callAllNode(nowNode, errBack, topicRecode, topic, val)
			break
		} else {
			sub := &Subscriber{
				Topic:     topic,
				CallTopic: topic,
				Publisher: that,
				isWorker:  true,
			}

			that.callErrBack(errBack, sub, ErrTopicNotFound)
			break
		}
	}
}

func (that *Publisher) callNode(node *Node, topic, callTopic string, errBack ErrCall, val interface{}) {
	for c_id, call := range node.Calls {
		sub := &Subscriber{
			Topic:     topic,
			CallTopic: callTopic,
			CId:       c_id,
			Node:      node,
			Publisher: that,
			isWorker:  true,
		}
		e := that.callSubscriber(call, sub, val)
		if e != nil && errBack != nil {
			that.callErrBack(errBack, sub, e)
		}
	}
}

func (that *Publisher) callAllNode(nowNode map[string]*Node, errBack ErrCall, topicInit []string, callTopic string, val interface{}) {
	var wg sync.WaitGroup
	for t, n := range nowNode {
		topic := topicInit
		topic = append(topic, t)

		that.callNode(n, strings.Join(topic, "/"), callTopic, errBack, val)

		wg.Add(1)
		go func(nextNode map[string]*Node, topic []string, val interface{}) {
			defer wg.Done()
			that.callAllNode(nextNode, errBack, topic, callTopic, val)
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
		return ErrSubscriberTimeout
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
