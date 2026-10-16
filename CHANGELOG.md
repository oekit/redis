# Changelog — oekit/redis

由 ekit v0.34.0（git.enjoye.top/enjoydream/ekit，冻结维护）拆分而来。所有 notable 变更记录如下。

## [v0.4.0] — 2026-10-16

RedisCache 后端契约化：接口断言 + cachetest 契约套件接入 + null 哨兵碰撞修复。

### Breaking（数据格式）
- **null 哨兵变更**：`__NULL__` → `"\x00oekit:null"`（以 `\x00` 开头的值保留为
  null 标记专用前缀）。旧哨兵可被 `SetRaw` 合法写入精确碰撞（写入字面
  `__NULL__` 即被永久误判为"数据不存在"，GetOrLoad 不再回源）；新哨兵带
  `\x00` 前缀，合法 JSON 序列化输出不会自然撞上。**存量影响**：旧缓存中由
  GetOrLoad 写入的 `__NULL__` 标记不再被识别——读回按普通值反序列化（默认
  JSON 下报错），该 key 的防穿透短暂失效，标记自身按 NullTTL（默认 1 分钟）
  自然过期后恢复；零已知消费方，无实际影响面。走 JSON 序列化的 Set/GetOrLoad
  正常值不受影响（json.Marshal 字符串带引号，从未与裸哨兵精确相等）
- `SetRaw` 新增契约约束：data 不得以 `\x00` 开头（保留为 null 标记前缀；
  精确等于哨兵的字节按 null marker 处理，其余 `\x00` 开头值不认作 null）

### 新增
- 编译期接口断言：`RedisCache` 显式满足 `data/cache` 门面 `Cache[T]` /
  `CacheAside[T]`（int/string 双类型实例化）
- 接入 `github.com/oekit/data/cachetest` 契约测试套件（8 场景 × string/
  结构体双值类型，每场景独立 miniredis 并开启 `ClockTTL` 对齐真实时间过期
  语义，`-race` 全绿）——通用缓存语义不再在本包重复编写

### 修复
- `DelBatch` 空 keys 由报错（Redis 裸 `DEL` 要求至少 1 个 key，返回
  wrong number of arguments）改为幂等 no-op——契约套件 `del_and_delbatch`
  场景（空 DelBatch 不应报错）发现的违约，与 `Del` 幂等语义对齐

### 变更
- import 路径正名：`github.com/oekit/data/data/cache` → `github.com/oekit/data/cache`
  （跟随 data v0.4.0 目录正名）

### 依赖
- `github.com/oekit/data` v0.3.3 → v0.4.0（v0.4.0 发布前本地经 go.work
  replace 解析）

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
