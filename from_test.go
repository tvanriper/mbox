package mbox

import (
	"testing"
	"time"
)

func TestParseFrom(t *testing.T) {
	from := "From pi@rpi.cu Mon Jul 04 19:23:45 2022 Gads, more crap from this guy?"
	addr, date, moreinfo, err := ParseFrom(from)
	if err != nil {
		t.Errorf("expected success but it failed: %s", err)
	}
	expectedAddr := "pi@rpi.cu"
	if addr != expectedAddr {
		t.Errorf("expected %s but got %s", expectedAddr, addr)
	}
	expectedTime, err := time.Parse(TimeFormat, "Mon Jul 04 19:23:45 2022")
	if err != nil {
		t.Error(err)
	}
	if expectedTime != date {
		t.Errorf("expected %s but got %s", expectedTime.String(), date.String())
	}
	expectedMoreinfo := "Gads, more crap from this guy?"
	if moreinfo != expectedMoreinfo {
		t.Errorf("expected %s but got %s", expectedMoreinfo, moreinfo)
	}
}

func TestBuildFrom(t *testing.T) {
	addr := "pi@rpi.cu"
	date, err := time.Parse(TimeFormat, "Fri Jul 04 19:23:45 2022")
	if err != nil {
		t.Error(err)
	}
	moreinfo := "Gads, more crap from this guy?"

	expectedFrom := "From pi@rpi.cu Mon Jul 04 19:23:45 2022 Gads, more crap from this guy?"

	from := BuildFrom(addr, date, moreinfo)
	if from != expectedFrom {
		t.Errorf("expected [%s] but got [%s]", expectedFrom, from)
	}
}
