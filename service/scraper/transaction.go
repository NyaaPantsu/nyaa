package scraperService

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util/log"
)

// TransactionTimeout 30 second timeout for transactions
const TransactionTimeout = time.Second * 30

const stateSendID = 0
const stateRecvID = 1
const stateTransact = 2

const actionError = 3
const actionScrape = 2

//const actionAnnounce = 1
const actionConnect = 0

// Transaction a scrape transaction on a udp tracker
type Transaction struct {
	TransactionID uint32
	ConnectionID  uint64
	bucket        *Bucket
	state         uint8
	swarms        []model.Torrent
	lastData      time.Time
}

// Done marks this transaction as done and removes it from parent
func (t *Transaction) Done() {
	t.bucket.Forget(t.TransactionID)
}

func (t *Transaction) handleScrapeReply(data []byte) {
	data = data[8:]
	now := time.Now()
	for idx := range t.swarms {
		t.swarms[idx].Scrape = &model.Scrape{}
		t.swarms[idx].Scrape.Seeders = binary.BigEndian.Uint32(data)
		data = data[4:]
		t.swarms[idx].Scrape.Completed = binary.BigEndian.Uint32(data)
		data = data[4:]
		t.swarms[idx].Scrape.Leechers = binary.BigEndian.Uint32(data)
		data = data[4:]
		t.swarms[idx].Scrape.LastScrape = now
		idx++
	}
}


var pgQuery = "INSERT INTO " + config.Conf.Models.ScrapeTableName + " (torrent_id, seeders, leechers, completed, last_scrape) VALUES ($1, $2, $3, $4, $5) "+
	"ON CONFLICT (torrent_id) DO UPDATE SET seeders=EXCLUDED.seeders, leechers=EXCLUDED.leechers, completed=EXCLUDED.completed, last_scrape=EXCLUDED.last_scrape"
var sqliteQuery = "REPLACE INTO " + config.Conf.Models.ScrapeTableName + " (torrent_id, seeders, leechers, completed, last_scrape) VALUES (?, ?, ?, ?, ?)"

// Sync syncs models with database
func (t *Transaction) Sync() (err error) {
	q := pgQuery
	if db.IsSqlite {
		q = sqliteQuery
	}
	tx, e := db.ORM.DB().Begin()
	err = e
	if err == nil {
		for idx := range t.swarms {
			_, err = tx.Exec(q, t.swarms[idx].ID, t.swarms[idx].Scrape.Seeders, t.swarms[idx].Scrape.Leechers, t.swarms[idx].Scrape.Completed, t.swarms[idx].Scrape.LastScrape)
		}
		tx.Commit()
	}
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
		t.state = stateTransact
	} else if t.state == stateSendID {
		ev.Data = make([]byte, 16)
		binary.BigEndian.PutUint64(ev.Data, InitialConnectionID)
		binary.BigEndian.PutUint32(ev.Data[8:], 0)
		binary.BigEndian.PutUint32(ev.Data[12:], t.TransactionID)
		t.state = stateRecvID
	}
	t.lastData = time.Now()
	return
}

func (t *Transaction) handleError(msg string) {
	log.Infof("scrape failed: %s", msg)
}

// handle data for transaction
func (t *Transaction) GotData(data []byte) (done bool) {
	t.lastData = time.Now()
	if len(data) > 4 {
		cmd := binary.BigEndian.Uint32(data)
		switch cmd {
		case actionConnect:
			if len(data) == 16 {
				if t.state == stateRecvID {
					t.ConnectionID = binary.BigEndian.Uint64(data[8:])
				}
			}
		case actionScrape:
			if len(data) == (12*len(t.swarms))+8 && t.state == stateTransact {
				t.handleScrapeReply(data)
			}
			done = true
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

func (t *Transaction) IsTimedOut() bool {
	return t.lastData.Add(TransactionTimeout).Before(time.Now())

}
