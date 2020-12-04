package addressBook

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/*
	数据库错误的定义
*/
type DBErrorCode int

const DBErr_NotFound DBErrorCode = 1
const DBErr_NoError DBErrorCode = 0
const DBErr_Other DBErrorCode = -1

type DBError struct {
	Error error
	Code  DBErrorCode
}

func newErrorWithString(code DBErrorCode, msg string) *DBError {
	return &DBError{errors.New(msg), code}
}

func (c *DBError) IsNotFound() bool {
	return c.Code == DBErr_NotFound
}

type DataStore struct {
	Total     int64       `json:"total"`
	TotalPage int         `json:"totalPage"`
	Data      interface{} `json:"data"`
}

func (d *DataStore) BuildData(data interface{}) *DataStore {
	d.Data = data
	return d
}
func (d *DataStore) BuildTotal(total int) *DataStore {
	d.Total = int64(total)
	return d
}
func (d *DataStore) Build(pageSize int) *DataStore {
	d.TotalPage = (d.TotalPage + pageSize - 1) / pageSize
	return d
}

/*
	数据链接
*/

type DBAccess struct {
	db *sql.DB
	tx *sql.Tx
}

func GetDBAccess() *DBAccess {
	return &DBAccess{nil, nil}
}

func (c *DBAccess) BeginTrans() *DBError {
	db, e := openDB()
	if e != nil {
		return newErrorWithString(DBErr_Other, e.Error())
	}

	c.db = db

	var e1 error
	c.tx, e1 = db.Begin()
	if e1 != nil {
		db.Close()
		return newErrorWithString(DBErr_Other, e1.Error())
	}

	return nil
}

func (c *DBAccess) Commit() *DBError {
	defer c.db.Close()

	if c.tx == nil {
		return newErrorWithString(DBErr_Other, "没有begin")
	}
	e := c.tx.Commit()
	if e != nil {
		return newErrorWithString(DBErr_Other, e.Error())
	}

	return nil
}

func (c *DBAccess) Rollback() *DBError {
	defer c.db.Close()

	if c.tx == nil {
		return newErrorWithString(DBErr_Other, "没有begin")
	}
	e := c.tx.Rollback()
	if e != nil {
		return newErrorWithString(DBErr_Other, e.Error())
	}
	return nil
}

var sql_str string

func SetDBAStr(user, passwd, host, port, database string) {
	sql_str = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local", user, passwd, host, port, database)
}

//将来这个地方做成可配置
func openDB() (*sql.DB, error) {
	//r, e := sql.Open("mysql", "webadmin:WebAdmin123!@tcp(127.0.0.1:3306)/cardmgr?charset=utf8&parseTime=true")
	r, e := sql.Open("mysql", sql_str)
	if e != nil {
		return nil, e
	}
	return r, nil
}

var gormDb *gorm.DB
var once sync.Once

func OpenGorm() *gorm.DB {
	if gormDb == nil || gormDb.Value == nil {
		once.Do(func() {
			var err error
			gormDb, err = gorm.Open("mysql", sql_str)

			if err != nil {
				fmt.Println("数据库连接异常%w", err)
			}
			if gormDb != nil {
				gormDb.LogMode(false)
			}
			gormDb.DB().SetMaxOpenConns(40)               //设置最大连接数
			gormDb.DB().SetMaxIdleConns(5)                //设置最大空闲连接数
			gormDb.DB().SetConnMaxLifetime(time.Hour * 4) //设置连接的最大生命周期
			gormDb.SingularTable(true)
		})
	}
	return gormDb
}

/*
var prikey_key = "codenai.com"

func encrptPriKey(prikey string) string {
	return tea.Encrypt(prikey, prikey_key, 16)
}

func decrptPriKey(prikey string) string {
	r, _ := tea.Decrypt(prikey, prikey_key, 16)
	return r
}
*/

func GetSqlResultMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	var resultMap = make([]map[string]interface{}, 0)
	columns, _ := rows.Columns()
	for rows.Next() {
		var pointslice = make([]interface{}, len(columns))
		var valslice = make([]interface{}, len(columns))
		for i, _ := range columns {
			pointslice[i] = &valslice[i]
		}
		if err := rows.Scan(pointslice...); err != nil {
			return nil, err
		}
		var tempMap = make(map[string]interface{})
		for i, ele := range columns {
			var v interface{}
			if val, ok := valslice[i].([]byte); ok {
				v = string(val)
			} else if val, ok := valslice[i].(time.Time); ok {
				v = val.Format("2006-01-02 15:04:05")
			} else {
				v = valslice[i]
			}
			tempMap[strings.ToLower(ele)] = v
		}
		resultMap = append(resultMap, tempMap)
	}
	return resultMap, nil
}

//upcase
func Rows2Map(rows *sql.Rows) ([]map[string]interface{}, error) {
	var resultMap = make([]map[string]interface{}, 0)
	columns, _ := rows.Columns()
	for rows.Next() {
		var pointslice = make([]interface{}, len(columns))
		var valslice = make([]interface{}, len(columns))
		for i, _ := range columns {
			pointslice[i] = &valslice[i]
		}
		if err := rows.Scan(pointslice...); err != nil {
			return nil, err
		}
		var tempMap = make(map[string]interface{})
		for i, ele := range columns {
			var v interface{}
			if val, ok := valslice[i].([]byte); ok {
				v = string(val)
			} else if val, ok := valslice[i].(time.Time); ok {
				v = val.Format("2006-01-02 15:04:05")
			} else {
				v = valslice[i]
			}
			tempMap[strings.ToUpper(ele)] = v
		}
		resultMap = append(resultMap, tempMap)
	}
	return resultMap, nil
}
