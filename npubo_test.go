package npubo_test

import (
	"fmt"
	"npubo"
	"testing"
)

func BenchmarkSub(b *testing.B) {
	n := npubo.New[bool]()

	for i := 0; i < b.N; i++ {
		n.Subscribe("sub_more/*", func(c *npubo.Context[bool]) error {
			return nil
		})
	}

}

func BenchmarkPub(b *testing.B) {
	n := npubo.New[bool]()

	n.Subscribe("sub_more/*", func(c *npubo.Context[bool]) error {
		return nil
	})

	for i := 0; i < b.N; i++ {
		n.Publish(fmt.Sprintf("sub_more/%d", i), true)
	}
}
