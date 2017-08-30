package test

import "testing"

func AssertEqualString(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got :%v, want: %v", got, want)
	}
}
