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

func (*cmd) Sql() string {
	return "select * from task_info"
}

func (*cmd) Args() []interface{} {
	return nil
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
