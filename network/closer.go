package network

import (
	"net"
	"net/http"
)

// implements io.Closer that gracefully closes an http server
type GracefulHttpCloser struct {
	Server   *http.Server
	Listener net.Listener
}

func (c *GracefulHttpCloser) Close() error {
	c.Listener.Close()
	return c.Server.Shutdown(nil)
}
