package mbox

import (
	"fmt"
	"strings"
	"time"
)

// TimeFormat describes the way mbox files format the 'From ' date/time field.
var TimeFormat string = "Mon Jan 02 15:04:05 2006"

// ParseFrom parses a from string to its component parts.
// It helpfully translates the date/time to a time.Time.  A mailer might use
// this information in some way, if needed.
func ParseFrom(from string) (addr string, date time.Time, moreinfo string, err error) {
	splitted := strings.Split(from, " ")
	l := len(splitted)
	if l > 1 {
		addr = splitted[1]
	}
	if l > 6 {
		strDate := ""
		for i := 2; i <= 6; i++ {
			strDate = fmt.Sprintf("%s %s", strDate, splitted[i])
		}
		date, err = time.Parse(TimeFormat, strings.TrimSpace(strDate))
	}
	if l > 7 {
		for i := 7; i < l; i++ {
			moreinfo = fmt.Sprintf("%s %s", moreinfo, splitted[i])
		}
	}
	return addr, date, strings.TrimSpace(moreinfo), err
}

// BuildFrom creates a from string based on the provided data.
// A mailer might build this to add to an mbox that it creates.
func BuildFrom(addr string, date time.Time, moreinfo string) (result string) {
	return fmt.Sprintf("From %s %s %s", addr, date.Format(TimeFormat), moreinfo)
}
