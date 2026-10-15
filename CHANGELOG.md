# Changelog — oekit/redis

由 ekit v0.34.0（git.enjoye.top/enjoydream/ekit，冻结维护）拆分而来。所有 notable 变更记录如下。

## [v0.2.0] — 2026-10-15

可观测性统一轮。

### Breaking
- `redisx.NewClient` / `NewClientFromAddr` / `NewClientFromConfig` **不再自动注入 OTel tracing**
  （原先强制 `redisotel.InstrumentTracing`，无出口）；go.mod 不再依赖 go-redis extra

### 迁移
- 需要 Redis 命令 span 的服务在构造后一行显式接入（无 SDK 时 noop）：
  `redisotel.InstrumentTracing(c.Client())`（`github.com/oekit/tracing/redisotel`）；
  不需要的无需任何改动

## [v0.1.5] — 2026-10-15

- v0.1.5

## [v0.1.4] — 2026-10-15

- v0.1.4

## [v0.1.3] — 2026-10-15

- v0.1.3

## [v0.1.2] — 2026-10-15

- v0.1.2

## [v0.1.1] — 2026-10-15

- v0.1.1

## [v0.1.0] — 2026-10-15

- v0.1.0
