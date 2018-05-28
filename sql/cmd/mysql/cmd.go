package mysql

import "github.com/Kathent/mem-from-db/sql/cmd/base"

const (
	ExecTypeInsert = iota
	ExecTypeUpdate
	ExecTypeDelete
)

type DbCmd interface {
	base.DbCmd
	Sql() string
	Args() []interface{}
}

type ExecCmd interface {
	DbCmd
	ExecType() int
}

type QueryCmd struct {
	SqlStr string
	ArgArr []interface{}
}

func (q *QueryCmd) Sql() string {
	return q.SqlStr
}

func (q *QueryCmd) Args() []interface{} {
	return q.ArgArr
}
