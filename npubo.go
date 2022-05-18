package npubo

import (
	"fmt"
	"sync"
)

type Npubo struct {
	node   *Node
	rwLock *sync.RWMutex
}

type Subscriber struct {
	npubo *Npubo
	n     *Node
	id    int
}

func New() *Npubo {
	return &Npubo{
		node:   newNode(nil),
		rwLock: &sync.RWMutex{},
	}
}

func (n *Npubo) Subscribe(topic string, h HandlerFunc) *Subscriber {
	n.rwLock.Lock()
	defer n.rwLock.Unlock()

	node, id := n.node.addNode(topic, h)
	return &Subscriber{n, node, id}
}

func (n *Subscriber) Unsubscribe() {
	n.npubo.rwLock.Lock()
	defer n.npubo.rwLock.Unlock()

	n.n.removeNode(n.id)
}

func (n *Npubo) Publish(topic string, data interface{}) error {
	n.rwLock.RLock()
	defer n.rwLock.RUnlock()

	h := n.node.getRoute(topic)
	if h == nil {
		return fmt.Errorf("no handler for topic: %s", topic)
	}

	c := newContext(topic, data)

	for i := 0; i < len(h); i++ {
		if e := func() (e error) {
			defer func() {
				if err := recover(); err != nil {
					e = fmt.Errorf("%v", err)
				}
			}()

			return h[i](c)
		}(); e != nil {
			return e
		}
	}

	return nil
}

func (n *Npubo) NepoPublish(topic string, data interface{}) {
	n.rwLock.RLock()
	defer n.rwLock.RUnlock()

	h := n.node.getRoute(topic)
	if h == nil {
		return
	}

	c := newContext(topic, data)

	for i := 0; i < len(h); i++ {
		func() {
			defer func() {
				recover()
			}()

			h[i](c)
		}()
	}
}
