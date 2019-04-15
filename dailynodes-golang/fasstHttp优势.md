### Http和fasstHttp的区别

fasthttp 的性能可以达到Http标准库的 10 倍，他魔性的实现方式,主要的点在于四个方面：

- `net/http` 的实现是一个连接新建一个 goroutine；`fasthttp` 是利用一个 worker 复用 goroutine，减轻 runtime 调度 goroutine 的压力
- `net/http` 解析的请求数据很多放在 `map[string]string`(http.Header) 或 `map[string][]string`(http.Request.Form)，有不必要的 []byte 到 string 的转换，是可以规避的
- `net/http` 解析 HTTP 请求每次生成新的 `*http.Request` 和 `http.ResponseWriter`; `fasthttp` 解析 HTTP 数据到 `*fasthttp.RequestCtx`，然后使用 `sync.Pool` 复用结构实例，减少对象的数量
- `fasthttp` 会延迟解析 HTTP 请求中的数据，尤其是 Body 部分。这样节省了很多不直接操作 Body 的情况的消耗但是因为 `fasthttp` 的实现与标准库差距较大，所以 API 的设计完全不同。使用时既需要理解 HTTP 的处理过程，又需要注意和标准库的差别。



### fasthttp 的不足

两个比较大的不足：

- HTTP/2.0 不支持
- WebSocket 不支持

严格来说 Websocket 通过 `Hijack()` 是可以支持的，但是 `fasthttp` 想自己提供直接操作的 API。那还需要等待开发。

### 总结

比较标准库的粗犷，`fasthttp` 有更精细的设计，对 Go 网络并发编程的主要痛点做了很多工作，达到了很好的效果。目前，[iris](https://github.com/kataras/iris) 和 [echo](https://github.com/labstack/echo) 支持 `fasthttp`，性能上和使用 `net/http` 的别的 Web 框架对比有明显的优势。如果选择 Web 框架，支持 `fasthttp` 可以看作是一个真好的卖点。

