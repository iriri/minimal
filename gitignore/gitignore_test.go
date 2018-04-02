package gitignore

import (
	"os"
	"testing"
)

func TestEverything(t *testing.T) {
	ign, err := FromAll("testgitignore")
	if err != nil {
		panic(err)
	}
	err = ign.AppendGlob("aaa")
	if err != nil {
		panic(err)
	}
	expected := [...]string{
		"testfs",
		"testfs/eee",
		"testfs/eee/ggg",
		"testfs/test.ou",
		"testfs/testdir",
	}
	actual := make([]string, 0, 3)
	ign.Walk(
		"testfs",
		false,
		func(path string, info os.FileInfo, err error) error {
			actual = append(actual, path)
			return nil
		})
	if len(actual) != len(expected) {
		t.Fail()
		return
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Fail()
		}
	}

	ign, err = From("testgitignore")
	if err != nil {
		panic(err)
	}
	if !ign.Match("Makefile") {
		t.Fail()
	}
	if !ign.Match("../../test/test/../../test/Makefile") {
		t.Fail()
	}

	ign, err = FromGit()
	ign.Walk(
		".",
		false,
		func(path string, info os.FileInfo, err error) error {
			if path == "ignoredfile" {
				t.Fail()
			}
			return nil
		})
}
