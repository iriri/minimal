package flag

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

type ftype int

const (
	notFlag ftype = iota
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
	*v = boolVal(b.(bool))
}

func (v *intVal) set(i interface{}) {
	*v = intVal(i.(int))
}

func (v *stringVal) set(s interface{}) {
	*v = stringVal(s.(string))
}

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

var shortFlags = make(map[rune]flag)
var longFlags = make(map[string]flag)

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
		switch t := f.val.(type) {
		case *boolVal:
			t.set(!*t)
		case *intVal:
			if j != len(os.Args[i])-2 || len(os.Args[i:]) < 2 {
				fmt.Fprintf(os.Stderr, "%d %d", j, len(os.Args[i:]))
				fmt.Fprintf(os.Stderr,
					"-%c must precede integer\n", r)
				printUsageAndExit()
			}
			n, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"-%c must precede integer\n", r)
				printUsageAndExit()
			}
			t.set(n)
			return 1
		case *stringVal:
			if j != len(os.Args[i])-1 || len(os.Args[i:]) < 2 ||
				isFlag(os.Args[i+1]) != notFlag {
				fmt.Fprintf(os.Stderr,
					"-%c must precede string\n", r)
				printUsageAndExit()
			}
			t.set(os.Args[i+1])
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
	switch t := f.val.(type) {
	case *boolVal:
		t.set(!*t)
	case *intVal:
		if len(os.Args[i:]) < 2 {
			fmt.Fprintf(os.Stderr,
				"%c must precede integer\n", os.Args[i])
			printUsageAndExit()
		}
		n, err := strconv.Atoi(os.Args[i+1])
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"%c must precede integer\n", os.Args[i])
			printUsageAndExit()
		}
		t.set(n)
		return 1
	case *stringVal:
		if len(os.Args[i:]) < 2 || isFlag(os.Args[i+1]) != notFlag {
			fmt.Fprintf(os.Stderr, "%s must precede string\n",
				os.Args[i])
			printUsageAndExit()
		}
		t.set(os.Args[i+1])
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
