package onig

/*
#include <stdlib.h>
#include <string.h>
#include <oniguruma.h>
*/
import "C"
import "unsafe"

func offsetInt(i *C.int, offset int) C.int {
	return *(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(i)) + uintptr(offset)*unsafe.Sizeof(C.int(0))))
}
