package reconfig

import (
	"context"
	"encoding/json"
	"time"

	ripc "github.com/stormi-li/Ripc"
)

type Client struct {
	RipcClient *ripc.Client
	ctx        context.Context
	namespace  string
}

func NewClient(addr string) (*Client, error) {
	client := Client{}
	client.ctx = context.Background()
	ripcC, err := ripc.NewClient(addr)
	if err != nil {
		return nil, err
	}
	client.RipcClient = ripcC
	return &client, nil
}

func (c *Client) SetNameSpace(str string) {
	c.namespace = str + ":"
	c.RipcClient.SetNameSpace(str)
}

func (c *Client) NewConfig(name string, addr string) *Config {
	config := Config{
		name:       name,
		Info:       &ConfigInfo{Addr: addr},
		ripcClient: c.RipcClient,
		ctx:        c.ctx,
		namespace:  c.namespace,
	}
	return &config
}

const configPrefix = "stormi:config:"

const updateConfig = "updateConfig"

func (c *Client) GetConfig(name string) *ConfigInfo {
	configStr, _ := c.RipcClient.RedisClient.Get(c.ctx, c.namespace+configPrefix+name).Result()
	var cfg ConfigInfo
	json.Unmarshal([]byte(configStr), &cfg)
	return &cfg
}

func (c *Client) GetTTL(name string) time.Duration {
	ttl, _ := c.RipcClient.RedisClient.TTL(c.ctx, c.namespace+configPrefix+name).Result()
	return ttl
}

func (c *Client) Connect(name string, handler func(configInfo *ConfigInfo)) {
	listener := c.RipcClient.NewListener(c.ctx, c.namespace+configPrefix+name)
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
