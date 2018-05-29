package manager

import (
	"fmt"
	"testing"
)

type indexCo int

func (i indexCo) Compare(comparator Comparator) int {
	if val, ok := comparator.(indexCo); ok {
		return int(i) - int(val)
	}

	return 1
}

func TestBinarySearch(t *testing.T) {
	val := make([]Comparator, 4)
	val[0] = indexCo(1)
	val[1] = indexCo(2)
	val[2] = indexCo(5)
	val[3] = indexCo(7)

	fmt.Println(val)
	search := binarySearch(val, indexCo(3))
	fmt.Println(search)

	search = binarySearch(val, indexCo(11))
	fmt.Println(search)

	search = binarySearch(val, indexCo(6))
	fmt.Println(search)

	search = binarySearch(val, indexCo(9))
	fmt.Println(search)
}
