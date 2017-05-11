package scraperService

import (
	"errors"
)

var ErrShortPacket = errors.New("short udp packet")
