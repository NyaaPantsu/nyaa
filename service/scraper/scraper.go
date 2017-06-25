package scraperService

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util/log"
)

// MTU yes this is the ipv6 mtu
const MTU = 1500

// max number of scrapes per packet
const ScrapesPerPacket = 74

// bittorrent scraper
type Scraper struct {
	done             chan int
	sendQueue        chan *SendEvent
	recvQueue        chan *RecvEvent
	errQueue         chan error
	trackers         map[string]*Bucket
	ticker           *time.Ticker
	cleanup          *time.Ticker
	interval         time.Duration
	PacketsPerSecond uint
}

func New(conf *config.ScraperConfig) (sc *Scraper, err error) {
	sc = &Scraper{
		done:      make(chan int),
		sendQueue: make(chan *SendEvent, 1024),
		recvQueue: make(chan *RecvEvent, 1024),
		errQueue:  make(chan error),
		trackers:  make(map[string]*Bucket),
		ticker:    time.NewTicker(time.Second * 10),
		interval:  time.Second * time.Duration(conf.IntervalSeconds),
		cleanup:   time.NewTicker(time.Minute),
	}

	if sc.PacketsPerSecond == 0 {
		sc.PacketsPerSecond = 10
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
		if err != nil {
			log.Warnf("failed: %s", err)
			break
		}
		action, err := ev.Action()
		if err != nil {
			log.Warnf("failed: %s", err)
			break
		}
		log.Debugf("transaction = %d action = %d", tid, action)
		bucket, ok = sc.trackers[ev.From.String()]
		if !ok || bucket == nil {
			log.Warnf("bucket not found for %s", ev.From)
			break
		}

		bucket.VisitTransaction(tid, func(t *Transaction) {
			if t == nil {
				log.Warnf("no transaction %d", tid)
				return
			}
			if t.GotData(ev.Data) {
				err := t.Sync()
				if err != nil {
					log.Warnf("failed to sync swarm: %s", err)
				}
				t.Done()
				log.Debugf("transaction %d done", tid)
			} else {
				sc.sendQueue <- t.SendEvent(ev.From)
			}
		})
	}
	return
}

func (sc *Scraper) Run() {
	for {
		select {
		case <-sc.ticker.C:
			sc.Scrape(sc.PacketsPerSecond)
			break
		case <-sc.cleanup.C:
			sc.removeStale()
			break
		}
	}
}

func (sc *Scraper) removeStale() {

	for k := range sc.trackers {
		sc.trackers[k].ForEachTransaction(func(tid uint32, t *Transaction) {
			if t == nil || t.IsTimedOut() {
				sc.trackers[k].Forget(tid)
			}
		})
	}
}

func (sc *Scraper) Scrape(packets uint) {
	now := time.Now().Add(0 - sc.interval)
	// only scrape torretns uploaded within 90 days
	oldest := now.Add(0 - (time.Hour * 24 * 90))

	query := fmt.Sprintf(
		"SELECT * FROM ("+

		// previously scraped torrents that will be scraped again:
		"SELECT %[1]s.torrent_id, torrent_hash FROM %[1]s, %[2]s WHERE "+
		"date > ? AND "+
		"%[1]s.torrent_id = %[2]s.torrent_id AND "+
		"last_scrape < ?"+

		// torrents that weren't scraped before:
		" UNION "+
		"SELECT torrent_id, torrent_hash FROM %[1]s WHERE "+
		"date > ? AND "+
		"torrent_id NOT IN (SELECT torrent_id FROM %[2]s)"+

		") AS x ORDER BY torrent_id DESC LIMIT ?",
		config.Conf.Models.TorrentsTableName, config.Conf.Models.ScrapeTableName)
	rows, err := db.ORM.Raw(query, oldest, now, oldest, packets*ScrapesPerPacket).Rows()

	if err == nil {
		counter := 0
		var scrape [ScrapesPerPacket]model.Torrent
		for rows.Next() {
			idx := counter % ScrapesPerPacket
			rows.Scan(&scrape[idx].ID, &scrape[idx].Hash)
			counter++
			if counter%ScrapesPerPacket == 0 {
				for _, b := range sc.trackers {
					t := b.NewTransaction(scrape[:])
					sc.sendQueue <- t.SendEvent(b.Addr)
				}
			}
		}
		idx := counter % ScrapesPerPacket
		if idx > 0 {
			for _, b := range sc.trackers {
				t := b.NewTransaction(scrape[:idx])
				sc.sendQueue <- t.SendEvent(b.Addr)
			}
		}
		log.Infof("scrape %d", counter)
		rows.Close()
	} else {
		log.Warnf("failed to select torrents for scrape: %s", err)
	}
}

func (sc *Scraper) Wait() {
	<-sc.done
}
