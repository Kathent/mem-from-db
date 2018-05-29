package bptree

import "sort"

type leafNode struct {
	kvs
	pre *leafNode
	next *leafNode
	p *indexNode

	m int
}

func (ln *leafNode) search(key comparator) (int, bool) {
	search := sort.Search(len(ln.kvs), func(i int) bool {
		return key.compare(ln.kvs[i].key) < 0
	})
	return search, true
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

type node interface {
	search(key comparator) (int, bool)
	parent() *indexNode
	setParent(p *indexNode)
	full() bool
}

type comparator interface {
	compare(c comparator) int
}

type kv struct {
	key comparator
	val interface{}
}

type kvs []kv

type indexNode struct {
	kis
	p *indexNode
	m int
}

func (in *indexNode) search(c comparator) (int, bool) {
	search := sort.Search(len(in.kis), func(i int) bool {
		return c.compare(in.kis[i].key) < 0
	})

	return search, in.kis[search].key == c
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

type ki struct {
	key comparator
	node node
}

type kis []ki

type BpTree struct {
	root *indexNode
	leaf *leafNode
	m int
}

func newLeafNode(p *indexNode, m int) *leafNode {
	return &leafNode{
		kvs: make([]kv, 0),
		p: p,
		m: m,
	}
}

func newIndexNode(p *indexNode, m int) *indexNode {
	return &indexNode{
		kis: make([]ki, 0),
		p: p,
		m: m,
	}
}

func NewBpTree(m int) *BpTree {
	root := newIndexNode(nil, m)
	return &BpTree{
		root: root,
		leaf: newLeafNode(root, m),
		m: m,
	}
}


func (b *BpTree) Search(c comparator) interface{}{
	var n node = b.root
	for {
		if val, ok := n.(*leafNode); ok {
			search, ok := val.search(c)
			if !ok {
				return nil
			}

			return val.kvs[search].val
		}else if val, ok := n.(*indexNode); ok {
			in, _ := val.search(c)
			n = val.kis[in].node
			continue
		}else {
			return nil
		}

	}

	return nil
}
