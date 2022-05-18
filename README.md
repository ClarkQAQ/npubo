### Npubo (脸滚键盘)

> 基于前缀树的发布/订阅

性能....就那样....

GoTest:

```raw
[clark@ArchOwO npubo]$ go test -benchmem -bench .
goos: linux
goarch: amd64
pkg: npubo
cpu: AMD Ryzen 7 5800H with Radeon Graphics         
BenchmarkSub-16          8911771               138.6 ns/op           100 B/op          2 allocs/op
BenchmarkPub-16          4401676               242.6 ns/op            88 B/op          4 allocs/op
PASS
ok      npubo   2.726s
```

示例:

```go

    // 初始化
	n := npubo.New()

	n.Subscribe("/user/*", func(c *npubo.Context) error {
		fmt.Printf("topic: %s, data: %s\n",
			c.Topic(), c.String())
		return nil
	})

	n.Publish("/user/1231/dwdw", "qaq")
	n.Publish("/user/123", n)

```