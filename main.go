package main

import (
	"fmt"
)

type Cell_BTree struct {
	max_size int
	cur_size int
	keys     []int
	children []*Cell_BTree
	parent   *Cell_BTree
	height   int
	is_leaf  bool
}

type Set_BTree struct {
	degree  int
	root    *Cell_BTree
	maximum *int
	minimum *int
}

type Value_Location struct {
	cell    *Cell_BTree
	key_idx int
	value   int
}

func NewCell_BTree(size int, is_leaf bool) *Cell_BTree {
	keys := make([]int, size-1)
	children := make([]*Cell_BTree, size)
	cell := &Cell_BTree{size, 0, keys, children, nil, 0, is_leaf}
	return cell
}

func PrintSet_AVL(s *Set_BTree) {
	var print func(*Cell_BTree, string)
	print = func(c *Cell_BTree, prefix string) {
		if c != nil {
			to_right := c.cur_size / 2
			to_left := c.cur_size - to_right
			// print right half of children
			for i := 0; i < to_right; i++ {
				print(c.children[i], prefix+"   ")
			}
			// print current self
			fmt.Printf("%s%d : %d\n", prefix, c.keys, c.height)
			// print left hafl of children
			for i := 0; i < to_left; i++ {
				print(c.children[to_right+i], prefix+"   ")
			}
		} else {
			fmt.Printf("%s-\n", prefix)
		}
	}

	fmt.Println(" /")
	print(s.root, " |  ")
	fmt.Println(" \\")

	min := Minimum_BTree(s)
	max := Maximum_BTree(s)
	minStr := "min = -"
	maxStr := "max = -"
	if min != nil {
		minStr = fmt.Sprintf("min = %d", min.value)
	}
	if max != nil {
		maxStr = fmt.Sprintf("max = %d", max.value)
	}
	fmt.Printf("  %s  %s\n", minStr, maxStr)
}

func NewSet_BTree(degree int) *Set_BTree {
	tree := &Set_BTree{degree, nil, nil, nil}
	return tree
}

func Insert_BTree(s *Set_BTree, v int) {
	fmt.Println("Not implemented yet")
}

func Search_BTree(s *Set_BTree, v int) *Value_Location {
	fmt.Println("Not implemented yet")
	return nil
}

func Delete_BTree(s *Set_BTree, v *Value_Location) {
	fmt.Println("Not implemented yet")
}

func Minimum_BTree(s *Set_BTree) *Value_Location {
	fmt.Println("Not implemented yet")
	return nil
}

func Maximum_BTree(s *Set_BTree) *Value_Location {
	fmt.Println("Not implemented yet")
	return nil
}

func test_BTree() {
	// B-Tree degree
	size := 5
	// First, insert all of 10, 20, 30, ..., 150 in order into an empty B-tree.
	fmt.Println("***** INSERTING RIGHT **************************************************")
	s := NewSet_BTree(size)
	for i := 1; i < 16; i += 1 {
		v := i * 10
		fmt.Printf("Inserting %d\n", v)
		Insert_BTree(s, v)
		PrintSet_AVL(s)
	}

	// Similarly, insert all of 350, 340, 330, ..., 210 in order into an empty B-tree.
	fmt.Println("***** INSERTING LEFT **************************************************")
	t := NewSet_BTree(size)
	for i := 15; i > 0; i -= 1 {
		v := i*10 + 200
		fmt.Printf("Inserting %d\n", v)
		Insert_BTree(t, v)
		PrintSet_AVL(t)
	}
	// Delete the inserted values from the preceding tree, in order from 210, 220, ..., 350.
	fmt.Println("***** DELETING **************************************************")
	for i := 1; i < 16; i += 1 {
		c := Search_BTree(t, i*10+200)
		fmt.Printf("Deleting %d\n", i*10+200)
		Delete_BTree(t, c)
		PrintSet_AVL(t)
	}
}

func main() {
	test_BTree()

}
