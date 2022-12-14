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
Inserting a key `v` into a b-tree works by the following proceedure. First, find a leaf node where `v` should be inserted into. All values inserted into a b-tree are always inserted into leaves. 

To find the leaf where the key should be inserted, we use the following function:
```go
func (c *Cell_BTree) FindLeaf(v int) *Cell_BTree {
	// if at leaf, return
	if c.IsLeaf() {
		return c
	}

	search_idx := 0
	for search_idx < c.cur_size && v > c.keys[search_idx] {
		search_idx += 1
	}

	// value is already in tree, return to that cell
	if search_idx != c.cur_size && c.keys[search_idx] == v {
		return c
	}

	// keep going deeper to find that leaf
	return c.children[search_idx].FindLeaf(v)
}
```
The function above works by recursively moving down the tree. If the current cell is a leaf, return the adress of the cell. If it is not a leaf, we find which child of the current cell is the root of a subtree that should hold `v`. 

Once we reach the insertion leaf, we add `v` into the leaf at the end then sort it using a quicksort function to make sure the keys in the leaf are still ordered. If the leaf node does not contain too many leaves, we finish the insertion process here. If the leaf contains `m` keys, we need to fix things. 

To fix a leaf that has too many keys, we split the current array of keys into three groups. The first group is the median number. The second group is the set of keys less than the median, and the third group is the set of keys that are greater than the median. We split the current cell into two cells. One cell holds group two, and the other holds group three. The median value is pushed upwards into the parent cell. 

We can see this operation in the following code snippet:
```go
// new left node
M := c.cur_size - 1
left_node := NewCell_BTree(s.degree, c.ID)
for idx := 0; idx < (M / 2); idx++ {
	left_node.keys[idx] = c.keys[idx]
	left_node.cur_size += 1
}

// new right node
right_node := NewCell_BTree(s.degree, c.ID+1)
for idx := (M / 2) + 1; idx < c.cur_size; idx++ {
	right_node.keys[idx-(M/2)-1] = c.keys[idx]
	right_node.cur_size += 1
}

// new middle value
new_mid_value := c.keys[M/2]
```
With these new nodes, we still have to make sure everything is properly linked. The new cells correctly have their `parent` property set to the original parent, and the parent has to add these two new children in the correct location within its `children` array.

Now that the parent cell has one more key than when it started, it is possible that the parent now also has too many keys. If this is the case, we again preform the same operaton that we did on the leaf to the parent. We split the parent into two new nodes and push the median value to its parent. We preform this recursive operaton untill all nodes in the tree have a valid number of keys. 

This whole operatoin is covered by our `FixTreeUpwardsInsert` function: 
```go
func (s *Set_BTree) FixTreeUpwardsInsert(c *Cell_BTree){
	// lots of code here ...


	/// lots of code here ...
}
```

### Search
The search function of our b-tree implementation is very simple. We recursivly search down the tree, heading in the correct direction of our search key `v` until we either reach a leaf or we find the value along the way. If we find `v` in a cell, we return the location of that cell. Otherwise, we search deeper into the tree. If we reach a leaf cell, we iterate through the leaf looking for `v`. If we do not find it, we return `nill` otherwise, we return the location of the leaf. This operation is done with the following recurrsive function:

```go
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

	// // if value is not found and current cell is a leaf, return nil
	if c.IsLeaf() {
		return nil
	}

	// // if value is not found and current cell is not leaf, search child cell for value
	return Search_Cell_BTree(c.children[search_idx], v)
}
```

Rather than returning the adress of the cell we find `v` in, we created a `Value_Location` struct that looks like the following:

``` go
type Value_Location struct {
	cell    *Cell_BTree
	key_idx int
	value   int
}
```

The `Value_Location` struct contains the cell that contains `v`, the index of the `keys` array where `v` is found, and the value of `v` itself. 


### Delete

### Min

### Max


## Sources
For insert, I used the following guide: [programiz.com](https://www.programiz.com/dsa/insertion-into-a-b-tree)

For delete, I used the follwing guides: 


