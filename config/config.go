package config

import (
	"flag"
	"os"
	"path"
)

var (
	Cwd, _  = os.Getwd()
	TempDir = path.Join(Cwd, "/temp")

	defaultRedisAddr = "127.0.0.1:6379"
	RedisAddr        = *flag.String("redisAddr", defaultRedisAddr, "Listen address for Redis service")
	HTTPAddr         = "127.0.0.1:8000"
	TaskQueue        = "default"
)

func init() {
	_ = os.MkdirAll(TempDir, 0755)
}
