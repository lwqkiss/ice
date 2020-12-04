package dbaccess

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"sync"

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
	db, e := OpenDB()
	if e != nil {
		return e
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
	sql_str = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&multiStatements=true", user, passwd, host, port, database)
}

//将来这个地方做成可配置
func OpenDB() (*sql.DB, *DBError) {
	//r, e := sql.Open("mysql", "webadmin:WebAdmin123!@tcp(127.0.0.1:3306)/cardmgr?charset=utf8&parseTime=true")
	r, e := sql.Open("mysql", sql_str)
	if e != nil {
		return nil, newErrorWithString(DBErr_Other, e.Error())
	}
	return r, nil
}

var gormDb *gorm.DB
var once sync.Once

func OpenGorm() *gorm.DB {
	if gormDb == nil {
		once.Do(func() {
			var err error
			gormDb, err = gorm.Open("mysql", sql_str)
			gormDb.LogMode(true)
			if err != nil {
				log.Println("数据库连接异常%w", err)
			}
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
