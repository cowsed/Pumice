package config

import "testing"

func Test(t *testing.T) {
	testCases := []struct {
		s string
		v Version
	}{
		{s: "3.2.1-DEV", v: Version{
			major:   3,
			minor:   2,
			patch:   1,
			comment: "DEV",
		}},
		{s: "0.0.0", v: Version{
			major:   0,
			minor:   0,
			patch:   0,
			comment: "",
		}},
	}
	for _, tC := range testCases {
		t.Run(tC.s, func(t *testing.T) {
			if tC.s != tC.v.String() {
				t.Errorf("Expected %v, got %v", tC.s, tC.v.String())
			}
		})
	}
}
