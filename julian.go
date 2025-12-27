package julian

import (
	"math"
	"time"
)

type Julian float64

const (
	seconds        = 86400
	nanoseconds    = seconds * 1000000000
	julian_unix    = 2440587.5 // 1/1/1970
	days_p_century = 36525
	epoch_j2000    = 2451545
)

// Time returns a julian date version of the time.
func Time(t time.Time) Julian {
	j := float64(t.UnixNano())/nanoseconds + julian_unix
	return Julian(j)
}

// Date returns the julian time corresponding to
//
// yyyy-mm-dd hh:mm:ss + nsec nanoseconds
// in the appropriate zone for that time in the given location.
//
// The month, day, hour, min, sec, and nsec values may be outside their usual
// ranges and will be normalized during the conversion. For example, October
// 32 converts to November 1.
//
// A daylight savings time transition skips or repeats times. For example, in
// the United States, March 13, 2011 2:15am never occurred, while November 6,
// 2011 1:15am occurred twice. In such cases, the choice of time zone, and
// therefore the time, is not well-defined. Date returns a time that is correct
// in one of the two zones involved in the transition, but it does not
// guarantee which.
//
// Date panics if loc is nil.
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Julian {
	t := time.Date(year, month, day, hour, min, sec, nsec, loc)
	j := float64(t.UnixNano())/nanoseconds + julian_unix
	return Julian(j)
}

// Gregorian
func (j Julian) Gregorian() time.Time {
	return time.Unix(0, j.UnixNano())
}

// Unix returns the local Time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC. It is
// valid to pass nsec outside the range [0, 999999999]. Not all sec
// values have a corresponding time value. One such value is 1<<63-1
// (the largest int64 value).
func (j Julian) Unix() int64 {
	return int64((j - julian_unix) * seconds)
}

// UnixNano returns j as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC. The result is undefined if the Unix time
// in nanoseconds cannot be represented by an int64 (a date before the
// year 1678 or after 2262). Note that this means the result of calling
// UnixNano on the zero Time is undefined. The result does not depend on
// the location associated with j.
func (j Julian) UnixNano() int64 {
	return int64((j - julian_unix) * nanoseconds)
}

// Time returns j time fraction.
func (j Julian) Time() float64 {
	return math.Mod(float64(j), 1)
}

// Day returns the float
func (j Julian) Day() float64 {
	return float64(j)
}

// DayNumber
func (j Julian) DayNumber() int {
	return int(j)
}

// Century
func (j Julian) Century() float64 {
	return float64(j-epoch_j2000) / days_p_century
}
