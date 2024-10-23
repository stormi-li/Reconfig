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

func (client *Client) NewConfig(name string, addr string) *Config {
	config := Config{
		name:       name,
		Info:       &ConfigInfo{Addr: addr},
		ripcClient: client.RipcClient,
		ctx:        client.ctx,
	}
	return &config
}

const configPrefix = "stormi:config:"

const updateConfig = "updateConfig"

func (client *Client) getConfig(name string) *ConfigInfo {
	configStr, _ := client.RipcClient.RedisClient.Get(client.ctx, configPrefix+name).Result()
	var c ConfigInfo
	json.Unmarshal([]byte(configStr), &c)
	return &c
}

func (client *Client) Connect(name string, handler func(configInfo *ConfigInfo)) {
	listener := client.RipcClient.NewListener(client.ctx, configPrefix+name)
	config := client.getConfig(name)
	handler(config)
	go func() {
		listener.Listen(func(msg string) {
			if msg == updateConfig {
				config = client.getConfig(name)
				handler(config)
			}
		})
	}()
	for {
		time.Sleep(10 * time.Second)
		cfg := client.getConfig(name)
		if cfg.ToString() != config.ToString() {
			config = cfg
			handler(config)
		}
	}
}
