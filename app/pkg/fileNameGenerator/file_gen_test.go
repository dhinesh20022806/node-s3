package filenamegenerator

import (
	"testing"
)

func TestFileName(t *testing.T) {

	filenamegen := NewFileNameGenerator(0)
	alpha := []string{
		"a", "b", "c", "d", "e", "f", "g", "h",
		"i", "j", "k", "l", "m", "n", "o", "p",
		"q", "r", "s", "t", "u", "v", "w", "x",
		"y", "z", "a-1", "b-1", "c-1", "d-1",
	}
	for i := 0; i < 30; i++ {
		want := alpha[i]
		got := filenamegen.Next()
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}

}
