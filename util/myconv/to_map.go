package myconv

import (
	"fmt"
	"reflect"
)

func Struct2Map(a any) (map[string]any, error) {
	a = DirectObj(a)

	v := reflect.ValueOf(a)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input type isn't struct: %T", a)
	}

	t := reflect.TypeOf(a)
	m := make(map[string]any, 10)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() {
			m[t.Field(i).Name] = v.Field(i).Interface()
		} else {
			fmt.Printf("field: name= %v,v= %v \n", t.Field(i).Name, v.Field(i).String())
		}
	}

	return m, nil
}

func Map2Slice[K comparable, V any](mp map[K]V) (ks []K, vs []V) {
	ks = make([]K, 0, len(mp))
	vs = make([]V, 0, len(mp))

	for k, v := range mp {
		ks = append(ks, k)
		vs = append(vs, v)
	}

	return
}
