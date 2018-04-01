// Copyright 2018 iriri. All rights reserved. Use of this source code is
// governed by a BSD-style license which can be found in the LICENSE file.

// Package gitignore can be used to parse .gitignore-style files into lists of
// globs that can be used to test against paths or selectively walk a file
// tree. Gobwas's glob package is used for matching because it is faster than
// using regexp, which is overkill, and supports globstars (**), unlike
// filepath.Match.
package gitignore

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
	rr "github.com/iriri/minimal/replaceright"
)

type ignore struct {
	globs []glob.Glob
	root  string
}

type IgnoreList []ignore

// New creates a new ignore list.
func New() IgnoreList {
	return make([]ignore, 0, 2)
}

// From creates a new ignore list and populates the first entry with the
// contents of the specified file.
func From(path string) (IgnoreList, error) {
	ign := New()
	err := ign.Append(path)
	return ign, err
}

// FromAll creates a new ignore list and populates it with the contents of
// every file in the file tree rooted at the current working directory with
// the specified filename.
func FromAll(fname string) (IgnoreList, error) {
	ign := IgnoreList(make([]ignore, 0, 4))
	err := ign.AppendAll(fname)
	return ign, err
}

func clean(s string) string {
	i := len(s) - 1
	for ; i >= 0; i-- {
		if s[i] != ' ' || i > 0 && s[i-1] == '\\' {
			break
		}
	}
	return s[:i+1]
}

// AppendGlob appends a single glob as a new entry in the ignore list. The root
// (relevant for matching patterns that begin with "/") is assumed to be the
// current working directory.
func (ign *IgnoreList) AppendGlob(s string) error {
	g, err := glob.Compile(clean(s))
	if err != nil {
		return err
	}
	*ign = append(*ign, ignore{[]glob.Glob{g}, ""})
	return nil
}

// Append appends the globs in the specified file to the ignore list. Files are
// expected to have the same format as .gitignore files.
func (ign *IgnoreList) Append(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	root := filepath.Dir(path)
	*ign = append(*ign, ignore{make([]glob.Glob, 0, 8), root})
	globs := &(*ign)[len(*ign)-1].globs
	scn := bufio.NewScanner(bufio.NewReader(f))
	for scn.Scan() {
		s := scn.Text()
		if s == "" || s[0] == '#' {
			continue
		}
		g, err := glob.Compile(clean(s))
		if err != nil {
			log.Printf("Invalid glob in %s: %s", path, s)
			continue
		}
		*globs = append(*globs, g)
	}
	return nil
}

func (ign *IgnoreList) AppendAll(fname string) error {
	err := filepath.Walk(
		".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Base(path) == fname {
				ign.Append(path)
			}
			return nil
		})
	return err
}

func (ign IgnoreList) match(path string, info os.FileInfo) bool {
	base := filepath.Base(path)
	for _, i := range ign {
		for _, g := range i.globs {
			if g.Match(base) ||
				g.Match(path) ||
				g.Match(rr.Replace(path, i.root, "", 1)) ||
				(i.root == "." && g.Match("/"+path)) {
				return true
			}
			if info != nil && info.IsDir() && (g.Match(base+"/") ||
				g.Match(path+"/") ||
				g.Match(rr.Replace(path, i.root, "", 1)+"/") ||
				(i.root == "." && g.Match("/"+path+"/"))) {
				return true
			}
		}
	}
	return false
}

// Match returns whether any of the globs in the ignore list match the
// specified path. Uses the same matching rules as .gitignore files.
func (ign IgnoreList) Match(path string) bool {
	return ign.match(path, nil)
}

// Walk walks the file tree with the specified root and calls fn on each file
// or directory. Files and directories that match any of the globs in the
// ignore list are skipped. This behavior is inverted (i.e. non-matching files
// and directories are skipped instead) if inv is true.
func (ign IgnoreList) Walk(root string, inv bool, fn filepath.WalkFunc) error {
	return filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if ign.match(path, info) != inv {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return err
			}
			return fn(path, info, err)
		})
}
