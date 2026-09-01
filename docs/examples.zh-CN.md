# 案例

[English](examples.md) | [文档索引](README.zh-CN.md)

## 案例 1：将模块加入应用

```bash
mkdir gnalloy-app && cd gnalloy-app
go mod init example.com/gnalloy-app
go get gnalloy.org/handler-cors@dev
go doc gnalloy.org/handler-cors
```

## 案例 2：查看当前包

当前源码树暴露这些 package 导入路径：
- `gnalloy.org/handler-cors`

按需要的行为对对应 package 执行 `go doc`：

```bash
go doc gnalloy.org/handler-cors
```

精选当前导出入口：
- `var ErrInvalidConfig = errors.New("gnalloy/handler/cors: invalid config") ...`
- `type Config struct{ ... }`
- `type Handler struct{ ... }`

## 案例 3：将可执行测试作为行为示例

仓库测试是受支持行为的可执行示例。先从下面的精选名称开始，再阅读对应 `_test.go` 文件中的完整 setup 和断言。完整发现列表见 [测试与性能](testing.zh-CN.md)。

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run Test -count=1
```

精选当前 test、benchmark、fuzz 与 example 入口：
- `TestHandlerAddsHeadersToSimpleResponse`
- `TestHandlerAnswersPreflight`
- `TestHandlerKeepsResponseOrderForMixedOrigins`
- `TestHandlerShortCircuitsForbiddenOrigin`
- `TestNewHandlerRejectsMissingOrigins`

## 案例 4：跨模块装配

本模块的直接 Gnalloy 依赖：
- `gnalloy.org/codec-http1`
- `gnalloy.org/gnalloy`

装配说明：
- handler 在 Gnalloy pipeline 内提供策略、生命周期、防护、日志、指标或流量行为。
- 协议解析留给 codec 模块，socket 所有权留给 transport 模块。
- handler 的位置应匹配所选协议栈的 inbound/outbound 语义。

## 案例 5：压测 Harness

持续负载测试时，如果该模块参与网络流量路径，将它接入 `gnalloy.org/benchmarks` 的场景，或接入 `gnalloy.org/examples` 的可运行客户端。报告中记录 host、OS、CPU、Go version、protocol、payload、concurrency、warmup、repetitions、throughput 和 p99 latency。
