package config

import (
	"os"
	"path"
)

const (
	RedisAddr = "127.0.0.1:6379"
	HTTPAddr  = "127.0.0.1:8000"
	TaskQueue = "default"
)

var (
	Cwd, _  = os.Getwd()
	TempDir = path.Join(Cwd, "/temp")
)

func init() {
	_ = os.MkdirAll(TempDir, 0755)
}
