package onig

/*
#include "regex.h"
*/
import "C"
import (
	"runtime"
)

// Region represents a set of capture groups found in a search or match.
type Region struct {
	raw   *C.region
	regex *Regex
}

// newRegion creates a new empty Region.
func newRegion(regex *Regex, raw *C.region) *Region {
	region := &Region{
		raw:   raw,
		regex: regex,
	}
	runtime.SetFinalizer(region, func(region *Region) {
		if region.raw != nil {
			C.freeRegion(region.raw)
			region.raw = nil
		}
	})
	return region
}

// Len returns the number of registers in the region.
func (r *Region) Len() int {
	return int(r.raw.groupCount)
}

// Pos returns the start and end positions of the Nth capture group.
// Returns nil if the capture group did not match anything or if index is not a valid capture group.
// The positions returned are always byte indices with respect to the original string matched.
func (r *Region) Pos(index int) *Range {
	if index >= r.Len() {
		return nil
	}
	begin := offsetInt(r.raw.groupStartIndices, index)
	end := offsetInt(r.raw.groupEndIndices, index)
	if begin == -1 || end == -1 {
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
		begin := offsetInt(r.raw.groupStartIndices, groupIndex)
		end := offsetInt(r.raw.groupEndIndices, groupIndex)
		return NewRange(int(begin), int(end))
	}
	return nil
}
