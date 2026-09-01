# handler-cors

[简体中文](README.zh-CN.md) | [Documentation](docs/README.md)

HTTP/1 CORS pipeline handler for Gnalloy with preflight handling and credential-safe wildcard behavior.

This module provides Pipeline handlers. A handler observes, transforms, rejects, delays, records, or protects messages after a Channel already exists. It does not own listening sockets or application protocols unless explicitly named.

## Status

- Import path: `gnalloy.org/handler-cors`
- Repository: `github.com/gnalloy/handler-cors`
- Default branch: `dev`
- Preview install: `go get gnalloy.org/handler-cors@dev`
- License: Apache-2.0

## Install
```bash
go get gnalloy.org/handler-cors@dev
go doc gnalloy.org/handler-cors
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
```

## Documentation
- [Overview](docs/overview.md) ([中文](docs/overview.zh-CN.md))
- [Usage](docs/usage.md) ([中文](docs/usage.zh-CN.md))
- [Examples](docs/examples.md) ([中文](docs/examples.zh-CN.md))
- [Configuration](docs/configuration.md) ([中文](docs/configuration.zh-CN.md))
- [Testing and Performance](docs/testing.md) ([中文](docs/testing.zh-CN.md))
- [API Reference](docs/api.md) ([中文](docs/api.zh-CN.md))
- [Notes and Caveats](docs/notes.md) ([中文](docs/notes.zh-CN.md))
- [ADR-001 Module Boundary](docs/decisions/0001-module-boundary.md) ([中文](docs/decisions/0001-module-boundary.zh-CN.md))

## Module Boundary

This repository owns: HTTP/1 CORS pipeline handler for Gnalloy with preflight handling and credential-safe wildcard behavior.

It does not absorb neighboring module responsibilities. Core primitives stay in `gnalloy.org/gnalloy`; protocol codecs, transports, handlers, resolvers, examples, and benchmarks stay in their own repositories.

## Packages
- `gnalloy.org/handler-cors` (`cors`)

## Gnalloy Dependencies
- `gnalloy.org/gnalloy`
- `gnalloy.org/codec-http1`

## Common Integration Pattern
- Handler constructors usually carry the policy: limits, timeouts, match rules, logging level, recorder, or traffic budget.
- Handler order matters. Place protocol decoders before handlers that inspect protocol objects, and place outbound encoders after handlers that write protocol objects.
- Handlers must keep backpressure and message ownership explicit; never retain ByteBuf values without clear lifetime control.

## Current Public Entry Points

The generated API reference lists the full public surface. Common constructors or option types currently include:
- `var ErrInvalidConfig = errors.New("gnalloy/handler/cors: invalid config") ...`
- `type Config struct{ ... }`

## Verification

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
GOWORK=off GOTOOLCHAIN=local go test ./... -run '^$' -bench . -benchmem -count=1
```

For pressure tests, assemble this module with the relevant transport, codec, and handler stack and run the scenario from `gnalloy.org/benchmarks` or `gnalloy.org/examples`. Keep host, operating system, payload, concurrency, warmup, and repetitions in the report.

## Caveats
- This repository is intentionally narrow. Cross-module behavior should be assembled in applications, recipes, examples, or benchmark harnesses.
- Public APIs should remain Go-native and explicit; avoid runtime scanning, hidden global registries, and reflection-heavy behavior in hot paths.
- Treat network input as untrusted. Configure parser limits and return typed errors instead of panics.
- Keep benchmark claims tied to a concrete host, operating system, protocol, payload, concurrency, warmup, and repetition count.
- Handler modules are reusable only when their position in the pipeline is correct for the messages they observe.
