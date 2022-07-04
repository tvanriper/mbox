package mbox_test

import (
	"bytes"
	"fmt"
	"net/mail"

	"github.com/tvanriper/mbox"
)

// Imagine this is a file on your filesystem instead of a variable in your code.
const mboxrd string = `From bubbles@bubbletown.com Mon Jul 04 14:23:45 2022
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation

>From all of us, to all of you, be happy!
From mrspam@corporate.corp.com Mon Jul 04 15:02:15 2022
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
>From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
`

func ExampleMboxReader() {
	// Imagine you used os.Open instead of bytes.NewBuffer here.
	file := bytes.NewBuffer([]byte(mboxrd))
	mailReader := mbox.NewReader(file)
	mailReader.Type = mbox.MBOXRD

	var err error
	mailBytes := bytes.NewBuffer([]byte{})
	for err == nil {
		_, err = mailReader.NextMessage(mailBytes)
		msg, e := mail.ReadMessage(bytes.NewBuffer(mailBytes.Bytes()))
		if e != nil {
			panic(e)
		}
		fmt.Println("From:")
		fmt.Println(msg.Header.Get("From"))
		fmt.Println("Subject:")
		fmt.Println(msg.Header.Get("Subject"))
		mailBytes.Reset()
	}

	// Output:
	// From:
	// bubbles@bubbletown.com
	// Subject:
	// To interpretation
	// From:
	// mrspam@corporate.corp.com
	// Subject:
	// Bestest offer in the universe!!11!!
}
