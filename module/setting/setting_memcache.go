package setting

import (
	_ "github.com/go-macaron/cache/memcache"
)

func init() {
	EnableMemcache = true
}
