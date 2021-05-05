package npubo_test

import (
	"errors"
	"fmt"
	"npubo"
	"testing"
	"time"
)

var pub *npubo.Publisher = npubo.NewPublisher(500)

func TestSub(t *testing.T) {
	pub.Subscribe("sub_one/one", "QwQ", func(sub *npubo.Subscriber, val interface{}) error {
		fmt.Println("sub", sub, " message", val)
		return nil
	})

	pub.Subscribe("sub_one/timeout", "QwQ", func(sub *npubo.Subscriber, val interface{}) error {
		time.Sleep(time.Second)
		return nil
	})

	pub.Subscribe("sub_one/error", "QwQ", func(sub *npubo.Subscriber, val interface{}) error {
		return errors.New("a error")
	})
}

func TestPub(t *testing.T) {
	pub.Publish("sub_one/one", "QwQ", func(sub *npubo.Subscriber, e error) {
		fmt.Println("sub", sub, " error", e)
	})

	pub.Publish("sub_one/timeout", "Message", func(sub *npubo.Subscriber, e error) {
		fmt.Println("sub", sub, " error", e)
	})

	pub.Publish("sub_one/error", "Message", func(sub *npubo.Subscriber, e error) {
		fmt.Println("sub", sub, " error", e)
	})
}

func BenchmarkSub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pub.Subscribe(fmt.Sprintf("sub_more/%v", i), "QwQ", func(sub *npubo.Subscriber, val interface{}) error {
			return nil
		})
	}
}

func BenchmarkPub(b *testing.B) {
	pub.Publish("sub_more/*", "Message", func(sub *npubo.Subscriber, e error) {
		return
	})
}
