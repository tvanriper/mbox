package mbox_test

import (
	"bytes"
	"fmt"

	"github.com/tvanriper/mbox"
)

// Imagine this is a file on your filesystem instead of a variable in your code.
const mboxrd2 string = `From bubbles@bubbletown.com Mon Jul 04 14:23:45 2022
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

func ExampleDetectType() {
	// Imagine you used os.Open instead of bytes.NewBuffer here.
	reader := bytes.NewReader([]byte(mboxrd2))
	mbType, err := mbox.DetectType(reader)
	if err != nil {
		fmt.Printf("failed to detect mbox type: %s\n", err)
		return
	}
	switch mbType {
	case mbox.MBOXO:
		fmt.Println("MBOXO")
	case mbox.MBOXRD:
		fmt.Println("MBOXRD")
	case mbox.MBOXCL:
		fmt.Println("MBOXCL")
	case mbox.MBOXCL2:
		fmt.Println("MBOXCL2")
	default:
		fmt.Println("Unknown")
	}

	// Output:
	// MBOXRD
}
