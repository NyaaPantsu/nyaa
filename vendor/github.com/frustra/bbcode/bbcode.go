// Copyright 2015 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package bbcode implements a parser and HTML generator for BBCode.
package bbcode

import "sort"

type BBOpeningTag struct {
	Name  string
	Value string
	Args  map[string]string
	Raw   string
}

type BBClosingTag struct {
	Name string
	Raw  string
}

func (t *BBOpeningTag) String() string {
	str := t.Name
	if len(t.Value) > 0 {
		str += "=" + t.Value
	}
	keys := make([]string, len(t.Args))
	i := 0
	for key := range t.Args {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		v := t.Args[key]
		str += " " + key
		if len(v) > 0 {
			str += "=" + v
		}
	}
	return str
}
