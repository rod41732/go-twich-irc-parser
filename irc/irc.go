package irc

import "strings"

type TagEntry struct {
	Key   string
	Value string
}

type IRCMessage struct {
	// Tag     map[string]string
	Tag     []TagEntry
	RawTags string
	Prefix  string
	Command string
	Params  string
}

func hash(s string, m int) int {
	var c byte
	for i := 0; i < len(s); i++ {
		c = s[i]
	}
	return (int(c)%m + m) % m
}

func realUnescapeValue(s string) string {
	nextEscaped := false
	unescaper := func(r rune) rune {
		if nextEscaped {
			nextEscaped = false
			switch r {
			case ':':
				return ';'
			case 'r':
				return '\r'
			case 'n':
				return '\n'
			case 's':
				return ' '
			default:
				return r
			}
		} else if r == '\\' {
			nextEscaped = true
			return -1
		} else {
			return r
		}
	}
	return strings.Map(unescaper, s)
}

func unescapeValue(s string) string {
	hasEscape := false
	for _, c := range s {
		if c == '\\' {
			hasEscape = true
		}
	}
	if !hasEscape {
		return s
	} else {
		return realUnescapeValue(s)
	}
}

func parseTags(raw string) []TagEntry {
	idx := 0
	tags := make([]TagEntry, 0, 31)

	var keyStart int
	var valueStart int
	// var lc byte

wholeParsing:
	for ; idx < len(raw); idx++ {
		// keyHash
		keyStart = idx
		tags = append(tags, TagEntry{raw[keyStart:idx], ""})
		for ; idx < len(raw); idx++ {
			// lc = raw[idx]
			if raw[idx] == ';' {
				tags[len(tags)-1].Key = raw[keyStart:idx]
				continue wholeParsing
			} else if raw[idx] == '=' {
				tags[len(tags)-1].Key = raw[keyStart:idx]
				idx++
				break
			}
		}
		// traliing word?
		if raw[idx-1] != '=' {
			tags[len(tags)-1].Key = raw[keyStart:idx]
		}

		// value
		valueStart = idx
		for ; idx < len(raw); idx++ {
			// raw[idx] = raw[idx]
			if raw[idx] == ';' {
				tags[len(tags)-1].Value = unescapeValue(raw[valueStart:idx])
				continue wholeParsing
			}
		}
		// value EOF
		tags[len(tags)-1].Value = unescapeValue(raw[valueStart:idx])
	}
	return tags
}

func (m *IRCMessage) ParseTags() {
	m.Tag = parseTags(m.RawTags)
}

func NewIRCMessage(raw string) IRCMessage {
	parsed := IRCMessage{}

	n := len(raw)
	idx := 0
	// tag: @foo=bar;foo2=bar2...
	if raw[idx] == '@' {
		start := idx + 1
		for idx = start; idx < n; idx++ {
			if raw[idx] == ' ' {
				break
			}
		}
		parsed.RawTags = (raw[start:idx])
		parsed.Tag = parseTags(raw[start:idx])
		idx++
	}

	// something like :userlogin!userlogin@userlogin.tmi.twitch.tv
	if raw[idx] == ':' {
		start := idx + 1
		for idx = start; idx < n; idx++ {
			if raw[idx] == ' ' {
				break
			}
		}
		parsed.Prefix = (raw[start:idx])
		idx++
	}

	// command e.g. PRIVMSG
	start := idx
	for idx = start; idx < n; idx++ {
		if raw[idx] == ' ' {
			break
		}
	}
	parsed.Command = (raw[start:idx])
	idx++

	// params: everything after command e.g. "#pajlada :this is my message"
	parsed.Params = (raw[idx:])

	return parsed
}
