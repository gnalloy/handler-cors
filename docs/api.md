# API Reference

[简体中文](api.zh-CN.md) | [Docs Index](README.md)

This inventory is generated from `go doc -short` for the packages in this repository. It is a quick public-surface map; source files and tests remain the authority for exact semantics.

## Packages

### `gnalloy.org/handler-cors`

Package name: `cors`

```text
var ErrInvalidConfig = errors.New("gnalloy/handler/cors: invalid config") ...
type Config struct{ ... }
type Handler struct{ ... }
    func NewHandler(cfg Config) (*Handler, error)
```
