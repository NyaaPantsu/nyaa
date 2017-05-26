package network

import (
	"net"
	"net/http"
)

// GracefulHttpCloser : implements io.Closer that gracefully closes an http server
type GracefulHttpCloser struct {
	Server   *http.Server
	Listener net.Listener
}

// Close method
func (c *GracefulHttpCloser) Close() error {
	c.Listener.Close()
	return c.Server.Shutdown(nil)
}
