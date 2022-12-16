# B-tree
Golang implementation of a B-tree

## To Run
Running the main file creates a B-tree of degree 7 and populates it with the values 0-100, printing out the tree at each step. It then removes the values 0-100, also printing at each step. Lastly, it generates a b-tree filled with 10000 random values, prints it, then deltes each value. 
```console
go run b-tree.go
```


## About B-tree
A B-tree data structure is a self-balancing tree data structure. Like a binary search tree, values stored in a b-tree are ordered. This allows for insertion, deletion, and search operations to happen in logarithmic time. However, unlike a binary search tree, each node in a b-tree may hold more than one value and two children.

### Structure
A b-tree with degree *m* has the following properties:
- Every node has at most *m* children
- Every internatl node (not leaf or root) has at least ⌈m/2⌉ children and ⌈m/2⌉-1 keys
- Every non-leaf node has at least two children
- All leaves appear on the same level. The tree is ballenced 
- A non-leaf node with *k* keys contains *k+1* children

Each node has a set of keys. These keys act as seperation values which devide it's children which are roots of sub-trees. If a tree has two keys *a* and *b*, it will have three children. Child one will have all values less than *a*. Child two will have all the values between *a* and *b*. Child three will have all the values greater than *b*. In this manner, the tree is organized. 

### Time Complexity

| **Algorithm** | **Asymptotic Running Times** |
|---------------|------------------------------|
|Insert||
|Search||
|Delete||
|Maximum||
|Minimum||

### Use Cases
A B-tree is well suited for storage systems that need to access large amount of data, such as databases. This is epspecially true if accessing the data is a time intensive process, for example if the data in a b-tree were stored in external drives. Since B-trees store large chunks of data in each node, to get to any piece of data, one needs to traverse less nodes, which means less disk read operaitons. In memory, each node will have all of its data stored physically close to eachother. Furthermore, B-trees don't reballance as often as other trees, so the time intensive task of moving data around on a disk happens less often. 

Historically, B-trees have been used for databases where data is stores on disk drives. Data on disk drives is stores in disk blocks. The time needed for the disk to acess the correct disk block is far greater than the time to read the data once it is already in the correct disk block. Therefore, it makes sense to use a b-tree. The fewer times, the disk has re-seek a disk block, the quicker one can access the data. Typically, a b-tree will set the size of the node to be the size of a disk block.


## Implementation
### Structure
A b-tree is represented by the following struct:

```go
type Set_BTree struct {
	degree int
	root   *Cell_BTree
}
```
The `degree` of a tree is equal to the *m* order of the tree as seen in the definition of a b-tree. The degree is then the maximum number of children any given cell of the tree can hold. The `root` of the tree is the address of the root cell. 


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
Inserting a key `v` into a b-tree works by the following procedure. First, find a leaf node where `v` should be inserted into. All values inserted into a b-tree are always inserted into leaves. 

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
The function above works by recursively moving down the tree. If the current cell is a leaf, return the address of the cell. If it is not a leaf, we find which child of the current cell is the root of a sub-tree that should hold `v`. 

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

Now that the parent cell has one more key than when it started, it is possible that the parent now also has too many keys. If this is the case, we again preform the same operation that we did on the leaf to the parent. We split the parent into two new nodes and push the median value to its parent. We preform this recursive operation until all nodes in the tree have a valid number of keys. 

This whole operation is covered by our `FixTreeUpwardsInsert` function: 
```go
func (s *Set_BTree) FixTreeUpwardsInsert(c *Cell_BTree){
	// lots of code here ...


	/// lots of code here ...
}
```

### Search
The search function of our b-tree implementation is very simple. We recursively search down the tree, heading in the correct direction of our search key `v` until we either reach a leaf or we find the value along the way. If we find `v` in a cell, we return the location of that cell. Otherwise, we search deeper into the tree. If we reach a leaf cell, we iterate through the leaf looking for `v`. If we do not find it, we return `nill` otherwise, we return the location of the leaf. This operation is done with the following recursive function:

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

Rather than returning the address of the cell we find `v` in, we created a `Value_Location` struct that looks like the following:

``` go
type Value_Location struct {
	cell    *Cell_BTree
	key_idx int
	value   int
}
```

The `Value_Location` struct contains the cell that contains `v`, the index of the `keys` array where `v` is found, and the value of `v` itself. 


### Delete
In order to delete a value `v` from the b-tree, we first try and find it using our `search(v)` function. If we cannot find it, we return. If we do find it, there are two options: the value we are trying to delete is in a leaf or it is not. Depending on this, there are different steps we preform. 

```go
func (s *Set_BTree) delete(v int) {
	// find location of the value
	location := s.search(v)

	// check to see if value is found
	if location == nil {
		return
	}

	// determine type of location and remove appropriatly
	if location.cell.IsLeaf() {
		s.DeleteFromLeaf(location)
	} else {
		s.DeleteFromLNonLeaf(location)
	}
}
```
### Deleting from a leaf
First off, if we are trying to delete a value from a leaf, there are two scenarios. If the leaf has at least *⌈m/2⌉* keys, we can safely remove a the value `v` from the cell. The cell will still have a minimum of *⌈m/2⌉-1* keys, the minimum number for a b-tree. We are then done. The second scenario is if the tree has *⌈m/2⌉-1* or less keys. After we remove a key, it will have an insufficient number of keys, so we will need to fix it by expanding the size of the current cell. There are four ways to fix this issue. 

1) Borrow a key from the cell directly to the right
2) Borrow a key from the cell directly to the left
3) Merge with the cell directly to the right
4) Merge with the cell directly to the left

The first option if we have too few keys in a cell is to borrow from an adjacent cell. To get the cell to the right, we go to the parent of the cell, then move over to the next child cell. Conversely, to get to the cell directly to the left, we go to the parent of the cell, and then move over to the previous child. In order to borrow from an adjacent cell, the adjacent cell needs to have a minimum of ⌈m/2⌉ + 1 keys. Therefore, we can safely remove a key from it and move a key into our current cell. 

To borrow a key from the cell to the left, we use the following code:

```go
// if we are not at the beginning, try borrowing from the left
if pred_cell != nil {
   	if pred_cell.cur_size > min_size {
	   	c.BorrowFromLeft()
   		// we only borrow one value at a time, so see if we need to do operation again
   		if c.cur_size < min_size {
   			s.FixTreeUpwardsDelete(c)
   		}
   		return
   	}
}
```
The `BorrowFromLeft()` function does the following:
- Find the cell directly to the left
- Find and remove the last key in the left_cell
- Move the last key into the parent
- Move a value from the parent into the beginning of the current cell

If we cannot borrow from the left cell, we can alternatively try to borrow from the right cell as follows: 

```go
// if we are not at the end, try borrowing from the right
if next_cell != nil {
	if next_cell.cur_size > min_size {
  		c.BorrowFromRight()
  		// we only borrow one value at a time, so see if we need to do operation again
  		if c.cur_size < min_size {
  			s.FixTreeUpwardsDelete(c)
  		}
  	return
  	}
}
```
The `BorrorFromRight()` function does the following:
- Find the cell directly to the right
- Find and remove the first key in the right_cell
- Move the first key into the parent
- Move a value from the parent into the end of the current cell

If we cannot borrow from the left cell or the right cell, we can then try to merge with the left cell. This happens if neither the cell to the left nor right have sufficient number of keys such that they can safely lose one. Therefore, they have at most *⌈m/2⌉-1* keys. Given that the current cell has *⌈m/2⌉ - 2* keys, if we were to merge with another cell, and steal a key from the parent, the resulting cell would have at maximum *M-1* keys, the maximum number of keys. Therefore, if we cannot borrow, we can safely merge with adjacent cells. 

To merge with a cell to the left, we use the following code:
```go
// if we cannot borrow but are not at the beginning, merge left
if pred_cell != nil {
	c.MergeWithLeft()
		if c.parent.cur_size < min_size {
			s.FixTreeUpwardsDelete(c.parent)
		}
	return
}
```

The `MergeWithleft()` function does the following:
- Find the cell directly to the left
- Find and remove a key from the parent
- Add the parent key to the end of the left cell
- Move all the keys from the current cell to the left cell

If we cannot, for some reason, merge with the left cell, we can alternatively merge with the right cell. To do so, we use the following code:
```go
// if we cannot borrow but are not at the end, merge right
if next_cell != nil {
	c.MergeWithRight()
	if c.parent.cur_size < min_size {
		s.FixTreeUpwardsDelete(c.parent)
	}
	return
}
```

The `MergeWithRight()` function does the following:
- Find the cell directly to the right
- Call `MergeWithLeft` upon the right cell, effectively merging the two cells

After we do one of the four operations above, we have removed `v` from the leaf, and we are guaranteed to have a sufficient number of keys in the cell. However, we may have forced the parent cell to now have too few keys. To fix this issue, the parent cell may do any one of the four operations to gain an extra key. We do this process recursively until all the cells in the tree have enough keys or until we reach the root of the tree. If we reach the root of the tree, we shrink the height of the tree by one. 

### Deleting from a non-leaf
If we want to delete a value that is not in a leaf, we can do the following. 
- Delete the key from somewhere in the middle of the tree
- Copy a leaf value into the deletion node so that the node maintains the same number of keys
- Delete the leaf value using the methods mentioned above

Since we know the location of the node where we want to delete a key, we need to find a value to replace it. We can replace the key with the next smallest key in the whole tree which we know must be in a leaf. To find this key, we use the following code:

```go
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
```

In the code above, we first find the location of the key to delete. We then find the child sub-tree of all values smaller than the key. We then iterate down the tree until we reach a leaf. For each iteration, we go to the furthest right sub-tree, the sub-tree that holds the greatest values. Once we reach the end, we can take the last value of the leaf, and we know it must be the greatest key that is still smaller than the deletion key. 

After swapping the leaf key into the deletion cell, we need to delete the key from the child. We do that with the following code:
```go
// remove the key from the leaf
child.keys[child.cur_size-1] = 0
child.cur_size -= 1

// check if we need to fix the leaf if it is too small
min_size := int(math.Ceil(float64(s.degree)/2.0)) - 1
if child.cur_size < min_size {
	s.FixTreeUpwardsDelete(child)
}
```
After again recursively fixing the tree upwards, we know that the key `v` is deleted and that each cell has a valid number of keys in it. 

### Min

In order to find the minimum value in the tree, we just keep moving down the tree using the first child until we reach a leaf. We then return the value of the first key. The code to do so look like the following:
```go
func (s *Set_BTree) min() *Value_Location {
	// if tree is empty, return nil
	if s.root == nil {
		return nil
	}

	// start at root
	current_node := s.root

	// keep going towards first child until reach leaf
	for current_node.children[0] != nil {
		current_node = current_node.children[0]
	}

	// return first value of the leaf
	return &Value_Location{current_node, 0, current_node.keys[0]}
}
```

### Max
In order to find the maximum value in a tree, we just keep moving down the tree using the last child until we reach a leaf. we then return the value of the last key. The code to do so looks like the following:
```go
func (s *Set_BTree) max() *Value_Location {
	// if tree is empty, return nil
	if s.root == nil {
		return nil
	}

	// start at root
	current_node := s.root

	// keep going towards last child until reach leaf
	for current_node.children[0] != nil {
		current_node = current_node.children[current_node.cur_size]
	}

	// return last value of the leaf
	return &Value_Location{current_node, current_node.cur_size - 1, current_node.keys[current_node.cur_size-1]}
}
```


## Sources
For insert, I used the following guide: [programiz.com](https://www.programiz.com/dsa/insertion-into-a-b-tree)

For delete, I used the following guide: [webdocs.cs.ualberta.ca](https://webdocs.cs.ualberta.ca/~holte/T26/del-b-tree.html)

