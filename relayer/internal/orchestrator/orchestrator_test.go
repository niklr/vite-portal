package orchestrator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
	wstest "github.com/vitelabs/vite-portal/shared/pkg/ws/test"
)

var timeout = 1000 * time.Millisecond

func TestInit(t *testing.T) {
	rpc := wstest.NewTestWsRpc(timeout)
	rpc.Start()
	o := NewOrchestrator(rpc.Url, timeout)
	require.NotNil(t, o)
	require.Equal(t, ws.Unknown, o.GetStatus())
	o.Start()
	commonutil.WaitFor(timeout, o.StatusChanged, func(status ws.ConnectionStatus) bool {
		return status == ws.Connected
	})
	require.Equal(t, ws.Connected, o.GetStatus())
	require.True(t, ws.CanConnect(rpc.Url, timeout))
	rpc.Stop()
	require.False(t, ws.CanConnect(rpc.Url, timeout))
	commonutil.WaitFor(timeout, o.StatusChanged, func(status ws.ConnectionStatus) bool {
		return status == ws.Disconnected
	})
	require.Equal(t, ws.Disconnected, o.GetStatus())
}

func TestMockInit(t *testing.T) {
	mock := ws.NewMockWsRpc(0)
	require.NotNil(t, mock)
	go mock.Serve(timeout)
	o := NewOrchestrator(mock.Url, timeout)
	require.NotNil(t, o)
	require.Equal(t, ws.Unknown, o.GetStatus())
	o.Start()
	commonutil.WaitFor(timeout, o.StatusChanged, func(status ws.ConnectionStatus) bool {
		return status == ws.Connected
	})
	require.Equal(t, ws.Connected, o.GetStatus())
	require.True(t, ws.CanConnect(mock.Url, timeout))
	mock.Close()
	require.False(t, ws.CanConnect(mock.Url, timeout))
	commonutil.WaitFor(timeout, o.StatusChanged, func(status ws.ConnectionStatus) bool {
		return status == ws.Disconnected
	})
	// Connections are not closed -> use TestWsRpc
	require.Equal(t, ws.Connected, o.GetStatus())
}

func TestInitInvalidUrl(t *testing.T) {
	t.Parallel()
	o := NewOrchestrator("http://localhost:1234", timeout)
	require.NotNil(t, o)
}