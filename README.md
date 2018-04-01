# minimal
Worse versions of packages that already exist.

## minimal/color
Just the consts.

## minimal/flag [![GoDoc](https://godoc.org/github.com/iriri/minimal/flag?status.svg)](https://godoc.org/github.com/iriri/minimal/flag)
Package flag provides a very minimal command line flag parser. Both short flags
(-s) and long flags (--long) are supported. Short flags can be chained (-xvzf)
and "--" is treated as the end of flags marker. Boolean and integer values have
first-class support; strings values are intended to serve as a catch-all for
anything else.

## minimal/gitignore [![GoDoc](https://godoc.org/github.com/iriri/minimal/gitignore?status.svg)](https://godoc.org/github.com/iriri/minimal/gitignore)
Package gitignore can be used to parse .gitignore-style files into globs that
can be used to test against a certain string or selectively walk a file tree.
Gobwas's glob package is used for matching because it is faster than using
regexp, which is overkill, and supports globstars (**), unlike filepath.Match.

## minimal/replaceright [![GoDoc](https://godoc.org/github.com/iriri/minimal/replaceright?status.svg)](https://godoc.org/github.com/iriri/minimal/replaceright)
Package replaceright provides similar functionality to some of the standard
library string replacement routines except they start from the right. Uses
Boyer-Moore on longer strings.
