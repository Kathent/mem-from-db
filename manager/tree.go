package manager

import (
	"math"
)

const (
	nodeTypeRoot = iota
	nodeTypeInternal
	nodeTypeLeaf
)

type Comparator interface {
	Compare(comparator Comparator) int
}

type treeNode struct {
	indexVal   []Comparator
	root       *treeNode
	subs       []*treeNode
	indexInArr int
}

func (node *treeNode) nodeType() int {
	if node.root == nil {
		return nodeTypeRoot
	} else if node.subs == nil {
		return nodeTypeLeaf
	} else {
		return nodeTypeInternal
	}
}

func (node *treeNode) add(comparator Comparator, tree *BpTree, flag bool) {
	switch node.nodeType() {
	case nodeTypeLeaf:
		searchIndex := binarySearch(node.indexVal, comparator)
		newArr := make([]Comparator, len(node.indexVal)+1)
		if searchIndex < 0 {
			copy(newArr[1:], node.indexVal)
			newArr[0] = comparator
		} else if searchIndex >= len(node.indexVal) {
			newArr[:len(newArr)-1], newArr[len(newArr)-1] = node.indexVal, comparator
		} else {
			newArr[:searchIndex], newArr[searchIndex], newArr[searchIndex+1:] =
				node.indexVal[:searchIndex], comparator, node.indexVal[searchIndex:]
		}

		// full, split node.
		if node.full(tree.m) {
			// split node to 2 node. and add mid val to parent node.
			mid := len(node.indexVal) / 2
			tmpArr := node.indexVal
			// clear this node
			node.indexVal = make([]Comparator, 0)
			// add element in parent node.
			node.root.add(node.indexVal[mid], tree, true)
			for k, v := range tmpArr {
				if k != mid {
					node.root.add(v, tree, false)
				}
			}
		} else {
			node.indexVal = newArr
		}
		break
	case nodeTypeRoot:
		if node.indexVal == nil {
			// first add
			node.indexVal = make([]Comparator, 0)
			node.indexVal = append(node.indexVal, comparator)
			node.subs = make([]*treeNode, 0)
			node.subs = append(node.subs, &treeNode{
				root:     node,
				indexVal: make([]Comparator, 0),
			})
		} else {
			// has index value already.
			for i := 0; i < len(node.indexVal); i++ {
				if comparator.Compare(node.indexVal[0]) < 0 {
					node.subs[i].add(comparator, tree, false)
				} else if comparator.Compare(node.indexVal[0]) == 0 {
					continue
				}
			}
		}
		break
	case nodeTypeInternal:
		break
	}
}

func binarySearch(comparators []Comparator, comparator Comparator) int {
	length := len(comparators)
	if length <= 0 {
		return 0
	}

	for {
		mid := len(comparators) / 2
		midVal := comparators[mid]
		compare := comparator.Compare(midVal)
		if compare != 0 {
			compare = int(compare / int(math.Abs(float64(compare))))
		}
		switch compare {
		case 0:
			return mid
		case 1:
			if mid+1 >= length {
				return mid + 1
			}
			search := binarySearch(comparators[mid+1:], comparator)
			return search + mid + 1
		case -1:
			if mid-1 < 0 {
				return 0
			}
			return binarySearch(comparators[:mid], comparator)
		}
	}
}

func (node *treeNode) remove(comparator Comparator, tree *BpTree) {

}
func (node *treeNode) full(count int) bool {
	return len(node.subs) >= count
}

type Tree interface {
	Add(comparator Comparator)
	Remove(comparator Comparator)
	Find(comparator Comparator)
}

type BpTree struct {
	Tree

	root *treeNode
	arr  []Comparator
	m    int
}

func NewBpTree(m int) *BpTree {
	return &BpTree{
		root: &treeNode{},
		arr:  make([]Comparator, 0),
		m:    m,
	}
}

func (b *BpTree) Add(comparator Comparator) {
	b.root.add(comparator, b, false)
}

func (b *BpTree) Remove(comparator Comparator) {
	b.root.remove(comparator, b)
}
