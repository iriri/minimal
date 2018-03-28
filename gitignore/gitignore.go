// Copyright 2018 iriri. All rights reserved. Use of this source code is
// governed by a BSD-style license which can be found in the LICENSE file.

// Package gitignore can be used to parse .gitignore-style files into globs
// that can be used to test against a certain string or selectively walk a file
// tree. Gobwas's glob package is used for matching because it is faster than
// using regexp, which is overkill, and supports globstars (**), unlike
// filepath.Match.
package gitignore

import (
	"bufio"
	"github.com/gobwas/glob"
	"log"
	"os"
	"path/filepath"
)

type Ignore []*glob.Glob

// New creates a new list of globs to ignore (or not ignore).
func New() Ignore {
	return make([]*glob.Glob, 0, 16)
}

// From creates a new list of globs and populates it with the contents of the
// specified file.
func From(fname string) (Ignore, error) {
	ign := New()
	err := ign.Append(fname)
	return ign, err
}

// AppendGlob appends a single glob to the ignore list.
func (ign *Ignore) AppendGlob(s string) error {
	g, err := glob.Compile(s)
	if err != nil {
		return err
	}
	*ign = append(*ign, &g)
	return nil
}

// Append appends the globs in the specified file to the ignore list. Files are
// expected to have the same format as .gitignore files: every non-empty line
// is expected to be a valid glob unless it starts with "#:.
func (ign *Ignore) Append(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	scn := bufio.NewScanner(bufio.NewReader(f))
	for scn.Scan() {
		s := scn.Text()
		if s == "" || s[0] == '#' {
			continue
		}
		err = ign.AppendGlob(s)
		if err != nil {
			log.Printf("Invalid glob in %s: %s", fname, s)
			continue
		}
	}
	return nil
}

// Match returns whether any of the globs in the ignore list match the
// specified path.
func (ign Ignore) Match(path string) bool {
	for _, g := range ign {
		if (*g).Match(path) {
			return true
		}
	}
	return false
}

// Walk walks the file tree with the specified root and calls walkFn on each
// file or directory. Files and directories that match any of the globs in the
// ignore list are skipped. This behavior can be inverted (i.e. non-matching
// files and directories are skipped instead) by setting the argument inv to
// true.
func (ign Ignore) Walk(root string, inv bool, walkFn filepath.WalkFunc) error {
	return filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			b := filepath.Base(path)
			if info.IsDir() &&
				((ign.Match(b) || ign.Match(b+"/")) != inv) {
				return filepath.SkipDir
			}
			if ign.Match(b) != inv {
				return err
			}
			return walkFn(path, info, err)
		})
}
