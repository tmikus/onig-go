package onig

/*
#include <oniguruma.h>
*/
import "C"
import "runtime"

// Region represents a set of capture groups found in a search or match.
type Region struct {
	raw   *C.OnigRegion
	regex *Regex
}

// NewRegion creates a new empty Region.
func NewRegion() *Region {
	region := &Region{
		raw: C.onig_region_new(),
	}
	runtime.SetFinalizer(region, func(region *Region) {
		C.onig_region_free(region.raw, 0)
	})
	return region
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
// Returns nil if the capture group did not match anything or if index is not a valid capture group.
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

// PosByGroupName returns the start and end positions of the named capture group.
// Returns nil if the capture group did not match anything or if groupName is not a valid capture group.
// The positions returned are always byte indices with respect to the original string matched.
func (r *Region) PosByGroupName(groupName string) *Range {
	groupIndices := r.regex.GetGroupNumbersForGroupName(groupName)
	for _, groupIndex := range groupIndices {
		begin := offsetInt(r.raw.beg, groupIndex)
		end := offsetInt(r.raw.end, groupIndex)
		if begin == C.ONIG_REGION_NOTPOS {
			continue
		}
		return NewRange(int(begin), int(end))
	}
	return nil
}
