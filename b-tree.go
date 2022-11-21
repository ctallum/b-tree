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
		new_node.cur_size += 1
		s.root = new_node
		return
	}

	// Step 1) Find the leaf node where the value should be inserted
	var Find_Leaf func(c *Cell_BTree, v int) *Cell_BTree
	Find_Leaf = func(c *Cell_BTree, v int) *Cell_BTree {
		// if at leaf, return
		if IsLeaf(c) {
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
		// Step 5) if tree is node is fine, return
		if c.cur_size != c.max_size {
			return
		}

		// Step 6) setup new left node
		M := c.cur_size - 1
		left_node := NewCell_BTree(c.max_size, true, c.ID)
		for idx := 0; idx <= (M/2)-1; idx++ {
			left_node.keys[idx] = c.keys[idx]
			left_node.cur_size += 1
		}

		// Step 7) setup new right node
		right_node := NewCell_BTree(c.max_size, true, c.ID+1)
		for idx := (M / 2) + 1; idx <= c.cur_size-1; idx++ {
			right_node.keys[idx-(M/2)-1] = c.keys[idx]
			right_node.cur_size += 1
		}

		// Step 8) find mid value
		new_mid_value := c.keys[M/2]

		// Step 9) if not leaf, fix all child linking
		if c.children[0] != nil {
			// fix left node
			for idx := 0; idx < left_node.cur_size+1; idx++ {
				left_node.children[idx] = c.children[idx]
				left_node.children[idx].parent = left_node
				left_node.children[idx].ID = idx

			}

			// fix right node
			for idx := 0; idx < right_node.cur_size+1; idx++ {
				right_node.children[idx] = c.children[idx+left_node.cur_size+1]
				right_node.children[idx].parent = right_node
				right_node.children[idx].ID = idx
			}
		}

		// Step 10) push new value to node above if root
		if c.parent == nil {
			new_root := NewCell_BTree(c.max_size, false, 0)
			new_root.keys[0] = new_mid_value
			new_root.cur_size += 1
			new_root.children[1] = right_node
			right_node.parent = new_root
			new_root.children[0] = left_node
			left_node.parent = new_root
			s.root = new_root
			return
		}

		// Step 11) shift values in array to make space
		ShiftCellItems(c.parent, c.ID)

		// Step 12) push new value to node above if not root
		// Adjust parent so that it connects to new left and right
		c.parent.keys[c.ID] = new_mid_value
		c.parent.cur_size += 1

		left_node.parent = c.parent
		c.parent.children[c.ID] = left_node

		right_node.parent = c.parent
		c.parent.children[c.ID+1] = right_node

		// Step 13) Recursivly fix tree
		FixTreeUpwards(s, c.parent)

	}
	FixTreeUpwards(s, insert_node)
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
		if IsLeaf(c) {
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
	size := 4
	// First, insert all of 1, 2, 3, ..., 100 in order into an empty B-tree.
	fmt.Println("***** INSERTING RIGHT **************************************************")
	s := NewSet_BTree(size)
	for i := 1; i < 100; i += 1 {
		fmt.Printf("Inserting %d\n", i)
		Insert_BTree(s, i)
		PrintSet_BTree(s)
	}

	// // Similarly, insert all of 350, 340, 330, ..., 210 in order into an empty B-tree.
	// fmt.Println("***** INSERTING LEFT **************************************************")
	// t := NewSet_BTree(size)
	// for i := 100; i > 0; i -= 1 {
	// 	fmt.Printf("Inserting %d\n", i)
	// 	Insert_BTree(t, i)
	// 	PrintSet_BTree(t)
	// }
	// // Delete the inserted values from the preceding tree, in order from 210, 220, ..., 350.
	// fmt.Println("***** DELETING **************************************************")
	// for i := 1; i < 16; i += 1 {
	// 	c := Search_BTree(t, i*10+200)
	// 	fmt.Printf("Deleting %d\n", i*10+200)
	// 	Delete_BTree(t, c)
	// 	PrintSet_BTree(t)
	// }
}

func main() {
	test_BTree()


	fmt.Println()
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
	msort(arr, 0, partitian-1, make([]int, partitian))
}

func PrintSet_BTree(s *Set_BTree) {
	var print func(*Cell_BTree, string)
	print = func(c *Cell_BTree, prefix string) {
		if c != nil {
			to_right := (c.cur_size + 1) / 2
			to_left := (c.cur_size) - to_right + 1
			// print right half of children
			if c.children[0] != nil {
				for i := 0; i < to_right; i++ {
					print(c.children[i], prefix+"   ")
				}
			}
			// print current self
			//fmt.Printf("%s%d : %d : %d\n", prefix, c.keys, c.cur_size, c.ID)
			fmt.Printf("%s%d\n", prefix, c.keys)
			// print left hafl of children
			if c.children[0] != nil {
				for i := 0; i < to_left; i++ {
					print(c.children[to_right+i], prefix+"   ")
				}
			}

		} else {
			fmt.Printf("%s-\n", prefix)
		}
	}

	fmt.Println(" /")
	print(s.root, " |  ")
	fmt.Println(" \\")

	fmt.Println("")

	// min := Minimum_BTree(s)
	// max := Maximum_BTree(s)
	// minStr := "min = -"
	// maxStr := "max = -"
	// if min != nil {
	// 	minStr = fmt.Sprintf("min = %d", min.value)
	// }
	// if max != nil {
	// 	maxStr = fmt.Sprintf("max = %d", max.value)
	// }
	// fmt.Printf("  %s  %s\n", minStr, maxStr)
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
		if c.children[idx] == nil {
			continue
		}
		c.keys[idx] = c.keys[idx-1]

		c.children[idx+1] = c.children[idx]
		c.children[idx+1].ID += 1
	}

	c.keys[free_idx] = 0
	c.children[free_idx] = nil
	c.children[free_idx+1] = nil
}

func IsLeaf(c *Cell_BTree) bool {
	return c.children[0] == nil
}
