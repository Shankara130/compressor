package optimizer

import "testing"

func TestParseSeconds(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"00:00:00", 0},
		{"00:00:05", 5},
		{"00:01:30", 90},
		{"01:00:00", 3600},
		{"00:00:05.500000", 5.5},
		{"00:00:12.345000", 12.345},
		{"00:02:03.250000", 123.25},
	}
	for _, c := range cases {
		got, ok := parseSeconds(c.in)
		if !ok {
			t.Errorf("parseSeconds(%q) returned not-ok", c.in)
			continue
		}
		if got != c.want {
			t.Errorf("parseSeconds(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseSeconds_Invalid(t *testing.T) {
	for _, in := range []string{"bad", "1:2:3:4", "00:99", ""} {
		if _, ok := parseSeconds(in); ok {
			t.Errorf("parseSeconds(%q) expected not-ok", in)
		}
	}
}
