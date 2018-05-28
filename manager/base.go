package manager

import (
	"github.com/Kathent/mem-from-db/db/mysql"
	"github.com/orcaman/concurrent-map"
)

type Manager struct {
	m    cmap.ConcurrentMap
	conf TableConfig

	dbImpl *mysql.DbImpl
}

type TableConfig struct {
	DbName  string
	Name    string
	InitArr []interface{}
}

type columnInfo struct {
	ColumnName string
	DataType   string
}

type indexInfo struct {
}
