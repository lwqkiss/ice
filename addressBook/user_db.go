package addressBook

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

const (
	ACTIVE     = 0 //记录有效
	DELETED    = 1 //记录删除
	TIMEFORMAT = "2006-01-02 15:04:05"
)

type User struct {
	Id          int         `json:"id" form:"id" gorm:"column:ID;primary_key"`
	UserName    string      `json:"username" form:"username" gorm:"column:USER_NAME"`
	Mobilephone string      `json:"mobilephone" form:"mobilephone" gorm:"column:MOBILEPHONE"`
	Telphone    string      `json:"telphone" form:"telphone" gorm:"column:TELPHONE"`
	Address     string      `json:"address" form:"address" gorm:"column:ADDRESS"`
	IsDelete    int         `json:"isdelete" form:"isdelete" gorm:"column:IS_DELETE"`
	UpdateTime  string      `json:"updatetime" form:"updatetime" gorm:"column:UPDATE_TIME"`
	IsEmPerson  string      `json:"isemperson" form:"isemperson" gorm:"column:IS_EMERGENCY_PERSON"`
	WorkNum     string      `gorm:column:"WORK_NUM" json:"work_num"`
	CropNum     string      `gorm:column:"CROP_NUM" json:"crop_num"`
	DpUserRel   []*DeptUser `json:"dpuserrel" gorm:"-"`
	DpList      []int       `json:"dplist" gorm:"-"`
}

/*type DeptUser struct {
	Id  int  `json:"id" form:"id" gorm:"column:ID"`
	DeptId  int  `json:"deptid" form:"deptid" gorm:"column:DEPT_ID"`
	UserId  string  `json:"userid" form:"userid" gorm:"column:USER_ID"`
	IsDelete  string  `json:"isdelete" form:"isdelete" gorm:"column:IS_DELETE"`
	UpdateTime  string  `json:"updatetime" form:"updatetime" gorm:"column:UPDATE_TIME"`
}*/

func (u *User) TableName() string {
	return "USER"
}

func AddUser(b *User) (int, error) {
	db := OpenGorm()
	b.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	er := db.Table("USER").Create(b).Error
	if er != nil {
		return 0, er
	}

	for _, v := range b.DpList {
		tmp := DeptUser{UserId: b.Id, UpdateTime: time.Now().Format("2006-01-02 15:04:05"), DeptId: v}

		err := db.Table("DEPT_USER").Create(&tmp).Error
		if err != nil {
			return 0, fmt.Errorf("添加关联关系失败")
		}
	}
	return b.Id, er
}

func GetUserByDpId(id, pageNo, pageSize int) ([]*User, error) {
	var offset = (pageNo - 1) * pageSize
	db := OpenGorm()
	bis := make([]*User, 0)
	sql1 := `SELECT a.* FROM USER a,DEPT_USER where a.ID = DEPT_USER.USER_ID and DEPT_USER.DEPT_ID = ?
			and a.IS_DELETE =0 and DEPT_USER.IS_DELETE =0 limit ?,?`
	err := db.Raw(sql1, id, offset, pageSize).Find(&bis).Error

	for _, v := range bis {
		tmp := v
		tmprel := make([]*DeptUser, 0)
		db.Table("DEPT_USER").Where("USER_ID =? and IS_DELETE =0 and DEPT_ID = ?", tmp.Id, id).Find(&tmprel)
		tmp.DpUserRel = append(tmp.DpUserRel, tmprel...)

	}

	return bis, err
}

//remember: all of the fields will be updated ,although  you do not want to , be carefull
func UpdateUser(b *User) (int64, error) {
	db := OpenGorm()
	err := db.Model(&User{}).Where("ID = ?", b.Id).Update(b).Error
	if err != nil {
		return 0, err
	}
	for _, v := range b.DpUserRel {
		err = db.Table("DEPT_USER").Model(&DeptUser{}).Where("ID =?", v.Id).Update(v).Error
	}
	return db.RowsAffected, err
}

func GetUserById(id int) (*User, error) {
	db := OpenGorm()
	b := &User{}
	err := db.Where("ID=? and IS_DELETE = 0", id).First(b).Error
	return b, err
}

//warning: the condition must be suited for your own business by cxq
func ListUser(b *User, pageNo, pageSize int) ([]*User, error) {
	var offset = (pageNo - 1) * pageSize
	db := OpenGorm()

	results := make([]*User, 0)
	err := db.Where("IS_DELETE = 0").Offset(offset).Limit(pageSize).Find(&results).Error
	return results, err
}

func CountUser(b *User) (int64, error) {
	db, err := openDB()

	if err != nil {
		return 0, fmt.Errorf("数据库连接异常：%w", err)
	}
	defer db.Close()
	sql := `select count(1) from USER  IS_DELETE=0 `
	var count int64
	er := db.QueryRow(sql).Scan(&count)
	if er != nil {
		return 0, fmt.Errorf("sql执行错误：%w", er)
	}
	return count, nil
}

func DeleteUser(id int64) (int64, error) {
	db, err := openDB()

	if err != nil {
		return 0, fmt.Errorf("数据库连接异常：%w", err)
	}
	defer db.Close()
	sql := `update  USER set IS_DELETE=? where  ID=? `
	result, er := db.Exec(sql, DELETED, id)
	if er != nil {
		return 0, fmt.Errorf("sql执行错误：%w", er)
	}
	return result.RowsAffected()
}

func GetUserByDeptId(deptid int) (*User, error) {
	var user User
	err := OpenGorm().Raw(`
select USER.* from USER 
inner join DEPT_USER 
on USER.ID =DEPT_USER.USER_ID
inner join DEPARTMENT 
on DEPT_USER.DEPT_ID = DEPARTMENT.ID
where
USER.IS_DELETE=0
and DEPT_USER.IS_DELETE=0
and DEPARTMENT.IS_DELETE=0
and DEPT_USER.DEPT_ID =?`, deptid).First(&user).Error
	return &user, err
}

type GetUserByDeptNameAndUserNameVO struct {
	UserName string `json:"username" form:"username" `
	DeptName string `json:"deptname" form:"deptname"`
}

func GetUserByDeptNameAndUserName(vo *GetUserByDeptNameAndUserNameVO) (*User, error) {
	var user User
	err := OpenGorm().Raw(`
        select USER.* from USER 
inner join DEPT_USER 
on USER.ID =DEPT_USER.USER_ID
inner join DEPARTMENT 
on DEPT_USER.DEPT_ID = DEPARTMENT.ID
where
USER.IS_DELETE=0
and DEPT_USER.IS_DELETE=0
and DEPARTMENT.IS_DELETE=0
and DEPARTMENT.DEPT_NAME =?
and DEPARTMENT.DEPT_TYPE=?
and USER.USER_NAME=?`, vo.DeptName, vo.DeptName, vo.UserName).First(&user).Error
	return &user, err
}

func GetContentByDeptName(deptname string, responseid int) (interface{}, error) {
	var planContent []*struct {
		ContentID int    `json:"contentid" gorm:"column:CONTENT_ID"`
		DeptName  string `json:"deptname" gorm:"column:DEPT_NAME"`
	}
	err := OpenGorm().Raw(`
SELECT
	ORGANIZATION_COMMAND_SYSTEM_CONTENT.ID CONTENT_ID,
	DEPARTMENT.DEPT_NAME 
FROM
	EMERGENCY_RESPONSE
	INNER JOIN PLAN_CONTENT ON EMERGENCY_RESPONSE.PLAN_CONTENT_ID = PLAN_CONTENT.ID
	INNER JOIN ORGANIZATION_COMMAND_SYSTEM ON ORGANIZATION_COMMAND_SYSTEM.PLAN_CONTENT_ID = PLAN_CONTENT.Id
	INNER JOIN ORGANIZATION_COMMAND_SYSTEM_CONTENT ON ORGANIZATION_COMMAND_SYSTEM_CONTENT.ORGANIZATION_ID = ORGANIZATION_COMMAND_SYSTEM.ID
	INNER JOIN DEPARTMENT ON DEPARTMENT.ID = ORGANIZATION_COMMAND_SYSTEM_CONTENT.DEPT_ID 
WHERE
	PLAN_CONTENT.IS_DELETE = 0 
	AND ORGANIZATION_COMMAND_SYSTEM.IS_DELETE = 0 
	AND ORGANIZATION_COMMAND_SYSTEM_CONTENT.IS_DELETE = 0 
	AND DEPARTMENT.IS_DELETE = 0 
	AND DEPARTMENT.DEPT_TYPE = "dept" 
	AND ORGANIZATION_COMMAND_SYSTEM.ORGANIZATION_NAME = "成员单位" 
	AND DEPARTMENT.DEPT_NAME LIKE CONCAT( "%", ?, "%" ) 
	AND EMERGENCY_RESPONSE.ID = ?
`, deptname, responseid).Find(&planContent).Error
	return planContent, err
}

func DelUser(id int) error {
	db := OpenGorm()
	err := db.Model(&User{}).Where("ID = ?", id).Update("IS_DELETE", 1).Error
	if err != nil {
		return fmt.Errorf("删除失败 ")
	}
	err = db.Table("DEPT_USER").Model(&DeptUser{}).Where("USER_ID = ?", id).Update("IS_DELETE", 1).Error
	if err != nil {
		return fmt.Errorf("删除失败 ")
	}
	return nil

}

func GetUserInfobySysUserId(userid string) (map[string]interface{}, error) {
	rows, err := OpenGorm().Raw(`select USER.MOBILEPHONE,
SYS_USER.USER_NAME username,
SYS_USER.TYPES_OF typesof,
SYS_ROLE.ROLE_RNAME rolename
from SYS_USER 
inner join  USER on SYS_USER.USER_ID=USER.Id
inner join SYS_ROLE 
on SYS_ROLE.ID=SYS_USER.ROLE
where SYS_USER.ID=?`, userid).Rows()
	if err != nil {
		return nil, err
	}
	data, err := GetSqlResultMap(rows)
	if err != nil {
		return nil, err
	}
	if len(data) >= 1 {
		return data[0], nil
	}
	return nil, gorm.ErrRecordNotFound
}

func ReportCameraMikeStatus(userID, camera, mike string) error {
	return OpenGorm().Exec("update SYS_USER set HAS_CAMERA =?,HAS_MIKE=? where ID=?", camera, mike, userID).Error
}

type CameraMike struct {
	ID        string `json:"id" gorm:"column:ID"`
	UserName  string `json:"username" gorm:"column:USER_NAME"`
	HasCamera string `json:"hascamera" gorm:"column:HAS_CAMERA"`
	HasMike   string `json:"hasmike" gorm:"column:HAS_MIKE"`
}

func ListCameraMikeByUserIds(userids []string) ([]*CameraMike, error) {
	var result []*CameraMike
	err := OpenGorm().Raw(fmt.Sprintf(`select ID,USER_NAME,HAS_CAMERA,HAS_MIKE from SYS_USER where ID in %s`, "("+strings.Join(userids, ",")+")")).Find(&result).Error
	return result, err
}
