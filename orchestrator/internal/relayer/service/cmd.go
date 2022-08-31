package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

func (s *Service) Handle(timeout time.Duration, c *rpc.Client, peerInfo rpc.PeerInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var resp sharedtypes.RpcAppInfoResponse
	err := c.CallContext(ctx, &resp, "core_getAppInfo")
	if err != nil {
		logger.Logger().Error().Err(err).Msg("calling context failed")
		return err
	}
	logger.Logger().Debug().Str("resp", fmt.Sprintf("%#v", resp)).Msg("onconnect result")
	return nil
}