package onig

import "iter"

// Captures represents a group of captured strings for a single match.
//
// The 0th capture always corresponds to the entire match.
// Each subsequent index corresponds to the next capture group in the regex.
// Positions returned from a capture group are always byte indices.
type Captures struct {
	Offset uint
	Region *Region
	Text   string
}

// All returns all the capture groups in order of appearance in the regular expression.
func (c *Captures) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, text := range c.AllWithIndex() {
			if !yield(text) {
				return
			}
		}
	}
}

// AllWithIndex returns all the capture groups in order of appearance in the regular expression.
func (c *Captures) AllWithIndex() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i := 0; i < c.Len(); i++ {
			if !yield(i, c.At(i)) {
				return
			}
		}
	}
}

// AllPos returns all the capture group positions in order of appearance in the regular expression.
// Positions are byte indices in terms of the original string matched.
func (c *Captures) AllPos() iter.Seq[*Range] {
	return func(yield func(*Range) bool) {
		for _, pos := range c.AllPosWithIndex() {
			if !yield(pos) {
				return
			}
		}
	}
}

// AllPosWithIndex returns all the capture group positions in order of appearance in the regular expression.
// Positions are byte indices in terms of the original string matched.
func (c *Captures) AllPosWithIndex() iter.Seq2[int, *Range] {
	return func(yield func(int, *Range) bool) {
		for i := 0; i < c.Len(); i++ {
			if !yield(i, c.Pos(i)) {
				return
			}
		}
	}
}

// At returns the matched string for the capture group i.
// If it isn’t a valid capture group or didn’t match anything, then an empty string is returned.
func (c *Captures) At(i int) string {
	r := c.Pos(i)
	if r == nil {
		return ""
	}
	return c.Text[r.From:r.To]
}

// IsEmpty returns true if and only if there are no captured groups.
func (c *Captures) IsEmpty() bool {
	return c.Len() == 0
}

// Len returns the number of captured groups.
func (c *Captures) Len() int {
	return c.Region.Len()
}

// Pos returns the start and end positions of the Nth capture group.
// Returns nil if it is not a valid capture group or if the capture group did not match anything.
// The positions returned are always byte indices with respect to the original string matched.
func (c *Captures) Pos(i int) *Range {
	return c.Region.Pos(i)
}
