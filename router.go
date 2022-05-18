package npubo

import (
	"strings"
)

type Node struct {
	field   *Node  // 上级节点
	part    string // 当前节点路径
	handler []HandlerFunc
	node    map[string]*Node
}

type HandlerFunc func(*Context) error

func newNode(field *Node) *Node {
	return &Node{
		field:   field,
		part:    "",
		handler: []HandlerFunc{},
		node:    make(map[string]*Node),
	}
}

func (p *Node) addNode(topic string, h HandlerFunc) (*Node, int) {
	parts := strings.Split(topic, "/")

	for i := 0; i < len(parts); i++ {
		t := parts[i]

		// 路由结构是否存在
		if _, ok := p.node[t]; !ok {
			p.node[t] = newNode(p)
		}

		// 指向
		p = p.node[t]
		p.part = parts[i]

		if t == "*" {
			break
		}
	}

	// 添加请求
	p.handler = append(p.handler, h)
	return p, len(p.handler) - 1
}

func (p *Node) getRoute(topic string) []HandlerFunc {
	parts := strings.Split(topic, "/")
	handler := []HandlerFunc{}

	for i := 0; i < len(parts); i++ {
		t := parts[i]

		// 查找以及通配
		if v, ok := p.node["*"]; ok {
			handler = append(handler, v.handler...)
		}

		if _, ok := p.node[t]; !ok {
			break
		}

		p = p.node[t]
		handler = append(handler, p.handler...)
	}

	return handler
}

func (p *Node) removeNode(id int) {
	if len(p.handler) > id {
		p.handler = append(p.handler[:id], p.handler[id+1:]...)
	}
}
