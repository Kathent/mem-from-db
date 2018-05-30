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
		return key.compare(ln.kvs[i].key) < 0
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

func (ln *leafNode) full() bool {
	return len(ln.kvs) >= ln.m
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
	mid := len(ln.kvs) / 2
	nl.kvs = make(kvs, len(ln.kvs)-mid)
	copy(nl.kvs, ln.kvs[mid:])
	nl.p = ln.p

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
	ln.p.kis[index + 1].node = nl
}

func (ln *leafNode) isNil() bool {
	return ln == nil
}

type node interface {
	search(key comparator) (int, bool)
	parent() *indexNode
	setParent(p *indexNode)
	full() bool
	split()
	isNil() bool
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
		return c.compare(in.kis[i].key) < 0
	})
	return search, true
}

func (in *indexNode) parent() *indexNode {
	return in.p
}

func (in *indexNode) setParent(p *indexNode) {
	in.p = p
}

func (in *indexNode) full() bool {
	return len(in.kis) >= in.m
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
	index := in.p.insert(in.kis[mid - 1].key)
	in.p.kis[index].node = nl
}

func (in *indexNode) isNil() bool {
	return in == nil
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
		if !n.isNil() && n.full() {
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
				fmt.Print(fmt.Sprintf("%s", val.kvs))
				fmt.Print("----")
			} else if val, ok := v.(*indexNode); ok {
				fmt.Print(fmt.Sprintf("%s", val.kis))
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
