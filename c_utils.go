package onig

/*
#include <stdlib.h>
#include <string.h>
#include <oniguruma.h>
*/
import "C"
import "unsafe"

func getUChar(s *C.char) *C.UChar {
	return (*C.UChar)(unsafe.Pointer(s))
}

func getUCharEnd(s *C.char) *C.UChar {
	return offsetUChar(getUChar(s), int(C.strlen(s)))
}

func getPointer[T any](s *T) uintptr {
	return uintptr(unsafe.Pointer(s))
}

func offsetInt(i *C.int, offset int) C.int {
	return *(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(i)) + uintptr(offset)*unsafe.Sizeof(C.int(0))))
}

func offsetUChar(s *C.UChar, offset int) *C.UChar {
	return (*C.UChar)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + uintptr(offset)))
}
