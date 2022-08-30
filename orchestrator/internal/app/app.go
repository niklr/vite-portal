package app

import (
	"sync"
	"time"

	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
)

const (
	initializingState = iota
	runningState
	closedState
)

type OrchestratorApp struct {
	config        types.Config
	startStopLock sync.Mutex // Start/Stop are protected by an additional lock
	state         int        // Tracks state of node lifecycle

	lock          sync.Mutex
	rpcAPIs       []rpc.API // List of APIs currently provided by the app
	rpc           *rpc.HTTPServer
	rpcAuth       *rpc.HTTPServer
	inprocHandler *rpc.Server // In-process RPC request handler to process the API requests
	context       *Context
}

func NewOrchestratorApp(cfg types.Config) *OrchestratorApp {
	c := NewContext(cfg)
	a := &OrchestratorApp{
		config:        cfg,
		inprocHandler: rpc.NewServer(),
		context:       c,
	}

	// Register built-in APIs.
	a.rpcAPIs = append(a.rpcAPIs, a.apis()...)

	defaultTimeout := time.Duration(cfg.RpcTimeout) * time.Millisecond
	timeouts := rpc.HTTPTimeouts{
		ReadTimeout:       defaultTimeout,
		ReadHeaderTimeout: defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultTimeout * 2,
	}

	logger := logger.Logger()

	// Configure RPC servers.
	a.rpc = rpc.NewHTTPServer(logger, timeouts)
	a.rpcAuth = rpc.NewHTTPServer(logger, timeouts)

	return a
}

func (a *OrchestratorApp) Start() error {
	logger.Logger().Info().Msg("Start called")
	a.startStopLock.Lock()
	defer a.startStopLock.Unlock()

	// start RPC endpoints
	err := a.startRPC()
	if err != nil {
		a.stopRPC()
	}
	return err
}

func (a *OrchestratorApp) Shutdown() {
	logger.Logger().Info().Msg("Shutdown called")
	a.startStopLock.Lock()
	defer a.startStopLock.Unlock()

	a.stopRPC()
}