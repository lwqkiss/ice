package addressBook

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Department struct {
	Id         int           `json:"id" form:"id" gorm:"column:ID;primary_key"`
	DeptName   string        `json:"deptname" form:"deptname" gorm:"column:DEPT_NAME"`
	DeptType   string        `json:"depttype" form:"depttype" gorm:"column:DEPT_TYPE"`
	ParentId   int           `json:"parentid" form:"parentid" gorm:"column:PARENT_ID"`
	IsDelete   int           `json:"isdelete" form:"isdelete" gorm:"column:IS_DELETE"`
	UpdateTime string        `json:"updatetime" form:"updatetime" gorm:"column:UPDATE_TIME"`
	PageNo     int           `json:"pageNo" form:"pageNo" gorm:"-"`
	PageSize   int           `json:"pageSize" form:"pageSize" gorm:"-"`
	ChildRen   []*Department `json:"children" gorm:"-"`
}

func (dept *Department) TableName() string {
	return "DEPARTMENT"
}

const (
	JOB  string = "job"
	DEPT string = "dept"
)

func AddDepartment(b *Department) (int, error) {
	db := OpenGorm()
	b.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	err := db.Create(b).Error
	fmt.Println(b.Id)
	return b.Id, err
}

//remember: all of the fields will be updated ,although  you do not want to , be carefull
func UpdateDepartment(b *Department) (int64, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return 0, fmt.Errorf("数据库连接异常：%w", err)
	}
	sql := `update DEPARTMENT set DEPT_NAME=?,DEPT_TYPE=?,PARENT_ID=?,UPDATE_TIME=?  where ID=?`

	result, er := db.Exec(sql, b.DeptName, b.DeptType, b.ParentId, time.Now().Format(TIMEFORMAT), b.Id)
	if er != nil {
		return 0, fmt.Errorf("sql执行错误：%w", err)
	}
	return result.RowsAffected()
}

func GetDepartmentById(id int) (*Department, error) {
	var dept Department
	err := OpenGorm().Model(Department{}).Where("ID=?", id).First(&dept).Error
	return &dept, err
}

//warning: the condition must be suited for your own business by cxq
type OrgRes struct {
	Org     Department    `json:"org"`
	HasNext int           `json:"hasnext"`
	OrgList []*Department `json:"orglist"`
}

func ListDepartment(parentid int) ([]OrgRes, error) {

	db := OpenGorm()
	bis := make([]*Department, 0)
	res := make([]OrgRes, 0)
	err := db.Where("IS_DELETE = 0 and PARENT_ID = ?", parentid).Find(&bis).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, v := range bis {
		tmpvv := v
		tmplist := make([]*Department, 0)
		tmp := OrgRes{}
		err := db.Where("PARENT_ID = ? and IS_DELETE =0", tmpvv.Id).Find(&tmplist).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if len(tmplist) > 0 {
			tmp.HasNext = 1
		}
		tmp.Org = *tmpvv
		tmp.OrgList = append(tmp.OrgList, tmplist...)
		res = append(res, tmp)
	}
	return res, nil
}

func CountDepartment(b *Department) (int64, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return 0, fmt.Errorf("数据库连接异常：%w", err)
	}
	sql := `select count(1) from DEPARTMENT where ID=? and DEPT_NAME=? and DEPT_TYPE=? and PARENT_ID=? and IS_DELETE=? and UPDATE_TIME=?`
	var count int64
	er := db.QueryRow(sql, b.Id, b.DeptName, b.DeptType, b.ParentId, b.IsDelete, b.UpdateTime).Scan(&count)
	if er != nil {
		return 0, fmt.Errorf("sql执行错误：%w", er)
	}
	return count, nil
}

func DeleteDepartment(id int64) (int64, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return 0, fmt.Errorf("数据库连接异常：%w", err)
	}
	er3 := CheckDepartmentCanDelete(id)
	if er3 != nil {
		return 0, er3
	}
	sql := `update   DEPARTMENT set IS_DELETE=1 where ID=? `
	result, er := db.Exec(sql, id)
	if er != nil {
		return 0, fmt.Errorf("sql执行错误：%w", er)
	}
	return result.RowsAffected()
}

func CheckDepartmentCanDelete(id int64) error {
	db := OpenGorm()
	con1 := 0
	err := db.Table("DEPARTMENT").Where("PARENT_ID = ? and IS_DELETE = 0", id).Count(&con1).Error
	if con1 > 0 {
		return fmt.Errorf("请先删除父节点")
	}
	cou := 0
	err = db.Table("DEPT_USER").Where("DEPT_ID=? and IS_DELETE = 0", id).Count(&cou).Error

	if err != nil {
		fmt.Println(err)
	}
	if cou > 0 {
		return fmt.Errorf("有人员属于该节点，不能删除")
	}
	return nil
}

func GetParentDepartmentById(id int) (Department, error) {
	db := OpenGorm()
	dept := Department{}
	err := db.Raw("select DEPARTMENT.* from DEPARTMENT  inner join DEPARTMENT d2 on DEPARTMENT.ID=d2.PARENT_ID where d2.ID=?", id).First(&dept).Error
	return dept, err
}

func ListDepartmentsByParentId(parentId int) ([]*Department, error) {
	var depts = make([]*Department, 0)
	err := OpenGorm().Model(Department{}).Select("ID,DEPT_NAME,DEPT_TYPE,IS_DELETE").Where("PARENT_ID=? and IS_DELETE = 0", parentId).Find(&depts).Error
	return depts, err
}
func GetRootDept() (dept Department, err error) {
	err = OpenGorm().Model(Department{}).Where("PARENT_ID=?", -1).First(&dept).Error
	return
}

func ListDeptIdByParent(parentId int) []int {
	var deptIds []int
	rows, _ := OpenGorm().Raw("select ID from DEPARTMENT where PARENT_ID=? where IS_DELETE=0", parentId).Rows()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		deptIds = append(deptIds, id)
	}
	return deptIds
}

func ListAllDept() ([]*Department, error) {
	var depts []*Department
	err := OpenGorm().Select("ID,DEPT_NAME,DEPT_TYPE,PARENT_ID").Find(&depts).Error
	return depts, err
}

func AddEmergencyParticipants(planContentId int, responseId int) error {
	return OpenGorm().Exec(`INSERT INTO RESCUE_PARTICIPANTS ( CATEGORY, RESPONSE_ID, IS_DELETE ) SELECT
IF
	(
		d0.DEPT_TYPE = "job",
		d1.DEPT_NAME,
		d0.DEPT_NAME 
	) DEPT_NAME,
	`+strconv.Itoa(responseId)+`,
	0 
FROM
	ORGANIZATION_COMMAND_SYSTEM_CONTENT syscontent
	INNER JOIN ORGANIZATION_COMMAND_SYSTEM sys ON syscontent.ORGANIZATION_ID = sys.ID
	INNER JOIN DEPARTMENT d0 ON d0.ID = syscontent.DEPT_ID
	INNER JOIN DEPARTMENT d1 ON d1.ID = d0.PARENT_ID 
WHERE
	sys.PLAN_CONTENT_ID = ? 
	AND sys.IS_DELETE = 0 
	AND syscontent.IS_DELETE = 0 
	AND d0.IS_DELETE = 0 
	AND d1.IS_DELETE = 0 UNION
SELECT
IF
	(
		d2.DEPT_TYPE = "job",
		d3.DEPT_NAME,
		d2.DEPT_NAME 
	) DEPT_NAME,
	`+strconv.Itoa(responseId)+`,
	0 
FROM
	EMERGENCY_TEAM_CONTENT teamcontent
	INNER JOIN RESCUE_EMERGENCY_TEAM team ON teamcontent.TEAM_ID = team.ID
	INNER JOIN DEPARTMENT d2 ON d2.ID = teamcontent.DEPT_ID
	INNER JOIN DEPARTMENT d3 ON d3.ID = d2.PARENT_ID 
WHERE
	team.PLAN_CONTENT_ID = ? 
	AND team.IS_DELETE = 0 
	AND teamcontent.IS_DELETE = 0 
	AND d2.IS_DELETE = 0 
	AND d3.IS_DELETE = 00`, planContentId, planContentId).Error

}
