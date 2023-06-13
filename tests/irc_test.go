package gotwichircparser_test

import (
	"go-twitch-irc-parser/irc"
	"os"
	"strings"
	"testing"
)

func TestSimpleIRC(t *testing.T) {
	msg := "@foo=bar :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed := irc.NewIRCMessage(msg)

	expRawTags := "foo=bar"
	if string(parsed.RawTags) != expRawTags {
		t.Errorf("RawTags: expected [%s], got: [%s]", expRawTags, parsed.RawTags)
	}

	expPrefix := "user!user@user.tmi.twitch.tv"
	if string(parsed.Prefix) != expPrefix {
		t.Errorf("Prefix: expected [%s], got: [%s]", expPrefix, parsed.Prefix)
	}

	expCommand := "PRIVMSG"
	if string(parsed.Command) != expCommand {
		t.Errorf("Command: expected [%s], got: [%s]", expCommand, parsed.Command)
	}

	expParams := "#pajlada :this is a test"
	if string(parsed.Params) != expParams {
		t.Errorf("Params: expected [%s], got: [%s]", expParams, parsed.Params)
	}

}

func BenchmarkParsingSingleMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := irc.NewIRCMessage("@badge-info=subscriber/22;badges=subscriber/3012;color=#FFFF00;display-name=FELYP8;emote-only=1;emotes=521050:0-6,8-14,16-22,24-30,32-38,40-46,48-54,56-62,64-70,72-78,80-86,88-94,96-102,104-110,148-154,156-162,164-170,172-178,180-186,188-194,196-202,204-210,212-218,220-226,228-234,236-242,244-250,252-258,260-266/302827730:112-119/302827734:121-128/302827735:130-137/302827737:139-146;first-msg=0;flags=;id=1844235a-c24e-4e18-937b-805d6601aebe;mod=0;returning-chatter=0;room-id=22484632;subscriber=1;tmi-sent-ts=1685664001040;turbo=0;user-id=162760707;user-type= :felyp8!felyp8@felyp8.tmi.twitch.tv PRIVMSG #forsen :forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE1 forsenE2 forsenE3 forsenE4 forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE")
		_ = m
		// m.ParseTags()
	}
}

// 1.80
// move to func 1.40 - 1.64
// forget to reset timer ->  1.37 - 1.64
// parse at same time, N 1.27 - 1.31
func BenchmarkParsing1000Messages(b *testing.B) {
	f, err := os.ReadFile("../data.txt")
	if err != nil {
		b.Fatalf("Read file failed %s", err)
	}

	lines := strings.Split(string(f), "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) != 1000 {
		b.Fatalf("Not 1000 lines")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			m := irc.NewIRCMessage(line)
			_ = m
			// m.ParseTags()
		}
	}
}
