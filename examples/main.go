package main

import (
	"fmt"
	"npubo"
)

func main() {
	n := npubo.New()

	n.Subscribe("/user/*", func(c *npubo.Context) error {
		fmt.Printf("topic: %s, data: %s\n",
			c.Topic(), c.String())
		return nil
	})

	n.Publish("/user/1231/dwdw", "qaq")
	n.Publish("/user/123", n)
}
