package manager

import (
	"testing"
	"git.icsoc.net/paas/voice-robot-dispatch/manage"
	"github.com/Kathent/mem-from-db/db/mysql"
)

func TestNewManager(t *testing.T) {
	tableConfig := TableConfig{
		DbName:  "robot",
		Name:    "task_info",
		InitArr: make([]manage.TaskInfo, 0)}

	d, err := mysql.NewDbImpl(mysql.DbConfig{
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
		t.Errorf("create db impl fail. %v", err)
		t.FailNow()
	}

	NewManager(tableConfig, d)
}