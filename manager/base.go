package manager

import (
	"github.com/Kathent/mem-from-db/db/mysql"
	"github.com/orcaman/concurrent-map"
	"github.com/Kathent/mem-from-db/manager/bptree"
	"github.com/Kathent/mem-from-db/manager/comparator"
)

type Manager struct {
	m    cmap.ConcurrentMap
	conf TableConfig

	dbImpl *mysql.DbImpl
	ir []*indexReserve
}

type TableConfig struct {
	DbName  string
	Name    string
	InitArr interface{}
}

type columnInfo struct {
	ColumnName string
	DataType   string
	FieldIndex int
}

type indexInfo struct {
	Columns string
}

type IndexTree interface {
	Insert(k bptree.KV)
	Delete(k comparator.Comparator) bool
	Search(c comparator.Comparator) interface{}
}
