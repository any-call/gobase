package myconv

import (
	"bytes"
	"encoding/gob"
)

func Obj2Stream(obj any) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func Stream2Obj(stream []byte, obj any) error {
	dec := gob.NewDecoder(bytes.NewReader(stream))
	return dec.Decode(obj)
}
