package main

import (
	"fmt"

	reconfig "github.com/stormi-li/Reconfig"
)

func main() {
	client, _ := reconfig.NewClient("118.25.196.166:6379")
	client.SetNameSpace("f1sf")
	client.Connect("mysql", func(configInfo *reconfig.ConfigInfo) {
		fmt.Println(configInfo.Addr)
	})
}
