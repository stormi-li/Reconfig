package main

import (
	"time"

	"github.com/go-redis/redis/v8"
	reconfig "github.com/stormi-li/Reconfig"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "118.25.196.166:6379",
	})
	client := reconfig.NewClient(redisClient)
	client.SetNamespace("123")
	config := client.NewConfig("mysql", "118.25.196.166:3306")
	config.Upload(1000 * time.Second)
	time.Sleep(1 * time.Second)
	config.Delete()
}
