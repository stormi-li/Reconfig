package main

import (
	reconfig "github.com/stormi-li/Reconfig"
)

func main() {
	client, _ := reconfig.NewClient("localhost:6379")
	client.Connect("mysql", func(configInfo *reconfig.ConfigInfo) {
	})
}
