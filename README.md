# RECONFIG Guides

Simple and practical distributed configuration management library.

# Overview

- Support configuration sharing
- Support configuration listening
- Every feature comes with tests
- Developer Friendly

# Install

```shell
go get -u github.com/stormi-li/Reconfig
```

# Quick Start

```go
package main

import (
	"github.com/go-redis/redis/v8"
	reconfig "github.com/stormi-li/Reconfig"
)

var redisAddr = “localhost:6379”
var password = “your password”

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	client := reconfig.NewClient(redisClient, "reconfig-namespace")
	cfg := client.NewConfig("mysql", "localhost:12234")
	cfg.Upload()
}
```

# Interface - reconfig

## NewClient

### Create reconfig client
```go
package main

import (
	"github.com/go-redis/redis/v8"
	reconfig "github.com/stormi-li/Reconfig"
)

var redisAddr = “localhost:6379”
var password = “your password”

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	client := reconfig.NewClient(redisClient, "reconfig-namespace")
}
```

The first parameter is a redis client of successful connection, the second parameter is a unique namespace.

# Interface - reconfig.Client

## GetConfig

### Get the configuration information through the configuration name
```go
client.GetConfig("mysql")
```
The parameter is configuration name. The returned value is a resync.Config struct.

## GetConfigNames

### Get all configuration names
```go
client.GetConfigNames("mysql")
```
The parameter is configuration name. The returned value is a name array.

## Listen

### Listen for configuration information changes
```go
client.Listen("mysql", func(config *reconfig.Config) {
	fmt.Println(config.Name, config.Addr)
})
```
The parameter is a handler for received configuration.

## NewConfig 

### Create Config
```go
cfg := client.NewConfig("mysql", "localhost:3306”)
```
The first parameter is the configuration name, the second parameter is the node address.

# Interface - reconfig.Config

## Upload

### Upload configuration information
```go
cfg.Data["username"]="root"
cfg.Data["password"]="123456"
cfg.Upload()
```

## Delete

### Delete configuration information
```go
cfg.Delete()
```

#  Community

## Ask

### How do I ask a good question?
- Email - 2785782829@qq.com
- Github Issues - https://github.com/stormi-li/Reconfig/issues