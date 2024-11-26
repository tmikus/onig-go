package onig

/*
#include <oniguruma.h>
*/
import "C"

type Region struct {
	raw *C.OnigRegion
}

func NewRegion() *Region {
	return &Region{
		raw: C.onig_region_new(),
	}
}

func (r *Region) Len() int {
	return int(r.raw.num_regs)
}

func (r *Region) Pos(index int) *Range {
	if index >= r.Len() {
		return nil
	}
	begin := offsetInt(r.raw.beg, index)
	end := offsetInt(r.raw.end, index)
	if begin == C.ONIG_REGION_NOTPOS {
		return nil
	}
	return NewRange(int(begin), int(end))
}
