# minimal
Worse versions of packages that already exist.

## minimal/color
Just the consts.

## minimal/flag [![GoDoc](https://godoc.org/github.com/iriri/minimal/flag?status.svg)](https://godoc.org/github.com/iriri/minimal/flag)
Package flag provides a very minimal command line flag parser.

## minimal/gitignore
Package gitignore is used to parse .gitignore-style files into globs that can be used to test against a certain string or selectively walk a file tree. Gobwas's glob package is used for matching because it is faster than using regexp, which is overkill, and supports globstars (**), unlike filepath.Match. Not thoroughly tested so you probably shouldn't use this.
