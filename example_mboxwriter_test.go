package mbox_test

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/tvanriper/mbox"
)

// PrintWithLineCount displays the mbox in a clean way for godoc.
func PrintWithLineCount(text string) {
	scanner := bufio.NewScanner(bytes.NewBuffer([]byte(text)))
	lineCount := 0
	for scanner.Scan() {
		lineCount += 1
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		l := len(trimmed)
		if l == 0 {
			fmt.Printf("%02d:\n", lineCount)
			continue
		}
		fmt.Printf("%02d: %s\n", lineCount, trimmed)
	}
}

func ExampleMboxWriter() {
	// Imagine these two emails came from some larger piece of software that
	// provides email.

	fromTime1, _ := time.Parse(time.RFC3339, "2022-07-04T14:03:04Z")
	from1 := mbox.BuildFrom("bubbles@bubbletown.com", fromTime1, "")
	email1 := `From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation

From all of us, to all of you, be happy!
`
	fromTime2, _ := time.Parse(time.RFC3339, "2022-07-04T13:12:34Z")
	from2 := mbox.BuildFrom("mrspam@corporate.corp.com", fromTime2, "")
	email2 := `From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!
From X: Quoi?

You won't believe these prices!
`

	// Imagine this is on a filesystem instead of bytes in a buffer.
	file := bytes.NewBuffer([]byte{})
	mailWriter := mbox.NewWriter(file)

	var err error
	err = mailWriter.WriteMail(from1, bytes.NewBuffer([]byte(email1)))
	if err != nil {
		panic(err)
	}
	err = mailWriter.WriteMail(from2, bytes.NewBuffer([]byte(email2)))
	if err != nil {
		panic(err)
	}

	fmt.Println("File contents:")
	PrintWithLineCount(file.String())

	// Output:
	// File contents:
	// 01: From bubbles@bubbletown.com Mon Jul 04 14:03:04 2022
	// 02: From: bubbles@bubbletown.com
	// 03: To: mrmxpdstk@lazytown.com
	// 04: Subject: To interpretation
	// 05:
	// 06: >From all of us, to all of you, be happy!
	// 07:
	// 08: From mrspam@corporate.corp.com Mon Jul 04 13:12:34 2022
	// 09: From: mrspam@corporate.corp.com
	// 10: To: mrmxpdstk@lazytown.com
	// 11: Subject: Bestest offer in the universe!!11!!
	// 12: >From X: Quoi?
	// 13:
	// 14: You won't believe these prices!
	// 15:
}
