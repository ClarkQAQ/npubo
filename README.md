### Npubo (脸滚键盘)

> 基于前缀树的发布/订阅


##### 性能: 

```
目前是订阅: 5256131/s, 发布: 2152751/s (AMD Ryzen 7 5800H)
优化了gc以及内存占用, 但感觉还是有很大的优化空间....
不过从数据看起来就只是在速度和消耗上面做平衡罢了
```


##### GoTest:

```raw
[clark@ArchOwO npubo]$ go test -benchmem -bench .
goos: linux
goarch: amd64
pkg: npubo
cpu: AMD Ryzen 7 5800H with Radeon Graphics         
 - BenchmarkSub-16          8911771               138.6 ns/op           100 B/op          2 allocs/op
 + BenchmarkSub-16          7450473               168.0 ns/op            90 B/op          1 allocs/op
 - BenchmarkPub-16          4401676               242.6 ns/op            88 B/op          4 allocs/op
 + BenchmarkPub-16          3306672               353.7 ns/op            56 B/op          3 allocs/op
PASS
 - ok      npubo   2.726s
 + ok      npubo   2.959s
```

示例:

```go

	// 支持泛型, 这里的泛型是 string
	n := npubo.New[string]()

	// 哪怕手抖多打了一个 /，也不会影响订阅以及发布
	// 但是如果你订阅了 sub_more/* ，那么发布的时候就必须是 sub_more/xxx 这样的格式
	// '*' 号代表通配
	n.Subscribe("/sub_more//zhiyin/*", func(c *npubo.Context[string]) error {
		fmt.Printf("topic: %s, data: %s\n",
			c.Topic(), c.Data())
		return nil
	})

	// subscribe
	n.Publish("/sub_more/zhiyin/nitaimei/music", "唱跳rap篮球") // 这段是 Github Copilot 自动补全的, 笑死hhhhh
	n.Publish("/sub_more/zhiyin//cxk", "Music~")

```