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
	ID       int
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

func NewCell_BTree(size int, is_leaf bool, ID int) *Cell_BTree {
	keys := make([]int, size)
	children := make([]*Cell_BTree, size+1)
	cell := &Cell_BTree{size, 0, keys, children, nil, ID, 0, is_leaf}
	return cell
}

func NewSet_BTree(degree int) *Set_BTree {
	tree := &Set_BTree{degree, nil, nil, nil}
	return tree
}

func Insert_BTree(s *Set_BTree, v int) {
	// Step 0, if the tree is empty, initialize it
	if s.root == nil {
		new_node := NewCell_BTree(s.degree, true, 0)
		new_node.keys[0] = v
		s.root = new_node
		return
	}

	// Step 1) Find the leaf node where the value should be inserted
	var Find_Leaf func(c *Cell_BTree, v int) *Cell_BTree
	Find_Leaf = func(c *Cell_BTree, v int) *Cell_BTree {
		// if at leaf, return
		if c.is_leaf {
			return c
		}

		// find first key that is smaller than or equal to v
		search_idx := 0
		for search_idx < c.cur_size && v > c.keys[search_idx] {
			search_idx += 1
		}

		// value is already in tree, return to that cell
		if c.keys[search_idx] == v {
			return c
		}

		// keep going deeper to find that leaf
		return Find_Leaf(c.children[search_idx], v)
	}

	// Step 2) Insert value into leaf node keys
	insert_node := Find_Leaf(s.root, v)

	// check if value is already in tree
	if ListContains(insert_node.keys, v) {
		return
	}

	insert_node.keys[insert_node.cur_size] = v
	insert_node.cur_size += 1

	// Step 3) Sort list
	PartialSort(insert_node.keys, insert_node.cur_size)

	// Step 4) Check if we are done
	if insert_node.cur_size != insert_node.max_size {
		return
	}

	var FixTreeUpwards func(s *Set_BTree, c *Cell_BTree)
	FixTreeUpwards = func(s *Set_BTree, c *Cell_BTree) {
		if c.cur_size != c.max_size {
			return
		}
		c.is_leaf = false

		// setup new left node
		M := c.cur_size
		left_node := NewCell_BTree(c.max_size, true, c.ID)
		for idx := 0; idx < (M/2)-1; idx++ {
			left_node.keys[idx] = c.keys[idx]
			left_node.cur_size += 1
		}

		// setup new right node
		right_node := NewCell_BTree(c.max_size, true, c.ID+1)
		for idx := (M / 2) + 1; idx < c.cur_size-1; idx++ {
			right_node.keys[idx] = c.keys[idx]
			right_node.cur_size += 1
		}

		new_mid_value := c.keys[M/2]

		// push new value to node above if root
		if c.parent == nil {
			new_root := NewCell_BTree(c.max_size, false, 0)
			new_root.keys[0] = new_mid_value
			new_root.cur_size += 1
			new_root.children[0] = left_node
			new_root.children[1] = right_node
			s.root = new_root
			return
		}

		// push new value to node above if not root
		// shift values in array to make space
		ShiftCellItems(c.parent, c.ID)

		// Adjust parent so that it connects to new left and right
		c.parent.keys[c.ID] = new_mid_value
		c.parent.cur_size += 1
		c.parent.children[c.ID] = left_node
		c.parent.children[c.ID+1] = right_node

		// Recursivly fix tree
		FixTreeUpwards(s, c.parent)


	}
	FixTreeUpwards(s, insert_node)

	// Step 5) adjust tree to free up space in leafs
	// Step 5.1) split into two nodes

}

func Search_BTree(s *Set_BTree, v int) *Value_Location {
	// if the tree is empty, return nil
	if s.root == nil {
		return nil
	}

	// setup a function to recursivly search via cell instead of whole tree
	var Search_Cell_BTree func(c *Cell_BTree, v int) *Value_Location
	Search_Cell_BTree = func(c *Cell_BTree, v int) *Value_Location {

		// find first key that is smaller than or equal to v
		search_idx := 0
		for search_idx < c.cur_size && v > c.keys[search_idx] {
			search_idx += 1
		}

		// if that key is equal to v, value has been found
		if c.keys[search_idx] == v {
			return &Value_Location{c, search_idx, v}
		}

		// if value is not found and current cell is a leaf, return nil
		if c.is_leaf {
			return nil
		}

		// if value is not found and current cell is not leaf, search child cell for value
		return Search_Cell_BTree(c.children[search_idx], v)
	}

	return Search_Cell_BTree(s.root, v)
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

func test_random() {
	node := NewCell_BTree(5, true, 0)
	a := NewCell_BTree(5, true, 0)
	b := NewCell_BTree(5, true, 1)
	c := NewCell_BTree(5, true, 2)
	d := NewCell_BTree(5, true, 3)
	e := NewCell_BTree(5, true, 4)

	fmt.Println(node.keys)
	node.keys[0] = 5

	node.children[0] = a

	node.keys[1] = 2
	node.children[1] = b
	node.keys[2] = -3
	node.children[2] = c
	node.children[3] = d

	node.cur_size = 3
	node.keys[3] = 10
	node.children[4] = e

	node.cur_size = 4

	PartialSort(node.keys, 2)
	fmt.Println(node.keys)

	fmt.Println(node.children)

	ShiftCellItems(node, 3)
	fmt.Println(node.keys)
	fmt.Println(node.children)
}

func main() {
	test_BTree()
	//test_random()

}

// ************** HELPER CODE **********************************

func PartialSort(arr []int, partitian int) {
	var msort func(arr []int, i int, j int, temp []int)
	var merge func(arr []int, i int, k int, j int, temp []int)

	// merge sort definition
	msort = func(arr []int, i int, j int, temp []int) {
		if i == j {
			return
		}
		half_point := (j-i)/2 + i
		msort(arr, i, half_point, temp)
		msort(arr, half_point+1, j, temp)
		merge(arr, i, half_point, j, temp)
	}

	// merge definition
	merge = func(arr []int, i int, k int, j int, temp []int) {
		first_end := k
		second_end := j
		var largest int
		for idx := j; idx >= i; idx-- {
			if first_end == i-1 {
				largest = arr[second_end]
				second_end -= 1
			} else if second_end == k {
				largest = arr[first_end]
				first_end -= 1
			} else if arr[second_end] > arr[first_end] {
				largest = arr[second_end]
				second_end -= 1
			} else {
				largest = arr[first_end]
				first_end -= 1
			}
			temp[idx] = largest
		}
		for idx2 := j; idx2 >= i; idx2-- {
			arr[idx2] = temp[idx2]
		}
	}

	// actuall merge sort here
	msort(arr, 0, partitian, make([]int, partitian+1))
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

func ListContains(list []int, v int) bool {
	for _, val := range list {
		if val == v {
			return true
		}
	}
	return false
}

func ShiftCellItems(c *Cell_BTree, free_idx int) {
	for idx := c.cur_size; idx > free_idx; idx-- {
		c.keys[idx] = c.keys[idx-1]
		c.children[idx+1] = c.children[idx]
		c.children[idx+1].ID += 1
	}
	c.keys[free_idx] = 0
	c.children[free_idx] = nil
	c.children[free_idx+1] = nil

}