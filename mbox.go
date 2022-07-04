// Package mbox provides a flexible mbox reader and writer for four file types.
package mbox

/*
Package mbox implements a reader and writer for working with mbox files.

The package supports four types of mbox files:

* mboxo
* mboxrd
* mboxcl
* mboxcl2

Type mboxo is the original mbox format.

Type mboxrd tries to address lines starting with 'From ' in a way to avoid
conflicts by prepending such lines with '>', removing those characters when
reading the mail.

Type mboxcl tries to address lines starting with 'From ' by doing what mboxrd
does, but also adding a 'Content-Length' header to the mail that provides the
size of the mail's body.

Type mboxcl2 tries to address the lines starting with 'From ' by doing what
mboxcl does, except it doesn't add '>' characters at all.

You will need to know which type to use when reading or writing an mbox, for
best results.

NOTE: These routines do not concern themselves with file locking.  You may want
to consider that while working with mbox files on systems that might actively
write to the file.  These simply use the golang writer/reader interfaces.

*/

const (
	MBOXO   int = iota // Specifies the mboxo mail box file type.
	MBOXRD             // Specifies the mboxrd mail box file type.
	MBOXCL             // Specifies the mboxcl mail box file type.
	MBOXCL2            // Specifies the mboxcl2 mail box file type.
)
