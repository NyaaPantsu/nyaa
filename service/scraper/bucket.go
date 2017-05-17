package scraperService

import (
	"math/rand"
	"net"
	"sync"

	"github.com/NyaaPantsu/nyaa/model"
)

const InitialConnectionID = 0x41727101980

type Bucket struct {
	Addr         net.Addr
	access       sync.Mutex
	transactions map[uint32]*Transaction
}

func (b *Bucket) NewTransaction(swarms []model.Torrent) (t *Transaction) {
	id := rand.Uint32()
	// get good id
	b.access.Lock()
	_, ok := b.transactions[id]
	for ok {
		id = rand.Uint32()
		_, ok = b.transactions[id]
	}
	t = &Transaction{
		TransactionID: id,
		bucket:        b,
		swarms:        make([]model.Torrent, len(swarms)),
		state:         stateSendID,
	}
	copy(t.swarms[:], swarms[:])
	b.transactions[id] = t
	b.access.Unlock()
	return

}

func (b *Bucket) ForEachTransaction(v func(uint32, *Transaction)) {

	clone := make(map[uint32]*Transaction)

	b.access.Lock()

	for k := range b.transactions {
		clone[k] = b.transactions[k]
	}

	b.access.Unlock()

	for k := range clone {
		v(k, clone[k])
	}
}

func (b *Bucket) Forget(tid uint32) {
	b.access.Lock()
	_, ok := b.transactions[tid]
	if ok {
		delete(b.transactions, tid)
	}
	b.access.Unlock()
}

func (b *Bucket) VisitTransaction(tid uint32, v func(*Transaction)) {
	b.access.Lock()
	t, ok := b.transactions[tid]
	b.access.Unlock()
	if ok {
		v(t)
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
