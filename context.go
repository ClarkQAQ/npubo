package npubo

type Context[T any] struct {
	topic string
	data  T
}

func newContext[T any](topic string, data T) *Context[T] {
	return &Context[T]{
		topic: topic,
		data:  data,
	}
}

// 获取发布的主题
// @return topic string 发布数据的主题路径
func (c *Context[T]) Topic() string {
	return c.topic
}

// 获取发布的数据
// @return data T 发布的数据
func (c *Context[T]) Data() T {
	return c.data
}
