package network

import (
	"errors"
	"net"
)

var ErrListenerStopped = errors.New("listener was stopped")

// GracefulListener provides safe and graceful net.Listener wrapper that prevents error on graceful shutdown
type GracefulListener struct {
	listener net.Listener
	stop     chan int
}

func (l *GracefulListener) Accept() (net.Conn, error) {
	for {
		c, err := l.listener.Accept()
		select {
		case <-l.stop:
			if c != nil {
				c.Close()
			}
			close(l.stop)
			l.stop = nil
			return nil, ErrListenerStopped
		default:

		}
		if err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() && neterr.Temporary() {
				continue
			}
		}
		return c, err
	}
}

func (l *GracefulListener) Close() (err error) {
	l.listener.Close()
	if l.stop != nil {
		l.stop <- 0
	}
	return
}

func (l *GracefulListener) Addr() net.Addr {
	return l.listener.Addr()
}

// WrapListener wraps a net.Listener such that it can be closed gracefully
func WrapListener(l net.Listener) net.Listener {
	return &GracefulListener{
		listener: l,
		stop:     make(chan int),
	}
}
