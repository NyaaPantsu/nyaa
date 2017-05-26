package i2p

import (
	"net"
)

// i2p network session
type Session interface {

	// get session name
	Name() string
	// open a new control socket
	// does handshaske
	OpenControlSocket() (net.Conn, error)

	// get printable b32.i2p address
	B32Addr() string

	// implements network.Network
	Addr() net.Addr

	// implements network.Network
	Accept() (net.Conn, error)

	// implements network.Session
	Lookup(name string, port int) (net.Addr, error)

	// lookup an i2p address
	LookupI2P(name string) (I2PAddr, error)

	// implements network.Network
	Dial(n, a string) (net.Conn, error)

	// dial out to a remote destination
	DialI2P(a I2PAddr) (net.Conn, error)

	// open the session, generate keys, start up destination etc
	Open() error
	// close the session
	Close() error
}

// create a new i2p session
func NewSession(name, addr, keyfile string) Session {
	return &samSession{
		name:       name,
		addr:       addr,
		minversion: "3.0",
		maxversion: "3.0",
		keys:       NewKeyfile(keyfile),
	}
}
