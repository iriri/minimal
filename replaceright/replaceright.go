// Copyright 2018 iriri. All rights reserved. Use of this source code is
// governed by a BSD-style license which can be found in the LICENSE file.

// Package replaceright provides similar functionality to some of the standard
// library string replacement routines except they start from the right. Uses
// Boyer-Moore on longer strings.
package replaceright

import "strings"

type Replacer struct {
	old         string
	new         string
	badCharSkip [256]int
	goodPfxSkip []int
}

// Similar to strings.Replace except it starts from the right, immediately
// returns if s or old is the empty string, and n cannot be < 0 as there's no
// reason to replace from the right if you're going to replace everything.
func Replace(s, old, new string, n int) string {
	if s == "" || old == "" || old == new || n <= 0 {
		return s
	}
	if len(s) > 32 {
		rep := NewReplacer(old, new)
		return rep.Replace(s, n)
	}

	end := len(old) - 1
	var buf []byte
	if len(new) > len(old) {
		buf = make([]byte, len(s)+int(n)*(len(new)-len(old)))
	} else {
		buf = make([]byte, len(s))
	}
	wr := len(buf)
	for i := len(s) - 1; i >= 0; i-- {
		if i < end {
			goto next
		}
		if s[i] == old[end] {
			for j, k := i-1, end-1; k >= 0; j, k = j-1, k-1 {
				if s[j] != old[k] {
					goto next
				}
			}
			for j := len(new) - 1; j >= 0; j-- {
				wr--
				buf[wr] = new[j]
			}
			i -= end
			continue
		}
	next:
		wr--
		buf[wr] = s[i]
	}
	return string(buf[wr:])
}

func maxSharedPfxLen(a, b string) int {
	for i := 0; i < len(b); i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return len(b)
}

// Creates a new Replacer that can be used to repeatedly replace a single
// old-new pair. This implementation of Boyer-Moore is largely copied from the
// Go standard library.
func NewReplacer(old, new string) *Replacer {
	rep := &Replacer{old: old, new: new}
	if old == new {
		return rep
	}

	for i := range rep.badCharSkip {
		rep.badCharSkip[i] = len(old)
	}
	for i := 1; i < len(old); i++ {
		rep.badCharSkip[old[i]] = i
	}
	var lastSfx int
	rep.goodPfxSkip = make([]int, len(old))
	for i := 0; i < len(old); i++ {
		if strings.HasSuffix(old, old[:i]) {
			lastSfx = i - 1
		}
		rep.goodPfxSkip[i] = (len(old) - 1 - lastSfx) + i
	}
	for i := len(old) - 1; i > 0; i-- {
		lenPfx := maxSharedPfxLen(old, old[i:len(old)-1])
		if old[i+lenPfx] != old[lenPfx] {
			rep.goodPfxSkip[lenPfx] = lenPfx + i
		}
	}
	return rep
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (rep *Replacer) find(text string) int {
	i := len(text) - len(rep.old)
	for i >= 0 {
		j := 0
		for j < len(rep.old) && text[i] == rep.old[j] {
			i++
			j++
		}
		if j == len(rep.old) {
			return i - 1
		}
		i -= max(rep.badCharSkip[text[i]], rep.goodPfxSkip[j])
	}
	return -1
}

// Performs up to n replacements starting from the right.
func (rep *Replacer) Replace(s string, n int) string {
	if s == "" || rep.old == "" || rep.old == rep.new || n <= 0 {
		return s
	}

	var buf []byte
	if len(rep.new) > len(rep.old) {
		buf = make([]byte, len(s)+int(n)*(len(rep.new)-len(rep.old)))
	} else {
		buf = make([]byte, len(s))
	}
	i := len(s)
	wr := len(buf)
	for ; n > 0; n-- {
		j := rep.find(s[:i])
		if j < 0 {
			for i > 0 {
				i--
				wr--
				buf[wr] = s[i]
			}
			break
		}
		for i-1 > j {
			i--
			wr--
			buf[wr] = s[i]
		}
		for k := len(rep.new) - 1; k >= 0; k-- {
			wr--
			buf[wr] = rep.new[k]
		}
		i -= len(rep.old)
	}
	for i > 0 {
		i--
		wr--
		buf[wr] = s[i]
	}
	return string(buf[wr:])
}
