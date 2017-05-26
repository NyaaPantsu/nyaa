package i2p

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type samSession struct {
	addr       string
	minversion string
	maxversion string
	name       string
	keys       *Keyfile
	// control connection
	c   net.Conn
	mtx sync.RWMutex
}

func (s *samSession) Close() error {
	return s.c.Close()
}

func (s *samSession) B32Addr() string {
	return s.keys.Addr().Base32Addr().String()
}

func (s *samSession) Name() string {
	return s.name
}

func (s *samSession) Addr() net.Addr {
	return s.keys.Addr()
}

func (s *samSession) OpenControlSocket() (n net.Conn, err error) {
	n, err = net.Dial("tcp", s.addr)
	if err == nil {
		// make the connection never time out
		err = n.(*net.TCPConn).SetKeepAlive(true)
		if err == nil {
			// send keepalive every 5 seconds
			err = n.(*net.TCPConn).SetKeepAlivePeriod(time.Second * 5)
		}
		if err != nil {
			err = nil
		}
		_, err = fmt.Fprintf(n, "HELLO VERSION MIN=%s MAX=%s\n", s.minversion, s.maxversion)
		r := bufio.NewReader(n)
		var line string
		line, err = r.ReadString(10)
		if err == nil {
			sc := bufio.NewScanner(strings.NewReader(line))
			sc.Split(bufio.ScanWords)
			for sc.Scan() {
				text := strings.ToUpper(sc.Text())
				if text == "HELLO" {
					continue
				}
				if text == "REPLY" {
					continue
				}
				if text == "RESULT=OK" {
					// we good
					return
				}
				err = errors.New(line)
			}
		}
		n.Close()
	}
	return
}

func (s *samSession) DialI2P(addr I2PAddr) (c net.Conn, err error) {
	var nc net.Conn
	nc, err = s.OpenControlSocket()
	if err == nil {
		// send connect
		_, err = fmt.Fprintf(nc, "STREAM CONNECT ID=%s DESTINATION=%s SILENT=false\n", s.Name(), addr.String())

		r := bufio.NewReader(nc)
		var line string
		// read reply
		line, err = r.ReadString(10)
		if err == nil {
			// parse reply
			sc := bufio.NewScanner(strings.NewReader(line))
			sc.Split(bufio.ScanWords)
			for sc.Scan() {
				txt := sc.Text()
				upper := strings.ToUpper(txt)
				if upper == "STREAM" {
					continue
				}
				if upper == "STATUS" {
					continue
				}
				if upper == "RESULT=OK" {
					// we are connected
					c = &I2PConn{
						c:     nc,
						laddr: s.keys.Addr(),
						raddr: addr,
					}
					return
				}
				err = errors.New(line)
				nc.Close()
			}
		}
	}
	return
}

func (s *samSession) Dial(n, a string) (c net.Conn, err error) {
	var addr I2PAddr
	addr, err = s.LookupI2P(a)
	if err == nil {
		c, err = s.DialI2P(addr)
	}
	return
}

func (s *samSession) LookupI2P(name string) (a I2PAddr, err error) {
	var n string
	n, _, err = net.SplitHostPort(name)
	if err == nil {
		name = n
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.c == nil {
		// no session socket
		err = errors.New("session not open")
		return
	}
	_, err = fmt.Fprintf(s.c, "NAMING LOOKUP NAME=%s\n", name)
	r := bufio.NewReader(s.c)
	var line string
	line, err = r.ReadString(10)
	if err == nil {
		// okay
		sc := bufio.NewScanner(strings.NewReader(line))
		sc.Split(bufio.ScanWords)
		for sc.Scan() {
			txt := sc.Text()
			upper := strings.ToUpper(txt)
			if upper == "NAMING" {
				continue
			}
			if upper == "REPLY" {
				continue
			}
			if upper == "RESULT=OK" {
				continue
			}
			if strings.HasPrefix(upper, "NAME=") {
				continue
			}
			if strings.HasPrefix(txt, "VALUE=") {
				// we got it
				a = I2PAddr(txt[6:])
				return
			}
			err = errors.New(line)
		}
	}
	return
}

func (s *samSession) Lookup(name string, port int) (a net.Addr, err error) {
	a, err = s.LookupI2P(name)
	return
}

func (s *samSession) createStreamSession() (err error) {
	// try opening if this session isn't already open
	_, err = fmt.Fprintf(s.c, "SESSION CREATE STYLE=STREAM ID=%s DESTINATION=%s\n", s.Name(), s.keys.privkey)
	if err == nil {
		// read response line
		r := bufio.NewReader(s.c)
		var line string
		line, err = r.ReadString(10)
		if err == nil {
			// parse response line
			sc := bufio.NewScanner(strings.NewReader(line))
			sc.Split(bufio.ScanWords)
			for sc.Scan() {
				text := sc.Text()
				upper := strings.ToUpper(text)
				if upper == "SESSION" {
					continue
				}
				if upper == "STATUS" {
					continue
				}
				if upper == "RESULT=OK" {
					// we good
					return
				}
				err = errors.New(line)
			}
		}
	}
	return
}

func (s *samSession) Open() (err error) {
	s.c, err = s.OpenControlSocket()
	if err == nil {
		err = s.keys.ensure(s.c)
	}
	if err == nil {
		err = s.createStreamSession()
		if err == nil {
			var a I2PAddr
			a, err = s.LookupI2P("ME")
			if err == nil {
				s.keys.pubkey = a.String()
			}
		}
	}
	if err != nil {
		s.Close()
	}
	return
}

func (s *samSession) Accept() (c net.Conn, err error) {
	l := &i2pListener{
		session: s,
		laddr:   I2PAddr(s.keys.pubkey),
	}
	c, err = l.Accept()
	return
}
