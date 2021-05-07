package npubo

import "strings"

// 订阅
func (that *Publisher) Subscribe(topic string, c_id string, call Call) (*Subscriber, error) {
	that.rwLock.Lock()
	defer that.rwLock.Unlock()

	nowNode := that.Node.NextNode
	cals := strings.Split(topic, "/")
	sub := &Subscriber{
		Topic:     topic,
		Publisher: that,
		CId:       c_id,
		Node:      nil,
		C:         make(chan ChanCall),
	}

	if strings.Contains(topic, "*") || cals[0] == "" {
		return sub, ErrInvaildTopic
	}

	for i, v := range cals {
		if _, ok := nowNode[v]; !ok {
			nowNode[v] = newNode()
		}
		if i == len(cals)-1 {
			nowNode[v].Calls[c_id] = &call

			if that.openChan {
				chanCallBack := func(subData *Subscriber, val interface{}) error {
					sub.C <- ChanCall{
						Subscriber: subData,
						Content:    val,
					}
					return nil
				}
				nowNode[v].Calls["__chan__"+c_id] = (*Call)(&chanCallBack)
			}

			sub.Node = nowNode[v]
		}
		nowNode = nowNode[v].NextNode
	}
	return sub, nil
}

func (that *Publisher) RootSubscribe(c_id string, call Call) *Subscriber {
	that.rwLock.Lock()
	defer that.rwLock.Unlock()

	that.Root.Calls[c_id] = &call
	return &Subscriber{
		Topic: "*",
		CId:   c_id,
		Node:  that.Root,
	}
}
