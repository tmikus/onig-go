package onig

/*
#include <oniguruma.h>
*/
import "C"

type MatchParam struct {
	raw *C.OnigMatchParam
}
