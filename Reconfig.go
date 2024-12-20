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
	namespace   string
	ctx         context.Context
}

func NewClient(redisClient *redis.Client, namespace string) *Client {
	return &Client{
		redisClient: redisClient,
		ripcClient:  ripc.NewClient(redisClient, namespace),
		namespace:   namespace + const_splitChar + const_prefix,
		ctx:         context.Background(),
	}
}

func (c *Client) NewConfig(name string, addr string) *Config {
	return newConfig(c.redisClient, c.ripcClient, name, addr, c.namespace)
}

func (c *Client) GetConfig(name string) *Config {
	//---------------------------------------------------redis代码
	configStr, _ := c.redisClient.Get(c.ctx, c.namespace+name).Result()
	var cfg Config
	json.Unmarshal([]byte(configStr), &cfg)
	cfg.ctx = c.ctx
	cfg.namespace = c.namespace
	cfg.redisClient = c.redisClient
	cfg.ripcClient = c.ripcClient
	return &cfg
}

func (c *Client) GetConfigNames() []string {
	return getKeysByNamespace(c.redisClient, c.namespace)
}

func (c *Client) Listen(name string, handler func(config *Config)) {
	listener := c.ripcClient.NewListener(c.namespace + name)
	config := c.GetConfig(name)
	handler(config)
	go func() {
		listener.Listen(func(msg string) {
			if msg == const_updateConfig {
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

func getKeysByNamespace(redisClient *redis.Client, namespace string) []string {
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
