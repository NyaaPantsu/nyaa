package scraperService

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"time"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
)

const stateSendID = 0
const stateRecvID = 1
const stateTransact = 2

const actionError = 3
const actionScrape = 2
const actionAnnounce = 1
const actionConnect = 0

// Transaction a scrape transaction on a udp tracker
type Transaction struct {
	TransactionID uint32
	ConnectionID  uint64
	bucket        *Bucket
	state         uint8
	swarms        []model.Torrent
}

// Done marks this transaction as done and removes it from parent
func (t *Transaction) Done() {
	delete(t.bucket.transactions, t.TransactionID)
}

func (t *Transaction) handleScrapeReply(data []byte) {
	data = data[8:]
	now := time.Now().Unix()
	for idx := range t.swarms {
		t.swarms[idx].Seeders = binary.BigEndian.Uint32(data[:idx*12])
		t.swarms[idx].Completed = binary.BigEndian.Uint32(data[:(idx*12)+4])
		t.swarms[idx].Leechers = binary.BigEndian.Uint32(data[:(idx*12)+8])
		t.swarms[idx].LastScrape = now
	}
}

// Sync syncs models with database
func (t *Transaction) Sync() (err error) {
	err = db.ORM.Update(t.swarms).Error
	return
}

// create send event
func (t *Transaction) SendEvent(to net.Addr) (ev *SendEvent) {
	ev = &SendEvent{
		To: to,
	}
	if t.state == stateRecvID {
		l := len(t.swarms) * 20
		l += 16

		ev.Data = make([]byte, l)

		binary.BigEndian.PutUint64(ev.Data[:], t.ConnectionID)
		binary.BigEndian.PutUint32(ev.Data[8:], 2)
		binary.BigEndian.PutUint32(ev.Data[12:], t.TransactionID)
		for idx := range t.swarms {
			ih, err := hex.DecodeString(t.swarms[idx].Hash)
			if err == nil && len(ih) == 20 {
				copy(ev.Data[16+(idx*20):], ih)
			}
		}
	} else if t.state == stateSendID {
		ev.Data = make([]byte, 16)
		binary.BigEndian.PutUint64(ev.Data, InitialConnectionID)
		binary.BigEndian.PutUint32(ev.Data[8:], 0)
		binary.BigEndian.PutUint32(ev.Data[12:], t.TransactionID)
		t.state = stateRecvID
	}
	return
}

func (t *Transaction) handleError(msg string) {
	log.Infof("scrape failed: %s", msg)
}

// handle data for transaction
func (t *Transaction) GotData(data []byte) (done bool) {

	if len(data) > 4 {
		cmd := binary.BigEndian.Uint32(data)
		switch cmd {
		case actionConnect:
			if len(data) == 16 {
				if t.state == stateRecvID {
					t.ConnectionID = binary.BigEndian.Uint64(data[8:])
				}
			}
			break
		case actionScrape:
			if len(data) == (12*len(t.swarms))+8 && t.state == stateTransact {
				t.handleScrapeReply(data[8:])
			}
			done = true
			break
		case actionError:
			if len(data) == 12 {
				t.handleError(string(data[4:12]))

			}
		default:
			done = true
		}
	}
	return
}
