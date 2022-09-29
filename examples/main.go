package main

import (
	"fmt"
	"npubo"
)

func main() {
	// 支持泛型, 这里的泛型是 string
	n := npubo.New[string]()

	// 哪怕手抖多打了一个 /，也不会影响订阅以及发布
	// 但是如果你订阅了 sub_more/* ，那么发布的时候就必须是 sub_more/xxx 这样的格式
	n.Subscribe("/sub_more//zhiyin/*", func(c *npubo.Context[string]) error {
		fmt.Printf("topic: %s, data: %s\n",
			c.Topic(), c.Data())
		return nil
	})

	// subscribe
	n.Publish("/sub_more/zhiyin/nitaimei/music", "唱跳rap篮球") // 这段是 Github Copilot 自动补全的, 笑死hhhhh
	n.Publish("/sub_more/zhiyin//cxk", "Music~")
}
