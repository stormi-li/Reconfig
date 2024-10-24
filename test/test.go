package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	reconfig "github.com/stormi-li/Reconfig"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "118.25.196.166:3934",
		Password: "12982397StrongPassw0rd",
	})
	ctx := context.Background()

	// 测试连接
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
	} else {
		fmt.Println("Redis connected successfully!")
	}
	client := reconfig.NewClient(redisClient)
	client.SetNamespace("123")
	config := client.NewConfig("mysql", "118.25.196.166:3306")
	config.Upload(1000 * time.Second)
	// time.Sleep(1 * time.Second)
	// config.Delete()
}
