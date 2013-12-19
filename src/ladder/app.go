package main

import (
	"fmt"
	"unicode/utf8"
)

type WordConnections []WordNode

type WordNode struct {
	value       string
	connections WordConnections
}

func NewWordNode(value string) *WordNode {
	return &WordNode{value: value, connections: make(WordConnections, 0)}
}

// Iterate over w.value and check other.value at the
// point and see if the characters are the same
// if we find more than one diff then return false
func (w *WordNode) oneDiff(other *WordNode) bool {
	diffs := 0
	for i, width := 0, 0; i < len(w.value); i += width {
		wRuneVal, size := utf8.DecodeRuneInString(w.value[i:])
		otherRuneVal, _ := utf8.DecodeRuneInString(other.value[i:])
		if wRuneVal != otherRuneVal {
			diffs++
		}
		if diffs > 1 {
			return false
		}
		width = size
	}
	return true
}

func main() {
	fmt.Println("Hey!")
	w1 := NewWordNode("mooo")
	fmt.Printf("My word is: %s\n", w1.value)
}
