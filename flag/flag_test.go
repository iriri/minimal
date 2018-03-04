package flag

import (
	"os"
	"strconv"
	"testing"
)

type flagSet struct {
	b    bool
	f    float64
	fStr string
	i    int
	j    int
	n    uint64
	nStr string
	s    string
}

var opt flagSet

func initFlags() {
	Bool(&opt.b, false, "bool flag", "bool", 'b')
	String(&opt.fStr, "", "float64 flag", "f64", 'f')
	Int(&opt.i, 0, "int flag", "int", 'i')
	Int(&opt.j, 0, "int flag", "", 'j')
	String(&opt.nStr, "", "uint64 flag", "", 'n')
	String(&opt.s, "", "string flag", "str", 0)
}

func verify(expected flagSet) bool {
	return expected == opt
}

func TestEverything(t *testing.T) {
	initFlags()
	args := os.Args
	os.Args = []string{
		"test",
		"-bf",
		"1234.5678",
		"--int",
		"-12345678",
		"-j",
		"0",
		"-n",
		"12345678901",
		"--str",
		"lastFlag",
		"firstArg",
	}
	expected := flagSet{
		true,
		1234.5678,
		"1234.5678",
		-12345678,
		0,
		12345678901,
		"12345678901",
		"lastFlag",
	}
	if Parse(1) != 11 {
		t.Fail()
	}
	opt.f, _ = strconv.ParseFloat(opt.fStr, 64)
	opt.n, _ = strconv.ParseUint(opt.nStr, 10, 64)
	if !verify(expected) {
		t.Fail()
	}
	os.Args = args
}

func TestDeclareInvalidLongFlags(t *testing.T) {
	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}
	var b bool
	Bool(&b, false, "", "b", 0)
	if exitCode != 1 {
		t.Fail()
	}

	var i int
	exitCode = 0
	Int(&i, 0, "", "i", 0)
	if exitCode != 1 {
		t.Fail()
	}

	var s string
	exitCode = 0
	String(&s, "", "", "s", 0)
	if exitCode != 1 {
		t.Fail()
	}
	osExit = os.Exit
}

func TestShortCircuit(t *testing.T) {
	initFlags()
	args := os.Args
	os.Args = []string{"test", "a"}
	if Parse(1) != 1 {
		t.Fail()
	}

	os.Args = []string{"test", "--a"}
	if Parse(1) != 1 {
		t.Fail()
	}
	os.Args = args
	os.Args = args
}

func TestInvalidFlag(t *testing.T) {
	initFlags()
	args := os.Args
	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}
	os.Args = []string{"test", "-h"}
	Parse(1)
	if exitCode != 1 {
		t.Fail()
	}

	initFlags()
	exitCode = 0
	os.Args = []string{"test", "--help"}
	Parse(1)
	if exitCode != 1 {
		t.Fail()
	}
	osExit = os.Exit
	os.Args = args
}

func TestFlagAfterIntOrStrFlag(t *testing.T) {
	initFlags()
	args := os.Args
	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}
	os.Args = []string{"test", "-ib"}
	Parse(1)
	if exitCode != 1 {
		t.Fail()
	}

	initFlags()
	exitCode = 0
	os.Args = []string{"test", "-fb"}
	Parse(1)
	if exitCode != 1 {
		t.Fail()
	}
	osExit = os.Exit
	os.Args = args
}

func TestInvalidIntAfterIntFlag(t *testing.T) {
	initFlags()
	args := os.Args
	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}
	os.Args = []string{"test", "-i", "1234.5689"}
	Parse(1)
	if exitCode != 1 {
		t.Fail()
	}

	initFlags()
	exitCode = 0
	os.Args = []string{"test", "--int", "1234.5689"}
	Parse(1)
	if exitCode != 1 {
		t.Fail()
	}
	osExit = os.Exit
	os.Args = args
}

func TestNothingAfterIntOrStrFlag(t *testing.T) {
	initFlags()
	args := os.Args
	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}
	os.Args = []string{"test", "-b", "--int"}
	Parse(1)
	if exitCode != 1 {
		t.Fail()
	}

	initFlags()
	exitCode = 0
	os.Args = []string{"test", "--bool", "--f64"}
	Parse(1)
	if exitCode != 1 {
		t.Fail()
	}
	osExit = os.Exit
	os.Args = args
}
