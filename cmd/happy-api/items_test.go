package main

import (
	"sort"
	"testing"
)

// Q: why do this?
// A: inserting records randomly is the only way to minimize merge conflicts. random insertion
//    is easiest if we'll use a random ID and keep the slice sorted - otherwise if there's
//    no system, everyone would just insert their records as the last one, resulting in conflicts
func TestSortedHappiness(t *testing.T) {
	if !sort.SliceIsSorted(happiness, func(i, j int) bool { return happiness[i].Id < happiness[j].Id }) {
		panic("items in happiness slice are not properly sorted!")
	}
}
