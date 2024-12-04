package onig

import (
	"fmt"
	"strconv"
)

// PythonRegexReplacer is a RegexReplacer that replaces capture groups with the corresponding captured text.
type PythonRegexReplacer struct {
	repl    []byte
	replLen int
}

// NewPythonRegexReplacer creates a new PythonRegexReplacer given a replacement pattern.
// See https://docs.python.org/3/library/re.html#re.sub for more information.
func NewPythonRegexReplacer(pattern string) RegexReplacer {
	repl := []byte(pattern)
	return &PythonRegexReplacer{
		repl:    repl,
		replLen: len(repl),
	}
}

// Replace applies the replacement pattern to the captures.
func (r *PythonRegexReplacer) Replace(captures *Captures) (string, error) {
	newReplacement := make([]byte, 0, r.replLen*3)
	inEscapeMode := false
	inGroupNameMode := false
	groupName := make([]byte, 0, r.replLen)
	for index := 0; index < r.replLen; index++ {
		ch := r.repl[index]
		if inGroupNameMode && ch == byte('<') {
		} else if inGroupNameMode && ch == byte('>') {
			inGroupNameMode = false
			groupNameStr := string(groupName)
			groupIndex, err := strconv.ParseInt(groupNameStr, 10, 32)
			var capture string
			if err == nil {
				capture = captures.At(int(groupIndex))
			} else {
				capture = captures.AtGroupName(groupNameStr)
			}
			newReplacement = append(newReplacement, []byte(capture)...)
			groupName = groupName[:0]
		} else if inGroupNameMode {
			groupName = append(groupName, ch)
		} else if inEscapeMode && ch >= byte('0') && ch <= byte('9') {
			capNumStr := string(ch)
			capNum, err := strconv.ParseInt(capNumStr, 10, 32)
			if err != nil {
				return "", fmt.Errorf("invalid capture group number: %s", capNumStr)
			}
			capBytes := []byte(captures.At(int(capNum)))
			newReplacement = append(newReplacement, capBytes...)
		} else if inEscapeMode && ch == byte('g') && (index+1) < r.replLen && r.repl[index+1] == byte('<') {
			inGroupNameMode = true
			inEscapeMode = false
			index++
		} else if inEscapeMode {
			newReplacement = append(newReplacement, '\\')
			newReplacement = append(newReplacement, ch)
		} else if ch != byte('\\') {
			newReplacement = append(newReplacement, ch)
		}
		if ch == byte('\\') || inEscapeMode {
			inEscapeMode = !inEscapeMode
		}
	}
	return string(newReplacement), nil
}
