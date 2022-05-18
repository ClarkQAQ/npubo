package npubo

import (
	"fmt"
	"reflect"
)

type Context struct {
	topic string
	data  interface{}
}

func newContext(topic string, data interface{}) *Context {
	return &Context{
		topic: topic,
		data:  data,
	}
}

func (c *Context) Topic() string {
	return c.topic
}

func (c *Context) Data() interface{} {
	return c.data
}

func (c *Context) Int64() int64 {
	if v, ok := c.data.(int64); ok {
		return v
	}

	return 0
}

func (c *Context) String() string {
	if v, ok := c.data.(string); ok {
		return v
	}

	return fmt.Sprint(c.data)
}

func (c *Context) Type() string {
	return reflect.TypeOf(c.data).Name()
}
