// Package ratelimit 提供基于 Redis 的分布式限流器。
//
// 包含两种限流器，Lua 脚本与三态返回语义忠实移植 go-zero core/limit
// （periodscript.lua / tokenscript.lua），可与其混部共用同一份 Redis 配额：
//
//   - [PeriodLimit]：固定窗口计数限流（INCRBY + EXPIRE 原子脚本）。
//     Take 返回三态（[Allowed] / [HitQuota] / [OverQuota]）。
//   - [TokenLimiter]：令牌桶限流（双 key 惰性回补）。
//     Allow/AllowN 返回 bool；可选进程内兜底（[WithTokenFallback]）。
//
// # 存储抽象
//
// 构造器只依赖 go-redis v9 的 redis.Scripter 接口（Client / ClusterClient /
// Ring 均已实现），不绑定具体客户端类型，便于单测注入 mock。
//
// key 布局：
//
//   - PeriodLimit：keyPrefix + ":" + key（keyPrefix 为空时直接用 key）
//   - TokenLimiter：桶标识即 keyPrefix，实际 key 为 "{keyPrefix}.tokens" /
//     "{keyPrefix}.ts"，花括号为 Redis Cluster hash-tag，保证同桶两 key
//     落在同一 slot（沿用 go-zero 惯例）
//
// # 固定窗口 PeriodLimit
//
//	l := ratelimit.NewPeriodLimit(time.Hour, 5, rdb, "sms:verify")
//	switch code, err := l.Take(ctx, "13800000000"); {
//	case err != nil:
//	    // Redis 故障：返回 Unknown + err，由调用方决定放行或拒绝
//	case code == ratelimit.OverQuota:
//	    // 拒绝
//	default:
//	    // 放行（HitQuota 表示恰好用完配额，本次仍放行）
//	}
//
// 可用 [WithPeriodAlign] 把窗口对齐到自然周期（本地时区的天/小时边界），
// 典型场景"每天 5 条短信"。
//
// 与 TokenLimiter 不同，PeriodLimit **不做**本地兜底：固定窗口计数是跨实例的
// 全局精确配额，故障窗口内进程无法在不破坏全局限额的前提下本地近似——本地
// 放行多少都会放大全局限额。因此 Redis 故障时 Take 直接返回 (Unknown, err)，
// fail-open 还是 fail-closed 由调用方业务决定。
//
// # 令牌桶 TokenLimiter
//
//	l := ratelimit.NewTokenLimiter(100, 200, rdb, "api:ingress",
//	    ratelimit.WithTokenFallback())
//	if l.Allow(ctx) {
//	    // 放行
//	}
//
// 默认（未开启兜底）Redis 故障时一律拒绝（fail-closed）。开启
// [WithTokenFallback] 后，Redis 故障降级到进程内 golang.org/x/time/rate
// 兜底限流（同 rate/burst 参数），后台 goroutine 按 [WithTokenPingInterval]
// （默认 100ms）周期探活，成功一次即自动回切 Redis；探活期间兜底持续生效。
// 单次探活的超时与间隔解耦（[WithTokenProbeTimeout]，默认 1s）：须大于部署
// 环境的正常 EVAL RT（如跨 region），否则探测会因超时误判故障、降级永不回切。
// 兜底参数与 Redis 桶一致，故障窗口内的实际放行速率不变，仅失去跨实例全局
// 精确性（各进程各自持有一个等参本地桶）。
//
// # 设计边界
//
//   - ctx-only API：所有判定方法首参收 context.Context，无历史包袱
//   - 零可观测性注入：不内置 metrics/logx，故障切换对调用方透明；需要
//     观测时在业务层包装判定结果
//   - 单机限流请直接用 golang.org/x/time/rate，本包解决的是跨实例全局配额
package ratelimit
