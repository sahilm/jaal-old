package test

import "testing"

func AssertEqualString(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got :%v, want: %v", got, want)
	}
}

func AssertEqualInt(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got :%v, want: %v", got, want)
	}
}

func AssertNotNil(t *testing.T, got interface{}) {
	if got == nil {
		t.Errorf("%v is nil. Expected not nil", got)
	}
}
