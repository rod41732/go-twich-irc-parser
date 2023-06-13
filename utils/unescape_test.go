package utils_test

import (
	"testing"

	"github.com/rod41732/go-twitch-irc-parser/utils"
)

func assertEqual(t *testing.T, exp string, act string) {
	if exp != act {
		t.Errorf("Expected %q, got %q", exp, act)
	}
}

func TestNoEscape(t *testing.T) {
	assertEqual(t, "forsen", utils.Unescape("forsen"))
	// all escapes
	assertEqual(t, "for\nsen", utils.Unescape("for\\nsen"))
	assertEqual(t, "for\rsen", utils.Unescape("for\\rsen"))
	assertEqual(t, "for sen", utils.Unescape("for\\ssen"))
	assertEqual(t, "for;sen", utils.Unescape("for\\:sen"))
	assertEqual(t, "for\\sen", utils.Unescape("for\\\\sen"))
	// invalid escape is same
	assertEqual(t, "xdxd", utils.Unescape("xd\\xd"))
	// trailing \ is dropped
	assertEqual(t, "forsen", utils.Unescape("forsen\\"))
	// multiple escapes
	assertEqual(t, "this is a test\\", utils.Unescape("this\\sis\\sa\\stest\\\\"))
	// many backslashs
	assertEqual(t, "this\\x", utils.Unescape("this\\\\x"))
	assertEqual(t, "this\\x", utils.Unescape("this\\\\\\x"))
}
