package bptree

import (
	"fmt"
	"sort"
)

type leafNode struct {
	kvs
	pre  *leafNode
	next *leafNode
	p    *indexNode

	m int
}

func (ln *leafNode) search(key comparator) (int, bool) {
	search := sort.Search(len(ln.kvs), func(i int) bool {
		return key.compare(ln.kvs[i].key) <= 0
	})

	if search >= len(ln.kvs) {
		return search, false
	}

	return search, ln.kvs[search].key == key
}

func (ln *leafNode) parent() *indexNode {
	return ln.p
}

func (ln *leafNode) setParent(p *indexNode) {
	ln.p = p
}

func (ln *leafNode) full(add int) bool {
	return len(ln.kvs)+add >= ln.m
}

func (ln *leafNode) hunger(minus int) bool {
	return len(ln.kvs)-minus <= ln.m/2-1
}

func (ln *leafNode) insert(k kv) int {
	index, _ := ln.search(k.key)
	ln.kvs = append(ln.kvs, nil)
	copy(ln.kvs[index+1:], ln.kvs[index:])
	ln.kvs[index] = &k

	if ln.p != nil {
		idx, _ := ln.p.search(k.key)
		ln.p.kis[idx].key = ln.kvs[len(ln.kvs)-1].key
	}

	return index
}

func (ln *leafNode) split() {
	// create new leaf node.
	nl := newLeafNode(ln.p, ln.m)
	nl.pre = ln
	nl.next = ln.next
	if nl.next != nil {
		nl.next.pre = nl
	}
	mid := len(ln.kvs) / 2
	nl.kvs = make(kvs, len(ln.kvs)-mid)
	copy(nl.kvs, ln.kvs[mid:])

	// resolve pre node.
	ln.next = nl
	ln.kvs = ln.kvs[:mid]

	if ln.p == nil {
		// no parent. create parent and link child.
		nn := newIndexNode(nil, ln.m)
		nn.kis = append(nn.kis, &ki{key: ln.kvs[len(ln.kvs)-1].key, node: ln},
			&ki{key: nl.kvs[len(nl.kvs)-1].key, node: nl})
		ln.p = nn
		nl.p = nn
		return
	}

	// has parent.
	index := ln.p.insert(ln.kvs[len(ln.kvs)-1].key)
	ln.p.kis[index].node = ln
	ln.p.kis[index+1].node = nl
}

func (ln *leafNode) isNil() bool {
	return ln == nil
}
func (ln *leafNode) delete(i int) {
	preKey := ln.kvs[i].key
	if i >= len(ln.kvs)-1 {
		ln.kvs = ln.kvs[:len(ln.kvs)-1]
	} else {
		copy(ln.kvs[i:], ln.kvs[i+1:])
		ln.kvs = ln.kvs[:len(ln.kvs)-1]
	}

	if ln.p == nil {
		return
	}

	search, _ := ln.p.search(preKey)
	if ln.p != nil {
		ln.p.kis[search].key = ln.kvs[len(ln.kvs)-1].key
	}

	if ln.hunger(0) {
		if ln.pre != nil && !ln.pre.hunger(1) {
			// lean from left sibling
			preVal := ln.pre.kvs[len(ln.pre.kvs)-1]
			ln.pre.kvs = ln.pre.kvs[:len(ln.pre.kvs)-1]
			ln.kvs = append(ln.kvs, nil)
			copy(ln.kvs[1:], ln.kvs[:len(ln.kvs)-1])
			ln.kvs[0] = preVal

			// update left sibling key
			ln.p.kis[search-1].key = ln.pre.kvs[len(ln.pre.kvs)-1].key
		} else if ln.next != nil && !ln.next.hunger(1) {
			// lean from next sibling
			preVal := ln.next.kvs[0]
			ln.next.kvs = ln.next.kvs[1:]
			ln.kvs = append(ln.kvs, preVal)

			// update right sibling key
			ln.p.kis[search].key = preVal.key
		} else if ln.pre != nil && !ln.pre.full(len(ln.kvs)) {
			// merge left
			ln.kvs = append(ln.pre.kvs, ln.kvs...)
			ln.pre = ln.pre.pre
			if ln.pre != nil {
				ln.pre.next = ln
			}
			copy(ln.p.kis[search-1:], ln.p.kis[search:])
			ln.p.kis = ln.p.kis[:len(ln.p.kis)-1]
		} else if ln.next != nil && !ln.next.full(len(ln.kvs)) {
			// merge right
			ln.next.kvs = append(ln.kvs, ln.next.kvs...)
			copy(ln.p.kis[search:], ln.p.kis[search+1:])
			ln.p.kis = ln.p.kis[:len(ln.p.kis)-1]
			if ln.pre != nil {
				ln.pre.next = ln.next
			}
			ln.next.pre = ln.pre
		}
	}
}

type node interface {
	search(key comparator) (int, bool)
	parent() *indexNode
	setParent(p *indexNode)
	full(add int) bool
	split()
	isNil() bool
	hunger(minus int) bool
}

type comparator interface {
	compare(c comparator) int
}

type kv struct {
	key comparator
	val interface{}
}

type kvs []*kv

func (k kvs) String() string {
	val := ""
	for _, v := range k {
		val += fmt.Sprintf("%+v,", v)
	}

	return val
}

type indexNode struct {
	kis
	p *indexNode
	m int
}

func (k kis) String() string {
	val := ""
	for _, v := range k {
		val += fmt.Sprintf("%+v,", v)
	}

	return val
}

func (in *indexNode) search(c comparator) (int, bool) {
	search := sort.Search(len(in.kis)-1, func(i int) bool {
		return c.compare(in.kis[i].key) <= 0
	})
	return search, true
}

func (in *indexNode) parent() *indexNode {
	return in.p
}

func (in *indexNode) setParent(p *indexNode) {
	in.p = p
}

func (in *indexNode) full(add int) bool {
	return len(in.kis)+add >= in.m
}

func (in *indexNode) insert(c comparator) int {
	index, _ := in.search(c)
	in.kis = append(in.kis, nil)
	copy(in.kis[index+1:], in.kis[index:])
	in.kis[index] = &ki{key: c}
	return index
}

func (in *indexNode) split() {
	// create new index node.
	nl := newIndexNode(in.p, in.m)
	mid := len(in.kis) / 2
	nl.kis = make(kis, len(in.kis)-mid)
	copy(nl.kis, in.kis[mid:])

	in.kis = in.kis[:mid]

	if in.p == nil {
		// no parent. create parent and link child.
		nn := newIndexNode(nil, in.m)
		nn.kis = append(nn.kis, &ki{key: in.kis[len(in.kis)-1].key, node: in},
			&ki{key: nl.kis[len(nl.kis)-1].key, node: nl})
		in.p = nn
		nl.p = nn
		return
	}

	// has parent.
	nl.p = in.p
	index := in.p.insert(in.kis[mid-1].key)
	in.p.kis[index].node = nl
}

func (in *indexNode) isNil() bool {
	return in == nil
}

func (in *indexNode) hunger(minus int) bool {
	return len(in.kis)-minus <= in.m/2-1
}

func (in *indexNode) eat() {
	if in.p == nil {
		return
	}

	index, _ := in.p.search(in.kis[0].key)
	var preIn, nextIn *ki
	if index > 0 {
		preIn = in.p.kis[index-1]
	}

	if index < len(in.p.kis)-1 {
		nextIn = in.p.kis[index+1]
	}

	if preIn != nil && !preIn.node.hunger(1) {
		preNode := preIn.node.(*indexNode)
		preLastKi := preNode.kis[len(preNode.kis)-1]
		preNode.kis = preNode.kis[:len(preNode.kis)-1]
		preIn.key = preNode.kis[len(preNode.kis)-1].key

		in.kis = append(in.kis, nil)
		copy(in.kis[1:], in.kis[:len(in.kis)-1])
		in.kis[0] = preLastKi
	} else if nextIn != nil && !nextIn.node.hunger(1) {
		nextNode := nextIn.node.(*indexNode)
		nextFirstKi := nextNode.kis[0]
		nextNode.kis = nextNode.kis[1:]

		in.kis = append(in.kis, nextFirstKi)
		in.p.kis[index].key = nextFirstKi.key
	} else if preIn != nil && !preIn.node.full(len(in.kis)) {
		preNode := preIn.node.(*indexNode)
		in.kis = append(preNode.kis, in.kis...)
		copy(in.p.kis[index-1:], in.p.kis[index:])
		in.p.kis = in.p.kis[:len(in.p.kis)-1]
		for _, v := range in.kis {
			v.node.setParent(in)
		}
	} else if nextIn != nil && !nextIn.node.full(len(in.kis)) {
		nextNode := nextIn.node.(*indexNode)
		nextNode.kis = append(in.kis, nextNode.kis...)
		copy(in.p.kis[index:], in.p.kis[index+1:])
		in.p.kis = in.p.kis[:len(in.p.kis)-1]
		for _, v := range nextNode.kis {
			v.node.setParent(nextNode)
		}
	}
}

type ki struct {
	key  comparator
	node node
}

type kis []*ki

type BpTree struct {
	root node
	m    int
}

func newLeafNode(p *indexNode, m int) *leafNode {
	return &leafNode{
		kvs: make([]*kv, 0),
		p:   p,
		m:   m,
	}
}

func newIndexNode(p *indexNode, m int) *indexNode {
	return &indexNode{
		kis: make([]*ki, 0),
		p:   p,
		m:   m,
	}
}

func NewBpTree(m int) *BpTree {
	return &BpTree{
		root: newLeafNode(nil, m),
		m:    m,
	}
}

func (b *BpTree) Search(c comparator) interface{} {
	search, _, _ := b.search(c)
	return search
}

func (b *BpTree) search(c comparator) (interface{}, int, *leafNode) {
	var n = b.root
	for {
		if val, ok := n.(*leafNode); ok {
			search, ok := val.search(c)
			if !ok {
				return nil, 0, val
			}

			if search >= len(val.kvs) {
				return nil, search, val
			}

			return val.kvs[search].val, search, val
		} else if val, ok := n.(*indexNode); ok {
			in, _ := val.search(c)
			n = val.kis[in].node
			continue
		} else {
			return nil, 0, nil
		}

	}

	return nil, 0, nil
}

func (b *BpTree) Insert(k kv) {
	search, index, leaf := b.search(k.key)

	// if found
	if search != nil {
		leaf.kvs[index].val = k.val
		return
	}

	// not found
	leaf.insert(k)
	var n node = leaf
	for ; n != nil && !n.isNil(); n = n.parent() {
		if !n.isNil() && n.full(0) {
			n.split()
			if !b.root.parent().isNil() {
				b.root = b.root.parent()
			}
		}
	}
}

func (b *BpTree) print() {
	nodeArr := []node{b.root}
	for {
		if len(nodeArr) <= 0 {
			return
		}

		tmp := make([]node, 0)
		for _, v := range nodeArr {
			if val, ok := v.(*leafNode); ok {
				fmt.Print(fmt.Sprintf("%s, %p", val.kvs, val.p))
				fmt.Print("____")
			} else if val, ok := v.(*indexNode); ok {
				fmt.Print(fmt.Sprintf("%s, %p", val.kis, val.p))
				fmt.Print("----")
				for _, v := range val.kis {
					tmp = append(tmp, v.node)
				}
			}
		}

		fmt.Println()

		nodeArr = tmp
	}
}

func (b *BpTree) Delete(k kv) bool {
	search, index, leaf := b.search(k.key)
	if search == nil {
		return false
	}

	leaf.delete(index)

	var n = leaf.p
	for ; !n.isNil(); n = n.parent() {
		if b.root == n {
			if len(n.kis) <= 1 {
				b.root = n.kis[0].node
				b.root.setParent(nil)
			}
			return true
		}

		if n.hunger(0) {
			n.eat()
		}
	}
	return true
}
