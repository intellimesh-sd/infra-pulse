package server

import (
	"context"
	"github.com/clarechu/infra-pulse/src/server/router"
	routerv1 "github.com/clarechu/infra-pulse/src/server/router/v1"
)

type CmdbConfig struct {
	DataRoot  string `yaml:"data_root"`
	Port      int32  `yaml:"port"`
	ProxyPort int32  `yaml:"proxy_port"`
}

func NewCmdb(config *CmdbConfig) (Bootstrap, error) {

	ctx, cancel := context.WithCancel(context.Background())
	routeServer := routerv1.NewRouteServer()
	return &CMDB{
		port:      config.Port,
		proxyPort: config.ProxyPort,
		server:    router.NewServer(routeServer),
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}
