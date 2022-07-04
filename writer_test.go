package mbox

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

var email1 string = `From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
From X: WackyHeader

From all of us, to all of you, be happy!
`

var email2 string = `From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
`

var email3 string = `From: nobody@nowhere.man
To: mrmxpdstk@lazytown.com
Subject: Mysterious Jenkins
`

var email4 string = `From: corrupter@argh.net
To: mrmxpdstk@lazytown.com
Subject: Ah, ha ha ha ha!

I remember when you wrote:

>From then on, I was a genius.

Do you remember?
`

var from1 = "bubbles@bubbletown.com"
var from2 = "mrspam@corporate.corp.com"
var from3 = "nobody@nowhere.man"
var from4 = "corrupter@argh.net"

func CompareBodies(expected string, received string, t *testing.T) {
	if expected != received {
		check := diff.Diff(expected, received)
		t.Errorf("expected data doesn't match received data:\n%s\n", check)
	}
}

type BrokenFS struct {
	BreakWriter bool
	BreakReader bool
	BreakRemove bool
}

type Discarder struct{}

func (d *Discarder) Read([]byte) (n int, err error) {
	return 0, nil
}

func (d *Discarder) Write([]byte) (n int, err error) {
	return 0, nil
}

func (d *Discarder) Close() (err error) { return nil }

var Discard *Discarder = &Discarder{}

func (b *BrokenFS) OpenWriter(from string) (result io.WriteCloser, err error) {
	if b.BreakWriter {
		err = fmt.Errorf("never gonna give you up")
	}
	return Discard, err
}

func (b *BrokenFS) OpenReader(from string) (result io.ReadCloser, err error) {
	if b.BreakReader {
		err = fmt.Errorf("never gonna let you down")
	}
	return Discard, err
}

func (b *BrokenFS) Remove(from string) (err error) {
	if b.BreakRemove {
		err = fmt.Errorf("never gonna run around and desert you")
	}
	return err
}

func TestWriteMBOXO(t *testing.T) {
	var err error
	result := bytes.NewBuffer([]byte{})
	mbox := NewWriter(result)
	err = mbox.WriteMail(from1, bytes.NewBuffer([]byte(email1)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from2, bytes.NewBuffer([]byte(email2)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from3, bytes.NewBuffer([]byte(email3)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from4, bytes.NewBuffer([]byte(email4)))
	if err != nil {
		t.Error(err)
	}
	expectedFile := `From bubbles@bubbletown.com
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
>From X: WackyHeader

>From all of us, to all of you, be happy!

From mrspam@corporate.corp.com
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
>From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!

From nobody@nowhere.man
From: nobody@nowhere.man
To: mrmxpdstk@lazytown.com
Subject: Mysterious Jenkins

From corrupter@argh.net
From: corrupter@argh.net
To: mrmxpdstk@lazytown.com
Subject: Ah, ha ha ha ha!

I remember when you wrote:

>From then on, I was a genius.

Do you remember?

`
	found := result.Bytes()
	CompareBodies(expectedFile, string(found), t)
}

func TestWriteMBOXRD(t *testing.T) {
	var err error
	result := bytes.NewBuffer([]byte{})
	mbox := NewWriter(result)
	mbox.Type = MBOXRD
	err = mbox.WriteMail(from1, bytes.NewBuffer([]byte(email1)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from2, bytes.NewBuffer([]byte(email2)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from3, bytes.NewBuffer([]byte(email3)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from4, bytes.NewBuffer([]byte(email4)))
	if err != nil {
		t.Error(err)
	}
	expectedFile := `From bubbles@bubbletown.com
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
>From X: WackyHeader

>From all of us, to all of you, be happy!

From mrspam@corporate.corp.com
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
>From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!

From nobody@nowhere.man
From: nobody@nowhere.man
To: mrmxpdstk@lazytown.com
Subject: Mysterious Jenkins

From corrupter@argh.net
From: corrupter@argh.net
To: mrmxpdstk@lazytown.com
Subject: Ah, ha ha ha ha!

I remember when you wrote:

>>From then on, I was a genius.

Do you remember?

`
	found := result.Bytes()
	CompareBodies(expectedFile, string(found), t)
}

func TestWriteMBOXCL(t *testing.T) {
	var err error
	result := bytes.NewBuffer([]byte{})
	mbox := NewWriter(result)
	mbox.Type = MBOXCL
	err = mbox.WriteMail(from1, bytes.NewBuffer([]byte(email1)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from2, bytes.NewBuffer([]byte(email2)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from3, bytes.NewBuffer([]byte(email3)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from4, bytes.NewBuffer([]byte(email4)))
	if err != nil {
		t.Error(err)
	}
	expectedFile := `From bubbles@bubbletown.com
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
>From X: WackyHeader
Content-Length: 42

>From all of us, to all of you, be happy!
From mrspam@corporate.corp.com
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!
Content-Length: 130

You won't believe these prices!
>From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
From nobody@nowhere.man
From: nobody@nowhere.man
To: mrmxpdstk@lazytown.com
Subject: Mysterious Jenkins
Content-Length: 0

From corrupter@argh.net
From: corrupter@argh.net
To: mrmxpdstk@lazytown.com
Subject: Ah, ha ha ha ha!
Content-Length: 78

I remember when you wrote:

>>From then on, I was a genius.

Do you remember?
`
	found := result.Bytes()
	CompareBodies(expectedFile, string(found), t)
}

func TestWriteMBOXCL2(t *testing.T) {
	var err error
	result := bytes.NewBuffer([]byte{})
	mbox := NewWriter(result)
	mbox.Type = MBOXCL2
	err = mbox.WriteMail(from1, bytes.NewBuffer([]byte(email1)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from2, bytes.NewBuffer([]byte(email2)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from3, bytes.NewBuffer([]byte(email3)))
	if err != nil {
		t.Error(err)
	}
	err = mbox.WriteMail(from4, bytes.NewBuffer([]byte(email4)))
	if err != nil {
		t.Error(err)
	}
	expectedFile := `From bubbles@bubbletown.com
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
From X: WackyHeader
Content-Length: 41

From all of us, to all of you, be happy!
From mrspam@corporate.corp.com
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!
Content-Length: 129

You won't believe these prices!
From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
From nobody@nowhere.man
From: nobody@nowhere.man
To: mrmxpdstk@lazytown.com
Subject: Mysterious Jenkins
Content-Length: 0

From corrupter@argh.net
From: corrupter@argh.net
To: mrmxpdstk@lazytown.com
Subject: Ah, ha ha ha ha!
Content-Length: 77

I remember when you wrote:

>From then on, I was a genius.

Do you remember?
`
	found := result.Bytes()
	CompareBodies(expectedFile, string(found), t)
}

func TestBrokenWriteFS(t *testing.T) {
	result := bytes.NewBuffer([]byte{})
	mbox := NewWriter(result)
	mbox.FS = &BrokenFS{BreakWriter: true}
	mbox.Type = MBOXCL2
	err := mbox.WriteMail(from1, bytes.NewBuffer([]byte(email1)))
	if err == nil {
		t.Errorf("expected error, but it succeeded.")
	}
	mbox.Type = MBOXCL
	err = mbox.WriteMail(from1, bytes.NewBuffer([]byte(email1)))
	if err == nil {
		t.Errorf("expected error, but it succeeded.")
	}
}
func TestBrokenReadFS(t *testing.T) {
	result := bytes.NewBuffer([]byte{})
	mbox := NewWriter(result)
	mbox.FS = &BrokenFS{BreakReader: true}
	mbox.Type = MBOXCL2
	err := mbox.WriteMail(from1, bytes.NewBuffer([]byte(email1)))
	if err == nil {
		t.Errorf("expected error, but it succeeded.")
	}
	mbox.FS = &BrokenFS{BreakReader: true}
	mbox.Type = MBOXCL
	err = mbox.WriteMail(from1, bytes.NewBuffer([]byte(email1)))
	if err == nil {
		t.Errorf("expected error, but it succeeded.")
	}
}

func TestFileFromFSOrder(t *testing.T) {
	fs := NewFileFromFS("")
	reader, err := fs.OpenReader("moo")
	if err == nil {
		t.Errorf("expected error, but succeeded")
		reader.Close()
		fs.Remove("moo")
	}
	err = fs.Remove("arf")
	if err == nil {
		t.Errorf("expected error, but succeeded")
	}
}
