package myjson

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonSlice[V any] []V

func (s *JsonSlice[V]) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, s)
}

func (s JsonSlice[V]) Value() (driver.Value, error) {
	return json.Marshal(s)
}
