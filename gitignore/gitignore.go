package gitignore

import (
	"bufio"
	"fmt"
	"github.com/gobwas/glob"
	"os"
	"path/filepath"
)

type Ignore []*glob.Glob

func New() Ignore {
	return make([]*glob.Glob, 0, 16)
}

func From(fname string) (Ignore, error) {
	ign := Ignore(make([]*glob.Glob, 0, 16))
	err := ign.Append(fname)
	return ign, err
}

func (i Ignore) Match(path string) bool {
	for _, g := range i {
		if (*g).Match(path) {
			return true
		}
	}
	return false
}

func (i *Ignore) Append(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	scn := bufio.NewScanner(bufio.NewReader(f))
	for scn.Scan() {
		g, err := glob.Compile(scn.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid glob in %s: %s", fname,
				scn.Text())
			continue
		}
		*i = append(*i, &g)
	}
	return nil
}

func (i *Ignore) AppendStr(s string) error {
	g, err := glob.Compile(s)
	if err != nil {
		return err
	}
	*i = append(*i, &g)
	return nil
}

func Walk(root string, ign Ignore, b bool, fn filepath.WalkFunc) error {
	globFn := func(path string, info os.FileInfo, err error) error {
		if ign.Match(path) == b {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return err
		}
		return fn(path, info, err)
	}

	return filepath.Walk(root, globFn)
}
