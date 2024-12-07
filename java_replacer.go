package onig

import (
	"fmt"
	"strconv"
)

// JavaReplacer is a RegexReplacer that replaces capture groups with the corresponding captured text.
type JavaReplacer struct {
	repl    []byte
	replLen int
}

// NewJavaRegexReplacer creates a new JavaReplacer given a replacement pattern.
// See https://docs.oracle.com/javase/8/docs/api/java/util/regex/Matcher.html#appendReplacement-java.lang.StringBuffer-java.lang.String- for more information.
func NewJavaRegexReplacer(pattern string) RegexReplacer {
	repl := []byte(pattern)
	return &JavaReplacer{
		repl:    repl,
		replLen: len(repl),
	}
}

// Replace applies the replacement pattern to the captures.
func (r *JavaReplacer) Replace(captures *Captures) (string, error) {
	newReplacement := make([]byte, 0, r.replLen*3)
	inEscapeMode := false
	isGroupMode := false
	for index := 0; index < r.replLen; index++ {
		ch := r.repl[index]
		if isGroupMode {
			if ch == '{' {
				isGroupClosed := false
				groupName := make([]byte, 0, r.replLen)
				for index++; index < r.replLen; index++ {
					ch = r.repl[index]
					if ch == '}' {
						isGroupClosed = true
					} else {
						groupName = append(groupName, ch)
					}
				}
				if !isGroupClosed {
					return "", fmt.Errorf("missing closing brace in replacement pattern")
				}
				capBytes := []byte(captures.AtGroupName(string(groupName)))
				newReplacement = append(newReplacement, capBytes...)
			} else if ch >= '0' && ch <= '9' {
				groupNumber := make([]byte, 0, r.replLen)
				groupNumber = append(groupNumber, ch)
				for index++; index < r.replLen; index++ {
					ch = r.repl[index]
					if ch >= '0' && ch <= '9' {
						groupNumber = append(groupNumber, ch)
					} else {
						index--
						break
					}
				}
				capNumStr := string(groupNumber)
				capNum, err := strconv.ParseInt(capNumStr, 10, 32)
				if err != nil {
					return "", fmt.Errorf("invalid capture group number: %s", capNumStr)
				}
				capBytes := []byte(captures.At(int(capNum)))
				newReplacement = append(newReplacement, capBytes...)
			} else {
				return "", fmt.Errorf("unexpected token in replacement pattern: %s", string(ch))
			}
			isGroupMode = false
		} else if inEscapeMode {
			newReplacement = append(newReplacement, ch)
			inEscapeMode = false
		} else if ch == '\\' {
			inEscapeMode = true
		} else if ch == '$' {
			isGroupMode = true
		} else {
			newReplacement = append(newReplacement, ch)
		}
	}
	if isGroupMode {
		return "", fmt.Errorf("incomplete replacement pattern")
	}
	if inEscapeMode {
		return "", fmt.Errorf("incomplete escape pattern")
	}
	return string(newReplacement), nil
}
