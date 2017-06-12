package i2p

import (
	"net"
	"time"
)

// tcp/i2p connection
// implements net.Conn
type I2PConn struct {
	// underlying connection
	c net.Conn
	// our local address
	laddr I2PAddr
	// remote peer's address
	raddr I2PAddr
}

// implements net.Conn
func (c *I2PConn) Read(d []byte) (n int, err error) {
	n, err = c.c.Read(d)
	return
}

// implements net.Conn
func (c *I2PConn) Write(d []byte) (n int, err error) {
	n, err = c.c.Write(d)
	return
}

// implements net.Conn
func (c *I2PConn) Close() error {
	return c.c.Close()
}

// implements net.Conn
func (c *I2PConn) LocalAddr() net.Addr {
	return c.laddr
}

// implements net.Conn
func (c *I2PConn) RemoteAddr() net.Addr {
	return c.raddr
}

// implements net.Conn
func (c *I2PConn) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}

// implements net.Conn
func (c *I2PConn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

// implements net.Conn
func (c *I2PConn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}
