package mysql

import (
	"database/sql"
	"fmt"
	"github.com/Kathent/mem-from-db/sql/cmd/mysql"
	"testing"
)

var (
	db *DbImpl
)

type cmd struct {
	sql string
	arg []interface{}
	typ int
	mysql.DbCmd
}

type taskInfo struct {
	Id              int64
	TaskName        string
	TaskNum         int
	WordStrategyId  int64
	CommitTime      sql.NullInt64
	StartTime       sql.NullInt64
	EndTime         sql.NullInt64
	Status          int
	ClientStaticsId sql.NullString
	WorkTimeId      sql.NullInt64
	VccId           sql.NullString
	CallNumber      sql.NullString
}

func (c *cmd) Sql() string {
	return c.sql
}

func (c *cmd) Args() []interface{} {
	return c.arg
}

func (c *cmd) ExecType() int {
	return c.typ
}

func TestNewDbImpl(t *testing.T) {
	var taskInfo []taskInfo
	query := db.Query(&cmd{}, &taskInfo)
	if query != nil {
		fmt.Println(query)
		t.FailNow()
	}

	fmt.Println(taskInfo)
}

func TestDbImpl_Exec(t *testing.T) {
	ti := taskInfo{
		TaskName:       "haha",
		TaskNum:        1,
		WordStrategyId: 1,
	}

	exec, err := db.Exec(&cmd{sql: "insert into task_info (task_name, task_num, word_strategy_id) values (?, ?, ?)",
		arg: []interface{}{ti.TaskName, ti.TaskNum, ti.WordStrategyId},
		typ: mysql.ExecTypeInsert})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Println(exec)
}

func init() {
	d, err := NewDbImpl(DbConfig{
		Addr:        "192.168.96.204:3306",
		DB:          "robot",
		Username:    "root",
		Password:    "123456",
		MaxIdleConn: 10,
		MaxOpenConn: 10,
		Timeout:     30,
		ReadTimeout: 30,
	})
	if err != nil {
		panic(err)
	}

	db = d
}
