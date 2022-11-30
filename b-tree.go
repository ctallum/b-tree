package main

import (
	"fmt"
	"math"
)

type Cell_BTree struct {
	ID       int
	cur_size int
	keys     []int
	children []*Cell_BTree
	parent   *Cell_BTree
}

type Set_BTree struct {
	degree int
	root   *Cell_BTree
}

type Value_Location struct {
	cell    *Cell_BTree
	key_idx int
	value   int
}

func NewCell_BTree(size int, ID int) *Cell_BTree {
	keys := make([]int, size)
	children := make([]*Cell_BTree, size+1)
	cell := &Cell_BTree{ID, 0, keys, children, nil}
	return cell
}

func NewSet_BTree(degree int) *Set_BTree {
	tree := &Set_BTree{degree, nil}
	return tree
}

func (s *Set_BTree) insert(v int) {

	// Step 0, if the tree is empty, initialize it
	if s.root == nil {
		new_node := NewCell_BTree(s.degree, 0)
		new_node.keys[0] = v
		new_node.cur_size += 1
		s.root = new_node
		return
	}

	// Step 1) Find the leaf node where the value should be inserted
	var Find_Leaf func(c *Cell_BTree, v int) *Cell_BTree
	Find_Leaf = func(c *Cell_BTree, v int) *Cell_BTree {
		// if at leaf, return
		if c.IsLeaf() {
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
	if insert_node.Contains(v) {
		return
	}
	
	

	insert_node.keys[insert_node.cur_size] = v
	insert_node.cur_size += 1

	
	// Step 3) Sort list
	PartialSort(insert_node.keys, insert_node.cur_size)
	fmt.Println("HI")
	s.print()
	fmt.Println("HI")
	// Step 4) Check if we are done
	if insert_node.cur_size != s.degree {
		return
	}

	
	var FixTreeUpwards func(s *Set_BTree, c *Cell_BTree)
	FixTreeUpwards = func(s *Set_BTree, c *Cell_BTree) {
		// Step 5) if tree is node is fine, return
		if c.cur_size != s.degree {
			return
		}
		// Step 6) setup new left node
		M := c.cur_size - 1
		// fmt.Println(0,M/2 -1, M/2, M/2 + 1, c.cur_size )
		left_node := NewCell_BTree(s.degree, c.ID)
		for idx := 0; idx < (M/2); idx++ {
			left_node.keys[idx] = c.keys[idx]
			left_node.cur_size += 1
		}

		// Step 7) setup new right node
		right_node := NewCell_BTree(s.degree, c.ID+1)
		for idx := (M / 2) + 1; idx < c.cur_size; idx++ {
			right_node.keys[idx-(M/2)-1] = c.keys[idx]
			right_node.cur_size += 1
		}

		// Step 8) find mid value
		new_mid_value := c.keys[M/2]

		// Step 9) if not leaf, fix all child linking
		if !c.IsLeaf() {
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
			new_root := NewCell_BTree(s.degree, 0)
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
		c.parent.ShiftCellItems(c.ID)



		// Step 12) push new value to node above if not root
		// Adjust parent so that it connects to new left and right
		c.parent.keys[c.ID] = new_mid_value
		c.parent.cur_size += 1

		left_node.parent = c.parent
		c.parent.children[c.ID] = left_node

		right_node.parent = c.parent
		c.parent.children[c.ID+1] = right_node

		// Step 13) Recursivly fix tree)
		FixTreeUpwards(s, c.parent)

	}
	FixTreeUpwards(s, insert_node)
}

func (s *Set_BTree) search(v int) *Value_Location {
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
		if c.IsLeaf() {
			return nil
		}

		// if value is not found and current cell is not leaf, search child cell for value
		return Search_Cell_BTree(c.children[search_idx], v)
	}

	return Search_Cell_BTree(s.root, v)
}

func (s *Set_BTree) delete(v int) {

	// step 1) find location of the value
	location := s.search(v)

	// step 2) check to see if value is found
	if location == nil {
		fmt.Printf("can't find %d\n", v)
		return
	}

	// fmt.Printf("removing %d\n", v)
	// step 3) determine type of location and remove appropriatly
	if location.cell.IsLeaf() {
		s.DeleteFromLeaf(location)
	} else {
		s.DeleteFromLNonLeaf(location)
	}
}

func (s *Set_BTree) min() *Value_Location {
	if s.root == nil {
		return nil
	}

	current_node := s.root
	for current_node.children[0] != nil {
		current_node = current_node.children[0]
	}
	return &Value_Location{current_node, 0, current_node.keys[0]}
}

func (s *Set_BTree) max() *Value_Location {
	if s.root == nil {
		return nil
	}

	current_node := s.root
	for current_node.children[0] != nil {
		current_node = current_node.children[current_node.cur_size]
	}
	return &Value_Location{current_node, 0, current_node.keys[current_node.cur_size-1]}
}

func test_BTree() {
	// B-Tree degree
	size := 7
	// First, insert all of 1, 2, 3, ..., 100 in order into an empty B-tree.
	// fmt.Println("***** INSERTING RIGHT **************************************************")
	s := NewSet_BTree(size)
	for i := -7; i <= 0; i += 1 {
		fmt.Printf("Inserting %d\n", i)
		s.insert(i)
		s.print()
	}
	s.print()


	// // Similarly, insert all of 100, 99, 98, ..., 97 in order into an empty B-tree.
	// fmt.Println("***** INSERTING LEFT **************************************************")
	// s = NewSet_BTree(size)
	// for i := 100; i >= 0; i -= 1 {
	// 	fmt.Printf("Inserting %d\n", i)
	// 	s.insert(i)
	// 	s.print()
	// }

	// Delete the inserted values from the preceding tree, in order from 210, 220, ..., 350.
	// fmt.Println("***** DELETING **************************************************")
	// for i := 100; i > 0; i -= 1 {
	// 	fmt.Printf("Deleting %d\n", i)
	// 	s.delete(i)
	// 	s.print()
	// }
}

func main() {
	test_BTree()
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

func (s *Set_BTree) print() {
	ascii_tree := ""
	var print func(*Cell_BTree, string, bool, int)
	print = func(c *Cell_BTree, prefix string, signal bool, layer int) {
		if c != nil {
			to_right := (c.cur_size + 1) / 2
			to_left := (c.cur_size) - to_right + 1
			// print right half of children
			if c.children[0] != nil {
				for i := 0; i < to_right; i++ {
					print(c.children[c.cur_size-i], prefix+"  ", true, layer+1)
				}
			}
			// print current self
			//fmt.Printf("%s%d : %d : %d\n", prefix, c.keys, c.cur_size, c.ID)
			if signal && c.parent != nil {
				ascii_tree += fmt.Sprintf("%s%d\n", prefix+"/", c.keys[0:c.cur_size])
			} else if !signal && c.parent != nil {
				ascii_tree += fmt.Sprintf("%s%d\n", prefix+"\\", c.keys[0:c.cur_size])
			} else {
				ascii_tree += fmt.Sprintf("%s%d\n", prefix+"-", c.keys[0:c.cur_size])
			}
			//ascii_tree += fmt.Sprintf("%s%d\n", prefix, c.keys[0:c.cur_size])
			// print left hafl of children
			if c.children[0] != nil {
				for i := 0; i < to_left; i++ {
					print(c.children[c.cur_size-to_right-i], prefix+"  ", false, layer+1)
				}
			}

		} else {
			ascii_tree += fmt.Sprintf("%s-\n", prefix)
		}
	}

	print(s.root, " ", true, 0)

	fmt.Print(ascii_tree + "\n")

	// min := s.min()
	// max := s.max()
	// minStr := "min = -"
	// maxStr := "max = -"
	// if min != nil {
	// 	minStr = fmt.Sprintf("min = %d", min.value)
	// }
	// if max != nil {
	// 	maxStr = fmt.Sprintf("max = %d", max.value)
	// }

	// fmt.Printf("%s  %s\n\n", minStr, maxStr)
}

func (c *Cell_BTree) Contains(v int) bool {
	for _, val := range c.keys[0:c.cur_size] {
		if val == v {
			return true
		}
	}
	return false
}

func (c *Cell_BTree) ShiftCellItems(free_idx int) {
	// shift all keys and children down so that we have free spot at free_idx
	for idx := c.cur_size; idx > free_idx; idx-- {
		c.keys[idx] = c.keys[idx-1]
		if c.children[idx] == nil {
			continue
		}
		c.children[idx+1] = c.children[idx]
		c.children[idx+1].ID += 1
	}

	// clear out the key and the sourounding children the border it
	c.keys[free_idx] = 0
	c.children[free_idx] = nil
	c.children[free_idx+1] = nil
}

func (c *Cell_BTree) IsLeaf() bool {
	return c.children[0] == nil
}

func (s *Set_BTree) DeleteFromLeaf(loc *Value_Location) {
	// current cell
	c := loc.cell

	// shift everting down
	for i := loc.key_idx + 1; i < s.degree; i++ {
		c.keys[i-1] = c.keys[i]
	}
	// reduce size of cell
	c.cur_size -= 1

	// find min value of leaf
	min_size := int(math.Ceil(float64(s.degree)/2.0)) - 1

	// if current cell is now too small, fix tree
	if c.cur_size < min_size {
		s.FixTreeUpwards(c)
	}
}

func (s *Set_BTree) DeleteFromLNonLeaf(loc *Value_Location) {
	// get current cell
	c := loc.cell

	// get index of deletion
	idx := loc.key_idx

	// get child for values smaller than deleted key
	child := c.children[idx]

	// find the largest value in subtree of values less than the key
	for !child.IsLeaf() {
		child = child.children[child.cur_size]
	}

	// get the key that is just smaller than deleted key
	swap_value := child.keys[child.cur_size-1]

	// swap the deleted key with value just smaller than it
	c.keys[idx] = swap_value

	// remove the key from the leaf
	child.keys[child.cur_size-1] = 0
	child.cur_size -= 1

	// check if we need to fix the leaf if it is too small
	min_size := int(math.Ceil(float64(s.degree)/2.0)) - 1
	if child.cur_size < min_size {
		s.FixTreeUpwards(child)
	}
}

func (s *Set_BTree) FixTreeUpwards(c *Cell_BTree) {
	if c.parent == nil {
		if c.cur_size == 0 {
			s.FixRoot(c)
		}
		return
	}

	min_size := int(math.Ceil(float64(s.degree)/2.0)) - 1

	// see if we can manipulate cell to direct left
	var pred_cell *Cell_BTree = nil
	var next_cell *Cell_BTree = nil
	if c.ID != 0 {
		pred_cell = c.parent.children[c.ID-1]
	}
	if c.ID != c.parent.cur_size {
		next_cell = c.parent.children[c.ID+1]
	}

	if pred_cell != nil {
		if pred_cell.cur_size > min_size {
			// fmt.Println("borrow left")
			c.BorrowFromLeft()
			if c.cur_size < min_size {
				s.FixTreeUpwards(c)
			}
			return
		}
	}
	if next_cell != nil {
		if next_cell.cur_size > min_size {
			// fmt.Println("borrow right")
			c.BorrowFromRight()
			if c.cur_size < min_size {
				s.FixTreeUpwards(c)
			}
			return
		}
	}
	if pred_cell != nil {
		// fmt.Println("merge left")
		c.MergeWithLeft()
		// s.print()
		if c.parent.cur_size < min_size {
			s.FixTreeUpwards(c.parent)
		}
		return
	}
	if next_cell != nil {
		// fmt.Println("merge right")
		c.MergeWithRight()
		if c.parent.cur_size < min_size {
			s.FixTreeUpwards(c.parent)
		}
		return
	}
}

func (c *Cell_BTree) BorrowFromLeft() {
	// get cell to left
	left_cell := c.parent.children[c.ID-1]

	// remove last value from left
	last_val := left_cell.keys[left_cell.cur_size-1]
	left_cell.keys[left_cell.cur_size-1] = 0

	// remove last child from left
	last_child := left_cell.children[left_cell.cur_size]
	left_cell.children[left_cell.cur_size] = nil

	// decrease size of left cell by 1
	left_cell.cur_size -= 1

	// free up keys
	for idx := c.cur_size; idx > 0; idx-- {
		c.keys[idx] = c.keys[idx-1]
	}

	// free up children if not leaf
	if !c.IsLeaf() {
		for idx := c.cur_size + 1; idx > 0; idx-- {
			c.children[idx] = c.children[idx-1]
			c.children[idx].ID += 1
		}
	}

	// increase size of c by 1
	c.cur_size += 1

	// move child of left_cell if it's not nil
	if last_child != nil {
		c.children[0] = last_child
		last_child.parent = c
		last_child.ID = 0
	}

	// move parent into first spot
	c.keys[0] = c.parent.keys[c.ID-1]

	// move pred into parent
	c.parent.keys[c.ID-1] = last_val
}

func (c *Cell_BTree) BorrowFromRight() {
	// get cell to left
	right_cell := c.parent.children[c.ID+1]

	// remove first value from right
	right_val := right_cell.keys[0]
	right_cell.keys[0] = 0

	// remove first key from right
	right_child := right_cell.children[0]
	right_cell.children[0] = nil

	// shift all right_cell keys down by 1
	for i := 0; i < right_cell.cur_size; i++ {
		right_cell.keys[i] = right_cell.keys[i+1]
	}

	// shift all right_cell children down by 1
	for i := 0; i < right_cell.cur_size+1; i++ {
		right_cell.children[i] = right_cell.children[i+1]
	}

	// decrease right_cell size by 1
	right_cell.cur_size -= 1

	// move parent into first spot
	c.keys[c.cur_size] = c.parent.keys[c.ID]

	// move right_cell child into current cell if not nill
	if right_child != nil {
		c.children[c.cur_size+1] = right_child
		right_child.parent = c
		right_child.ID = c.cur_size + 1
	}

	// increase size of c by 1
	c.cur_size += 1

	// move pred into parent
	c.parent.keys[c.ID] = right_val
}

func (c *Cell_BTree) MergeWithLeft() {
	// get cell directly to left
	left_cell := c.parent.children[c.ID-1]

	// remove key from parent
	parent_key := c.parent.keys[c.ID-1]

	// clean up parent by shifting keys and children
	for i := c.ID - 1; i < c.parent.cur_size; i++ {
		c.parent.keys[i] = c.parent.keys[i+1]
	}

	// clean up parent by shifting children down
	for i := c.ID; i < c.cur_size+1; i++ {
		c.parent.children[i] = c.parent.children[i+1]
		c.parent.children[i].ID = i
	}

	// reduce size of parent by 1
	c.parent.cur_size -= 1

	// add parent key to left_cell
	left_cell.keys[left_cell.cur_size] = parent_key
	left_cell.cur_size += 1

	// store original size of left
	init_len := left_cell.cur_size

	// add keys from current cell to left_cell
	for idx := 0; idx < c.cur_size; idx++ {
		left_cell.keys[left_cell.cur_size] = c.keys[idx]
		left_cell.cur_size += 1
	}

	// add children from current cel to left_cell if not leaf
	if !c.IsLeaf() {
		for idx := 0; idx < c.cur_size+1; idx++ {
			left_cell.children[init_len+idx] = c.children[idx]
			c.children[idx].parent = left_cell
			c.children[idx].ID = init_len + idx
		}
	}
}

func (c *Cell_BTree) MergeWithRight() {
	// get cell directly to right
	right_cell := c.parent.children[c.ID+1]
	right_cell.MergeWithLeft()
}

func (s *Set_BTree) FixRoot(c *Cell_BTree) {
	// fmt.Println("Fixing root")
	if c.IsLeaf() {
		s.root = nil
		return
	}
	c.children[0].parent = nil
	s.root = c.children[0]
}

