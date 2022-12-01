package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func GetRandom_BTree(degree int, n int) (*Set_BTree, []int) {
	s := NewSet_BTree(degree)
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		v := rand.Intn(100000)
		if v%2 == 0 {
			// test negative values!
			v = -v
		}
		s.insert(v)
		vals[i] = v
	}
	return s, vals
}

func TestInsert_BTree(t *testing.T) {
	for degree := 3; degree < 10; degree++ {
		s, val := GetRandom_BTree(degree, 10*degree*degree)
		for _, v := range val {
			location := s.search(v)
			if location == nil {
				t.Errorf("Error trying to find value.")
				fmt.Printf("Tried to find %d. Found nil\n", v)
			}
			if location.value != location.cell.keys[location.key_idx] {
				t.Errorf("Incorrect location read when trying to find value.")
				fmt.Printf("location value is %d, cell value is %d\n", location.value, location.cell.keys[location.key_idx])
			}
			if location.value != v {
				t.Error("Error trying to find value.")
				fmt.Printf("Expected %d, found %d\n", v, location.value)
			}
			if location.cell.parent != nil && location.cell.parent.children[location.cell.ID] != location.cell {
				t.Error("ID of cell is wrong")
			}
		}
	}
}

func TestDelete_BTree(t *testing.T) {
	for degree := 3; degree < 10; degree++ {
		s, val := GetRandom_BTree(degree, 10*degree*degree)
		for _, v := range val {
			s.delete(v)
			location := s.search(v)
			if location != nil {
				t.Errorf("Error, value was not deleted.")
				fmt.Printf("Tried to delete %d, but value was still found\n", v)
			}
		}
	}
}
