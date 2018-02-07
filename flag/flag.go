// Copyright 2018 iriri. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package flag provides a very minimal command line flag parser. Both short
flags (-s) and long flags (--long) are supported and short flags can be chained
(-xvzf). Boolean and integer values can be set directly by this package;
strings values can serve as a catch-all for anything else.
*/
package flag

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

type flagType int

const (
	notFlag flagType = iota
	shortFlag
	longFlag
)

type flagVal interface {
	set(interface{})
}

type boolVal bool
type intVal int
type stringVal string

type flag struct {
	val   flagVal
	usage string
	long  string
	short rune
}

func (v *boolVal) set(b interface{}) {
	*v = b.(boolVal)
}

func (v *intVal) set(i interface{}) {
	*v = intVal(i.(int))
}

func (v *stringVal) set(s interface{}) {
	*v = stringVal(s.(string))
}

var shortFlags = make(map[rune]flag)
var longFlags = make(map[string]flag)

func (f flag) printUsageAndDelete() {
	if f.short != 0 && f.long != "" {
		fmt.Fprintf(os.Stderr, "    -%c --%s\t%s\n",
			f.short, f.long, f.usage)
		delete(shortFlags, f.short)
		delete(longFlags, f.long)
	} else if f.short != 0 {
		fmt.Fprintf(os.Stderr, "    -%c\t\t%s\n",
			f.short, f.usage)
		delete(shortFlags, f.short)
	} else if f.long != "" {
		fmt.Fprintf(os.Stderr, "    --%s\t%s\n",
			f.long, f.usage)
		delete(longFlags, f.long)
	}
}

// Bool defines a bool flag with the specified base value, usage text, and
// long and/or short flags. The argument val points to where the value is
// stored.
func Bool(val *bool, baseVal bool, usage string, long string, short rune) {
	f := flag{(*boolVal)(val), usage, long, short}
	*val = baseVal
	if short != 0 {
		shortFlags[short] = f
	}
	if long != "" {
		longFlags[long] = f
	}
}

// Int defines an int flag with the specified base value, usage text, and
// long and/or short flags. The argument val points to where the value is
// stored.
func Int(val *int, baseVal int, usage string, long string, short rune) {
	f := flag{(*intVal)(val), usage, long, short}
	*val = baseVal
	if short != 0 {
		shortFlags[short] = f
	}
	if long != "" {
		longFlags[long] = f
	}
}

// String defines a string flag with the specified base value, usage text, and
// long and/or short flags. The argument val points to where the value is
// stored.
func String(val *string, baseVal string, usage string, long string,
	short rune) {
	f := flag{(*stringVal)(val), usage, long, short}
	*val = baseVal
	if short != 0 {
		shortFlags[short] = f
	}
	if long != "" {
		longFlags[long] = f
	}
}

func isFlag(s string) flagType {
	if len(s) < 2 {
		return notFlag
	}
	if s[0] == '-' {
		if s[1] == '-' {
			if len(s) < 4 {
				return notFlag
			}
			return longFlag
		}
		return shortFlag
	}
	return notFlag
}

func parseShortFlag(i int) int {
	for j, r := range os.Args[i][1:] {
		f, ok := shortFlags[r]
		if !ok {
			fmt.Fprintf(os.Stderr, "invalid flag: -%c\n", r)
			PrintUsageAndExit()
		}
		switch t := f.val.(type) {
		case *boolVal:
			t.set(!*t)
		case *intVal:
			if j != len(os.Args[i])-2 || len(os.Args[i:]) < 2 {
				fmt.Fprintf(os.Stderr,
					"-%c must precede integer\n", r)
				PrintUsageAndExit()
			}
			n, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"-%c must precede integer\n", r)
				PrintUsageAndExit()
			}
			t.set(n)
			return 1
		case *stringVal:
			if j != len(os.Args[i])-2 || len(os.Args[i:]) < 2 ||
				isFlag(os.Args[i+1]) != notFlag {
				fmt.Fprintf(os.Stderr,
					"-%c must precede string\n", r)
				PrintUsageAndExit()
			}
			t.set(os.Args[i+1])
			return 1
		}
	}
	return 0
}

func parseLongFlag(i int) int {
	f, ok := longFlags[os.Args[i][2:]]
	if !ok {
		fmt.Fprintf(os.Stderr, "invalid flag: %s\n", os.Args[i])
		PrintUsageAndExit()
	}
	switch t := f.val.(type) {
	case *boolVal:
		t.set(!*t)
	case *intVal:
		if len(os.Args[i:]) < 2 {
			fmt.Fprintf(os.Stderr,
				"%c must precede integer\n", os.Args[i])
			PrintUsageAndExit()
		}
		n, err := strconv.Atoi(os.Args[i+1])
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"%c must precede integer\n", os.Args[i])
			PrintUsageAndExit()
		}
		t.set(n)
		return 1
	case *stringVal:
		if len(os.Args[i:]) < 2 || isFlag(os.Args[i+1]) != notFlag {
			fmt.Fprintf(os.Stderr, "%s must precede string\n",
				os.Args[i])
			PrintUsageAndExit()
		}
		t.set(os.Args[i+1])
		return 1
	}
	return 0
}

// Parse parses the command line flags from os.Args[firstFlag:] and returns
// the index of the first non-flag command line argument. If Parse encounters
// a flag that has not been defined PrintUsageAndExit will be called.
func Parse(firstFlag int) int {
	i := firstFlag
	for ; i <= len(os.Args[firstFlag:]); i++ {
		switch isFlag(os.Args[i]) {
		case shortFlag:
			i += parseShortFlag(i)
		case longFlag:
			i += parseLongFlag(i)
		default:
			return i
		}
	}
	return i
}

// PrintUsageAndExit prints usage text based on the defined flags and exits.
func PrintUsageAndExit() {
	fmt.Fprintf(os.Stderr, "usage of %s:\n", os.Args[0])
	shortKeys := make([]int, len(shortFlags))
	for r := range shortFlags {
		shortKeys = append(shortKeys, int(r))
	}
	sort.Ints(shortKeys)
	for _, r := range shortKeys {
		shortFlags[rune(r)].printUsageAndDelete()
	}
	longKeys := make([]string, len(longFlags))
	for s := range longFlags {
		longKeys = append(longKeys, s)
	}
	sort.Strings(longKeys)
	for _, s := range longKeys {
		longFlags[s].printUsageAndDelete()
	}
	os.Exit(0)
}
