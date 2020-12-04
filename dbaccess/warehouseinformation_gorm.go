package dbaccess

import (
	"fmt"
	"github.com/astaxie/beego/logs"

	"strconv"
)

type LoginType struct {
	Type string `json:"type" gorm:"-"`
}

const Dept string = "部门"
const Warehouse string = "库所"

type WarehouseInformation struct {
	Id            int     `json:"id" form:"id" gorm:"column:ID;primary_key"`
	Name          string  `json:"equipmentlibraryname" form:"equipmentlibraryname" gorm:"column:NAME"`
	Location      string  `json:"address" form:"address" gorm:"column:LOCATION"`
	Longtitude    float64 `json:"lon" form:"lon" gorm:"column:LONGITUDE"`
	Latitude      float64 `json:"lat" form:"lat" gorm:"column:LATITUDE"`
	ContactPerson string  `json:"contactperson" form:"contactperson" gorm:"column:CONTACT_PERSON"`
	Phone         string  `json:"contactnumber" form:"contactnumber" gorm:"column:PHONE"`
	XzqhName      string  `json:"xzqhname" form:"xzqhname" gorm:"column:XZQH_NAME"`
	XzqhId        int     `json:"xzqhid" form:"xzqhid" gorm:"column:XZQH_ID"`
	NodeCode      string  `json:"-" gorm:"column:NODE_CODE"`
	GeoInfo       string  `json:"-" gorm:"column:GEO_INFO"`
	CenterPoint   string  `json:"-" gorm:"column:CENTER_POINT"`
	ManageUnit    string  `json:"manageunit" form:"manageunit" gorm:"column:MANAGE_UNIT"`        // 维护责任部门
	ManageUnitID  int     `json:"manageunitid" form:"manageunitid" gorm:"column:MANAGE_UNIT_ID"` // 责任部门id
	LoginCode     string  `json:"logincode" form:"logincode" gorm:"column:LOGIN_CODE"`
	IsDelete      int     `json:"isdelete" form:"isdelete" gorm:"column:IS_DELETE"`
	PageNo        int     `json:"pageNo" form:"pageNo" gorm:"-"`
	PageSize      int     `json:"pageSize" form:"pageSize" gorm:"-"`
	Key           string  `json:"key" form:"key" gorm:"-"`
	LoginType
}

func (b *WarehouseInformation) TableName() string {
	return "WAREHOUSE_INFORMATION"
}
func AddWarehouseInformation(b []*WarehouseInformation) error {
	db := OpenGorm()
	for _, v := range b {
		tmp := v
		tmp.GeoInfo = "(" + strconv.FormatFloat(tmp.Longtitude, 'f', -1, 64) + " " + strconv.FormatFloat(tmp.Latitude, 'f', -1, 64) + ")"
		tmp.CenterPoint = tmp.GeoInfo
		tmp.XzqhId, _ = GetAdminIdByName(tmp.XzqhName)
		err := db.Create(tmp).Error
		if err != nil {
			logs.Error(err.Error())
			return err
		}
	}
	return nil
}
func AddWarehouseInformationByTeam(b *WarehouseInformation) error {
	db := OpenGorm()

	return db.Create(b).Error

}

func UpdateWarehouseInformation(b *WarehouseInformation) error {
	db := OpenGorm()
	b.XzqhId, _ = GetAdminIdByName(b.XzqhName)
	return db.Model(WarehouseInformation{}).Where("ID=?", b.Id).Update(b).Error
}

func GetWarehouseInformationById(id int) (*WarehouseInformation, error) {
	db := OpenGorm()
	b := &WarehouseInformation{}
	err := db.Where("ID=?", id).First(b).Error
	return b, err
}

func GetWarehouseInformationByName(name string) (*WarehouseInformation, error) {
	db := OpenGorm()
	b := &WarehouseInformation{}
	err := db.Where("NAME = ? and IS_DELETE =0", name).First(b).Error
	return b, err
}

func ListWarehouseInformation(b *WarehouseInformation) ([]*WarehouseInformation, error) {
	db := OpenGorm()
	bis := make([]*WarehouseInformation, 0)

	if b.Key == "" {
		err := db.Where("IS_DELETE = 0").Offset((b.PageNo - 1) * b.PageSize).Limit(b.PageSize).Find(&bis).Error
		for i, _ := range bis {

			tmp, _ := GetMaterialDepartmentById(bis[i].ManageUnitID)
			bis[i].ManageUnit = tmp.DepartmentName
			// 把验证码去掉
			bis[i].LoginCode = ""

		}
		return bis, err
	}

	err := db.Where("IS_DELETE = 0 and CONCAT_WS(',',NAME,LOCATION,CONTACT_PERSON,PHONE,XZQH_NAME) REGEXP ?", b.Key).Offset((b.PageNo - 1) * b.PageSize).Limit(b.PageSize).Find(&bis).Error
	return bis, err
}

func CountWarehouseInformation(b *WarehouseInformation) (int64, error) {
	db := OpenGorm()
	var count int64
	if b.Key == "" {
		err := db.Model(WarehouseInformation{}).Where("IS_DELETE = 0").Count(&count).Error
		return count, err
	}
	err := db.Model(WarehouseInformation{}).Where("IS_DELETE = 0 and CONCAT_WS(',',NAME,LOCATION,CONTACT_PERSON,PHONE,XZQH_NAME) REGEXP ?", b.Key).Count(&count).Error
	return count, err
}

func DeleteWarehouseInformation(id int64) error {
	db := OpenGorm()
	return db.Model(&WarehouseInformation{}).Where("ID=?", id).Update("IS_DELETE", 1).Error
}

func GetAdminIdByName(str string) (int, error) {
	var id int
	db, err := OpenDB()
	defer db.Close()
	if err != nil {
		return 0, fmt.Errorf("数据库连接异常%w", err)
	}
	er := db.QueryRow("select id from administration where name =?", str).Scan(&id)
	return id, er

}

func GetWarehouseInformationByPhone(phone string) (*WarehouseInformation, error) {
	db := OpenGorm()
	b := &WarehouseInformation{}
	err := db.Where("PHONE=?", phone).First(b).Error
	return b, err
}

func CheckDeptHasWarehouse(deptid int) bool {
	db := OpenGorm()
	var count int64

	err := db.Model(WarehouseInformation{}).Where("IS_DELETE = 0 and MANAGE_UNIT_ID = ?", deptid).Count(&count).Error
	if err != nil {
		return true
	}
	if count > 0 {
		return true
	} else {
		return false
	}

}

func ListWarehouseInformationAll() ([]*WarehouseInformation, error) {
	db := OpenGorm()
	bis := make([]*WarehouseInformation, 0)

	err := db.Where("IS_DELETE = 0 and ID >6342 ").Find(&bis).Error
	/*for i, _ := range bis {

		tmp, _ := GetMaterialDepartmentById(bis[i].ManageUnitID)
		bis[i].ManageUnit = tmp.DepartmentName
		// 把验证码去掉
		bis[i].LoginCode = ""

	}*/
	return bis, err

}
