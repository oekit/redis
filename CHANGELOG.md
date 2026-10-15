# Changelog — oekit/redis

由 ekit v0.34.0（git.enjoye.top/enjoydream/ekit，冻结维护）拆分而来。所有 notable 变更记录如下。

## [v0.3.0] — 2026-10-16

新增 `ratelimit/` 分布式限流子包（Lua 脚本与三态返回语义忠实移植 go-zero core/limit）。

### 新增
- `ratelimit.PeriodLimit`——Redis 固定窗口限流（INCRBY+EXPIRE 原子脚本）：
  `Take(ctx, key)` 返回三态 `Allowed`/`HitQuota`/`OverQuota`（Redis 故障返回
  `Unknown`+err，不消耗配额）；`WithPeriodAlign()` 对齐自然周期（如"每日 5 条"）；
  **不做本地兜底**——固定窗口计数无法本地近似，fail-open/closed 由调用方决定
- `ratelimit.TokenLimiter`——Redis 令牌桶限流（双 key 惰性回补，Cluster
  hash-tag 同 slot）：`Allow(ctx)`/`AllowN(ctx, n)`；`WithTokenFallback()` 开启
  Redis 故障时进程内 golang.org/x/time/rate 本地兜底 + 后台探活自动回切
  （`WithTokenPingInterval` 默认 100ms，探活超时 `WithTokenProbeTimeout`
  默认 1s、与间隔解耦）；未开启时故障 fail-closed
- 存储抽象收 go-redis v9 `redis.Scripter` 接口（Client/ClusterClient/Ring 均可用）；
  全部 ctx-only API；后台 goroutine 经 `concurrency/async.GoSafe` recover

### 依赖
- 新增直接依赖 `golang.org/x/time v0.16.0`（TokenLimiter 本地兜底，矩阵内与
  concurrency 同版本）；其余不变

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
