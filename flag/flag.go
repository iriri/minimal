package flag

import (
	"fmt"
	"os"
	//"strings"
)

type ftype int

const (
	notFlag ftype = iota
	shortFlag
	longFlag
)

type flag interface {
	set(interface{})
	printUsageAndDelete()
}

type boolFlag struct {
	val     *bool
	baseVal bool
	usage   string
	long    string
	short   rune
}

type stringFlag struct {
	val     *string
	baseVal string
	usage   string
	long    string
	short   rune
}

func (f *boolFlag) set(b interface{}) {
	*f.val = b.(bool)
}

func (f *boolFlag) printUsageAndDelete() {
	if f.short != 0 && f.long != "" {
		fmt.Fprintf(os.Stderr, "    -%c --%s\t%s [default: %t]\n",
			f.short, f.long, f.usage, f.baseVal)
		delete(shortFlags, f.short)
		delete(longFlags, f.long)
	} else if f.short != 0 {
		fmt.Fprintf(os.Stderr, "    -%c\t\t%s [default: %t]\n",
			f.short, f.usage, f.baseVal)
		delete(shortFlags, f.short)
	} else if f.long != "" {
		fmt.Fprintf(os.Stderr, "    --%s\t%s [default: %t]\n",
			f.long, f.usage, f.baseVal)
		delete(longFlags, f.long)
	}
}

func (f *stringFlag) set(s interface{}) {
	*f.val = s.(string)
}

func (f *stringFlag) printUsageAndDelete() {
	if f.short != 0 && f.long != "" {
		fmt.Fprintf(os.Stderr, "    -%c --%s\t%s [default: %s]\n",
			f.short, f.long, f.usage, f.baseVal)
		delete(shortFlags, f.short)
		delete(longFlags, f.long)
	} else if f.short != 0 {
		fmt.Fprintf(os.Stderr, "    -%c\t\t%s [default: %s]\n",
			f.short, f.usage, f.baseVal)
		delete(shortFlags, f.short)
	} else if f.long != "" {
		fmt.Fprintf(os.Stderr, "    --%s\t%s [default: %s]\n",
			f.long, f.usage, f.baseVal)
		delete(longFlags, f.long)
	}
}

var shortFlags = make(map[rune]flag)
var longFlags = make(map[string]flag)

func Bool(val *bool, baseVal bool, usage string, long string, short rune) {
	f := boolFlag{val, baseVal, usage, long, short}
	*val = baseVal

	if short != 0 {
		shortFlags[short] = &f
	}
	if long != "" {
		longFlags[long] = &f
	}
}

func String(val *string, baseVal string, usage string, long string,
	short rune) {
	f := stringFlag{val, baseVal, usage, long, short}
	*val = baseVal

	if short != 0 {
		shortFlags[short] = &f
	}
	if long != "" {
		longFlags[long] = &f
	}
}

func isFlag(s string) ftype {
	if s[0] == '-' {
		if s[1] == '-' {
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
			fmt.Fprintf(os.Stderr, "Invalid flag: -%c\n", r)
			printUsageAndExit()
		}
		switch f := f.(type) {
		case *boolFlag:
			f.set(!*f.val)
		case *stringFlag:
			if j != len(os.Args[i])-1 ||
				isFlag(os.Args[i+1]) != notFlag {
				fmt.Fprintf(os.Stderr,
					"-%c must precede string\n", r)
				printUsageAndExit()
			}
			f.set(os.Args[i+1])
			return 1
		default:
			fmt.Fprintf(os.Stderr, "Invalid flag: -%c\n", r)
			printUsageAndExit()
		}
	}
	return 0
}

func parseLongFlag(i int) int {
	f, ok := longFlags[os.Args[i][2:]]
	if !ok {
		fmt.Fprintf(os.Stderr, "Invalid flag: %s\n", os.Args[i])
		printUsageAndExit()
	}
	switch f := f.(type) {
	case *boolFlag:
		f.set(!*f.val)
	case *stringFlag:
		if isFlag(os.Args[i+1]) != notFlag {
			fmt.Fprintf(os.Stderr, "%s must precede string\n",
				os.Args[i])
			printUsageAndExit()
		}
		f.set(os.Args[i+1])
		return 1
	default:
		fmt.Fprintf(os.Stderr, "Invalid flag: %s\n", os.Args[i])
		printUsageAndExit()
	}
	return 0
}

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

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	for _, f := range shortFlags {
		f.printUsageAndDelete()
	}
	for _, f := range longFlags {
		f.printUsageAndDelete()
	}
	os.Exit(0)
}
