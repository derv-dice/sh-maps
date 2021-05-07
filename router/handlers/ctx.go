package handlers

import (
	"fmt"
	"sh-maps/config"
)

type BaseCtx struct {
	SelfAddr string                 `json:"self_addr"`
	Data     map[string]interface{} `json:"data"`
}

func (b *BaseCtx) Add(key string, val interface{}) *BaseCtx {
	if b.Data == nil {
		b.Data = map[string]interface{}{}
	}
	b.Data[key] = val
	return b
}

func NewBaseCtx() *BaseCtx {
	return &BaseCtx{
		SelfAddr: fmt.Sprintf(config.RemoteAddrTmpl(), config.Cfg.Server.Addr),
	}
}
