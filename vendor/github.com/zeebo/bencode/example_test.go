package bencode

import (
	"fmt"
	"io"
)

var (
	data string
	r    io.Reader
	w    io.Writer
)

func ExampleDecodeString() {
	var torrent interface{}
	if err := DecodeString(data, &torrent); err != nil {
		panic(err)
	}
}

func ExampleEncodeString() {
	var torrent interface{}
	data, err := EncodeString(torrent)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func ExampleDecodeBytes() {
	var torrent interface{}
	if err := DecodeBytes([]byte(data), &torrent); err != nil {
		panic(err)
	}
}

func ExampleEncodeBytes() {
	var torrent interface{}
	data, err := EncodeBytes(torrent)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func ExampleEncoder_Encode() {
	var x struct {
		Foo string
		Bar []string `bencode:"name"`
	}

	enc := NewEncoder(w)
	if err := enc.Encode(x); err != nil {
		panic(err)
	}
}

func ExampleDecoder_Decode() {
	dec := NewDecoder(r)
	var torrent struct {
		Announce string
		List     [][]string `bencode:"announce-list"`
	}
	if err := dec.Decode(&torrent); err != nil {
		panic(err)
	}
}
