package main

import (
	"github.com/go-redis/redis/v8"
	reconfig "github.com/stormi-li/Reconfig"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	client := reconfig.NewClient(redisClient, "reconfig-namespace")
	cfg := client.NewConfig("mysql", "localhost:12234")
	cfg.Data["username"]="root"
	cfg.Data["password"]="123456"
	cfg.Upload()
}
