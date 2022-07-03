package mbox

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"testing"
)

var mboxo string = `From someone
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation

>From all of us, to all of you, be happy!
From someone-else
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
>From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
`

var mboxcl string = `From someone
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
Content-Length: 42

>From all of us, to all of you, be happy!
From someone-else
>From mug: weird header
Content-Length: 130
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
>From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
From nobody
From: nobody@nowhere.man
To: mrmxpdstk@lazytown.com
Subject: Mysterious Jenkins
Content-Length: 0

`

var badmboxcl string = `From someone
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
Content-Length: ts

>From all of us, to all of you, be happy!
From someone-else
Content-Length: 130
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
>From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
`
var badlenmboxcl string = `From someone
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
Content-Length: 42

>From all of us, to all of you, be happy!
From someone-else
Content-Length: 130
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
line of jets this side of the Gobi Desert!
`

type MsgTest struct {
	Headers map[string]string
	Body    string
}

func CheckMessage(expected MsgTest, msg *mail.Message) (err error) {
	for header, value := range expected.Headers {
		found := msg.Header.Get(header)
		if found != value {
			return fmt.Errorf("expected %s but found %s", value, found)
		}
	}

	b, err := io.ReadAll(msg.Body)
	if err != nil {
		return err
	}
	if string(b) != expected.Body {
		return fmt.Errorf("body of email does not match expectations\nExpected:\n%s\n\nGot (len:%d):\n%s\n", expected.Body, len(b), string(b))
	}
	return err
}

func TestReadMBOXO(t *testing.T) {
	box := NewReader(bytes.NewBuffer([]byte(mboxo)))
	msgStream := bytes.NewBuffer([]byte{})
	from, err := box.NextMessage(msgStream)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}
	expectedFrom := "From someone"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}
	msg, err := mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Error(err)
	}
	if msg == nil {
		t.Fatalf("no message from mail.ReadMessage")
	}
	headers := map[string]string{}
	headers["Subject"] = "To interpretation"
	headers["From"] = "bubbles@bubbletown.com"
	headers["To"] = "mrmxpdstk@lazytown.com"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body:    ">From all of us, to all of you, be happy!\n",
	}, msg)
	if err != nil {
		t.Error(err)
	}

	msgStream = bytes.NewBuffer([]byte{})
	from, err = box.NextMessage(msgStream)
	if err == nil {
		t.Errorf("expected error, but got nil")
	}
	if err != io.EOF {
		t.Errorf("expected io.EOF but got %s", err)
	}
	expectedFrom = "From someone-else"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}
	msg, err = mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Errorf("ReadMessage failed: %s\nmsgStream held:\n%s\n", err, msgStream.Bytes())
	}
	if msg == nil {
		t.Fatalf("no message from mail.ReadMessage")
	}
	headers = map[string]string{}
	headers["Subject"] = "Bestest offer in the universe!!11!!"
	headers["From"] = "mrspam@corporate.corp.com"
	headers["To"] = "mrmxpdstk@lazytown.com"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body: `You won't believe these prices!
>From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
`,
	}, msg)
	if err != nil {
		t.Error(err)
	}
}

func TestReadMBOXRC(t *testing.T) {
	box := NewReader(bytes.NewBuffer([]byte(mboxo)))
	box.Type = MBOXRD

	msgStream := bytes.NewBuffer([]byte{})
	from, err := box.NextMessage(msgStream)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}
	expectedFrom := "From someone"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}
	msg, err := mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Error(err)
	}
	headers := map[string]string{}
	headers["Subject"] = "To interpretation"
	headers["From"] = "bubbles@bubbletown.com"
	headers["To"] = "mrmxpdstk@lazytown.com"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body:    "From all of us, to all of you, be happy!\n",
	}, msg)
	if err != nil {
		t.Error(err)
	}

	msgStream = bytes.NewBuffer([]byte{})
	from, err = box.NextMessage(msgStream)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}
	if err != io.EOF {
		t.Errorf("expected an io.EOF error but got %s", err)
	}
	expectedFrom = "From someone-else"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}

	msg, err = mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Error(err)
	}
	headers = map[string]string{}
	headers["Subject"] = "Bestest offer in the universe!!11!!"
	headers["From"] = "mrspam@corporate.corp.com"
	headers["To"] = "mrmxpdstk@lazytown.com"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body: `You won't believe these prices!
From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
`,
	}, msg)
	if err != nil {
		t.Error(err)
	}
}

func TestReadMBOXCL(t *testing.T) {
	box := NewReader(bytes.NewBuffer([]byte(mboxcl)))
	box.Type = MBOXCL
	msgStream := bytes.NewBuffer([]byte{})
	from, err := box.NextMessage(msgStream)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}
	expectedFrom := "From someone"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}
	msg, err := mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Error(err)
	}
	headers := map[string]string{}
	headers["Subject"] = "To interpretation"
	headers["From"] = "bubbles@bubbletown.com"
	headers["To"] = "mrmxpdstk@lazytown.com"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body:    "From all of us, to all of you, be happy!\n",
	}, msg)
	if err != nil {
		t.Error(err)
	}

	msgStream = bytes.NewBuffer([]byte{})
	from, err = box.NextMessage(msgStream)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}
	expectedFrom = "From someone-else"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}

	msg, err = mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Error(err)
	}
	headers = map[string]string{}
	headers["Subject"] = "Bestest offer in the universe!!11!!"
	headers["From"] = "mrspam@corporate.corp.com"
	headers["To"] = "mrmxpdstk@lazytown.com"
	headers["From mug"] = "weird header"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body: `You won't believe these prices!
From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
`,
	}, msg)
	if err != nil {
		t.Error(err)
	}

	msgStream = bytes.NewBuffer([]byte{})
	from, err = box.NextMessage(msgStream)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}
	if err != io.EOF {
		t.Errorf("expected an io.EOF error but got %s", err)
	}
	expectedFrom = "From nobody"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}

	msg, err = mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Error(err)
	}
	headers = map[string]string{}
	headers["Subject"] = "Mysterious Jenkins"
	headers["From"] = "nobody@nowhere.man"
	headers["To"] = "mrmxpdstk@lazytown.com"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body:    ``,
	}, msg)
	if err != nil {
		t.Error(err)
	}
}

func TestReadMBOXCL2(t *testing.T) {
	box := NewReader(bytes.NewBuffer([]byte(mboxcl)))
	box.Type = MBOXCL2
	msgStream := bytes.NewBuffer([]byte{})
	from, err := box.NextMessage(msgStream)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}
	expectedFrom := "From someone"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}
	msg, err := mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Error(err)
	}
	headers := map[string]string{}
	headers["Subject"] = "To interpretation"
	headers["From"] = "bubbles@bubbletown.com"
	headers["To"] = "mrmxpdstk@lazytown.com"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body:    ">From all of us, to all of you, be happy!\n",
	}, msg)
	if err != nil {
		t.Error(err)
	}

	msgStream = bytes.NewBuffer([]byte{})
	from, err = box.NextMessage(msgStream)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}
	expectedFrom = "From someone-else"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}

	msg, err = mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Error(err)
	}
	headers = map[string]string{}
	headers["Subject"] = "Bestest offer in the universe!!11!!"
	headers["From"] = "mrspam@corporate.corp.com"
	headers["To"] = "mrmxpdstk@lazytown.com"
	headers[">From mug"] = "weird header"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body: `You won't believe these prices!
>From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
`,
	}, msg)
	if err != nil {
		t.Error(err)
	}

	msgStream = bytes.NewBuffer([]byte{})
	from, err = box.NextMessage(msgStream)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}
	if err != io.EOF {
		t.Errorf("expected an io.EOF error but got %s", err)
	}
	expectedFrom = "From nobody"
	if from != expectedFrom {
		t.Errorf("expected %s but got %s", expectedFrom, from)
	}

	msg, err = mail.ReadMessage(bytes.NewBuffer(msgStream.Bytes()))
	if err != nil {
		t.Error(err)
	}
	headers = map[string]string{}
	headers["Subject"] = "Mysterious Jenkins"
	headers["From"] = "nobody@nowhere.man"
	headers["To"] = "mrmxpdstk@lazytown.com"
	err = CheckMessage(MsgTest{
		Headers: headers,
		Body:    ``,
	}, msg)
	if err != nil {
		t.Error(err)
	}
}

func TestReadMBOXCLBadContentLength(t *testing.T) {
	box := NewReader(bytes.NewBuffer([]byte(badmboxcl)))
	box.Type = MBOXCL
	msgStream := bytes.NewBuffer([]byte{})
	_, err := box.NextMessage(msgStream)
	if err == nil {
		t.Errorf("expected an error, but it worked")
	}
}

func TestReadMBOXCLBadLength(t *testing.T) {
	box := NewReader(bytes.NewBuffer([]byte(badlenmboxcl)))
	box.Type = MBOXCL
	msgStream := bytes.NewBuffer([]byte{})
	_, err := box.NextMessage(msgStream)
	if err != nil {
		t.Errorf("expected this to work, but we got an error: %s", err)
	}
	msgStream = bytes.NewBuffer([]byte{})
	_, err = box.NextMessage(msgStream)
	if err == nil {
		t.Errorf("expected an error, but it worked")
	}
}

func TestReadMBOXCL2BadContentLength(t *testing.T) {
	box := NewReader(bytes.NewBuffer([]byte(badmboxcl)))
	box.Type = MBOXCL2
	msgStream := bytes.NewBuffer([]byte{})
	_, err := box.NextMessage(msgStream)
	if err == nil {
		t.Errorf("expected an error, but it worked")
	}
}
