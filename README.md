# mbox

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/github.com/tvanriper/mbox)
[![Coverage Status](https://coveralls.io/repos/github/tvanriper/mbox/badge.svg?branch=main)](https://coveralls.io/github/tvanriper/mbox?branch=main)

Supporting four different mbox file formats in pure golang.

Package mbox implements a reader and writer for working with mbox files.

The package supports four types of mbox files:

- mboxo
- mboxrd
- mboxcl
- mboxcl2

Use `mboxo` for the original mbox format.

Use `mboxrd` to handle lines starting with 'From ' in a way to avoid
conflicts by prepending such lines with '>', removing those characters when
reading the mail.

Use `mboxcl` to address lines starting with 'From ' by doing what mboxrd
does, but also adding a 'Content-Length' header to the mail that provides the
size of the mail's body.

Use `mboxcl2` to address the lines starting with 'From ' by doing what
mboxcl does, except it doesn't add '>' characters at all.

You may need to know which type to use when reading or writing an mbox, for
best results.

NOTE: These routines do not concern themselves with file locking. You may want
to consider that while working with mbox files on systems that might actively
write to the file, such the mbox for a Linux account on a local system. This
library simply use the golang writer/reader interfaces.

Hopefully, one may use a new structure to work with file locking once golang
exposes a standardized, tested file locking API.  Currently, one must work with
golang's internal API, or write their own code, for proper file locking.

## Installation

```bash
go get github.com/tvanriper/mbox
```
