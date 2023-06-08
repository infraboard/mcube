package rpc

import "github.com/infraboard/mcenter/clients/rpc"

var (
	client *ClientSet
)

func C() *ClientSet {
	if client == nil {
		panic("mpaas rpc client config not load")
	}
	return client
}

func HasLoaded() bool {
	return client != nil
}

func LoadClientFromEnv() error {
	cs, err := NewClientSetFromEnv()
	if err != nil {
		return err
	}
	client = cs
	return nil
}

func LoadClientFromConfig(conf *rpc.Config) error {
	c, err := NewClientSetFromConfig(conf)
	if err != nil {
		return err
	}
	client = c
	return nil
}