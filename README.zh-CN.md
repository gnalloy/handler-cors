# handler-cors

Gnalloy HTTP/1 CORS Pipeline Handler，支持 preflight 与 credential-safe wildcard 行为。

本仓库属于 Gnalloy 模块化网络栈。默认分支是 `dev`；本次初始化不创建 release，不打 tag。

## 安装

```bash
go get gnalloy.org/handler-cors@dev
```

## 模块边界

- 模块路径：`gnalloy.org/handler-cors`
- 职责：Gnalloy HTTP/1 CORS Pipeline Handler，支持 preflight 与 credential-safe wildcard 行为
- 核心依赖：需要 Gnalloy ByteBuf、Channel、EventLoop 或 Bootstrap 契约时依赖 `gnalloy.org/gnalloy`。

## Gnalloy 依赖

- `gnalloy.org/gnalloy`
- `gnalloy.org/codec-http1`

## 开发验证

```bash
go test ./... -count=1
go vet ./...
go test ./... -run '^$' -bench . -benchmem -benchtime=100ms -count=1
```

跨仓库开发使用 `G:\opensource\gnalloy\go.work`。独立模块复验时设置 `GOWORK=off`。

## 许可证

Apache-2.0。
