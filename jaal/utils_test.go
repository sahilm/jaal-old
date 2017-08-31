package jaal

import (
	"testing"
)

func TestToSHA256(t *testing.T) {
	t.Run("it return SHA256 of input", func(t *testing.T) {
		got, err := ToSHA256("hello")
		if err != nil {
			t.Fatal(err)
		}

		want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})
}
