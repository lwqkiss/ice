package dbaccess

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"strings"
	"unicode/utf8"
)

type MaterialNameStandardization struct {
	Id           int     `json:"id" form:"id" gorm:"column:ID;primary_key"`
	StandardName string  `json:"standardname" form:"standardname" gorm:"column:STANDARD_NAME"`
	Alias        string  `json:"alias" form:"alias" gorm:"column:ALIAS"`
	Demand       float64 `json:"demand" form:"demand" gorm:"column:DEMAND"`
	MeasureUnit  string  `json:"measureunit" form:"measureunit" gorm:"column:MEASURE_UNIT;comment:单位"`
	IsDelete     int     `json:"isdelete" form:"isdelete" gorm:"column:IS_DELETE"`
	PageNo       int     `json:"pageNo" form:"pageNo" gorm:"-"`
	PageSize     int     `json:"pageSize" form:"pageSize" gorm:"-"`
	//Category     []*MaterialNameKindRel `json:"category" gorm:"-"`
	CategoryName string `json:"categoryname" form:"categoryname" gorm:"column:CATEGORY_NAME"`
	CategoryId   int    `json:"categoryid" form:"categoryid" gorm:"column:CATEGORY_ID"`
	StandardNo   string `json:"standardno" form:"standardno" gorm:"column:STANDARD_NO"`
	Key          string `json:"key" form:"key" gorm:"-"`
}

func (b *MaterialNameStandardization) TableName() string {
	return "MATERIAL_NAME_STANDARDIZATION"
}
func AddMaterialNameStandardization(b *MaterialNameStandardization) error {
	db := OpenGorm()
	if b.StandardName == "" {
		return fmt.Errorf("%s", "标准化名称不能为空")
	}
	if CheckStandardName(b) {
		return fmt.Errorf("标准化名称 [%s] 重复", b.StandardName)
	}
	/*b.Alias = strings.ReplaceAll(b.Alias, "，", ",")
	aliasSlice := strings.Split(b.Alias, ",")
	b.Alias = ""
	sliceOk := make([]string, 0)
	for _, v := range aliasSlice {
		tmp := v
		if tmp != "" {
			sliceOk = append(sliceOk, tmp)
			b.Alias += tmp + ","
		}
	}*/
	/*if flag, repeat := CheckAlias(sliceOk, b.Id); flag {
		return fmt.Errorf("别名 [%s] 重复", repeat)
	}

	b.Alias = TrimLastChar(b.Alias)*/

	err := db.Create(b).Error
	if err != nil {
		fmt.Println(err)
		logs.Error(err)
		return err
	}

	return nil

}

func UpdateMaterialNameStandardization(b *MaterialNameStandardization) error {

	if b.StandardName == "" {
		return fmt.Errorf("%s", "标准化名称不能为空")
	}
	if CheckStandardName(b) {
		return fmt.Errorf("标准化名称 [%s] 重复", b.StandardName)
	}

	//处理别名
	b.Alias = strings.ReplaceAll(b.Alias, "，", ",")
	aliasSlice := strings.Split(b.Alias, ",")
	b.Alias = ""
	sliceOk := make([]string, 0)
	for _, v := range aliasSlice {
		tmp := v
		if tmp != "" {
			sliceOk = append(sliceOk, tmp)
			b.Alias += tmp + ","
		}
	}
	/*if flag, repeat := CheckAlias(sliceOk, b.Id); flag {
		return fmt.Errorf("别名 [%s] 重复", repeat)
	}*/
	//去掉最后一个逗号
	b.Alias = TrimLastChar(b.Alias)

	db := OpenGorm()

	//删除别名的关系

	tx := db.Begin()

	err := db.Model(MaterialNameStandardization{}).Where("ID=?", b.Id).Update(b).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func UpdateMaterialNameStandardizationV2(b *MaterialNameStandardization) error {
	db := OpenGorm()
	return db.Model(MaterialNameStandardization{}).Where("ID=?", b.Id).Update(b).Error
}

func GetMaterialNameStandardizationById(id int) (*MaterialNameStandardization, error) {
	db := OpenGorm()
	b := &MaterialNameStandardization{}
	err := db.Where("ID=?", id).First(b).Error
	return b, err
}

func GetMaterialNameStandardizationByName(name string) (*MaterialNameStandardization, error) {
	db := OpenGorm()
	b := &MaterialNameStandardization{}
	err := db.Where("IS_DELETE = 0 and STANDARD_NAME = ?", name).First(b).Error
	return b, err
}

//这里，别名应该不用处理
func ListMaterialNameStandardization(b *MaterialNameStandardization) ([]*MaterialNameStandardization, error) {
	db := OpenGorm()
	bis := make([]*MaterialNameStandardization, 0)
	if b.Key == "" {
		err := db.Where("IS_DELETE = 0").Offset((b.PageNo - 1) * b.PageSize).Limit(b.PageSize).Find(&bis).Error
		return bis, err
	}
	err := db.Where("IS_DELETE = 0 and CONCAT_WS(',',STANDARD_NAME,ALIAS) REGEXP ?", b.Key).Offset((b.PageNo - 1) * b.PageSize).Limit(b.PageSize).Find(&bis).Error
	return bis, err
}

func CountMaterialNameStandardization(b *MaterialNameStandardization) (int64, error) {
	db := OpenGorm()
	var count int64
	if b.Key == "" {
		err := db.Model(MaterialNameStandardization{}).Where("IS_DELETE = 0").Count(&count).Error
		return count, err
	}
	err := db.Model(MaterialNameStandardization{}).Where("IS_DELETE = 0 and CONCAT_WS(',',STANDARD_NAME,ALIAS) REGEXP ?", b.Key).Count(&count).Error
	return count, err
}

/*func DeleteMaterialNameStandardization(id int) error {
	db := OpenGorm()
	//删除和别名有关的
	db.Table("MATERIAL_NAME_ALIAS").Where("NAME_ID = ?", id).Delete(&MaterialNameAlias{})
	return db.Model(&MaterialNameStandardization{}).Where("ID=?", id).Update("IS_DELETE", 1).Error
}*/

/*//检查别名
func CheckAlias(alias []string, nameid int) (bool, string) {
	db := OpenGorm()

	for _, v := range alias {
		tmp := v
		var count int64
		err := db.Model(MaterialNameAlias{}).Where("ALIAS = ? and IS_DELETE = 0 and NAME_ID!=?", tmp, nameid).Count(&count).Error
		if err != nil {
			return true, tmp
		}
		if count > 0 {
			return true, tmp
		} else {
			return false, ""
		}
	}
	return false, ""
}*/

func CheckStandardName(b *MaterialNameStandardization) bool {
	db := OpenGorm()
	var count int64

	err := db.Model(MaterialNameStandardization{}).Where("IS_DELETE = 0 and STANDARD_NAME = ? and ID != ?", b.StandardName, b.Id).Count(&count).Error
	if err != nil {
		logs.Error(err.Error())
		fmt.Println(err.Error())
		return true
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

//根据物资名称返回单位
func GetMeasureUnitByStandardName(name string) (string, error) {
	db := OpenGorm()
	bis := MaterialNameStandardization{}
	err := db.Where("STANDARD_NAME = ?", name).First(&bis).Error
	return bis.MeasureUnit, err
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
