package main

import (
	"fmt"
)

type Cell_BTree struct{
	max_size int
	cur_size int
	keys []int
	children []*Cell_BTree
	parent *Cell_BTree
	is_leaf bool
}

type Set_BTree struct{
	degree int
	root *Cell_BTree
	maximum *int
	minimum *int
}

type Value_Location struct{
	cell *Cell_BTree
	key_idx int
}

func NewCell_BTree(size int, is_leaf bool) *Cell_BTree{
	keys := make([]int, size-1)
	children := make([]*Cell_BTree, size)
	cell := &Cell_BTree{size, 0, keys, children, nil, is_leaf}
	return cell
}

func NewSet_BTree(degree int) *Set_BTree{
	tree := &Set_BTree{degree, nil, nil, nil}
	return tree
}

func Insert_BTree(s *Set_BTree, v int){
	fmt.Println("Not implemented yet")
} 

func Search_BTree(s *Set_BTree, v int) *Value_Location{
	fmt.Println("Not implemented yet")
	return nil
}

func Delete_BTree(s *Set_BTree, v int){
	fmt.Println("Not implemented yet")
}

func Minimum_BTree(s *Set_BTree) *Cell_BTree {
	fmt.Println("Not implemented yet")
	return nil
}



func main(){
	fmt.Println("B-Tree")
}