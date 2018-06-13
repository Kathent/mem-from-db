package manager

import (
	"github.com/Kathent/mem-from-db/db/mysql"
	"github.com/Kathent/mem-from-db/manager/bptree"
	"github.com/Kathent/mem-from-db/manager/comparator"
	"github.com/orcaman/concurrent-map"
)

type Manager struct {
	m    cmap.ConcurrentMap
	conf TableConfig

	dbImpl *mysql.DbImpl
	IR     map[string]*IndexReserve
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
	Name    string
	Columns string
}

type IndexTree interface {
	Insert(k bptree.KV)
	Delete(k comparator.Comparator) bool
	Search(c comparator.Comparator) interface{}
}
