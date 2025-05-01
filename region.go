package onig

/*
#include "regex.h"
*/
import "C"

// Region represents a set of capture groups found in a search or match.
type Region struct {
	groupIndicesMap map[string][]int
	positions       []*Range
}

// A newRegion creates a new empty Region.
func newRegion(regex *Regex, raw C.region) *Region {
	positions := make([]*Range, raw.groupCount)
	for i := 0; i < int(raw.groupCount); i++ {
		begin := offsetInt(raw.groupStartIndices, i)
		end := offsetInt(raw.groupEndIndices, i)
		if begin == -1 || end == -1 {
			positions[i] = nil
		} else {
			positions[i] = NewRange(int(begin), int(end))
		}
	}
	return &Region{
		groupIndicesMap: regex.groupIndicesMap,
		positions:       positions,
	}
}

// Len returns the number of registers in the region.
func (r *Region) Len() int {
	return len(r.positions)
}

// Pos returns the start and end positions of the Nth capture group.
// Returns nil if the capture group did not match anything or if index is not a valid capture group.
// The positions returned are always byte indices with respect to the original string matched.
func (r *Region) Pos(index int) *Range {
	if index >= r.Len() {
		return nil
	}
	return r.positions[index]
}

// PosByGroupName returns the start and end positions of the named capture group.
// Returns nil if the capture group did not match anything or if groupName is not a valid capture group.
// The positions returned are always byte indices with respect to the original string matched.
func (r *Region) PosByGroupName(groupName string) *Range {
	indices, ok := r.groupIndicesMap[groupName]
	if !ok {
		return nil
	}
	for _, index := range indices {
		if index >= r.Len() {
			return nil
		}
		position := r.positions[index]
		if position != nil {
			return position
		}
	}
	return nil
}
