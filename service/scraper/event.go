package scraperService

import (
	"encoding/binary"
	"net"
)

type RecvEvent struct {
	From net.Addr
	Data []byte
}

// TID extract transaction id
func (ev *RecvEvent) TID() (id uint32, err error) {
	if len(ev.Data) < 8 {
		err = ErrShortPacket
	} else {
		id = binary.BigEndian.Uint32(ev.Data[4:])
	}
	return
}

// Action extract action
func (ev *RecvEvent) Action() (action uint32, err error) {
	if len(ev.Data) < 4 {
		err = ErrShortPacket
	} else {
		action = binary.BigEndian.Uint32(ev.Data)
	}
	return
}

type SendEvent struct {
	To   net.Addr
	Data []byte
}
