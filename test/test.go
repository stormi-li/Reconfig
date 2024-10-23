package main

import (
	"time"

	reconfig "github.com/stormi-li/Reconfig"
)

func main() {
	client, _ := reconfig.NewClient("localhost:6379")
	config := client.NewConfig("mysql", "localhost:3307")
	config.Upload(10 * time.Second)
}
