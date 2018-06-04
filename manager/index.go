package manager

import (
	"strings"
	"errors"
	"fmt"
	"github.com/Kathent/mem-from-db/manager/comparator"
	"reflect"
	"github.com/Kathent/mem-from-db/manager/bptree"
)

const(
	treeM = 100
)

func errorF(val string) error{
	return errors.New(val)
}

type indexReserve struct {
	tree *IndexTree
}


type indexReserveBuilder struct {
	res []interface{}
	i indexInfo
	c []columnInfo
}

type keyValueComparator struct {
	keys []comparator.Comparator
}

func (kvc *keyValueComparator) Compare(c comparator.Comparator) int {
	if val, ok := c.(*keyValueComparator); ok {
		for idx, v := range kvc.keys {
			compareVal := v.Compare(val.keys[idx])
			if compareVal != 0 {
				return compareVal
			}
		}

		return 0
	}

	return 1
}

func NewIndexReserveBuilder() *indexReserveBuilder {
	return &indexReserveBuilder{}
}

func (b *indexReserveBuilder) withRes(res []interface{}) *indexReserveBuilder {
	b.res = res
	return b
}

func (b *indexReserveBuilder) withIndexInfo(i indexInfo) *indexReserveBuilder {
	b.i = i
	return b
}

func (b *indexReserveBuilder) withColumnInfo(c []columnInfo) *indexReserveBuilder {
	b.c = c
	return b
}

func (b *indexReserveBuilder) build() (*indexReserve, error) {
	bpt := bptree.NewBpTree(treeM)
	ir := &indexReserve{
		tree: bpt,
	}

	columns := strings.Split(b.i.Columns, ",")
	kc := keyValueComparator{}

	indexColumnArr := make([]*columnInfo, 0)
	for idx, v := range columns {
		info := findInColumnInfo(v, b.c)
		if info == nil {
			return nil, errorF(fmt.Sprintf("find index column fail. v:%s", v))
		}
		indexColumnArr[idx] = info
	}

	for _, v := range b.res {
		val := reflect.ValueOf(v)
		for _, c := range indexColumnArr {
			kc.keys = append(kc.keys, comparator.NewComparator(c.DataType, val.FieldByName(c.ColumnName).Interface()))
		}

		bpt.Insert(bptree.KV{Key: &kc, Val: val.Interface()})
	}

	return ir, nil
}

func arrayMap(arr []columnInfo, f func(info columnInfo) bool) *columnInfo {
	for _, v := range arr {
		if f(v) {
			return &v
		}
	}

	return nil
}

func findInColumnInfo(s string, arr []columnInfo) *columnInfo{
	return arrayMap(arr, func(info columnInfo) bool {
		if info.ColumnName == s {
			return true
		}
		return false
	})
}