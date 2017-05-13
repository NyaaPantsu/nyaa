package network

import (
	"fmt"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/util/log"
	"github.com/majestrate/i2p-tools/lib/i2p"
	"net"
)

// CreateHTTPListener creates a net.Listener for main http webapp given main config
func CreateHTTPListener(conf *config.Config) (l net.Listener, err error) {
	if conf.I2P == nil {
		l, err = net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	} else {
		s := i2p.NewSession(conf.I2P.Name, conf.I2P.Addr, conf.I2P.Keyfile)
		err = s.Open()
		if s != nil {
			log.Infof("i2p address: %s", s.B32Addr())
			l = s
		}
	}
	if l != nil {
		l = WrapListener(l)
	}
	return
}

// CreateScraperSocket creates a UDP Scraper socket
func CreateScraperSocket(conf *config.Config) (pc net.PacketConn, err error) {
	if conf.I2P == nil {
		var laddr *net.UDPAddr
		laddr, err = net.ResolveUDPAddr("udp", conf.Scrape.Addr)
		if err == nil {
			pc, err = net.ListenUDP("udp", laddr)
		}
	} else {
		log.Fatal("i2p udp scraper not supported")
	}
	return
}
