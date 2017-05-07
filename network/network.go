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
	return
}
