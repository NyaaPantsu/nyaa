package i2p

import (
	"crypto/sha256"
)

// implements net.Addr
type I2PAddr string

func (a I2PAddr) Network() string {
	return "i2p"
}

func (a I2PAddr) String() string {
	return string(a)
}

// compute base32 address
func (a I2PAddr) Base32Addr() (b32 Base32Addr) {
	buf := make([]byte, i2pB64enc.DecodedLen(len(a)))
	if _, err := i2pB64enc.Decode(buf, []byte(a)); err != nil {
		return
	}
	h := sha256.New()
	h.Write(buf)
	d := h.Sum(nil)
	copy(b32[:], d)
	return
}

// i2p destination hash
type Base32Addr [32]byte

// get string version
func (b32 Base32Addr) String() string {
	b32addr := make([]byte, 56)
	i2pB32enc.Encode(b32addr, b32[:])
	return string(b32addr[:52]) + ".b32.i2p"
}
