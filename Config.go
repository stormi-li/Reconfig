package reconfig

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	ripc "github.com/stormi-li/Ripc"
)

type Config struct {
	ripcClient  *ripc.Client
	redisClient *redis.Client
	Name        string
	Addr        string
	ConfigId    int
	Desc        string
	Data        map[string]string
	namespace   string
	ctx         context.Context
}

func (c Config) ToString() string {
	bs, _ := json.MarshalIndent(c, " ", "  ")
	return string(bs)
}

func (c *Config) Upload() {
	//---------------------------------------------------redis代码
	c.redisClient.Set(c.ctx, c.namespace+ConfigPrefix+c.Name, c.ToString(), 0)
	c.ripcClient.Notify(ConfigPrefix+c.Name, updateConfig)
}

func (c *Config) Delete() {
	//---------------------------------------------------redis代码
	c.redisClient.Del(c.ctx, c.namespace+ConfigPrefix+c.Name)
	c.ripcClient.Notify(ConfigPrefix+c.Name, updateConfig)
}
