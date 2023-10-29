package mbox_test

import (
	"strings"
	"testing"

	"github.com/tvanriper/mbox"
)

func TestDetectRD(t *testing.T) {
	mb := `From someone
From: chuckles@funbunny.org
To: hhefner@playboy.com
Subject: Closed Captioning

>From the engineers: do we really need closed captioning on this material?
Please let me know, it'd save a lot of money if we could avoid it.

From hhefner
From: hhefner@playboy.com
To: chuckles@funbunny.org
Subject: RE: Closed Captioning

Yes, silly as it sounds, we broadcast this material, and it must therefore have
closed captioning.  The deaf will enjoy reading the material.
`
	mType, err := mbox.DetectType(strings.NewReader(mb))
	if err != nil {
		t.Error(err)
	}
	if mType != mbox.MBOXRD {
		t.Errorf("expected %d but got %d", mbox.MBOXRD, mType)
	}

	// Ensure we also detect MBOXRD if the second mail has the character.
	mb = `From someone
From: chuckles@funbunny.org
To: hhefner@playboy.com
Subject: Closed Captioning

Do we really need closed captioning on this material? Please let me know, it'd
save a lot of money if we could avoid it.

From hhefner
From: hhefner@playboy.com
To: chuckles@funbunny.org
Subject: RE: Closed Captioning

>From my lawyers: yes, silly as it sounds, we broadcast this material, and it
must therefore have closed captioning.  The deaf will enjoy reading the
material.
`
	mType, err = mbox.DetectType(strings.NewReader(mb))
	if err != nil {
		t.Error(err)
	}
	if mType != mbox.MBOXRD {
		t.Errorf("expected %d but got %d", mbox.MBOXRD, mType)
	}
}

func TestDetectCL(t *testing.T) {
	mb := `From someone
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
Content-Length: 33

We should all try to enjoy life!
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
	mType, err := mbox.DetectType(strings.NewReader(mb))
	if err != nil {
		t.Error(err)
	}
	if mType != mbox.MBOXCL {
		t.Errorf("expected %d but got %d", mbox.MBOXCL, mType)
	}
}

func TestDetectCL2(t *testing.T) {
	mb := `From someone
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation
Content-Length: 33

We should all try to enjoy life!
From someone-else
>From mug: weird header
Content-Length: 129
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
From nobody
From: nobody@nowhere.man
To: mrmxpdstk@lazytown.com
Subject: Mysterious Jenkins
Content-Length: 0

`
	mType, err := mbox.DetectType(strings.NewReader(mb))
	if err != nil {
		t.Error(err)
	}
	if mType != mbox.MBOXCL2 {
		t.Errorf("expected %d but got %d", mbox.MBOXCL2, mType)
	}
}

func TestDetectO(t *testing.T) {
	mb := `From someone
From: bubbles@bubbletown.com
To: mrmxpdstk@lazytown.com
Subject: To interpretation

From all of us, to all of you, be happy!
From someone-else
From: mrspam@corporate.corp.com
To: mrmxpdstk@lazytown.com
Subject: Bestest offer in the universe!!11!!

You won't believe these prices!
From 1 cent to 11 cents, we carry the least expensive
line of jets this side of the Gobi Desert!
`
	mType, err := mbox.DetectType(strings.NewReader(mb))
	if err != nil {
		t.Error(err)
	}
	if mType != mbox.MBOXO {
		t.Errorf("expected %d but got %d", mbox.MBOXO, mType)
	}
}

func TestDetectAltLinefeed(t *testing.T) {
	mb := "From someone\r\nFrom: chuckles@funbunny.org\r\nTo: hhefner@playboy.com\r\nSubject: Closed Captioning\r\n\r\n>From the engineers: do we really need closed captioning on this material?\r\nPlease let me know, it'd save a lot of money if we could avoid it.\r\n\r\nFrom hhefner\r\nFrom: hhefner@playboy.com\r\nTo: chuckles@funbunny.org\r\nSubject: RE: Closed Captioning\r\n\r\nYes, silly as it sounds, we broadcast this material, and it must therefore have\r\nclosed captioning.  The deaf will enjoy reading the material.\r\n"
	mType, err := mbox.DetectType(strings.NewReader(mb))
	if err != nil {
		t.Error(err)
	}
	if mType != mbox.MBOXRD {
		t.Errorf("expected %d but got %d", mbox.MBOXRD, mType)
	}

}
