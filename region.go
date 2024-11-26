package onig

/*
#include <oniguruma.h>
*/
import "C"

type Region struct {
	raw *C.OnigRegion
}
