package setting

import (
	_ "github.com/go-macaron/cache/redis"
	_ "github.com/go-macaron/session/redis"
)

func init() {
	EnableRedis = true
}
