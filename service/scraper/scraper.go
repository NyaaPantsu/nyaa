package scraperService

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
	"net"
	"net/url"
	"time"
)

// MTU yes this is the ipv6 mtu
const MTU = 1488

// bittorrent scraper
type Scraper struct {
	done      chan int
	sendQueue chan *SendEvent
	recvQueue chan *RecvEvent
	errQueue  chan error
	trackers  map[string]*Bucket
	ticker    *time.Ticker
	interval  time.Duration
}

func New(conf *config.ScraperConfig) (sc *Scraper, err error) {
	sc = &Scraper{
		done:      make(chan int),
		sendQueue: make(chan *SendEvent, 128),
		recvQueue: make(chan *RecvEvent, 1028),
		errQueue:  make(chan error),
		trackers:  make(map[string]*Bucket),
		ticker:    time.NewTicker(time.Minute),
		interval:  time.Second * time.Duration(conf.IntervalSeconds),
	}
	for idx := range conf.Trackers {
		err = sc.AddTracker(&conf.Trackers[idx])
		if err != nil {
			break
		}
	}
	return
}

func (sc *Scraper) AddTracker(conf *config.ScrapeConfig) (err error) {
	var u *url.URL
	u, err = url.Parse(conf.URL)
	if err == nil {
		var ips []net.IP
		ips, err = net.LookupIP(u.Hostname())
		if err == nil {
			// TODO: use more than 1 ip ?
			addr := &net.UDPAddr{
				IP: ips[0],
			}
			addr.Port, err = net.LookupPort("udp", u.Port())
			if err == nil {
				sc.trackers[addr.String()] = NewBucket(addr)
			}
		}
	}
	return
}

func (sc *Scraper) Close() (err error) {
	close(sc.sendQueue)
	close(sc.recvQueue)
	close(sc.errQueue)
	sc.ticker.Stop()
	sc.done <- 1
	return
}

func (sc *Scraper) runRecv(pc net.PacketConn) {
	for {
		var buff [MTU]byte
		n, from, err := pc.ReadFrom(buff[:])

		if err == nil {

			log.Debugf("got %d from %s", n, from)
			sc.recvQueue <- &RecvEvent{
				From: from,
				Data: buff[:n],
			}
		} else {
			sc.errQueue <- err
		}
	}
}

func (sc *Scraper) runSend(pc net.PacketConn) {
	for {
		ev, ok := <-sc.sendQueue
		if !ok {
			return
		}
		log.Debugf("write %d to %s", len(ev.Data), ev.To)
		pc.WriteTo(ev.Data, ev.To)
	}
}

func (sc *Scraper) RunWorker(pc net.PacketConn) (err error) {

	go sc.runRecv(pc)
	go sc.runSend(pc)
	for {
		var bucket *Bucket
		ev, ok := <-sc.recvQueue
		if !ok {
			break
		}
		tid, err := ev.TID()
		action, err := ev.Action()
		log.Debugf("transaction = %d action = %d", tid, action)
		if err == nil {
			bucket, ok = sc.trackers[ev.From.String()]
			if ok && bucket != nil {
				bucket.VisitTransaction(tid, func(t *Transaction) {
					if t == nil {
						log.Warnf("no transaction %d", tid)
					} else {
						if t.GotData(ev.Data) {
							err := t.Sync()
							if err != nil {
								log.Warnf("failed to sync swarm: %s", err)
							}
							t.Done()
						} else {
							sc.sendQueue <- t.SendEvent(ev.From)
						}
					}
				})
			} else {
				log.Warnf("bucket not found for %s", ev.From)
			}
		}

	}
	return
}

func (sc *Scraper) Run() {
	sc.Scrape()
	for {
		<-sc.ticker.C
		sc.Scrape()
	}
}

func (sc *Scraper) Scrape() {

	swarms := make([]model.Torrent, 0, 128)
	now := time.Now().Add(0 - sc.interval).Unix()
	err := db.ORM.Where("last_scrape < ?", now).Or("last_scrape IS NULL").Find(&swarms).Error
	if err == nil {
		for swarms != nil {
			var scrape []model.Torrent
			if len(swarms) > 74 {
				scrape = swarms[:74]
				swarms = swarms[74:]
			} else {
				scrape = swarms
				swarms = nil
			}
			log.Infof("scraping %d", len(scrape))
			if len(scrape) > 0 {
				for _, b := range sc.trackers {
					t := b.NewTransaction(scrape)
					log.Debugf("new transaction %d", t.TransactionID)
					sc.sendQueue <- t.SendEvent(b.Addr)
				}
			}
		}
	} else {
		log.Warnf("failed to select torrents for scrape: %s", err)
	}
}

func (sc *Scraper) Wait() {
	<-sc.done
}
