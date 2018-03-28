// Copyright 2018 iriri. All rights reserved. Use of this source code is
// governed by a BSD-style license which can be found in the LICENSE file.

// Package flag provides a very minimal command line flag parser. Both short
// flags (-s) and long flags (--long) are supported. Short flags can be chained
// (-xvzf) and "--" is treated as the end of flags marker. Boolean and integer
// values have first-class support; strings values are intended to serve as a
// catch-all for anything else.
package flag

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

type flagVal interface {
	set(interface{})
}

type boolVal bool
type int64Val int64
type stringVal string

type flag struct {
	val   flagVal
	usage string
	long  string
	short rune
}

type flagType uint

const (
	notFlag flagType = iota
	shortFlag
	longFlag
	endFlag
)

var osExit = os.Exit

func (v *boolVal) set(b interface{}) {
	*v = boolVal(b.(bool))
}

func (v *int64Val) set(i interface{}) {
	*v = int64Val(i.(int64))
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

// Bool defines a bool flag with the specified short and/or long variants, base
// value, and usage text. The argument val points to where the value is stored.
func Bool(val *bool, short rune, long string, base bool, usage string) {
	if len(long) == 1 {
		fmt.Fprintf(
			os.Stderr,
			"single character flags cannot be declared as long\n")
		osExit(1)
		return
	}
	f := flag{(*boolVal)(val), usage, long, short}
	*val = base
	if short != 0 {
		shortFlags[short] = f
	}
	if long != "" {
		longFlags[long] = f
	}
}

// Int64 defines an int64 flag with the specified short and/or long variants,
// base value, and usage text. The argument val points to where the value is
// stored.
func Int64(val *int64, short rune, long string, base int64, usage string) {
	if len(long) == 1 {
		fmt.Fprintf(
			os.Stderr,
			"single character flags cannot be declared as long\n")
		osExit(1)
		return
	}
	f := flag{(*int64Val)(val), usage, long, short}
	*val = base
	if short != 0 {
		shortFlags[short] = f
	}
	if long != "" {
		longFlags[long] = f
	}
}

// String defines a string flag with the specified short and/or long variants,
// base value, and usage text. The argument val points to where the value is
// stored.
func String(val *string, short rune, long string, base string, usage string) {
	if len(long) == 1 {
		fmt.Fprintf(
			os.Stderr,
			"single character flags cannot be declared as long\n")
		osExit(1)
		return
	}
	f := flag{(*stringVal)(val), usage, long, short}
	*val = base
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
			if len(s) == 2 {
				return endFlag
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
			return 1
		}
		switch t := f.val.(type) {
		case *boolVal:
			t.set(true)
		case *int64Val:
			if j != len(os.Args[i])-2 || len(os.Args[i:]) < 2 {
				fmt.Fprintf(os.Stderr,
					"-%c must precede integer\n", r)
				PrintUsageAndExit()
				return 1
			}
			n, err := strconv.ParseInt(os.Args[i+1], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"-%c must precede integer\n", r)
				PrintUsageAndExit()
				return 1
			}
			t.set(n)
			return 1
		case *stringVal:
			if j != len(os.Args[i])-2 || len(os.Args[i:]) < 2 {
				fmt.Fprintf(os.Stderr,
					"-%c must precede string\n", r)
				PrintUsageAndExit()
				return 1
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
		return 1
	}
	switch t := f.val.(type) {
	case *boolVal:
		t.set(true)
	case *int64Val:
		if len(os.Args[i:]) < 2 {
			fmt.Fprintf(os.Stderr,
				"%s must precede integer\n", os.Args[i])
			PrintUsageAndExit()
			return 1
		}
		n, err := strconv.ParseInt(os.Args[i+1], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"%s must precede integer\n", os.Args[i])
			PrintUsageAndExit()
			return 1
		}
		t.set(n)
		return 1
	case *stringVal:
		if len(os.Args[i:]) < 2 {
			fmt.Fprintf(os.Stderr, "%s must precede string\n",
				os.Args[i])
			PrintUsageAndExit()
			return 1
		}
		t.set(os.Args[i+1])
		return 1
	}
	return 0
}

// Parse parses the command line flags from os.Args[firstFlag:] and returns the
// index of the first non-flag command line argument. If present, "--" is
// treated as the end of flags marker and the index of the next argument is
// returned. If Parse encounters a flag that has not been defined
// PrintUsageAndExit will be called.
func Parse(firstFlag int) int {
	i := firstFlag
	for ; i <= len(os.Args[firstFlag:]); i++ {
		switch isFlag(os.Args[i]) {
		case shortFlag:
			i += parseShortFlag(i)
		case longFlag:
			i += parseLongFlag(i)
		case endFlag:
			return i + 1
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
	osExit(1)
}
