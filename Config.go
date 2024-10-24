package reconfig

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	ripc "github.com/stormi-li/Ripc"
)

type Config struct {
	name        string
	Info        *ConfigInfo
	ripcClient  *ripc.Client
	redisClient *redis.Client
	Context     context.Context
}

type ConfigInfo struct {
	Addr     string
	ConfigId int
	Desc     string
	Data     map[string]string
}

func (c ConfigInfo) ToString() string {
	bs, _ := json.MarshalIndent(c, " ", "  ")
	return string(bs)
}

func (c *Config) Upload(ttl time.Duration) {
	//---------------------------------------------------redis代码
	c.redisClient.Set(c.Context, c.ripcClient.Namespace+ConfigPrefix+c.name, c.Info.ToString(), ttl)
	c.ripcClient.Notify(ConfigPrefix+c.name, updateConfig)
}

func (c *Config) Delete() {
	//---------------------------------------------------redis代码
	c.redisClient.Del(c.Context, c.ripcClient.Namespace+ConfigPrefix+c.name)
	c.ripcClient.Notify(ConfigPrefix+c.name, updateConfig)
}
