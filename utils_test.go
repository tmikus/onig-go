package onig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertReplaceError(t *testing.T, r *Regex, text string, replacement string) {
	t.Helper()
	_, err := r.Replace(text, replacement)
	assert.Error(t, err)
}
