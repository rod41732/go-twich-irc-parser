package irc

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

func parseTags(raw string) []TagEntry {
	idx := 0
	tags := make([]TagEntry, 0, 31)

	var keyEnd int
	var keyStart int
	var valueStart int

wholeParsing:
	for idx < len(raw) {
		idx++
		// keyHash
		keyStart = idx
		for ; idx < len(raw); idx++ {
			if raw[idx] == ';' {
				// TODO
				tags = append(tags, TagEntry{raw[keyStart:keyEnd], ""})

				continue wholeParsing
			} else if raw[idx] == '=' {
				idx++
				break
			}
		}
		keyEnd = idx

		// value
		valueStart = idx
		for ; idx < len(raw); idx++ {
			if raw[idx] == ';' {
				tags = append(tags, TagEntry{raw[keyStart:keyEnd], raw[valueStart:idx]})
				continue wholeParsing
			}
		}
		// value EOF
		tags = append(tags, TagEntry{raw[keyStart:keyEnd], raw[valueStart:idx]})
	}
	return tags
}

func (m *IRCMessage) ParseTags() {
	m.Tag = parseTags(m.RawTags)
}

func NewIRCMessage(raw string) IRCMessage {
	parsed := IRCMessage{}

	rawBytes := raw
	idx := 0
	// tag: @foo=bar;foo2=bar2...
	if rawBytes[idx] == '@' {
		start := idx + 1
		for idx = start; idx < len(rawBytes); idx++ {
			if rawBytes[idx] == ' ' {
				break
			}
		}
		parsed.RawTags = (rawBytes[start:idx])
		idx++
	}

	// something like :userlogin!userlogin@userlogin.tmi.twitch.tv
	if rawBytes[idx] == ':' {
		start := idx + 1
		for idx = start; idx < len(rawBytes); idx++ {
			if rawBytes[idx] == ' ' {
				break
			}
		}
		parsed.Prefix = (rawBytes[start:idx])
		idx++
	}

	// command e.g. PRIVMSG
	start := idx
	for idx = start; idx < len(rawBytes); idx++ {
		if rawBytes[idx] == ' ' {
			break
		}
	}
	parsed.Command = (rawBytes[start:idx])
	idx++

	// params: everything after command e.g. "#pajlada :this is my message"
	parsed.Params = (rawBytes[idx:])

	return parsed
}
