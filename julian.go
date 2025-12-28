// Package julian implements functions for working with Julian dates.
//
// The Julian date is the continuous count of days since the beginning
// of the Julian Period on January 1, 4713 BC at noon Universal Time.
// For example, the Julian date for January 1, 2000 at noon UTC is 2451545.0.
// The fractional part of the day counts the time since the preceding noon UTC.
// For example, January 1, 2000 at 6 PM UTC is 2451545.25.
package julian

import (
	"math"
	"time"
)

type Date float64

const (
	day_seconds     = 86400
	day_nanoseconds = day_seconds * 1_000_000_000
	julian_unix     = 2440587.5 // 1/1/1970
	days_p_century  = 36525
	epoch_j2000     = 2451545
)

// Time returns a julian date version of the time.
func Time(t time.Time) Date {
	j := float64(t.UnixNano())/day_nanoseconds + julian_unix
	return Date(j)
}

// NewDate returns the julian date corresponding to yyyy-mm-dd hh:mm:ss + nsec
// nanoseconds in the appropriate zone for that time in the given location.
//
// The month, day, hour, min, sec, and nsec values may be outside their usual
// ranges and will be normalized during the conversion. For example, October
// 32 converts to November 1.
//
// A daylight savings time transition skips or repeats times. For example, in
// the United States, March 13, 2011 2:15am never occurred, while November 6,
// 2011 1:15am occurred twice. In such cases, the choice of time zone, and
// therefore the time, is not well-defined. NewDate returns a time that is correct
// in one of the two zones involved in the transition, but it does not
// guarantee which.
//
// NewDate panics if loc is nil.
func NewDate(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Date {
	t := time.Date(year, month, day, hour, min, sec, nsec, loc)
	jd := float64(t.UnixNano())/day_nanoseconds + julian_unix
	return Date(jd)
}

// Gregorian
func (jd Date) Gregorian() time.Time {
	return time.Unix(0, jd.UnixNano())
}

// Unix returns the Unix time corresponding to the julian date
func (jd Date) Unix() int64 {
	return int64((jd - julian_unix) * day_seconds)
}

// UnixNano returns julian date as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC.
//
// The result is undefined if the Unix time in nanoseconds cannot be represented
// by an int64 (a date before the year 1678 or after 2262). Note that this means
// the result of calling UnixNano on the zero Time is undefined. The result does
// not depend on the location associated with j.
func (jd Date) UnixNano() int64 {
	return int64((jd - julian_unix) * day_nanoseconds)
}

// Time returns the time fraction.
func (jd Date) Time() float64 {
	return math.Mod(float64(jd), 1)
}

// Duration returns the time fraction.
func (jd Date) Duration() time.Duration {
	return time.Duration(day_nanoseconds * math.Mod(float64(jd), 1))
}

// Day returns the float64 representation of the Julian day.
func (jd Date) Day() float64 {
	return float64(jd)
}

// DayNumber returns the integer part of the Julian day.
func (jd Date) DayNumber() int {
	return int(jd)
}

// Century returns the Julian century.
func (jd Date) Century() float64 {
	return float64(jd-epoch_j2000) / days_p_century
}
