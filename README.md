### Npubo (脸滚键盘)

> 基于前缀树的发布/订阅

Topic抄的MQTT的方法用路径来匹配
性能....就那样....

GoTest:

```raw
[clark@ArchLinux npubo]$ go test -benchmem -bench .
sub &{0xc00004c4d0 0xc000026040 sub_one/one QwQ true}  message QwQ
sub &{0xc00004c4e0 0xc000026040 sub_one/timeout QwQ true}  error subscriber timeout
sub &{0xc00004c4f0 0xc000026040 sub_one/error QwQ true}  error a error
goos: linux
goarch: amd64
pkg: npubo
cpu: Intel(R) Core(TM) i5-7300HQ CPU @ 2.50GHz
BenchmarkSub-4            754231              1427 ns/op             527 B/op          8 allocs/op
BenchmarkPub-4                 1        1780671142 ns/op        332174032 B/op   6788814 allocs/op
PASS
ok      npubo   3.190s
```

示例:

```go

    // 初始化
	pub := npubo.NewPublisher(500)

    // 订阅
	pub.Subscribe("sub/1", "QwQ", func(sub *Subscriber, val interface{}) error {
        fmt.Println(sub, val)
		return nil
	})
	
    // 发布  topic 支持星号通配 ("sub/*", "*")
	pub.Publish("sub/1", "消息", func(sub *Subscriber, e error) {
		fmt.Println(sub, e)
	})

    pub.Close()

```