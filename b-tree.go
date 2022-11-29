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
	if ListContains(insert_node.keys, v) {
		return
	}

	insert_node.keys[insert_node.cur_size] = v
	insert_node.cur_size += 1

	// Step 3) Sort list
	PartialSort(insert_node.keys, insert_node.cur_size)

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
		left_node := NewCell_BTree(s.degree, c.ID)
		for idx := 0; idx <= (M/2)-1; idx++ {
			left_node.keys[idx] = c.keys[idx]
			left_node.cur_size += 1
		}

		// Step 7) setup new right node
		right_node := NewCell_BTree(s.degree, c.ID+1)
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

		// Step 13) Recursivly fix tree
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

	fmt.Printf("removing %d\n", v)
	// step 3) determine type of location and remove appropriatly
	if location.cell.IsLeaf() {
		s.DeleteFromLeaf(location)
	} else {
		s.DeleteFromLNonLeaf(location)
	}
}

func (s *Set_BTree) min() *Value_Location {
	if s.root == nil{
		return nil
	}

	current_node := s.root
	for current_node.children[0] != nil {
		current_node = current_node.children[0]
	}
	return &Value_Location{current_node, 0, current_node.keys[0]}
}

func (s *Set_BTree) max() *Value_Location {
	if s.root == nil{
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
	size := 4
	// First, insert all of 1, 2, 3, ..., 100 in order into an empty B-tree.
	fmt.Println("***** INSERTING RIGHT **************************************************")
	s := NewSet_BTree(size)
	for i := 1; i < 30; i += 1 {
		fmt.Printf("Inserting %d\n", i)
		s.insert(i)
		s.print()
	}

	// Similarly, insert all of 100, 99, 98, ..., 97 in order into an empty B-tree.
	fmt.Println("***** INSERTING LEFT **************************************************")
	s = NewSet_BTree(size)
	for i := 99; i > 0; i -= 1 {
		fmt.Printf("Inserting %d\n", i)
		s.insert(i)
		s.print()
	}

	// Delete the inserted values from the preceding tree, in order from 210, 220, ..., 350.
	fmt.Println("***** DELETING **************************************************")
	for i := 99; i > 0; i -= 1 {
		fmt.Printf("Deleting %d\n", i)
		s.delete(i)
		s.print()
	}
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

	min := s.min()
	max := s.max()
	minStr := "min = -"
	maxStr := "max = -"
	if min != nil {
		minStr = fmt.Sprintf("min = %d", min.value)
	}
	if max != nil {
		maxStr = fmt.Sprintf("max = %d", max.value)
	}

	fmt.Printf("%s  %s\n\n", minStr, maxStr)
}

func ListContains(list []int, v int) bool {
	for _, val := range list {
		if val == v {
			return true
		}
	}
	return false
}

func (c *Cell_BTree) ShiftCellItems(free_idx int) {
	for idx := c.cur_size; idx > free_idx; idx-- {
		c.keys[idx] = c.keys[idx-1]
		if c.children[idx] == nil {
			continue
		}
		c.children[idx+1] = c.children[idx]
		c.children[idx+1].ID += 1
	}

	c.keys[free_idx] = 0
	c.children[free_idx] = nil
	c.children[free_idx+1] = nil
}

func (c *Cell_BTree) UnShiftCellItems(free_idx int) {
	for idx := free_idx; idx < c.cur_size; idx++ {
		c.keys[idx] = c.keys[idx+1]
		if c.children[idx] == nil {
			continue
		}
		c.children[idx+1] = c.children[idx+2]
		c.children[idx+1].ID -= 1
	}

	c.keys[c.cur_size] = 0
	c.children[free_idx+1] = nil
}

func (c *Cell_BTree) IsLeaf() bool {
	return c.children[0] == nil
}

func (s *Set_BTree) DeleteFromLeaf(loc *Value_Location) {
	// current cell
	c := loc.cell

	// shift everting donw
	for i := loc.key_idx + 1; i < s.degree; i++ {
		c.keys[i-1] = c.keys[i]
	}
	// reduce size of cell
	c.cur_size -= 1

	// find min value of leaf
	min_size := int(math.Ceil(float64(s.degree)/2.0)) - 1

	if c.cur_size >= min_size {
		return
	}
	s.FixTreeUpwards(c)
}

func (s *Set_BTree) FixTreeUpwards(c *Cell_BTree) {

	if c.parent == nil {
		s.FixRoot(c)
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

func (s *Set_BTree) DeleteFromLNonLeaf(loc *Value_Location) {
	// find value to replace current one with

	// get value just smaller
	c := loc.cell
	del_idx := loc.key_idx

	child := c.children[del_idx]

	for !child.IsLeaf(){
		child = child.children[child.cur_size]
	}

	swap_value := child.keys[child.cur_size-1]

	c.keys[del_idx] = swap_value

	child.keys[child.cur_size-1] = 0
	child.cur_size -= 1

	min_size := int(math.Ceil(float64(s.degree)/2.0)) - 1

	if child.cur_size < min_size{
		s.FixTreeUpwards(child)
	}

}

func (c *Cell_BTree) BorrowFromLeft() {
	// get cell to left
	left_cell := c.parent.children[c.ID-1]

	// remove last value from left
	pred := left_cell.keys[left_cell.cur_size-1]
	left_cell.keys[left_cell.cur_size-1] = 0
	left_cell.cur_size -= 1

	// free up space in current
	
	c.children[c.cur_size+1] = c.children[c.cur_size]
	for idx := c.cur_size; idx > 0; idx-- {
		c.keys[idx] = c.keys[idx-1]
		c.children[idx] = c.children[idx-1]
	}
	c.cur_size += 1

	// move children if nescessary
	if !left_cell.IsLeaf() {
		c.children[0] = left_cell.children[left_cell.cur_size+1]
		left_cell.children[left_cell.cur_size+1] = nil
	}

	// move parent into first spot
	c.keys[0] = c.parent.keys[c.ID-1]

	// move pred into parent
	c.parent.keys[c.ID-1] = pred
}

func (c *Cell_BTree) BorrowFromRight() {
	// get cell to left
	right_cell := c.parent.children[c.ID+1]

	// remove first value from right
	next := right_cell.keys[0]

	// move parent into first spot
	c.keys[c.cur_size] = c.parent.keys[c.ID]
	c.cur_size += 1
	if !right_cell.IsLeaf(){
		c.children[c.cur_size] = right_cell.children[0]
		right_cell.children[0].parent = c
	}

	// shfit values
	for i := 1; i <= right_cell.cur_size; i++ {
		right_cell.keys[i-1] = right_cell.keys[i]
		if !right_cell.IsLeaf() {
			right_cell.children[i-1] = right_cell.children[i]
			right_cell.children[i-1].ID -= 1
		}
	}

	right_cell.cur_size -= 1

	// move pred into parent
	c.parent.keys[c.ID] = next
}

func (c *Cell_BTree) MergeWithLeft() {
	// get prev cell
	left_cell := c.parent.children[c.ID-1]

	// add parent to left
	left_cell.keys[left_cell.cur_size] = c.parent.keys[c.ID-1]
	left_cell.cur_size += 1

	// fmt.Println(left_cell.keys)
	// fmt.Println(left_cell.children)

	// add keys from current cell to left
	// add curent cell to left

	if !c.IsLeaf() {
		left_cell.children[left_cell.cur_size] = c.children[0]
		c.children[0].parent = left_cell
		c.children[0].ID = left_cell.cur_size
	}

	for idx := 0; idx < c.cur_size; idx++ {
		left_cell.keys[left_cell.cur_size] = c.keys[idx]
		left_cell.children[left_cell.cur_size+1] = c.children[idx+1]
		left_cell.cur_size += 1
	}

	// shift parent keys and children down one
	for i := c.ID; i <= c.parent.cur_size; i++ {
		c.parent.keys[i-1] = c.parent.keys[i]
		c.parent.children[i] = c.parent.children[i+1]
		c.parent.children[i-1].ID = i - 1
	}

	//c.cur_size = c.parent.cur_size-1

	c.parent.cur_size -= 1

}

func (c *Cell_BTree) MergeWithRight() {

	// get next cell
	right_cell := c.parent.children[c.ID+1]

	// clear up space for values
	right_cell.children[right_cell.cur_size+1] = right_cell.children[right_cell.cur_size]
	for idx := right_cell.cur_size; idx > 0; idx-- {
		right_cell.keys[idx] = right_cell.keys[idx-1]
		right_cell.children[idx] = right_cell.children[idx-1]
	}

	// add curent cell to left

	if !c.IsLeaf() {
		right_cell.children[0] = c.children[0]
		c.children[0].parent = right_cell
		c.children[0].ID = 0
	}

	for idx := 0; idx < c.cur_size; idx++ {
		right_cell.keys[idx] = c.keys[idx]
		right_cell.children[idx+1] = c.children[idx+1]
		right_cell.cur_size += 1
	}

	// add parent value
	right_cell.keys[c.cur_size] = c.parent.keys[c.ID]
	right_cell.cur_size += 1

	//fmt.Println(right_cell.children[0])

	// shift parent keys and children down one
	for i := c.ID + 1; i <= c.parent.cur_size; i++ {
		c.parent.keys[i-1] = c.parent.keys[i]
		c.parent.children[i-1] = c.parent.children[i]
		c.parent.children[i-1].ID = i - 1
	}
	//c.parent.children[c.parent.cur_size] = nil

	c.parent.cur_size -= 1
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
