package onig

/*
#include <oniguruma.h>
*/
import "C"

// Region represents a set of capture groups found in a search or match.
type Region struct {
	raw *C.OnigRegion
}

// NewRegion creates a new empty Region.
func NewRegion() *Region {
	return &Region{
		raw: C.onig_region_new(),
	}
}

// Clear can be used to clear out a region so it can be used again. See [onig_sys::onig_region_clear](https://docs.rs/onig/latest/onig/onig_sys/fn.onig_region_clear.html)
func (r *Region) Clear() {
	C.onig_region_clear(r.raw)
}

// Len returns the number of registers in the region.
func (r *Region) Len() int {
	return int(r.raw.num_regs)
}

// Pos returns the start and end positions of the Nth capture group.
// Returns nil if index is not a valid capture group or if the capture group did not match anything.
// The positions returned are always byte indices with respect to the original string matched.
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
