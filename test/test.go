package main

import (
	"time"

	reconfig "github.com/stormi-li/Reconfig"
)

func main() {
	client, _ := reconfig.NewClient("118.25.196.166:6379")
	client.SetNameSpace("fsf")
	config := client.NewConfig("mysql", "118.25.196.166:3307")
	config.Upload(10 * time.Second)
}
