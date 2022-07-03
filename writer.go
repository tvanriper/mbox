package mbox

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// MboxWriter describes a writer for any of the mbox file types.
// Use NewWriter to instantiate.  Set Type to specify the type.  Type is set to
// MBOXO by default.
type MboxWriter struct {
	Type  int    // Specifies the type of MboxWriter, defaulting to MBOXO.
	FS    FromFS // A filesystem for working with temporary files that handle MBOXCL/MBOXCL2 mboxes. Defaults to a FileFromFS.
	write io.Writer
}

// FromFS describes an interface for providing a reader and writer independent
// of the underlying file system.  Replace MboxWriter.FS with your own
// implementation if needed.  MboxWriter uses this interface for working with
// MBOXCL and MBOXCL2 files.
type FromFS interface {
	OpenReader(from string) (result io.ReadCloser, err error)  // Opens a ReadCloser for the item specified by the 'from' field.
	OpenWriter(from string) (result io.WriteCloser, err error) // Opens a WriteCloser for the item specified by the 'from' field.
	Remove(from string) (err error)                            // Removes any information associated with the 'from' field.
}

// FileFromFS provides a structure for working with temporary files used while
// creating MBOXCL and MBOXCL2 files.
type FileFromFS struct {
	Base  string            // The base folder in which to write/read temporary files.
	names map[string]string // Previously created names
}

// NewWriter instantiates a new mbox file writer.
// Subsequent calls to MBoxWriter.WriteMail() will write mbox-formatted output
// to the writer provided to this function.
func NewWriter(write io.Writer) (result *MboxWriter) {
	return &MboxWriter{write: write, FS: NewFileFromFS("")}
}

// NewFileFromFS creates a new FileFromFS with the provided base folder.
// If the base is length 0, it will use os.TempDir() to determine the base
// folder location.  This is the default used when calling NewWriter().
func NewFileFromFS(base string) *FileFromFS {
	if len(base) == 0 {
		base = os.TempDir()
	}
	return &FileFromFS{Base: base, names: map[string]string{}}
}

// WriteMail adds new mail to the mbox-formatted file stream provided to
// NewWriter.
// The 'from' argument may come from a call to MboxReader.NextMessage(), or
// a tool delivering the mail to the box.  The 'mail' argument contains the
// bytes composing the mail (headers and body).
func (m *MboxWriter) WriteMail(from string, mail []byte) (err error) {
	switch m.Type {
	case MBOXCL2:
		err = m.writeMBOXCL2Mail(from, mail)
	case MBOXCL:
		err = m.writeMBOXCLMail(from, mail)
	case MBOXRD:
		err = m.writeMBOXRDMail(from, mail)
	default:
		err = m.writeMBOXOMail(from, mail)
	}
	return err
}

// writeMBOXOMail writes the email using mboxo formatting.
func (m *MboxWriter) writeMBOXOMail(from string, mail []byte) (err error) {
	reader := bufio.NewReader(bytes.NewBuffer(mail))
	m.write.Write([]byte(fmt.Sprintf("From %s\n", from)))
	b, err := reader.ReadBytes('\n')
	for err == nil {
		line := string(b)
		if strings.HasPrefix(line, "From ") {
			orig := b
			b = []byte{'>'}
			b = append(b, orig...)
		}
		_, err = m.write.Write(b)
		if err != nil {
			break
		}
		b, err = reader.ReadBytes('\n')
	}
	if err == io.EOF {
		m.write.Write([]byte{'\n'})
		err = nil
	}
	return err
}

func (m *MboxWriter) writeMBOXRDMail(from string, mail []byte) (err error) {
	re, err := regexp.Compile("^>*From ")
	if err != nil {
		// This should only happen if I didn't unit test and got the regexp wrong.
		panic(err)
	}
	reader := bufio.NewReader(bytes.NewBuffer(mail))
	m.write.Write([]byte(fmt.Sprintf("From %s\n", from)))
	b, err := reader.ReadBytes('\n')
	for err == nil {
		if re.Match(b) {
			orig := b
			b = []byte{'>'}
			b = append(b, orig...)
		}
		_, err = m.write.Write(b)
		if err != nil {
			break
		}
		b, err = reader.ReadBytes('\n')
	}
	if err == io.EOF {
		m.write.Write([]byte{'\n'})
		err = nil
	}
	return err
}

func (m *MboxWriter) writeMBOXCLMail(from string, mail []byte) (err error) {
	// We need to know, in advance, the size of the message to write, after we've
	// modified it to handle 'From '.  The safest way to do this requires writing
	// to a file first to get the resulting size, then authoring the Content-Length
	// header, and copying the contents of that file to the write stream.
	// Otherwise, the size of the body might consume available RAM.
	// But, I don't like depending on a file system at all, not knowing how someone
	// might want to use this library.  So, I chose an approach that uses helper
	// functions that one may replace to acquire the temporary io.Writer/io.Reader
	// objects.  The library will provide a reasonable default that most would
	// likely use.
	tmpWriter, err := m.FS.OpenWriter(from)
	if err != nil {
		return fmt.Errorf("unable to open temporary stream: %s", err)
	}
	defer func() {
		tmpWriter.Close()
		m.FS.Remove(from)
	}()

	re, err := regexp.Compile("^>*From ")
	if err != nil {
		// This should only happen if I didn't unit test and got the regexp wrong.
		panic(err)
	}
	inHeader := true
	reader := bufio.NewReader(bytes.NewBuffer(mail))
	m.write.Write([]byte(fmt.Sprintf("From %s\n", from)))
	b, err := reader.ReadBytes('\n')
	count := int64(0)
	for err == nil {
		if re.Match(b) {
			orig := b
			b = []byte{'>'}
			b = append(b, orig...)
		}
		l := len(b)
		if inHeader {
			if l <= 2 && len(strings.TrimSpace(string(b))) == 0 {
				// We are about to write the body.
				inHeader = false
				b, err = reader.ReadBytes('\n')
				continue
			}
			_, err = m.write.Write(b)
			if err != nil {
				return err
			}
			b, err = reader.ReadBytes('\n')
			continue
		}
		// If we made it here, we aren't in the header anymore.
		count += int64(l)
		_, err = tmpWriter.Write(b)
		if err != nil {
			return err
		}
		b, err = reader.ReadBytes('\n')
	}
	if err == io.EOF {
		err = nil
	}
	if err == nil {
		var tmpReader io.ReadCloser
		tmpWriter.Close()
		m.write.Write([]byte(fmt.Sprintf("Content-Length: %d\n\n", count)))
		tmpReader, err = m.FS.OpenReader(from)
		if err != nil {
			return err
		}
		defer tmpReader.Close()
		_, err = io.Copy(m.write, tmpReader)
	}
	return err
}

func (m *MboxWriter) writeMBOXCL2Mail(from string, mail []byte) (err error) {
	// We need to know, in advance, the size of the message to write, after we've
	// modified it to handle 'From '.  The safest way to do this requires writing
	// to a file first to get the resulting size, then authoring the Content-Length
	// header, and copying the contents of that file to the write stream.
	// Otherwise, the size of the body might consume available RAM.
	// But, I don't like depending on a file system at all, not knowing how someone
	// might want to use this library.  So, I chose an approach that uses helper
	// functions that one may replace to acquire the temporary io.Writer/io.Reader
	// objects.  The library will provide a reasonable default that most would
	// likely use.
	tmpWriter, err := m.FS.OpenWriter(from)
	if err != nil {
		return fmt.Errorf("unable to open temporary stream: %s", err)
	}
	defer func() {
		tmpWriter.Close()
		m.FS.Remove(from)
	}()

	inHeader := true
	reader := bufio.NewReader(bytes.NewBuffer(mail))
	m.write.Write([]byte(fmt.Sprintf("From %s\n", from)))
	b, err := reader.ReadBytes('\n')
	count := int64(0)
	for err == nil {
		l := len(b)
		if inHeader {
			// Might be \r\n or \n.
			if l <= 2 &&
				len(strings.TrimSpace(string(b))) == 0 &&
				!bytes.ContainsAny(b, " \t") {
				// We are about to write the body.
				inHeader = false
				b, err = reader.ReadBytes('\n')
				continue
			}
			_, err = m.write.Write(b)
			if err != nil {
				return err
			}
			b, err = reader.ReadBytes('\n')
			continue
		}
		// If we made it here, we aren't in the header anymore.
		count += int64(l)
		_, err = tmpWriter.Write(b)
		if err != nil {
			return err
		}
		b, err = reader.ReadBytes('\n')
	}
	if err == io.EOF {
		err = nil
	}
	if err == nil {
		var tmpReader io.ReadCloser
		tmpWriter.Close()
		m.write.Write([]byte(fmt.Sprintf("Content-Length: %d\n\n", count)))
		tmpReader, err = m.FS.OpenReader(from)
		if err != nil {
			return err
		}
		defer tmpReader.Close()
		_, err = io.Copy(m.write, tmpReader)
	}
	return err
}

func (f *FileFromFS) getPattern(from string) (result string) {
	for _, c := range from {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			result = fmt.Sprintf("%s_", result)
		} else {
			result = fmt.Sprintf("%s%c", result, c)
		}
	}
	result = fmt.Sprintf("%s_*.txt", result)
	return result
}

// OpenReader opens an io.ReadCloser based on the 'from' provided.
func (f *FileFromFS) OpenReader(from string) (result io.ReadCloser, err error) {
	filepath, ok := f.names[from]
	if !ok {
		return nil, fmt.Errorf("did not call OpenWriter first")
	}
	return os.Open(filepath)
}

// OpenWriter opens an io.WriteCloser based on the 'from' provided.
func (f *FileFromFS) OpenWriter(from string) (result io.WriteCloser, err error) {
	file, err := os.CreateTemp(f.Base, f.getPattern(from))
	f.names[from] = file.Name()
	result = file
	return result, err
}

// RemoveFile removes the temporary file created for working with the mbox.
func (f *FileFromFS) Remove(from string) (err error) {
	filepath, ok := f.names[from]
	if !ok {
		return fmt.Errorf("did not call OpenWriter first")
	}
	delete(f.names, from)
	return os.Remove(filepath)
}
