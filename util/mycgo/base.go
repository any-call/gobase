package mycgo

import "C"

func ToString(in *C.char) string {
	return C.GoString(in)
}

func ToPTRChar(in string) *C.char {
	return C.CString(in)
}
