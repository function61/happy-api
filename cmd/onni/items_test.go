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
	copied := append([]Happiness{}, happiness...)

	sort.Slice(copied, func(a, b int) bool { return copied[a].Id < copied[b].Id })

	for idx := range copied {
		if copied[idx].Id != happiness[idx].Id {
			panic("items in happiness slice are not properly sorted! offending Id = " + happiness[idx].Id)
		}
	}
}
