package ratelimit

import (
	"context"
	"time"

	"github.com/oekit/concurrency/async"
	xrate "golang.org/x/time/rate"
)

// defaultPingInterval 兜底期间 Redis 探活默认间隔（沿用 go-zero 同名默认值）。
const defaultPingInterval = 100 * time.Millisecond

// newRescueLimiter 构建与 Redis 桶同参数的进程内兜底限流器。
// rate 恒为正（NewTokenLimiter 已校验），不存在除零。
func newRescueLimiter(rate, burst int) *xrate.Limiter {
	return xrate.NewLimiter(xrate.Every(time.Second/time.Duration(rate)), burst)
}

// startMonitor 启动 Redis 探活监控（幂等）。
//
// 首个 Redis 故障触发：置 redisAlive=false（后续判定直接走本地桶）→
// 后台 goroutine 按 pingInterval 周期探活，成功一次即回切并结束监控；
// 探活期间再次故障不会重复拉起监控（monitorStarted 幂等）。
// goroutine 经 async.GoSafe 启动，panic 被 recover（本仓硬规则）。
func (l *TokenLimiter) startMonitor() {
	l.rescueLock.Lock()
	defer l.rescueLock.Unlock()

	if l.monitorStarted {
		return
	}
	l.monitorStarted = true
	l.redisAlive.Store(false)

	async.GoSafe(l.waitForRedis)
}

// waitForRedis 探活循环：先立即探活一次，之后按 pingInterval 周期进行，
// 成功一次即回切 Redis（redisAlive=true）并结束。退出时无论何种路径都
// 复位 monitorStarted，保证下一次故障能重新拉起监控。
func (l *TokenLimiter) waitForRedis() {
	defer func() {
		l.rescueLock.Lock()
		l.monitorStarted = false
		l.rescueLock.Unlock()
	}()

	ticker := time.NewTicker(l.pingInterval)
	defer ticker.Stop()

	if l.probe() {
		l.redisAlive.Store(true)
		return
	}
	for range ticker.C {
		if l.probe() {
			l.redisAlive.Store(true)
			return
		}
	}
}

// probe 探测 Redis 可用性：执行只含常量返回的最小 Lua 脚本（scripts.go probeLua）。
// redis.Scripter 接口无 Ping，EVAL 常量脚本等价于"连通 + 脚本执行能力"探测。
// 探活以独立 ctx 进行（不受单次请求 ctx 影响），超时取 pingInterval，
// 防止挂死的连接把探活 goroutine 永久卡住。
func (l *TokenLimiter) probe() bool {
	ctx, cancel := context.WithTimeout(context.Background(), l.pingInterval)
	defer cancel()
	return probeScript.Run(ctx, l.store, nil).Err() == nil
}
