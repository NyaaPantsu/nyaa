package bencode

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type testCase struct {
		in     string
		val    interface{}
		expect interface{}
		err    bool
	}

	type dT struct {
		X string
		Y int
		Z string `bencode:"zff"`
	}

	var decodeCases = []testCase{
		//integers
		{`i5e`, new(int), int(5), false},
		{`i-10e`, new(int), int(-10), false},
		{`i8e`, new(uint), uint(8), false},
		{`i8e`, new(uint8), uint8(8), false},
		{`i8e`, new(uint16), uint16(8), false},
		{`i8e`, new(uint32), uint32(8), false},
		{`i8e`, new(uint64), uint64(8), false},
		{`i8e`, new(int), int(8), false},
		{`i8e`, new(int8), int8(8), false},
		{`i8e`, new(int16), int16(8), false},
		{`i8e`, new(int32), int32(8), false},
		{`i8e`, new(int64), int64(8), false},
		{`i-2e`, new(uint), nil, true},

		//bools
		{`i1e`, new(bool), true, false},
		{`i0e`, new(bool), false, false},
		{`i8e`, new(bool), true, false},

		//strings
		{`3:foo`, new(string), "foo", false},
		{`4:foob`, new(string), "foob", false},
		{`6:short`, new(string), nil, true},

		//lists
		{`l3:foo3:bare`, new([]string), []string{"foo", "bar"}, false},
		{`li15ei20ee`, new([]int), []int{15, 20}, false},
		{`ld3:fooi0eed3:bari1eee`, new([]map[string]int), []map[string]int{
			{"foo": 0},
			{"bar": 1},
		}, false},

		//dicts
		{`d3:foo3:bar4:foob3:fooe`, new(map[string]string), map[string]string{
			"foo":  "bar",
			"foob": "foo",
		}, false},
		{`d1:X3:foo1:Yi10e3:zff3:bare`, new(dT), dT{"foo", 10, "bar"}, false},
		{`d1:X3:foo1:Yi10e1:Z3:bare`, new(dT), dT{"foo", 10, "bar"}, false},
		{`d1:X3:foo1:Yi10e1:h3:bare`, new(dT), dT{"foo", 10, ""}, false},
		{`d3:fooli0ei1ee3:barli2ei3eee`, new(map[string][]int), map[string][]int{
			"foo": []int{0, 1},
			"bar": []int{2, 3},
		}, false},
		{`de`, new(map[string]string), map[string]string{}, false},

		//into interfaces
		{`i5e`, new(interface{}), int64(5), false},
		{`li5ee`, new(interface{}), []interface{}{int64(5)}, false},
		{`5:hello`, new(interface{}), "hello", false},
		{`d5:helloi5ee`, new(interface{}), map[string]interface{}{"hello": int64(5)}, false},

		//malformed
		{`i53:foo`, new(interface{}), nil, true},
		{`6:foo`, new(interface{}), nil, true},
		{`di5ei2ee`, new(interface{}), nil, true},
		{`d3:fooe`, new(interface{}), nil, true},
		{`l3:foo3:bar`, new(interface{}), nil, true},
		{`d-1:`, new(interface{}), nil, true},
	}

	for i, tt := range decodeCases {
		err := DecodeString(tt.in, tt.val)
		if !tt.err && err != nil {
			t.Errorf("#%d: Unexpected err: %v", i, err)
			continue
		}
		if tt.err && err == nil {
			t.Errorf("#%d: Expected err is nil", i)
			continue
		}
		v := reflect.ValueOf(tt.val).Elem().Interface()
		if !reflect.DeepEqual(v, tt.expect) && !tt.err {
			t.Errorf("#%d: Val: %#v != %#v", i, v, tt.expect)
		}
	}
}

func TestRawDecode(t *testing.T) {
	type testCase struct {
		in     string
		expect []byte
		err    bool
	}

	var rawDecodeCases = []testCase{
		{`i5e`, []byte(`i5e`), false},
		{`5:hello`, []byte(`5:hello`), false},
		{`li5ei10e5:helloe`, []byte(`li5ei10e5:helloe`), false},
		{`llleee`, []byte(`llleee`), false},
		{`li5eli5eli5eeee`, []byte(`li5eli5eli5eeee`), false},
		{`d5:helloi5ee`, []byte(`d5:helloi5ee`), false},
	}

	for i, tt := range rawDecodeCases {
		var x RawMessage
		err := DecodeString(tt.in, &x)
		if !tt.err && err != nil {
			t.Errorf("#%d: Unexpected err: %v", i, err)
			continue
		}
		if tt.err && err == nil {
			t.Errorf("#%d: Expected err is nil", i)
			continue
		}
		if !reflect.DeepEqual(x, RawMessage(tt.expect)) && !tt.err {
			t.Errorf("#%d: Val: %#v != %#v", i, x, tt.expect)
		}
	}
}

func TestNestedRawDecode(t *testing.T) {
	type testCase struct {
		in     string
		val    interface{}
		expect interface{}
		err    bool
	}

	type message struct {
		Key string
		Val int
		Raw RawMessage
	}

	var cases = []testCase{
		{`li5e5:hellod1:a1:beli5eee`, new([]RawMessage), []RawMessage{
			RawMessage(`i5e`),
			RawMessage(`5:hello`),
			RawMessage(`d1:a1:be`),
			RawMessage(`li5ee`),
		}, false},
		{`d1:a1:b1:c1:de`, new(map[string]RawMessage), map[string]RawMessage{
			"a": RawMessage(`1:b`),
			"c": RawMessage(`1:d`),
		}, false},
		{`d3:Key5:hello3:Rawldedei5e1:ae3:Vali10ee`, new(message), message{
			Key: "hello",
			Val: 10,
			Raw: RawMessage(`ldedei5e1:ae`),
		}, false},
	}

	for i, tt := range cases {
		err := DecodeString(tt.in, tt.val)
		if !tt.err && err != nil {
			t.Errorf("#%d: Unexpected err: %v", i, err)
			continue
		}
		if tt.err && err == nil {
			t.Errorf("#%d: Expected err is nil", i)
			continue
		}
		v := reflect.ValueOf(tt.val).Elem().Interface()
		if !reflect.DeepEqual(v, tt.expect) && !tt.err {
			t.Errorf("#%d: Val:\n%#v !=\n%#v", i, v, tt.expect)
		}
	}
}
