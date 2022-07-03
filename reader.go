package mbox

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// MboxReader provides a reader for mbox files.
type MboxReader struct {
	Type int // Specifies the type of MboxReader, defaulting to MBOXO.
	from string
	read *bufio.Reader
}

// lineReader is a function you provide to MboxReader.nextMessageGeneric that either ignores or processes
// a line of text from the mail box.  If it returns false, the line is written as-is to the message.
// Otherwise, the line is not passed to the message unless the function does so on its own somehow.
// It may return an error if it encounters something unrecoverable, suggesting a malformed mbox file.
type lineReader func(string) (bool, error)

// NewReader creates a new MboxReader.
// You may wish to set the Type after instantiation if your mbox is anything other than MBOXO.
func NewReader(read io.Reader) *MboxReader {
	return &MboxReader{
		read: bufio.NewReader(read),
	}
}

// NextMessage writes the next message into the writer.
// This returns the 'From ' string that separates the mbox email, and an err for an error.
// This returns an io.EOF error when the last message is read.
func (m *MboxReader) NextMessage(write io.Writer) (from string, err error) {
	switch m.Type {
	case MBOXRD:
		return m.nextMBOXRDMessage(write)
	case MBOXCL:
		return m.nextMBOXCLMessage(write)
	case MBOXCL2:
		return m.nextMBOXCL2Message(write)
	default:
		return m.nextMBOXOMessage(write)
	}
}

// nextMessageGeneric is a common parser that handles most of the needs for parsing an mbox.
func (m *MboxReader) nextMessageGeneric(write io.Writer, fn lineReader) (from string, err error) {
	inMessage := false
	if len(m.from) > 0 {
		// We have already read the From line.
		from = m.from
		inMessage = true
	}
	for {
		b, err := m.read.ReadBytes('\n')
		if err != nil {
			return from, err
		}
		line := string(b[:len(b)-1])
		if strings.HasPrefix(line, "From ") {
			if inMessage {
				// We've finished the message... this starts a new one.
				m.from = line
				// Since we've finished, break out of scanning.
				break
			} else {
				// We're starting a new message.
				from = line
				continue
			}
		}

		if len(line) == 0 && !inMessage {
			inMessage = true
		}

		ok, err := fn(line)
		if err != nil {
			return from, err
		}
		if ok {
			continue
		}

		write.Write([]byte(fmt.Sprintf("%s\n", line)))
	}
	return from, err
}

// nextMBOXOMessage parses mboxo files.
func (m *MboxReader) nextMBOXOMessage(write io.Writer) (from string, err error) {
	return m.nextMessageGeneric(write, func(line string) (bool, error) { return false, nil })
}

// nextMBOXRDMessage parses mboxrd files.
func (m *MboxReader) nextMBOXRDMessage(write io.Writer) (from string, err error) {
	// Find all lines starting with any number of '>' characters followed by 'From '.
	// We need to change them to have one less '>' character before writing it.
	re, err := regexp.Compile("^>+From .*")
	if err != nil {
		// This should only happen if I didn't test the regular expression properly.
		panic(err)
	}
	return m.nextMessageGeneric(write, func(line string) (ok bool, err error) {
		if re.Match([]byte(line)) {
			// We need to remove the first character before writing.
			write.Write([]byte(fmt.Sprintf("%s\n", line[1:])))
			return true, nil
		}
		return false, nil
	})
}

// nextMBOXCLMessage parses mboxcl files.
func (m *MboxReader) nextMBOXCLMessage(write io.Writer) (from string, err error) {
	// A bit more complicated.
	// We want the 'Content-Length:' header in the email, which tells us the size of the
	// body of the email.  So we have to detect when we are no longer reading headers
	// (the first blank line), and write the whole body into the writer.
	// However... both header and body still have >From_ bits in it that need to be
	// corrected.  This is easily the most complicated MBOX type for processing purposes.
	re, err := regexp.Compile("^>+From .*")
	if err != nil {
		// This should only happen if I didn't test the regular expression properly.
		panic(err)
	}
	size := int64(0)
	return m.nextMessageGeneric(write, func(line string) (bool, error) {
		if re.Match([]byte(line)) {
			// We need to remove the first character before writing.
			write.Write([]byte(fmt.Sprintf("%s\n", line[1:])))
			return true, nil
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length: ") {
			size, err = strconv.ParseInt(line[16:], 10, 64)
			if err != nil {
				return false, fmt.Errorf("failed to parse Content-Length: %s", err)
			}
			return false, nil
		}
		if len(line) == 0 {
			// We are now in the body.
			// But, we still need to look for our regular expression.
			// We shouldn't read the whole message at once, but we should look for lines.
			// We must count bytes.
			write.Write([]byte{'\n'})
			if size == 0 {
				// We have the body already... it's empty.
				return true, nil
			}
			for {
				b, err := m.read.ReadBytes('\n')
				if err != nil {
					// Probably no more data.
					return true, err
				}
				// Decrementing the size we allocated earlier.
				size -= int64(len(b))
				if re.Match(b) {
					write.Write(b[1:])
				} else {
					write.Write(b)
				}
				if size <= 0 {
					break
				}
			}
			return true, nil
		}
		return false, nil
	})
}

// nextMBOXCL2Message parses mboxcl2 files.
func (m *MboxReader) nextMBOXCL2Message(write io.Writer) (from string, err error) {
	// A bit more complicated.
	// We want the 'Content-Length:' header in the email, which tells us the size of the
	// body of the email.  So we have to detect when we are no longer reading headers
	// (the first blank line), and write the whole body into the writer.
	size := int64(0)
	return m.nextMessageGeneric(write, func(line string) (bool, error) {
		if strings.HasPrefix(strings.ToLower(line), "content-length: ") {
			size, err = strconv.ParseInt(line[16:], 10, 64)
			if err != nil {
				return false, fmt.Errorf("failed to parse Content-Length: %s", err)
			}
			return false, nil
		}
		if len(line) == 0 {
			// We are now in the body.
			write.Write([]byte{'\n'})
			io.CopyN(write, m.read, size)
			return true, nil
		}
		return false, nil
	})
}
