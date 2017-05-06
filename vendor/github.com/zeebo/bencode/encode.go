package bencode

import (
	"bytes"

	"fmt"
	"io"
	"reflect"
	"sort"
)

type sortValues []reflect.Value

func (p sortValues) Len() int           { return len(p) }
func (p sortValues) Less(i, j int) bool { return p[i].String() < p[j].String() }
func (p sortValues) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type sortFields []reflect.StructField

func (p sortFields) Len() int { return len(p) }
func (p sortFields) Less(i, j int) bool {
	iName, jName := p[i].Name, p[j].Name
	if name, _ := parseTag(p[i].Tag.Get("bencode")); name != "" {
		iName = name
	}
	if name, _ := parseTag(p[j].Tag.Get("bencode")); name != "" {
		jName = name
	}
	return iName < jName
}
func (p sortFields) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

//An Encoder writes bencoded objects to an output stream.
type Encoder struct {
	w io.Writer
}

//NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

//Encode writes the bencoded data of val to its output stream.
//See the documentation for Decode about the conversion of Go values to
//bencoded data.
func (e *Encoder) Encode(val interface{}) error {
	return encodeValue(e.w, reflect.ValueOf(val))
}

//EncodeString returns the bencoded data of val as a string.
func EncodeString(val interface{}) (string, error) {
	buf := new(bytes.Buffer)
	e := NewEncoder(buf)
	if err := e.Encode(val); err != nil {
		return "", err
	}
	return buf.String(), nil
}

//EncodeBytes returns the bencoded data of val as a slice of bytes.
func EncodeBytes(val interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	e := NewEncoder(buf)
	if err := e.Encode(val); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeValue(w io.Writer, val reflect.Value) error {
	//inspect the val to check
	v := indirect(val)

	//send in a raw message if we have that type
	if rm, ok := v.Interface().(RawMessage); ok {
		_, err := io.Copy(w, bytes.NewReader(rm))
		return err
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err := fmt.Fprintf(w, "i%de", v.Int())
		return err

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, err := fmt.Fprintf(w, "i%de", v.Uint())
		return err

	case reflect.Bool:
		i := 0
		if v.Bool() {
			i = 1
		}
		_, err := fmt.Fprintf(w, "i%de", i)
		return err

	case reflect.String:
		_, err := fmt.Fprintf(w, "%d:%s", len(v.String()), v.String())
		return err

	case reflect.Slice, reflect.Array:
		// handle byte slices like strings
		if byteSlice, ok := val.Interface().([]byte); ok {
			_, err := fmt.Fprintf(w, "%d:", len(byteSlice))

			if err == nil {
				_, err = w.Write(byteSlice)
			}

			return err
		}

		if _, err := fmt.Fprint(w, "l"); err != nil {
			return err
		}

		for i := 0; i < v.Len(); i++ {
			if err := encodeValue(w, v.Index(i)); err != nil {
				return err
			}
		}

		_, err := fmt.Fprint(w, "e")
		return err

	case reflect.Map:
		if _, err := fmt.Fprint(w, "d"); err != nil {
			return err
		}
		var (
			keys sortValues = v.MapKeys()
			mval reflect.Value
		)
		sort.Sort(keys)
		for i := range keys {
			if err := encodeValue(w, keys[i]); err != nil {
				return err
			}
			mval = v.MapIndex(keys[i])
			if err := encodeValue(w, mval); err != nil {
				return err
			}
		}
		_, err := fmt.Fprint(w, "e")
		return err

	case reflect.Struct:
		t := v.Type()
		if _, err := fmt.Fprint(w, "d"); err != nil {
			return err
		}
		//put keys into keys
		var (
			keys       = make(sortFields, t.NumField())
			fieldValue reflect.Value
			rkey       reflect.Value
		)
		for i := range keys {
			keys[i] = t.Field(i)
		}
		sort.Sort(keys)
		for _, key := range keys {
			rkey = reflect.ValueOf(key.Name)
			fieldValue = v.FieldByIndex(key.Index)

			// filter out unexported values etc.
			if !fieldValue.CanInterface() {
				continue
			}

			/* Tags
			* Near identical to usage in JSON except with key 'bencode'

			* Struct values encode as BEncode dictionaries. Each exported
			  struct field becomes a set in the dictionary unless
			  - the field's tag is "-", or
			  - the field is empty and its tag specifies the "omitempty"
			    option.

			* The default key string is the struct field name but can be
			  specified in the struct field's tag value.  The "bencode"
			  key in struct field's tag value is the key name, followed
			  by an optional comma and options.
			*/
			tagValue := key.Tag.Get("bencode")
			if tagValue != "" {
				// Keys with '-' are omit from output
				if tagValue == "-" {
					continue
				}

				name, options := parseTag(tagValue)
				// Keys with 'omitempty' are omitted if the field is empty
				if options.Contains("omitempty") && isEmptyValue(fieldValue) {
					continue
				}

				// All other values are treated as the key string
				if isValidTag(name) {
					rkey = reflect.ValueOf(name)
				}
			}

			//encode the key
			if err := encodeValue(w, rkey); err != nil {
				return err
			}
			//encode the value
			if err := encodeValue(w, fieldValue); err != nil {
				return err
			}
		}
		_, err := fmt.Fprint(w, "e")
		return err
	}

	return fmt.Errorf("Can't encode type: %s", v.Type())
}
