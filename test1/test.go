package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	reconfig "github.com/stormi-li/Reconfig"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "118.25.196.166:6379",
	})
	client := reconfig.NewClient(redisClient)
	client.SetNamespace("123")

	names := client.GetConfigNames()
	fmt.Println(names)
	client.Connect("mysql", func(configInfo *reconfig.ConfigInfo) {
		fmt.Println(configInfo.Addr)
	})
}
