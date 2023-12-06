package mycgo

/*
#include <stdio.h>
#include <stdlib.h>
void print(char *str) {
    printf("%s\n", str);
}
*/
import "C"

func ToString(in *C.char) string {
	return C.GoString(in)
}

func ToPTRChar(in string) *C.char {
	return C.CString(in)
}

func CPrintln(char *C.char) {
	C.print(char)
}
