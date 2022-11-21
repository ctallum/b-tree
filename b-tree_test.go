package main

import (
	"fmt"
	"math/rand"

	//"sort"
	"testing"
)

func GetRandom_BTree(degree int, n int) (*Set_BTree, []int) {
	s := NewSet_BTree(degree)
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		v := rand.Intn(1000)
		if v%2 == 0 {
			// test negative values!
			v = -v
		}
		Insert_BTree(s, v)
		vals[i] = v
	}

	return s, vals
}

func TestInsert_BTree(t *testing.T) {
	s, val := GetRandom_BTree(5, 100)

	fmt.Println("Inserting values", val)

	for _, v := range val {
		location := Search_BTree(s, v)
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
	}

}
