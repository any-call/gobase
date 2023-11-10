package myjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

type (
	FieldValue struct {
		Key   string
		Value any
	}

	orderArray []any

	OrderMap struct {
		list []FieldValue
	}
)

func ToOrderMap(jsonStr string) (list []FieldValue, err error) {
	orderMap := &OrderMap{}
	if err := orderMap.UnmarshalJSON([]byte(jsonStr)); err != nil {
		return nil, err
	}

	return orderMap.Fields(), nil
}

func (self *OrderMap) Fields() []FieldValue {
	if self.list == nil {
		self.list = []FieldValue{}
	}

	return self.list
}

// implementation of interface Marshaler
func (self *OrderMap) MarshalJSON() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.WriteRune('{')

	if self.list == nil {
		self.list = []FieldValue{}
	}

	for i, val := range self.list {
		// write key
		b, e := json.Marshal(val.Key)
		if e != nil {
			return nil, e
		}
		buf.Write(b)

		// write delimiter
		buf.WriteRune(':')

		// write value
		b, e = json.Marshal(val.Value)
		if e != nil {
			return nil, e
		}
		buf.Write(b)

		// write delimiter
		if i+1 < len(self.list) {
			buf.WriteRune(',')
		}
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}

// implementation of interface Unmarshaler
func (self *OrderMap) UnmarshalJSON(b []byte) error {
	d := json.NewDecoder(bytes.NewReader(b))
	t, err := d.Token()
	if err == io.EOF {
		return nil
	}

	// 识别是否为对象
	if t != json.Delim('{') {
		// log.Print("unexpected start of object")
		return errors.New("unexpected start of object")
	}
	return self.unmarshalEmbededObject(d)
}

func (self *OrderMap) unmarshalEmbededObject(d *json.Decoder) error {
	for d.More() {
		kToken, err := d.Token()
		if err == io.EOF || (err == nil && kToken == json.Delim('}')) {
			// log.Print("unexpected EOF")
			return errors.New("unexpected EOF")
		}

		vToken, err := d.Token()
		if err == io.EOF {
			// log.Print("unexpected EOF")
			return errors.New("unexpected EOF")
		}

		var val interface{}
		switch vToken {
		case json.Delim('{'):
			var obj OrderMap
			if err = obj.unmarshalEmbededObject(d); err != nil {
				return err
			}
			val = obj
		case json.Delim('['):
			var arr orderArray
			err = arr.unmarshalEmbededArray(d)
			val = arr
		default:
			val = vToken
		}

		if err != nil {
			return err
		}

		if self.list == nil {
			self.list = []FieldValue{{
				Key:   kToken.(string),
				Value: val,
			}}
		} else {
			self.list = append(self.list, FieldValue{kToken.(string), val})
		}
	}

	// 读取对象结束token '}'
	kToken, err := d.Token()
	if err == io.EOF || kToken != json.Delim('}') {
		// log.Print("unexpected EOF")
		return errors.New("unexpected EOF")
	}

	return err
}

func (self *orderArray) unmarshalEmbededArray(d *json.Decoder) error {
	for d.More() {
		token, err := d.Token()
		if err == io.EOF || (err == nil && token == json.Delim(']')) {
			// log.Print("unexpected EOF")
			return errors.New("unexpected EOF")
		}

		var val interface{}
		switch token {
		case json.Delim('{'):
			var obj OrderMap
			if err = obj.unmarshalEmbededObject(d); err != nil {
				return err
			}
			val = obj
		case json.Delim('['):
			var arr orderArray
			err = arr.unmarshalEmbededArray(d)
			val = arr
		default:
			// 字面量 literial
			val = token
		}

		if err != nil {
			return err
		}

		*self = append(*self, val)
	}

	// 读取数组结束token ']'
	kToken, err := d.Token()
	if err == io.EOF || kToken != json.Delim(']') {
		return errors.New("unexpected EOF")
	}

	if *self == nil {
		*self = orderArray(make([]interface{}, 0))
	}

	return nil
}
