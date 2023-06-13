package utils

import "strings"

// Unescape string in tag according to IRCv3 (e.g. \s -> space)
// /	see: https://ircv3.net/specs/extensions/message-tags.html#escaping-values)
func Unescape(s string) string {
	var foundEscape bool
	var lastEscape bool
	var sb strings.Builder
	var idxToCopy int = 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && !lastEscape {
			if !foundEscape {
				sb.Grow(len(s))
				foundEscape = true
			}
			lastEscape = true
			sb.WriteString(s[idxToCopy:i])
			idxToCopy = i + 2
		} else if lastEscape {
			if s[i] == ':' {
				sb.WriteByte(';')
			} else if s[i] == 's' {
				sb.WriteByte(' ')
			} else if s[i] == 'r' {
				sb.WriteByte('\r')
			} else if s[i] == '\\' {
				sb.WriteByte('\\')
			} else if s[i] == 'n' {
				sb.WriteByte('\n')
			} else {
				sb.WriteByte(s[i])
			}
			lastEscape = false
		}
	}
	if !foundEscape {
		return s
	}

	if idxToCopy < len(s) {
		sb.WriteString(s[idxToCopy:])
	}
	return sb.String()
}
