package npubo

type Node[T any] struct {
	field   *Node[T] // 上级节点
	handler []HandlerFunc[T]
	node    map[rune]*Node[T]
}

type HandlerFunc[T any] func(*Context[T]) error

func newNode[T any](field *Node[T]) *Node[T] {
	return &Node[T]{
		field:   field,
		handler: []HandlerFunc[T]{},
		node:    make(map[rune]*Node[T]),
	}
}

func (p *Node[T]) addNode(topic string, h HandlerFunc[T]) (*Node[T], int) {
	parts := []rune(topic)

	for i := 0; i < len(parts); i++ {
		t := parts[i]
		if t == '/' {
			continue
		}

		// 路由结构是否存在
		if _, ok := p.node[t]; !ok {
			p.node[t] = newNode(p)
		}

		// 指向
		p = p.node[t]

		if t == '*' {
			break
		}
	}

	// 添加请求
	p.handler = append(p.handler, h)
	return p, len(p.handler) - 1
}

func (p *Node[T]) getRoute(topic string) []HandlerFunc[T] {
	parts := []rune(topic)
	handler := []HandlerFunc[T]{}

	for i := 0; i < len(parts); i++ {
		t := parts[i]
		if t == '/' {
			continue
		}

		// 查找以及通配
		if v, ok := p.node['*']; ok {
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

func (p *Node[T]) removeNode(id int) {
	if len(p.handler) > id {
		p.handler = append(p.handler[:id], p.handler[id+1:]...)
	}
}
