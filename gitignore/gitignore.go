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
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

type IgnoreList struct {
	globs []glob.Glob
	cwd   []string
}

func toSplit(path string) []string {
	return strings.Split(filepath.ToSlash(path), "/")
}

func fromSplit(path []string) string {
	return filepath.FromSlash(strings.Join(path, "/"))
}

// New creates a new ignore list.
func New() (*IgnoreList, error) {
	cwd, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}
	return &IgnoreList{
		make([]glob.Glob, 0, 16),
		toSplit(cwd),
	}, nil
}

// From creates a new ignore list and populates the first entry with the
// contents of the specified file.
func From(path string) (*IgnoreList, error) {
	ign, err := New()
	if err != nil {
		return nil, err
	}
	err = ign.Append(path)
	if err != nil {
		return nil, err
	}
	return ign, nil
}

// FromGit finds the root directory of the current git repository and creates a
// new ignore list with the contents of all .gitignore files in that git
// repository.
func FromGit() (*IgnoreList, error) {
	ign, err := New()
	if err != nil {
		return nil, err
	}
	err = ign.AppendGit()
	if err != nil {
		return nil, err
	}
	return ign, nil
}

func clean(s string) string {
	i := len(s) - 1
	for ; i >= 0; i-- {
		if s[i] != ' ' || i > 0 && s[i-1] == '\\' {
			return s[:i+1]
		}
	}
	return ""
}

func toRelpath(s string, root, cwd []string) string {
	if s != "" {
		if s[0] != '/' {
			return s
		}
		if root == nil || cwd == nil {
			return s[1:]
		}
		root = append(root, toSplit(s[1:])...)
	}

	i := 0
	min := len(cwd)
	if len(root) < min {
		min = len(root)
	}
	for ; i < min; i++ {
		if root[i] != cwd[i] {
			break
		}
	}
	ss := make([]string, (len(cwd)-i)+(len(root)-i))
	j := 0
	for ; j < len(cwd)-i; j++ {
		ss[j] = ".."
	}
	for k := 0; j < len(ss); j, k = j+1, k+1 {
		ss[j] = root[i+k]
	}
	return fromSplit(ss)
}

func (ign *IgnoreList) appendGlob(s string, root []string) error {
	g, err := glob.Compile(toRelpath(clean(s), root, ign.cwd))
	if err != nil {
		return err
	}
	ign.globs = append(ign.globs, g)
	return nil
}

// AppendGlob appends a single glob as a new entry in the ignore list. The root
// (relevant for matching patterns that begin with "/") is assumed to be the
// current working directory.
func (ign *IgnoreList) AppendGlob(s string) error {
	return ign.appendGlob(s, nil)
}

// Append appends the globs in the specified file to the ignore list. Files are
// expected to have the same format as .gitignore files.
func (ign *IgnoreList) Append(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return err
	}
	var root []string
	if dir != fromSplit(ign.cwd) {
		root = toSplit(dir)
	}
	scn := bufio.NewScanner(bufio.NewReader(f))
	for scn.Scan() {
		s := scn.Text()
		if s == "" || s[0] == '#' {
			continue
		}
		err = ign.appendGlob(s, root)
		if err != nil {
			return err
		}
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func findGitRoot(cwd []string) (string, error) {
	p := fromSplit(cwd)
	for !fileExists(p + "/.git") {
		if len(cwd) == 1 {
			return "", errors.New("not in a git repository")
		}
		cwd = cwd[:len(cwd)-1]
		p = fromSplit(cwd)
	}
	return p, nil
}

func (ign *IgnoreList) appendAll(fname, root string) error {
	err := filepath.Walk(
		root,
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

// AppendGit finds the root directory of the current git repository and appends
// the contents of all .gitignore files in that git repository to the ignore
// list.
func (ign *IgnoreList) AppendGit() error {
	root, err := findGitRoot(ign.cwd)
	if err != nil {
		return err
	}
	err = ign.appendAll(".gitignore", root)
	if err != nil {
		return err
	}
	return nil
}

func (ign *IgnoreList) match(path string, info os.FileInfo) bool {
	ss := make([]string, 0, 4)
	base := filepath.Base(path)
	ss = append(ss, path)
	if base != path {
		ss = append(ss, base)
	} else {
		ss = append(ss, "./"+path)
	}
	if info != nil && info.IsDir() {
		ss = append(ss, path+"/")
		if base != path {
			ss = append(ss, base+"/")
		} else {
			ss = append(ss, "./"+path+"/")
		}
	}

	for _, g := range ign.globs {
		for _, s := range ss {
			if g.Match(s) {
				return true
			}
		}
	}
	return false
}

// Match returns whether any of the globs in the ignore list match the
// specified path. Uses the same matching rules as .gitignore files.
func (ign *IgnoreList) Match(path string) bool {
	return ign.match(path, nil)
}

// Walk walks the file tree with the specified root and calls fn on each file
// or directory. Files and directories that match any of the globs in the
// ignore list are skipped.
func (ign *IgnoreList) Walk(root string, fn filepath.WalkFunc) error {
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	return filepath.Walk(
		toRelpath("", toSplit(abs), ign.cwd),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if ign.match(path, info) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return err
			}
			return fn(path, info, err)
		})
}
