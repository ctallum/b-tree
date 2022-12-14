# B-tree
Golang implementation of a B-tree

## About B-tree
A B-tree data structure is a self-ballencing tree data structure. Like a binary search tree, values stored in a b-tree are ordered. This allows for insertion, deleteion, and search operations to happen in logorithmic time. However, unlike a binary search tree, each node in a b-tree may hold more than one value and two children.

### Structure



### Use Cases



## Implementation
### Structure
A b-tree is represented by the following struct:

```go
type Set_BTree struct {
	degree int
	root   *Cell_BTree
}
```
The `degree` of a tree is equal to the *m* order of the tree as seen in the definition of a b-tree. The degree is then the maximum number of children any given cell of the tree can hold. The `root` of the tree is the adress of the root cell. 


Each cell in the b-tree is represented by the following struct:
```go
type Cell_BTree struct {
	ID       int
	cur_size int
	keys     []int
	children []*Cell_BTree
	parent   *Cell_BTree
}
```
`ID` refers to which child the current cell is with respect to its parent. If the current cell is the first child, it will have an `ID` of 0. If it is the *(m)* child, it will have an `ID` of *m-1*. The `cur_size` or current size of the cell is how many keys the cell is holding. `cur_size` can range from anywhere from *1* to *m-1*. `keys` is an array of integers that hold all of the keys held in the current cell. The size of the `keys` array is constant, extra spaces that aren't being used are filled by zeros, and the `cur_size` reflects the number of keys that the array is actually holding. `children` is an array of pointers to other cells. The size of `children` is also constant, extra spaces that aren't being used are filled with `nill`. At any given time there are `cur_size + 1` children in the array. `parent` refers to the parent of the current cell. If the current cell does not have a parent (if the cell is the root of the tree), parent is set to `nill`.

### Insert

### Search

### Delete

### Min

### Max


## Sources

