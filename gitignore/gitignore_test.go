package gitignore

import (
	"os"
	"strings"
	"testing"
)

func TestEverything(t *testing.T) {
	ign, err := New()
	if err != nil {
		panic(err)
	}
	err = ign.appendAll("testgitignore", ".")
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
		"../gitignore/testfs",
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
		"../../minimal",
		func(path string, info os.FileInfo, err error) error {
			if strings.Contains(path, "ignoredfile") {
				t.Fail()
			}
			return nil
		})
	ign.Walk(
		"../gitignore",
		func(path string, info os.FileInfo, err error) error {
			if strings.Contains(path, "ignoredfile") {
				t.Fail()
			}
			return nil
		})
	ign.Walk(
		"../gitignore/testfs",
		func(path string, info os.FileInfo, err error) error {
			if strings.Contains(path, "ignoredfile") {
				t.Fail()
			}
			return nil
		})
	err = ign.Walk(
		"../../../not_a_real_directory",
		func(path string, info os.FileInfo, err error) error {
			if strings.Contains(path, "ignoredfile") {
				t.Fail()
			}
			return nil
		})
	if err == nil {
		t.Fail()
	}
}
