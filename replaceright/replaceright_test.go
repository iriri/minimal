package replaceright

import (
	"testing"
)

func Test(t *testing.T) {
	if Replace("aaaaa", "a", "b", 5) != "bbbbb" {
		t.Fail()
	}
	if Replace("aaaaa", "a", "b", 6) != "bbbbb" {
		t.Fail()
	}
	if Replace("bbbbb", "a", "b", 5) != "bbbbb" {
		t.Fail()
	}
	if Replace("bbbbb", "a", "a", 5) != "bbbbb" {
		t.Fail()
	}
	if Replace("abcdefg", "fg", "12", 1) != "abcde12" {
		t.Fail()
	}
	if Replace("abcdefgabcdefg", "g", "123", 2) != "abcdef123abcdef123" {
		t.Fail()
	}
	if Replace("abcdefg", "cde", "12345", 1) != "ab12345fg" {
		t.Fail()
	}
	if Replace("abcdefg", "ab", "1", 1) != "1cdefg" {
		t.Fail()
	}
	if Replace("abcdefg", "xfg", "123", 1) != "abcdefg" {
		t.Fail()
	}
	if Replace("abcdefg", "efg", "", 3) != "abcd" {
		t.Fail()
	}
	if Replace(
		"daddy give me cummies (◔◡◔✿)",
		"(◔◡◔✿)",
		"uwu",
		69) != "daddy give me cummies uwu" {
		t.Fail()
	}

	rep := NewReplacer("a", "a")
	if rep.Replace("aaaaa", 3) != "aaaaa" {
		t.Fail()
	}

	rep = NewReplacer("a", "bbb")
	if rep.Replace("bbbbb", 3) != "bbbbb" {
		t.Fail()
	}
	if rep.Replace("aaaaa", 1) != "aaaabbb" {
		t.Fail()
	}
	if rep.Replace("aaaaa", 2) != "aaabbbbbb" {
		t.Fail()
	}
	if rep.Replace("aaaaa", 5) != "bbbbbbbbbbbbbbb" {
		t.Fail()
	}
	if rep.Replace("aaaaa", 6) != "bbbbbbbbbbbbbbb" {
		t.Fail()
	}

	rep = NewReplacer("efg", "")
	if rep.Replace("abcdefg", 5) != "abcd" {
		t.Fail()
	}
	if rep.Replace("abefgcefg", 5) != "abc" {
		t.Fail()
	}
}
