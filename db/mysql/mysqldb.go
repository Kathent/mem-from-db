package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Kathent/mem-from-db/sql/cmd/base"
	"github.com/Kathent/mem-from-db/sql/cmd/mysql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
)

const (
	// userName:password@tcp(addr)/db?timeout=30s&readTimeout=30s
	dataSourceFmt = "%s:%s@tcp(%s)/%s?timeout=%ds&readTimeout=%ds"
)

var (
	queryResTypeErr = errors.New("res type can not be map")

	intNullType    = reflect.TypeOf(sql.NullInt64{})
	floatNullType  = reflect.TypeOf(sql.NullFloat64{})
	stringNullType = reflect.TypeOf(sql.NullString{})
	boolNullType   = reflect.TypeOf(sql.NullBool{})
)

type DbConfig struct {
	Addr        string
	DB          string
	Username    string
	Password    string
	MaxIdleConn int
	MaxOpenConn int
	Timeout     int
	ReadTimeout int
}

type DbImpl struct {
	db *sql.DB
}

func NewDbImpl(conf DbConfig) (*DbImpl, error) {
	dataSource := fmt.Sprintf(dataSourceFmt, conf.Username, conf.Password,
		conf.Addr, conf.DB, conf.Timeout, conf.ReadTimeout)
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(conf.MaxIdleConn)
	db.SetMaxOpenConns(conf.MaxOpenConn)
	imp := DbImpl{
		db: db,
	}

	return &imp, nil
}

func (d *DbImpl) Exec(cmd base.DbCmd) (interface{}, error) {
	if val, ok := cmd.(mysql.ExecCmd); ok {
		result, err := d.db.Exec(val.Sql(), val.Args())
		if err != nil {
			return nil, err
		}

		if val.ExecType() == mysql.ExecTypeInsert {
			return result.LastInsertId()
		}

		return result.RowsAffected()
	}

	return nil, base.TypeErrF(fmt.Sprintf("type err..%v", cmd))
}

func (d *DbImpl) Query(cmd base.DbCmd, res interface{}) error {
	if val, ok := cmd.(mysql.DbCmd); ok {
		rows, err := d.db.Query(val.Sql(), val.Args()...)
		if err != nil {
			fmt.Println(fmt.Sprintf("query err:%v", err))
			return err
		}
		defer rows.Close()
		return mapperRes(rows, res)
	}

	return base.TypeErrF(fmt.Sprintf("type err..%v", cmd))
}

func mapperRes(rows *sql.Rows, res interface{}) error {
	tp := reflect.TypeOf(res)
	if tp.Kind() != reflect.Ptr {
		return queryResTypeErr
	}

	tp = tp.Elem()
	if tp.Kind() == reflect.Array || tp.Kind() == reflect.Slice {
		// 数组元素类型
		eleType := tp.Elem()
		// 新生成一个数组
		valArr := make([]reflect.Value, 0)
		for rows.Next() {
			val := reflect.New(eleType)
			valArr = append(valArr, val.Elem())
			err := resolveEle(rows, val)
			if err != nil {
				return err
			}
			fmt.Println(val)
		}
		slice := reflect.MakeSlice(tp, len(valArr), len(valArr))
		for k, val := range valArr {
			slice.Index(k).Set(val)
		}
		reflect.Copy(reflect.ValueOf(res).Elem(), slice)
	} else if tp.Kind() == reflect.Struct {
		err := resolveEle(rows, reflect.ValueOf(tp))
		if err != nil {
			return err
		}
	} else {
		return queryResTypeErr
	}

	return nil
}

func resolveEle(rows *sql.Rows, val reflect.Value) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	// 去除指针
	elem := val.Elem()
	dt := make([]interface{}, 0)
	for _, col := range cols {
		afterName := transferName(col)
		field := elem.FieldByName(afterName).Addr().Interface()
		dt = append(dt, field)
	}

	scan := rows.Scan(dt...)
	return scan
}

func nullTypeMap(field reflect.Type) reflect.Type {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return intNullType
	case reflect.Float32, reflect.Float64:
		return floatNullType
	case reflect.Bool:
		return boolNullType
	case reflect.String:
		return stringNullType
	default:
		return field

	}
	return field
}

func transferName(colName string) string {
	var needCap = true
	return strings.Map(func(r rune) rune {
		if needCap && r >= 'a' && r <= 'z' {
			r += 'A' - 'a'
			needCap = false
		} else if r == '_' {
			needCap = true
			return -1
		}
		return r
	}, colName)
}
