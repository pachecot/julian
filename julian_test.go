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
				t.Errorf("JulianDate.Gregorian() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.jd.Century(); got != tt.want {
				t.Errorf("JulianDate.Century() = %v, want %v", got, tt.want)
			}
		})
	}
}
