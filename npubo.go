package npubo

import (
	"fmt"
	"sync"
)

type Npubo[T any] struct {
	node   *Node[T]
	rwLock *sync.RWMutex
}

// 创建一个新的 Npubo, 支持泛型
// .e.g xxx := New[string]()
func New[T any]() *Npubo[T] {
	return &Npubo[T]{
		node:   newNode[T](nil),
		rwLock: &sync.RWMutex{},
	}
}

// Subscribe
// @param topic string 主题路径, '/' 分隔, '*' 号代表通配
// @param h HandlerFunc[T] 回调函数
func (n *Npubo[T]) Subscribe(topic string, h HandlerFunc[T]) (cancel func()) {
	n.rwLock.Lock()
	defer n.rwLock.Unlock()

	node, id := n.node.addNode(topic, h)
	return func() {
		n.rwLock.Lock()
		defer n.rwLock.Unlock()

		node.removeNode(id)
	}
}

// Publish
// @param topic string 主题路径, '/' 分隔
// @param data T 数据
func (n *Npubo[T]) Publish(topic string, data T) error {
	n.rwLock.RLock()
	defer n.rwLock.RUnlock()

	h := n.node.getRoute(topic)
	if h == nil {
		return fmt.Errorf("no handler for topic: %s", topic)
	}

	c := newContext(topic, data)

	for i := 0; i < len(h); i++ {
		if e := h[i](c); e != nil {
			return e
		}
	}

	return nil
}
