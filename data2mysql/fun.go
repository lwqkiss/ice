package data2mysql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lwq/dbaccess"
	"lwq/utils"
	"strings"
)

type HouseCommunity struct {
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

func (b *HouseCommunity) TableName() string {
	return "HOUSE_COMMUNITY"
}
func AddHouse(b *HouseCommunity) error {
	db := dbaccess.OpenGorm()

	return db.Create(b).Error
}

func Deal0808() {
	//学校导入到数据库
	by, err := ioutil.ReadFile("nature.json")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(by))

	school := School{}

	json.Unmarshal(by, &school)
	fmt.Println(school)

	for i, _ := range school.Data {
		code, err := dbaccess.GetClassMapInfo(school.Data[i].Type)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for ii, _ := range school.Data[i].Pois {

			vv := school.Data[i].Pois[ii]

			val := vv[1].(string)

			tmp := HouseCommunity{

				SchoolName:  vv[0].(string),
				Latitude:    0,
				Longtude:    0,
				GeoInfo:     "",
				Position:    "",
				NodeCode:    code.Nodecode,
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
			if code.ID < "523" {

				//tt := &Protection{
				//
				//	ProtectionName:    tmp.SchoolName,
				//	ProtectionAddress: tmp.Position,
				//	GeoInfo:           tmp.GeoInfo,
				//	NodeCode:          tmp.NodeCode,
				//	CenterPoint:       tmp.GeoInfo,
				//	XzqhId:            tmp.XzqhId,
				//	XzqhName:          tmp.XzqhName,
				//}
				//err := AddProtection(tt)
				//if err != nil {
				//	fmt.Println(err)
				//}
			} else if code.ID >= "523" && code.ID < "534" {

				//err := AddHouse(&tmp)
				//if err != nil {
				//	fmt.Println(err)
				//}
			} else {
				tt := &Party{

					DzName:      tmp.SchoolName,
					DzAddress:   tmp.Position,
					DzPhone:     tmp.Phone,
					DzPoint:     tmp.GeoInfo,
					GeoInfo:     tmp.GeoInfo,
					NodeCode:    tmp.NodeCode,
					CenterPoint: tmp.GeoInfo,
					XzqhId:      tmp.XzqhId,
					XzqhName:    tmp.XzqhName,
				}
				err := AddParty(tt)
				if err != nil {
					fmt.Println(err)
				}
			}

		}
	}
}

type Protection struct {
	ID                int    `gorm:"column:ID" json:"id" form:"id"`
	ProtectionName    string `gorm:"column:PROTECTION_NAME" json:"protection_name" form:"protection_name"`
	ProtectionAddress string `gorm:"column:PROTECTION_ADDRESS" json:"protection_address" form:"protection_address"`
	GeoInfo           string `gorm:"column:GEO_INFO" json:"geo_info" form:"geo_info"`
	NodeCode          string `gorm:"column:NODE_CODE" json:"node_code" form:"node_code"`
	CenterPoint       string `gorm:"column:CENTER_POINT" json:"center_point" form:"center_point"`
	XzqhId            int    `gorm:"column:XZQH_ID" json:"xzqh_id" form:"xzqh_id"`
	XzqhName          string `gorm:"column:XZQH_NAME" json:"xzqh_name" form:"xzqh_name"`
}

func (b *Protection) TableName() string {
	return "PROTECTION"
}
func AddProtection(b *Protection) error {
	db := dbaccess.OpenGorm()

	return db.Create(b).Error
}

type Party struct {
	ID          int    `gorm:"column:ID" json:"id" form:"id"`
	DzName      string `gorm:"column:DZ_NAME" json:"dz_name" form:"dz_name"`
	DzAddress   string `gorm:"column:DZ_ADDRESS" json:"dz_address" form:"dz_address"`
	DzPhone     string `gorm:"column:DZ_PHONE" json:"dz_phone" form:"dz_phone"`
	DzPoint     string `gorm:"column:DZ_POINT" json:"dz_point" form:"dz_point"`
	GeoInfo     string `gorm:"column:GEO_INFO" json:"geo_info" form:"geo_info"`
	NodeCode    string `gorm:"column:NODE_CODE" json:"node_code" form:"node_code"`
	CenterPoint string `gorm:"column:CENTER_POINT" json:"center_point" form:"center_point"`
	XzqhId      int    `gorm:"column:XZQH_ID" json:"xzqh_id" form:"xzqh_id"`
	XzqhName    string `gorm:"column:XZQH_NAME" json:"xzqh_name" form:"xzqh_name"`
}

func (b *Party) TableName() string {
	return "PARTY"
}
func AddParty(b *Party) error {
	db := dbaccess.OpenGorm()

	return db.Create(b).Error
}
