package bptree

import (
	"fmt"
	"reflect"
	"testing"
)

type IntComparator int

func (ic IntComparator) compare(c comparator) int {
	if val, ok := c.(IntComparator); ok {
		return int(ic) - int(val)
	}

	return 1
}

var bp *BpTree

func TestNewKvs(t *testing.T) {
	k := new(kvs)
	v := make(kvs, 0)
	fmt.Println(k, reflect.TypeOf(k))
	fmt.Println(v, reflect.TypeOf(v))
}

func TestBpTree_Insert(t *testing.T) {
	intNs := []int{81, 87, 45, 78, 98, 76, 100, 102, 101}
	for i := 0; i < len(intNs); i++ {
		intN := intNs[i]
		fmt.Println(fmt.Sprintf("input %d", intN))
		bp.Insert(kv{key: IntComparator(intN), val: 2})
		bp.print()

		fmt.Println("##################")
	}
}

func init() {
	bp = NewBpTree(5)
}
