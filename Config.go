package reconfig

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	ripc "github.com/stormi-li/Ripc"
)

type Config struct {
	redisClient *redis.Client
	ripcClient  *ripc.Client
	Name        string
	Addr        string
	ConfigId    int
	Desc        string
	Data        map[string]string
	namespace   string
	ctx         context.Context
}

func newConfig(redisClient *redis.Client, ripcClient *ripc.Client, name string, addr string, namespace string) *Config {
	return &Config{
		redisClient: redisClient,
		ripcClient:  ripcClient,
		Name:        name,
		Addr:        addr,
		Data:        map[string]string{},
		namespace:   namespace,
		ctx:         context.Background(),
	}
}

func (c Config) ToString() string {
	bs, _ := json.MarshalIndent(c, " ", "  ")
	return string(bs)
}

func (c *Config) Upload() {
	c.redisClient.Set(c.ctx, c.namespace+c.Name, c.ToString(), 0)
	c.ripcClient.Notify(c.namespace+c.Name, const_updateConfig)
}

func (c *Config) Delete() {
	c.redisClient.Del(c.ctx, c.namespace+c.Name)
	c.ripcClient.Notify(c.namespace+c.Name, const_updateConfig)
}
