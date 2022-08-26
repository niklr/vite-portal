package app

import (
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/configutil"
)

var (
	DefaultAllowedOrigins = []string{"*"}
	DefaultVhosts         = []string{"localhost"}
	DefaultModules        = []string{"nodes", "public"}
)

func InitApp(debug bool, configPath string) (*RelayerApp, error) {
	logger.Init(debug)
	p := configPath
	if p == "" {
		p = types.DefaultConfigFilename
	}
	cfg := types.NewDefaultConfig()
	err := configutil.InitConfig(&cfg, debug, p, types.DefaultConfigVersion)
	if err != nil {
		return nil, err
	}
	o := InitOrchestrator(cfg.OrchestratorWsUrl, cfg.RpcTimeout)
	c, err := InitContext(cfg)
	if err != nil {
		return nil, err
	}
	a := NewRelayerApp(cfg, o, c)
	return a, nil
}

