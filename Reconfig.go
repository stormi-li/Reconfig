package reconfig

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	ripc "github.com/stormi-li/Ripc"
)

type Client struct {
	redisClient *redis.Client
	ripcClient  *ripc.Client
	ctx         context.Context
	namespace   string
}

func NewClient(redisClient *redis.Client, namespace string) *Client {
	ripcClient := ripc.NewClient(redisClient, namespace)
	return &Client{redisClient: redisClient, ripcClient: ripcClient, ctx: context.Background(), namespace: namespace}
}

func (c *Client) NewConfig(name string, addr string) *Config {
	return &Config{
		Name:        name,
		Addr:        addr,
		Data:        map[string]string{},
		ripcClient:  c.ripcClient,
		redisClient: c.redisClient,
		ctx:         c.ctx,
		namespace:   c.namespace,
	}
}

const ConfigPrefix = "stormi:config:"

const updateConfig = "updateConfig"

func (c *Client) GetConfig(name string) *Config {
	//---------------------------------------------------redis代码
	configStr, _ := c.redisClient.Get(c.ctx, c.namespace+ConfigPrefix+name).Result()
	var cfg Config
	json.Unmarshal([]byte(configStr), &cfg)
	cfg.ctx = c.ctx
	cfg.namespace = c.namespace
	cfg.redisClient = c.redisClient
	cfg.ripcClient = c.ripcClient
	return &cfg
}

func (c *Client) GetConfigNames() []string {
	return GetKeysByNamespace(c.redisClient, c.namespace+ConfigPrefix)
}

func (c *Client) Listen(name string, handler func(config *Config)) {
	listener := c.ripcClient.NewListener(ConfigPrefix + name)
	config := c.GetConfig(name)
	handler(config)
	go func() {
		listener.Listen(func(msg string) {
			if msg == updateConfig {
				cfg := c.GetConfig(name)
				if cfg.ToString() != config.ToString() {
					config = cfg
					handler(config)
				}
			}
		})
	}()
	for {
		time.Sleep(10 * time.Second)
		cfg := c.GetConfig(name)
		if cfg.ToString() != config.ToString() {
			config = cfg
			handler(config)
		}
	}
}

func GetKeysByNamespace(redisClient *redis.Client, namespace string) []string {
	var keys []string
	cursor := uint64(0)

	for {
		// 使用 SCAN 命令获取键名
		res, newCursor, err := redisClient.Scan(context.Background(), cursor, namespace+"*", 0).Result()
		if err != nil {
			return nil
		}

		// 处理键名，去掉命名空间
		for _, key := range res {
			// 去掉命名空间部分
			keyWithoutNamespace := key[len(namespace):]
			keys = append(keys, keyWithoutNamespace)
		}

		cursor = newCursor

		// 如果游标为0，则结束循环
		if cursor == 0 {
			break
		}
	}

	return keys
}
