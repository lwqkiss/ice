package data2mysql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lwq/dbaccess"
	"lwq/utils"
	"strings"
)

type School struct {
	Class string `json:"class"`
	Data  []Data `json:"data"`
	Type  string `json:"type"`
}
type Data struct {
	Class string          `json:"class"`
	Pois  [][]interface{} `json:"pois"`
	Type  string          `json:"type"`
}

type SchoolTable struct {
	ID          int     `gorm:"column:ID" json:"id" form:"id"`
	SchoolName  string  `gorm:"column:SCHOOL_NAME" json:"school_name" form:"school_name"`
	Latitude    float64 `gorm:"column:LATITUDE" json:"latitude" form:"latitude"`
	Longtude    float64 `gorm:"column:LONGTUDE" json:"longtude" form:"longtude"`
	GeoInfo     string  `gorm:"column:GEO_INFO" json:"geo_info" form:"geo_info"`
	Position    string  `gorm:"column:POSITION" json:"position" form:"position"`
	NodeCode    string  `gorm:"column:NODE_CODE" json:"node_code" form:"node_code"`
	CenterPoint string  `gorm:"column:CENTER_POINT" json:"center_point" form:"center_point"`
	XzqhId      int     `gorm:"column:XZQH_ID" json:"xzqh_id" form:"xzqh_id"`
	XzqhName    string  `gorm:"column:XZQH_NAME" json:"xzqh_name" form:"xzqh_name"`
	Phone       string  `gorm:"column:PHONE" json:"phone" form:"phone"`
}

func (b *SchoolTable) TableName() string {
	return "SCHOOL"
}
func AddSchoo(b *SchoolTable) error {
	db := dbaccess.OpenGorm()

	return db.Create(b).Error
}

func DealSchoolData() {
	//学校导入到数据库
	by, err := ioutil.ReadFile("school1.json")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(by))

	school := School{}

	json.Unmarshal(by, &school)
	fmt.Println(school)

	for i, _ := range school.Data {
		code, _ := dbaccess.GetCodeByName(school.Data[i].Class)
		for ii, _ := range school.Data[i].Pois {

			vv := school.Data[i].Pois[ii]

			val := vv[1].(string)

			tmp := SchoolTable{

				SchoolName:  vv[0].(string),
				Latitude:    0,
				Longtude:    0,
				GeoInfo:     "",
				Position:    "",
				NodeCode:    code,
				CenterPoint: "",
				XzqhId:      0,
				XzqhName:    "",
			}
			tmp.Longtude, tmp.Latitude, tmp.GeoInfo = utils.DealPoint(val)
			tmp.CenterPoint = tmp.GeoInfo

			if cc, ok := vv[4].(string); ok {
				tmp.Position = cc
			} else {
				fmt.Println(school.Data[i].Pois[ii])
			}

			tt := vv[11].(string)
			if strings.Contains(tmp.Position, "经济技术开发区") && tt == "崇川区" {

				tt = "开发区"

				tmp.Position = "江苏省南通市" + tmp.Position

			} else {
				tmp.Position = "江苏省南通市" + tt + tmp.Position
			}
			tmp.XzqhName = tt
			tmp.XzqhId, tmp.XzqhName, _ = dbaccess.GetIdByAdministrationName(tt)
			if aa, ok := vv[5].(string); ok {
				slsl := strings.Split(aa, ";")
				if len(slsl) > 0 {
					tmp.Phone = slsl[0]
				}
			}

			err := AddSchoo(&tmp)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
