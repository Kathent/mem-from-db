package bptree

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

type IntComparator int

func (ic IntComparator) compare(c Comparator) int {
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

func BenchmarkBpTree_Insert(b *testing.B) {
	bp.Insert(KV{Key: IntComparator(rand.Intn(1000000)), Val: 2})
}

func TestBpTree_Delete(t *testing.T) {
	intNs := []int{81, 87, 45, 78, 98, 76, 100, 102, 101, 105, 106, 107}
	for i := 0; i < len(intNs); i++ {
		intN := intNs[i]
		bp.Insert(KV{Key: IntComparator(intN), Val: 2})
	}

	fmt.Println("after insert")
	bp.print()

	fmt.Println("start delete.")

	for i := 0; i < len(intNs); i++ {
		intN := intNs[i]
		fmt.Println(fmt.Sprintf("delete Key:%d", intN))
		bp.print()
		fmt.Println(fmt.Sprintf("delete res : %t", bp.Delete(KV{Key: IntComparator(intN)})))
		bp.print()
		fmt.Println("###########################")
	}
}

func TestBpTree_Insert(t *testing.T) {
	//intNs := []int{81, 87, 45, 78, 98, 76, 100, 102, 101, 105,106,107}
	for i := 0; i < 1000000; i++ {
		intN := rand.Intn(1000000)
		fmt.Println(fmt.Sprintf("input %d", intN))
		bp.Insert(KV{Key: IntComparator(intN), Val: 2})
		//bp.print()

		//fmt.Println("##################")
	}
}

func init() {
	bp = NewBpTree(5)
}
