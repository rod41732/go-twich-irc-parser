package utils

import "strings"

// Unescape string in tag according to IRCv3 (e.g. \s -> space)
// /	see: https://ircv3.net/specs/extensions/message-tags.html#escaping-values)
func Unescape(s string) string {
	var lastEscape bool
	var sb strings.Builder
	sb.Grow(len(s))
	var idxToCopy int = 0
	for i := 0; i < len(s); i++ {
		if !lastEscape && s[i] == '\\' {
			lastEscape = true
			sb.WriteString(s[idxToCopy:i])
			idxToCopy = i + 2
		} else if lastEscape {
			lastEscape = false
			if s[i] == 's' {
				sb.WriteByte(' ')
			} else if s[i] == '\\' {
				sb.WriteByte('\\')
			} else if s[i] == 'n' {
				sb.WriteByte('\n')
			} else if s[i] == 'r' {
				sb.WriteByte('\r')
			} else if s[i] == ':' {
				sb.WriteByte(';')
			} else {
				sb.WriteByte(s[i])
			}
		}
	}

	if idxToCopy < len(s) {
		sb.WriteString(s[idxToCopy:])
	}
	return sb.String()
}

func NeedsUnescape(s string) bool {
	for _, c := range s {
		if c == '\\' {
			return true
		}
	}
	return false
}
