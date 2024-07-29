package myjson

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonObject[V any] struct {
	Object V
}

func (s *JsonObject[V]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &(s.Object))
}

func (s JsonObject[V]) Value() (driver.Value, error) {
	return json.Marshal(s.Object)
}
