package mysql

import (
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
	CommitTime      int64
	StartTime       int64
	EndTime         int64
	Status          int
	ClientStaticsId string
	WorkTimeId      int64
	VccId           string
	CallNumber      string
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
		Addr:        "localhost:32770",
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
