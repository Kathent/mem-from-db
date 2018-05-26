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
