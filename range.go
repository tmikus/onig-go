package onig

// Range is a struct that contains the regex match start and end indices.
type Range struct {
	From int // the start index of the match
	To   int // the end index of the match
}

// NewRange creates a new Range given the from and to indices.
func NewRange(from int, to int) *Range {
	return &Range{
		From: from,
		To:   to,
	}
}
