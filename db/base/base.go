package base

import "github.com/Kathent/mem-from-db/sql/cmd/base"

type Db interface {
	Exec(cmd base.DbCmd) error
	Query(cmd base.DbCmd, res interface{}) error
}
