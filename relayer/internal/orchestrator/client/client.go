package client

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vitelabs/vite-portal/shared/pkg/crypto"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Client struct {
	url       string
	timeout   time.Duration
	jwtHanler crypto.JWTHandler
	Conn      *websocket.Conn
}

func NewClient(url, jwtSecret string, timeout time.Duration) *Client {
	return &Client{
		url:       url,
		timeout:   timeout,
		jwtHanler: *crypto.NewDefaultJWTHandler([]byte(jwtSecret)),
	}
}

func (c *Client) Connect() error {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: c.timeout,
	}
	headers := make(http.Header, 1)
	headers.Set("Authorization", c.jwtHanler.IssueDefaultToken(sharedtypes.JWTRelayerSubject))
	conn, _, err := dialer.Dial(c.url, headers)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}
