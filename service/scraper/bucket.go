package scraperService

import (
	"math/rand"
	"net"

	"github.com/ewhal/nyaa/model"
)

const InitialConnectionID = 0x41727101980

type Bucket struct {
	Addr         net.Addr
	transactions map[uint32]*Transaction
}

func (b *Bucket) NewTransaction(swarms []model.Torrent) (t *Transaction) {
	id := rand.Uint32()
	// get good id
	_, ok := b.transactions[id]
	for ok {
		id = rand.Uint32()
		_, ok = b.transactions[id]
	}
	t = &Transaction{
		TransactionID: id,
		swarms:        swarms,
		state:         stateSendID,
	}
	b.transactions[id] = t
	return

}

func (b *Bucket) VisitTransaction(tid uint32, v func(*Transaction)) {
	t, ok := b.transactions[tid]
	if ok {
		go v(t)
	} else {
		v(nil)
	}
}

func NewBucket(a net.Addr) *Bucket {
	return &Bucket{
		transactions: make(map[uint32]*Transaction),
		Addr:         a,
	}
}
