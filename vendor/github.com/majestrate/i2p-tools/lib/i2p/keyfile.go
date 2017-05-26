package i2p

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// a destination keypair file
type Keyfile struct {
	privkey string
	pubkey  string
	fname   string
}

// save to filesystem
func (k *Keyfile) Store() (err error) {
	if len(k.fname) > 0 {
		var f *os.File
		f, err = os.OpenFile(k.fname, os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			err = k.write(f)
			f.Close()
		}
	}
	return
}

// load from filesystem
func (k *Keyfile) Load() (err error) {
	if len(k.fname) > 0 {
		var f *os.File
		f, err = os.Open(k.fname)
		if err == nil {
			err = k.read(f)
			f.Close()
		}
	}
	return
}

func (k *Keyfile) write(w io.Writer) (err error) {
	_, err = fmt.Fprintf(w, "%s\n%s\n", k.privkey, k.pubkey)
	return
}

func (k *Keyfile) read(r io.Reader) (err error) {
	br := bufio.NewReader(r)
	k.privkey, err = br.ReadString(10)
	k.pubkey, err = br.ReadString(10)
	k.privkey = strings.Trim(k.privkey, "\n")
	k.pubkey = strings.Trim(k.pubkey, "\n")
	return
}

func (k *Keyfile) Addr() I2PAddr {
	return I2PAddr(k.pubkey)
}

// ensure keys are created using a control socket
func (k *Keyfile) ensure(nc net.Conn) (err error) {
	if len(k.fname) > 0 {
		_, err = os.Stat(k.fname)
	}
	if os.IsNotExist(err) || len(k.fname) == 0 {
		// no keyfile
		_, err = fmt.Fprintf(nc, "DEST GENERATE\n")
		r := bufio.NewReader(nc)
		var line string
		line, err = r.ReadString(10)
		if err == nil {
			sc := bufio.NewScanner(strings.NewReader(line))
			sc.Split(bufio.ScanWords)
			for sc.Scan() {
				txt := sc.Text()
				upper := strings.ToUpper(txt)
				if upper == "DEST" {
					continue
				}
				if upper == "REPLY" {
					continue
				}
				if strings.HasPrefix(upper, "PUB=") {
					k.pubkey = txt[4:]
					continue
				}
				if strings.HasPrefix(upper, "PRIV=") {
					k.privkey = txt[5:]
					continue
				}
			}
			// store new keys
			err = k.Store()
			return
		}
	}
	// load keys
	err = k.Load()
	return
}

// create new keyfile given filepath
func NewKeyfile(f string) *Keyfile {
	if strings.ToUpper(f) == "TRANSIENT" {
		f = ""
	}
	return &Keyfile{
		fname: f,
	}
}
