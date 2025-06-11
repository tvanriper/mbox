package mbox

import (
	"fmt"
	"strings"
	"time"
)

// TimeFormat describes the way mbox files format the 'From ' date/time field.
// One may use this with time.Format and time.Parse functions.
var TimeFormat string = "Mon Jan  2 15:04:05 2006"

// ParseFrom parses a from string to its component parts.
// It helpfully translates the date/time to a time.Time.  A mailer might use
// this information in some way, if needed.
func ParseFrom(from string) (addr string, date time.Time, moreinfo string, err error) {
	data, _ := strings.CutPrefix(from, "From ")
	addr, remainder, _ := strings.Cut(data, " ")
	remainder = strings.TrimSpace(remainder)
	if len(remainder) >= len(TimeFormat) {
		date, err = time.Parse(TimeFormat, strings.TrimSpace(remainder[:len(TimeFormat)]))
		moreinfo = remainder[len(TimeFormat):]
	}

	return addr, date, strings.TrimSpace(moreinfo), err
}

// BuildFrom creates a from string based on the provided data.
// A mailer might build this to add to an mbox that it creates.
func BuildFrom(addr string, date time.Time, moreinfo string) (result string) {
	return fmt.Sprintf("From %s %s %s", addr, date.Format(TimeFormat), moreinfo)
}
