package addressBook

import (
	"fmt"
)

type DeptUser struct {
	Id         int    `json:"id" form:"id" gorm:"column:ID"`
	DeptId     int    `json:"deptid" form:"deptid" gorm:"column:DEPT_ID"`
	UserId     int    `json:"userid" form:"userid" gorm:"column:USER_ID"`
	IsDelete   int    `json:"isdelete" form:"isdelete" gorm:"column:IS_DELETE"`
	UpdateTime string `json:"updatetime" form:"updatetime" gorm:"column:UPDATE_TIME"`
}

func (u *DeptUser) TableName() string {
	return "DEPT_USER"
}
func AddDeptUser(b *DeptUser) (int64, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return 0, fmt.Errorf("数据库连接异常：%w", err)
	}
	sql := `insert into DEPT_USER ( DEPT_ID,USER_ID,IS_DELETE,UPDATE_TIME ) values (?,?,?,?)`
	result, er := db.Exec(sql, b.DeptId, b.UserId, b.IsDelete, b.UpdateTime)
	if er != nil {
		return 0, er
	}
	return result.LastInsertId()
}

//remember: all of the fields will be updated ,although  you do not want to , be carefull
func UpdateDeptUser(b *DeptUser) (int64, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return 0, fmt.Errorf("数据库连接异常：%w", err)
	}
	sql := `update DEPT_USER set DEPT_ID=?,USER_ID=?,IS_DELETE=?,UPDATE_TIME=?  where ID=?`

	result, er := db.Exec(sql, b.DeptId, b.UserId, b.IsDelete, b.UpdateTime, b.Id)
	if er != nil {
		return 0, fmt.Errorf("sql执行错误：%w", err)
	}
	return result.RowsAffected()
}

func GetDeptUserById(id int) (*DeptUser, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return nil, fmt.Errorf("数据库连接异常：%w", err)
	}
	sql := `select ID,DEPT_ID,USER_ID,IS_DELETE,UPDATE_TIME from DEPT_USER where ID=?`
	var result = DeptUser{}
	er := db.QueryRow(sql, id).Scan(&result.Id, &result.DeptId, &result.UserId, &result.IsDelete, &result.UpdateTime)
	if er != nil {
		return nil, er
	}
	return &result, nil
}

//warning: the condition must be suited for your own business by cxq
func ListDeptUser(b *DeptUser, pageNo, pageSize int) ([]DeptUser, error) {
	var offset = (pageNo - 1) * pageSize
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return nil, fmt.Errorf("数据库连接异常：%w", err)
	}
	sql := `select ID,DEPT_ID,USER_ID,IS_DELETE,UPDATE_TIME from DEPT_USER where ID=? and DEPT_ID=? and USER_ID=? and IS_DELETE=? and UPDATE_TIME=? limit ?,? `
	rows, er := db.Query(sql, b.Id, b.DeptId, b.UserId, b.IsDelete, b.UpdateTime, offset, pageSize)
	if er != nil {
		return nil, er
	}
	results := make([]DeptUser, 0)
	for rows.Next() {
		var c = DeptUser{}
		er = rows.Scan(&c.Id, &c.DeptId, &c.UserId, &c.IsDelete, &c.UpdateTime)
		if er != nil {
			return nil, er
		}
		results = append(results, c)
	}
	return results, nil
}

func CountDeptUser(b *DeptUser) (int64, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return 0, fmt.Errorf("数据库连接异常：%w", err)
	}
	sql := `select count(1) from DEPT_USER where ID=? and DEPT_ID=? and USER_ID=? and IS_DELETE=? and UPDATE_TIME=?`
	var count int64
	er := db.QueryRow(sql, b.Id, b.DeptId, b.UserId, b.IsDelete, b.UpdateTime).Scan(&count)
	if er != nil {
		return 0, fmt.Errorf("sql执行错误：%w", er)
	}
	return count, nil
}

func DeleteDeptUser(id int64) (int64, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return 0, fmt.Errorf("数据库连接异常：%w", err)
	}
	sql := `delete from DEPT_USER where ID=? `
	result, er := db.Exec(sql, id)
	if er != nil {
		return 0, fmt.Errorf("sql执行错误：%w", er)
	}
	return result.RowsAffected()
}
