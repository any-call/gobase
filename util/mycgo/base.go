package mycgo

/*
#include <stdio.h>
#include <stdlib.h>
void print(char *str) {
    printf("%s\n", str);
}
*/
import "C"

type (
	PtrChar *C.char
)

func ToString(in PtrChar) string {
	return C.GoString(in)
}

func ToPTRChar(in string) PtrChar {
	return C.CString(in)
}

func CPrintln(char PtrChar) {
	C.print(char)
}
