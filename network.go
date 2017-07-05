package main

import (
	"fmt"
	"net"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/majestrate/i2p-tools/lib/i2p"
)

// CreateHTTPListener creates a net.Listener for main http webapp given main config
func CreateHTTPListener(conf *config.Config) (net.Listener, error) {
	if conf.I2P != nil {
		s := i2p.NewSession(conf.I2P.Name, conf.I2P.Addr, conf.I2P.Keyfile)
		err := s.Open()
		if s != nil {
			log.Infof("i2p address: %s", s.B32Addr())
		}
		return s, err
	}
	return net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
}

// CreateScraperSocket creates a UDP Scraper socket
func CreateScraperSocket(conf *config.Config) (net.PacketConn, error) {
	if conf.I2P != nil {
		log.Fatal("i2p udp scraper not supported")
	}
	var laddr *net.UDPAddr
	laddr, err := net.ResolveUDPAddr("udp", conf.Scrape.Addr)
	if err != nil {
		return nil, err
	}
	return net.ListenUDP("udp", laddr)
}
