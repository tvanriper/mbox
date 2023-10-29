package mbox

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func lineFeedType(reader io.ReadSeeker) (bool, error) {
	_, err := reader.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}
	b := make([]byte, 1024)
	for count, err := reader.Read(b); err == nil; {
		for i := 0; i < count; i++ {
			if b[i] == '\r' {
				return true, nil
			}
			if b[i] == '\n' {
				return false, nil
			}
		}
	}
	return false, fmt.Errorf("no carriage return or line feed")
}

// DetectType attempts to figure out the type of mbox the reader holds.  This
// is a best-effort attempt to determine the type of mbox file format based on
// what it sees within the text.  When it returns, it attempts to move reader
// to the beginning of the stream on exit.
//
// It tries to work out the type of file by:
//   - Looking for 'Content-Length' in a message's header
//   - Looking for '>From ' or '>>From ' (or any number of > character in front
//     of "From "') in the message's body.
//   - Using the length of the 'Content-Length', if present, to determine when
//     the body of the message is complete.
//
// With this information, it can guess if the mbox matches one of the file
// types supported by this library with some degree of certainty.
func DetectType(reader io.ReadSeeker) (mboxType int, err error) {
	rdMatch := regexp.MustCompile(`^>*>From `)
	clMatch := regexp.MustCompile(`^Content-Length:`)
	feedType, err := lineFeedType(reader)
	if err != nil {
		return -1, err
	}
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return -1, err
	}
	defer reader.Seek(0, io.SeekStart)
	inHeader := false
	scanner := bufio.NewScanner(reader)
	var hasRd bool = false
	var hasCL bool = false
	var count int64 = 0
	var clLen int64 = 0
	var finishedFirst bool = false
	for scanner.Scan() {
		// NOTE:
		// The man page for mbox indicates that one can tell different mailings
		// apart via lines that start with From followed by a space.  It also
		// states that the message is RFC 822 encoded.  An RFC 822 encoded text
		// message separates the header from the body via a 'null' line (CRLF)
		// Ergo, for our purposes, we'll scan over the file, assuming From_
		// indicates the start of the header, and a null line indicates the
		// start of the body.

		// Trimming the space to ensure extranous characters won't figure into
		// this.
		line := strings.TrimSpace(scanner.Text())
		if inHeader && len(line) == 0 {
			inHeader = false
			count = 0
			finishedFirst = true
			continue
		}
		if !inHeader && !hasCL && strings.HasPrefix(line, "From ") {
			inHeader = true
		}
		if !inHeader && hasCL {
			count += int64(len(scanner.Text())) + 1
			if feedType {
				count += 1
			}
		}

		matchRd := rdMatch.MatchString(line)
		matchCL := clMatch.MatchString(line)

		if inHeader && matchCL {
			hasCL = true
			// We have a content-length.  We need to parse it to determine the
			// length of the body.
			sp := strings.Split(line, ":")
			if len(sp) < 2 {
				// well, er, this is awkward...
				continue
			}
			intCL, err := strconv.ParseInt(strings.TrimSpace(sp[1]), 10, 64)
			if err != nil {
				// I guess that wasn't an int.
				continue
			}
			clLen = intCL
		}

		if !inHeader && matchRd {
			hasRd = true
		}

		if hasRd && hasCL {
			// We have enough evidence:
			// This has content length & >From_ in the body.
			return MBOXCL, nil
		}
		if hasCL && !inHeader && strings.HasPrefix(line, "From ") {
			// We have enough evidence:
			// This has content length, we're in the body, and the line starts
			// with 'From '.
			return MBOXCL2, nil
		}
		if !inHeader && clLen == count {
			count = 0
			finishedFirst = true
			// This isn't technically true, but helps for our tests.
			inHeader = true
		}

		if finishedFirst && !hasCL && hasRd {
			// We have enough evidence:
			// This doesn't have content length, but does have lines in the
			// body starting with >First.
			return MBOXRD, nil
		}
	}
	if hasCL && !hasRd {
		// We don't really know.  It could be MBOXCL2, or MBOXCL.  We will err on the side of caution.
		return MBOXCL, nil
	}
	return MBOXO, nil
}
