package manager

import (
	"fmt"
	"github.com/Kathent/mem-from-db/db/mysql"
	cmd "github.com/Kathent/mem-from-db/sql/cmd/mysql"
	"github.com/orcaman/concurrent-map"
)

var (
	queryTableColumnStr = "select column_name, " +
		"ordinal_position, " +
		"data_type " +
		"from information_schema.columns where " +
		"table_schema = ? && table_name = ?"

	queryIndexStr = "select key_name, column_name, index_type from %s.%s"
)

func NewManager(conf TableConfig, impl *mysql.DbImpl) *Manager {
	manager := &Manager{
		m:      cmap.New(),
		conf:   conf,
		dbImpl: impl,

		ir: make([]*indexReserve, 0),
	}

	manager.init()
	return manager
}

func (m *Manager) init() {
	queryCmd := cmd.QueryCmd{
		SqlStr: queryTableColumnStr,
		ArgArr: []interface{}{m.conf.DbName, m.conf.Name},
	}

	ca := make([]columnInfo, 0)
	err := m.dbImpl.Query(queryCmd, &ca)
	if err != nil {
		panic(err)
	}

	indexCmd := cmd.QueryCmd{
		SqlStr: fmt.Sprintf(queryIndexStr, m.conf.DbName, m.conf.Name),
	}

	ic := make([]indexInfo, 0)
	err = m.dbImpl.Query(indexCmd, &ic)
	if err != nil {
		panic(err)
	}

	for _, v := range ic {
		m.ir = append(m.ir, NewIndexReserveBuilder().withRes(m.conf.InitArr).withColumnInfo(ca).
			withIndexInfo(v).build())
	}
}
