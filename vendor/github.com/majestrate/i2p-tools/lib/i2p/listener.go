package i2p

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type i2pListener struct {
	// parent session
	session Session
	// local address
	laddr I2PAddr
}

// implements net.Listener
func (l *i2pListener) Addr() net.Addr {
	return l.laddr
}

// implements net.Listener
func (l *i2pListener) Close() error {
	l.session = nil
	return nil
}

// implements net.Listener
func (l *i2pListener) Accept() (c net.Conn, err error) {
	if l.session == nil {
		err = errors.New("session closed")
		return
	}
	var nc net.Conn
	nc, err = l.session.OpenControlSocket()
	if err == nil {
		_, err = fmt.Fprintf(nc, "STREAM ACCEPT ID=%s SILENT=false\n", l.session.Name())
		if err == nil {
			r := bufio.NewReader(nc)
			var line string
			// read response line
			line, err = r.ReadString(10)
			if err == nil {
				sc := bufio.NewScanner(strings.NewReader(line))
				sc.Split(bufio.ScanWords)
				for sc.Scan() {
					text := sc.Text()
					upper := strings.ToUpper(text)
					if upper == "STREAM" {
						continue
					}
					if upper == "RESULT" {
						continue
					}
					if upper == "RESULT=OK" {
						// we good
						break
					}
					// error
					err = errors.New(text)
				}
			}
			// read address line
			line, err = r.ReadString(10)
			if err == nil {
				// we got a new connection yeeeeh
				c = &I2PConn{
					c:     nc,
					laddr: l.laddr,
					raddr: I2PAddr(line[:len(line)-1]),
				}
			}
		}
		if c == nil {
			// we didn't get a connection
			nc.Close()
		}
	}
	return
}
