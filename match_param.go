package onig

/*
#include <oniguruma.h>
*/
import "C"

type MatchParam struct {
	raw *C.OnigMatchParam
}

func NewMatchParam() *MatchParam {
	raw := C.onig_new_match_param()
	C.onig_initialize_match_param(raw)
	return &MatchParam{
		raw: raw,
	}
}
