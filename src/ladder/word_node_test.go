package main

import (
	"testing"
)

func TestOneDiff(t *testing.T) {
	w1 := NewWordNode("cat")
	w2 := NewWordNode("rat")

	oneAway := w1.oneDiff(w2)
	if !oneAway {
		t.Errorf("Expected words to be one away")
	}
}
