package reconfig

import (
	"context"
	"encoding/json"
	"time"

	ripc "github.com/stormi-li/Ripc"
)

type Config struct {
	name       string
	Info       *ConfigInfo
	ripcClient *ripc.Client
	ctx        context.Context
	namespace  string
}

type ConfigInfo struct {
	Addr     string
	ConfigId int
	Desc     string
	Info     map[string]string
}

func (c ConfigInfo) ToString() string {
	bs, _ := json.MarshalIndent(c, " ", "  ")
	return string(bs)
}

func (c *Config) Upload(ttl time.Duration) {
	c.ripcClient.RedisClient.Set(c.ctx, c.namespace+configPrefix+c.name, c.Info.ToString(), ttl)
	c.ripcClient.Notify(c.ctx, c.namespace+configPrefix+c.name, updateConfig)
}

func (c *Config) Delete() {
	c.ripcClient.RedisClient.Del(c.ctx, c.namespace+configPrefix+c.name)
	c.Info.Addr = ""
	c.ripcClient.Notify(c.ctx, c.namespace+configPrefix+c.name, updateConfig)
}
