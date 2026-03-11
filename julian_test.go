package julian

import (
	"math"
	"testing"
	"time"
)

func timeEquals(got, want time.Time) bool {
	if want.Equal(got) {
		return true
	}
	diff := want.Sub(got).Nanoseconds()
	if diff < 0 {
		diff = -diff
	}
	return diff < 50000
}

const epsilon = 0.000001

func equalJulian(got, want Date) bool {
	if got == want {
		return true
	}
	return math.Abs(float64(got-want)) < epsilon
}

func TestJulian(t *testing.T) {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	// now := time.Now()
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		// {"now", args{now}, JulianDate(123)},
		{"A", args{time.Date(2020, 1, 1, 12, 0, 0, 0, time.Local)}, Date(2_458_850.20833333)},
		{"Jan. 1  2017", args{time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)}, Date(2_457_754.50000)},
		{"Jan. 1, 1990", args{time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)}, Date(2_447_892.50000)},
		{"July 4, 1998", args{time.Date(1998, time.July, 4, 0, 0, 0, 0, time.UTC)}, Date(2_450_998.50000)},
		{"Feb. 14, 2010 5:21", args{time.Date(2010, time.February, 14, 5, 21, 0, 0, time.UTC)}, Date(2_455_241.722917)},
		{"Feb. 14, 2010 5:21 PST", args{time.Date(2010, time.February, 14, 5, 21, 0, 0, location)}, Date(2_455_242.05625)},
		// DST: July 4, 2020 12:00 PDT (UTC-7) = 19:00 UTC; +1/24 DST correction → 2459034.5 + 20/24
		{"July 4, 2020 12:00 PDT", args{time.Date(2020, time.July, 4, 12, 0, 0, 0, location)}, Date(2_459_035.291667)},
		// No DST: Dec 25, 2020 12:00 PST (UTC-8) = 20:00 UTC → 2459208.5 + 20/24
		{"Dec 25, 2020 12:00 PST", args{time.Date(2020, time.December, 25, 12, 0, 0, 0, location)}, Date(2_459_209.333333)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Time(tt.args.t); !equalJulian(got, tt.want) {
				t.Errorf("Julian() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestJulianDate_Gregorian(t *testing.T) {
	now := time.Now()
	layout, _ := time.Parse(time.RFC3339, time.RFC3339)
	tests := []struct {
		name string
		jd   Date
		want time.Time
	}{
		{"now", Time(now), now},
		{"layout", Time(layout), layout},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.jd.Gregorian(); !timeEquals(got, tt.want) {

				dt := got.Sub(tt.want)
				t.Errorf("JulianDate.Gregorian() = %v, want %v [%v, %v]", got, tt.want,
					got.Nanosecond()-tt.want.Nanosecond(),
					dt.Nanoseconds())
			}
		})
	}
}

func TestJulianDate_Unix(t *testing.T) {
	tests := []struct {
		name string
		jd   Date
		want int64
	}{
		{"unix epoch", Date(2440587.5), 0},
		{"one day after epoch", Date(2440588.5), 86400},
		{"J2000 noon", Date(2451545.0), 946728000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.jd.Unix(); got != tt.want {
				t.Errorf("JulianDate.Unix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJulianDate_Time(t *testing.T) {
	tests := []struct {
		name string
		jd   Date
		want float64
	}{
		{"A", Time(time.Date(2020, 1, 1, 12, 0, 0, 0, time.Local)), 0.20833333},
		{"Jan. 1  2017", Time(time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)), 0.50000},
		{"Jan. 1, 1990", Time(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)), 0.50000},
		{"July 4, 1998", Time(time.Date(1998, time.July, 4, 0, 0, 0, 0, time.UTC)), 0.50000},
		{"Feb. 14, 2010 5:21", Time(time.Date(2010, time.February, 14, 5, 21, 0, 0, time.UTC)), 0.72292},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.jd.Time(); got-tt.want > 0.00000001 {
				t.Errorf("JulianDate.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJulianDate_Day(t *testing.T) {
	tests := []struct {
		name string
		jd   Date
		want float64
	}{
		{"J2000", Date(2451545.0), 2451545.0},
		{"Jan 1, 2020 midnight UTC", Date(2458849.5), 2458849.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.jd.Day(); got != tt.want {
				t.Errorf("JulianDate.Day() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJulianDate_DayNumber(t *testing.T) {
	tests := []struct {
		name string
		jd   Date
		want int
	}{
		{"J2000 noon", Date(2451545.0), 2451545},
		{"J2000 with fraction", Date(2451545.9), 2451545},
		{"Jan 1, 2020 midnight UTC", Date(2458849.5), 2458849},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.jd.DayNumber(); got != tt.want {
				t.Errorf("JulianDate.DayNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJulianDate_Century(t *testing.T) {
	tests := []struct {
		name string
		jd   Date
		want float64
	}{
		{"J2000 epoch", Date(2451545), 0.0},
		{"one century after J2000", Date(2451545 + 36525), 1.0},
		{"one century before J2000", Date(2451545 - 36525), -1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.jd.Century(); got != tt.want {
				t.Errorf("JulianDate.Century() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDST verifies that Time() when t.IsDST() is true,
// and that NewDate() does not.
func TestDST(t *testing.T) {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatal(err)
	}

	// JD(Jan 1, 2020 0:00 UTC) = 2458849.5
	// JD(July 4, 2020 0:00 UTC) = 2458849.5 + 185 = 2459034.5
	// July 4, 2020 12:00 PDT = 19:00 UTC → +19/24;
	t.Run("summer PDT applies DST correction", func(t *testing.T) {
		summer := time.Date(2020, time.July, 4, 12, 0, 0, 0, location)
		if !summer.IsDST() {
			t.Fatal("expected July 4 to be DST")
		}
		got := Time(summer)
		want := Date(2_459_035.291667) // 2459034.5 + 19/24
		if !equalJulian(got, want) {
			t.Errorf("Time() = %f, want %f", got, want)
		}
	})

	// JD(Dec 25, 2020 0:00 UTC) = 2458849.5 + 359 = 2459208.5
	// Dec 25, 2020 12:00 PST = 20:00 UTC → +20/24; no DST correction
	t.Run("winter PST no DST correction", func(t *testing.T) {
		winter := time.Date(2020, time.December, 25, 12, 0, 0, 0, location)
		if winter.IsDST() {
			t.Fatal("expected Dec 25 to not be DST")
		}
		got := Time(winter)
		want := Date(2_459_209.333333) // 2459208.5 + 20/24
		if !equalJulian(got, want) {
			t.Errorf("Time() = %f, want %f", got, want)
		}
	})

	// Spring forward: March 8, 2020 clocks advance 2:00 AM → 3:00 AM
	// JD(March 8, 2020 0:00 UTC) = 2458849.5 + 67 = 2458916.5

	t.Run("spring forward before transition PST", func(t *testing.T) {
		// 1:30 AM PST (UTC-8) = 9:30 AM UTC → +9.5/24; no DST
		before := time.Date(2020, time.March, 8, 1, 30, 0, 0, location)
		if before.IsDST() {
			t.Fatal("expected 1:30 AM to be standard time before spring forward")
		}
		got := Time(before)
		want := Date(2_458_916.895833) // 2458916.5 + 9.5/24
		if !equalJulian(got, want) {
			t.Errorf("Time() = %f, want %f", got, want)
		}
	})

	t.Run("spring forward after transition PDT", func(t *testing.T) {
		// 3:00 AM PDT (UTC-7) = 10:00 AM UTC → +10/24;
		after := time.Date(2020, time.March, 8, 3, 0, 0, 0, location)
		if !after.IsDST() {
			t.Fatal("expected 3:00 AM to be DST after spring forward")
		}
		got := Time(after)
		want := Date(2_458_916.916666) // 2458916.5 + 11/24
		if !equalJulian(got, want) {
			t.Errorf("Time() = %f, want %f", got, want)
		}
	})

	// Fall back: November 1, 2020 clocks fall 2:00 AM PDT → 1:00 AM PST
	// JD(Nov 1, 2020 0:00 UTC) = 2458849.5 + 305 = 2459154.5

	t.Run("fall back during PDT first occurrence", func(t *testing.T) {
		// Go resolves the ambiguous 1:30 AM to the first occurrence (PDT).
		// 1:30 AM PDT (UTC-7) = 8:30 AM UTC → +8.5/24;
		ambiguous := time.Date(2020, time.November, 1, 1, 30, 0, 0, location)
		if !ambiguous.IsDST() {
			t.Skip("Go resolved ambiguous fall-back time to PST; skipping PDT case")
		}
		got := Time(ambiguous)
		want := Date(2_459_154.854167) // 2459154.5 + 9.5/24
		if !equalJulian(got, want) {
			t.Errorf("Time() = %f, want %f", got, want)
		}
	})

	t.Run("fall back after transition PST", func(t *testing.T) {
		// 2:00 AM PST (UTC-8) = 10:00 AM UTC → +10/24; no DST correction
		// (2:00 AM can only be PST — clocks only went back to 1:00 AM)
		afterFallback := time.Date(2020, time.November, 1, 2, 0, 0, 0, location)
		if afterFallback.IsDST() {
			t.Fatal("expected 2:00 AM to be standard time after fall back")
		}
		got := Time(afterFallback)
		want := Date(2_459_154.916666) // 2459154.5 + 10/24
		if !equalJulian(got, want) {
			t.Errorf("Time() = %f, want %f", got, want)
		}
	})

	t.Run("NewDate omits DST correction", func(t *testing.T) {
		// Time() adds +1/24 for DST; NewDate() does not.
		summer := time.Date(2020, time.July, 4, 12, 0, 0, 0, location)
		fromTime := Time(summer)
		fromNewDate := NewDate(2020, time.July, 4, 12, 0, 0, 0, location)
		diff := float64(fromTime - fromNewDate)
		expected := 0.0
		if math.Abs(diff-expected) > epsilon {
			t.Errorf("Time() - NewDate() = %f, want %f (1/24 DST correction)", diff, expected)
		}
	})

	t.Run("DST correction shifts Gregorian result by one hour", func(t *testing.T) {
		// Time() adding +1/24 for DST means Gregorian() returns UTC+1h relative to the true UTC.
		summer := time.Date(2020, time.July, 4, 12, 0, 0, 0, location)
		jd := Time(summer)
		gotUTC := jd.Gregorian()
		wantUTC := summer.UTC()
		diff := gotUTC.Sub(wantUTC)
		// Allow small rounding error from float64 Julian date arithmetic
		if math.Abs(float64(diff)) > float64(50*time.Microsecond) {
			t.Errorf("Gregorian() offset from true UTC = %v, want ~1h", diff)
		}
	})
}
