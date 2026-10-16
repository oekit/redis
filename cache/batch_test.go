package cache

import (
	"context"
	"testing"
	"time"
)

// MGet 管道化批量读：命中真值进 map；不存在与 null marker 一律缺席；
// 纯读不回填、不触发 Source。
func TestRedisCacheMGet(t *testing.T) {
	client, _ := newTestRedisClient(t)
	sourceCalls := 0
	c := NewRedisCache[testUser](client, RedisConfig[testUser]{
		Prefix: "mget:",
		TTL:    time.Minute,
		Source: func(ctx context.Context, key string) (*testUser, error) {
			sourceCalls++
			if key == "ghost" {
				return nil, nil // 生成 null marker
			}
			return &testUser{ID: len(key), Name: key}, nil
		},
	})
	ctx := context.Background()

	// 预置：真值直写、marker 经 GetOrLoad 生成
	if err := c.Set(ctx, "hit", &testUser{ID: 1, Name: "alice"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, err := c.GetOrLoad(ctx, "ghost"); err != nil {
		t.Fatalf("GetOrLoad marker: %v", err)
	}

	got, err := c.MGet(ctx, "hit", "ghost", "absent")
	if err != nil {
		t.Fatalf("MGet: %v", err)
	}
	if len(got) != 1 || got["hit"] == nil || got["hit"].Name != "alice" {
		t.Fatalf("MGet = %v, want 只含 hit（marker/不存在缺席）", got)
	}
	if sourceCalls != 1 {
		t.Fatalf("source calls = %d, want 1", sourceCalls)
	}

	// 纯 miss：缺席且不回源
	got2, err := c.MGet(ctx, "m1", "m2")
	if err != nil {
		t.Fatalf("MGet second: %v", err)
	}
	if len(got2) != 0 {
		t.Fatalf("纯 miss MGet = %v, want 空", got2)
	}
	if sourceCalls != 1 {
		t.Fatalf("MGet 不得触发回源: source calls = %d", sourceCalls)
	}

	// 空键集
	empty, err := c.MGet(ctx)
	if err != nil || len(empty) != 0 {
		t.Fatalf("MGet() = %v, %v", empty, err)
	}
}

// 编译期能力契约：RedisCache 具备批量读取能力。
var _ interface {
	MGet(ctx context.Context, keys ...string) (map[string]*testUser, error)
} = (*RedisCache[testUser])(nil)
