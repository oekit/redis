# ⚠️ 已改名：oekit/redis → oekit/redisx

本模块已更名为 [oekit/redisx](https://github.com/oekit/redisx)（v0.4.0 起），
以消除与 go-redis（`package redis`）的包名撞车及 `/redis/redis` 路径口吃。

- 旧路径 `github.com/oekit/redis` 经 GitHub 重定向仍可拉取 ≤ v0.3.1 的版本
- 新用法见新仓 README；`redisx.Client` 等全部能力原样保留
  （原 `redis/` 子包更名为 `cache/`，`redisx` 子包上收为根包）

```bash
go get github.com/oekit/redisx@v0.4.0
```
