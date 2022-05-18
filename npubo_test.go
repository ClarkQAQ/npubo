package npubo_test

import (
	"fmt"
	"npubo"
	"testing"
)

func BenchmarkSub(b *testing.B) {
	n := npubo.New()

	for i := 0; i < b.N; i++ {
		n.Subscribe("sub_more/*id", func(c *npubo.Context) error {
			return nil
		})
	}

}

func BenchmarkPub(b *testing.B) {
	n := npubo.New()

	n.Subscribe("sub_more/*id", func(c *npubo.Context) error {
		return nil
	})

	for i := 0; i < b.N; i++ {
		n.NepoPublish(fmt.Sprintf("sub_more/%d", i), "1")
	}
}
