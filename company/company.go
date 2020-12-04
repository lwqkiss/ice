package company

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"lwq/dbaccess"
	"lwq/egis"
	"lwq/ship"
	"lwq/utils"
	"math"
	"net/http"
	"strings"
	"time"
)

type PlanEnterprise struct {
	EnterpriseId        int     `json:"enterpriseid" form:"enterpriseid" gorm:"column:ENTERPRISE_ID"`
	EnterpriseName      string  `json:"enterprisename" form:"enterprisename" gorm:"column:ENTERPRISE_NAME"`
	LegalRepresentative string  `json:"legalrepresentative" form:"legalrepresentative" gorm:"column:LEGAL_REPRESENTATIVE"`
	IndustryLargeClass  string  `json:"industrylargeclass" form:"industrylargeclass" gorm:"column:INDUSTRY_LARGE_CLASS"`    // da类
	IndustryMediumClass string  `json:"industrymediumclass" from:"industrymediumclass" gorm:"column:INDUSTRY_MEDIUM_CLASS"` // z类
	IndustrySmallClass  string  `json:"industrysmallclass" form:"industrysmallclass" gorm:"column:INDUSTRY_SMALL_CLASS"`
	WorkPhone           string  `json:"workphone" form:"workphone" gorm:"column:WORK_PHONE"`
	PlanSituation       string  `json:"plansituation" form:"plansituation" gorm:"column:PLAN_SITUATION"`
	Lon                 float64 `json:"lon" form:"lon" gorm:"column:LON"`
	Lat                 float64 `json:"lat" form:"lat" gorm:"column:LAT"`
	Address             string  `json:"address" form:"address" gorm:"column:Address"`
	Code                string  `json:"code" form:"code" gorm:"column:Code"`
	PageNo              int     `json:"pageNo" form:"pageNo" gorm:"-"`
	PageSize            int     `json:"pageSize" form:"pageSize" gorm:"-"`
}

func (e *PlanEnterprise) TableName() string {
	return "T06024_单位概况"
}

func ListPlanEnterprise() ([]*PlanEnterprise, error) {
	bis := make([]*PlanEnterprise, 0)
	var extense = ""
	err := dbaccess.OpenGorm().Raw(` select 单位名称 ENTERPRISE_NAME,
法定代表人 LEGAL_REPRESENTATIVE,
监管主行业大类 INDUSTRY_LARGE_CLASS,
监管主行业小类 INDUSTRY_SMALL_CLASS,
单位电话 WORK_PHONE,
经度 LON,
纬度 LAT,
通讯地址 Address,
单位代码 Code
from T06024_单位概况` + extense).Find(&bis).Error
	return bis, err
}

func DealCompanyLonAndLat() {
	data, err := ListPlanEnterprise()
	fmt.Println(err)
	for i := range data {
		lon := data[i].Lon
		lat := data[i].Lat
		// fmt.Println(data[i].EnterpriseName, data[i].Address)
		if (lon == 121 && lat == 31) || lon == 0 || lat == 0 || lon < 119 || lat <= 30 || (lon == 120.89078395928195 && lat == 31.98307056920128) {
			ss := ""
			if data[i].Address != "" {
				add := []rune(data[i].Address)
				ss = string(add[0:3])
			}
			/*if !strings.Contains(data[i].Address, "南通") {
				data[i].Address = "南通市" + data[i].Address
			}*/
			data[i].Address = ss + data[i].Address
			fmt.Println(data[i].Code, data[i].EnterpriseName, data[i].Address)
			da := egis.GetLonAndLat(ss + data[i].EnterpriseName)
			fmt.Println(da.Result.Level, da.Result.Location.Lng, da.Result.Location.Lat)
			if da.Result.Location.Lng != 0 && da.Result.Location.Lat != 0 {
				err := UpdateCompany(data[i].Code, da.Result.Location.Lng, da.Result.Location.Lat, "by addr")

				if err != nil {
					fmt.Println(err.Error())
				}
				time.Sleep(time.Millisecond * 100)
			}

		}

	}
}

func UpdateCompany(comId string, lon float64, lat float64, upflag string) error {
	db := dbaccess.OpenGorm()

	/*res := make([]*PlanEnterprise, 0)
	db.Find(&res)*/
	// err := db.Raw(`update T06024_单位概况 set 经度 = ? ,纬度 = ?,update_flag = ? where 单位代码 = ?`, lon, lat, comId, upflag).Debug().Error
	err := db.Table("T06024_单位概况").Where("单位代码 = ?", comId).Update("经度", lon).Update("纬度", lat).Update("update_flag", upflag).Error
	return err
}

type DangerCompanyDataList struct {
	Id              int     `json:"id" gorm:"column:ID"`
	Name            string  `json:"name" gorm:"column:COMPANY_NAME"`
	RegisterAddress string  `json:"registeraddress" gorm:"column:COMPANY_ADDRESS"`
	CompanyConatct  string  `json:"companyconatct" gorm:"column:COMPANY_CONTACT"`
	CompanyType     string  `json:"companytype" gorm:"column:COMPANY_TYPE"`
	Lon             float64 `json:"lon" gorm:"column:LON"`
	Lat             float64 `json:"lat" gorm:"column:LAT"`
	NodeCode        string  `json:"-" gorm:"column:NODE_CODE"`
	GeoInfo         string  `json:"-" gorm:"column:GEO_INFO"`
	CenterPoint     string  `json:"-" gorm:"column:CENTER_POINT"`
	XzqhId          int     `json:"-" gorm:"column:XZQH_ID"`
	XzqhName        string  `json:"-" gorm:"column:XZQH_NAME"`
	UpdateFlag      string  `json:"update_flag" gorm:"column:update_flag"`
	OLon            float64 `json:"olon" gorm:"column:O_LON"`
	OLat            float64 `json:"olat" gorm:"column:O_LAT"`
}

func GetDangerCompanyList() ([]*DangerCompanyDataList, error) {
	db := dbaccess.OpenGorm()
	res := make([]*DangerCompanyDataList, 0)
	err := db.Table("DANGER_COMPANY").Find(&res).Error
	if err != nil {
		logs.Error(err.Error())
	}
	return res, err
}

func UpdateDangerCompanyList(b *DangerCompanyDataList) error {
	db := dbaccess.OpenGorm()
	// 这样注册地址就不会更新
	b.RegisterAddress = ""
	return db.Table("DANGER_COMPANY").Where("id=?", b.Id).Update(b).Error
}

func DealDangerCompany() {
	data, err := GetDangerCompanyList()
	if err != nil {
		fmt.Println(err)
	}
	DealPoint()

	for i := range data {
		lon := data[i].Lon
		lat := data[i].Lat
		// fmt.Println(data[i].EnterpriseName, data[i].Address)
		if !CheckIsInGeo(lon, lat) {

			if !strings.Contains(data[i].RegisterAddress, "南通") {
				data[i].RegisterAddress = "南通市" + data[i].XzqhName + data[i].RegisterAddress
			}
			data[i].RegisterAddress = data[i].RegisterAddress
			fmt.Println(data[i].XzqhName, data[i].Name, data[i].RegisterAddress)
			da := egis.GetLonAndLat(data[i].RegisterAddress)
			// fmt.Println(data[i].Name, data[i].RegisterAddress)
			fmt.Println(da.Result.Level, da.Result.Location.Lng, da.Result.Location.Lat)
			if da.Result.Location.Lng != 0 && da.Result.Location.Lat != 0 {
				data[i].Lon = da.Result.Location.Lng
				data[i].Lat = da.Result.Location.Lat
				data[i].GeoInfo = utils.GenGeostr(data[i].Lon, data[i].Lat)
				data[i].CenterPoint = data[i].GeoInfo
				data[i].UpdateFlag = "byaddr"
				err := UpdateDangerCompanyList(data[i])

				if err != nil {
					fmt.Println(err.Error())
				}
				time.Sleep(time.Millisecond * 100)
			}

		}

	}
}

type CompanyDataList struct {
	Id              int     `json:"id" gorm:"column:id"`
	Name            string  `json:"name" gorm:"column:name"`
	RegisterAddress string  `json:"registeraddress" gorm:"column:registeredAddress"`
	Lon             float64 `json:"lon" gorm:"column:longitude"`
	Lat             float64 `json:"lat" gorm:"column:latitude"`
	NodeCode        string  `json:"-" gorm:"column:NODE_CODE"`
	GeoInfo         string  `json:"-" gorm:"column:GEO_INFO"`
	CenterPoint     string  `json:"-" gorm:"column:CENTER_POINT"`
	Level           string  `json:"-" gorm:"column:LEVEL"`
	XzqhId          int     `json:"-" gorm:"column:XZQH_ID"`
	XzqhName        string  `json:"-" gorm:"column:XZQH_NAME"`
	NodeCode2       string  `json:"-" gorm:"column:NODE_CODE2"`
	UpdateFlag      string  `json:"update_flag" gorm:"column:update_flag"`
	OLon            float64 `json:"olon" gorm:"column:O_LON"`
	OLat            float64 `json:"olat" gorm:"column:O_LAT"`
}

func GetCompanyList() ([]*CompanyDataList, error) {
	db := dbaccess.OpenGorm()
	res := make([]*CompanyDataList, 0)
	err := db.Table("COMPANY").Find(&res).Error
	if err != nil {
		logs.Error(err.Error())
	}
	return res, err
}

func UpdateCompanyList(b *CompanyDataList) error {
	db := dbaccess.OpenGorm()
	// 这样注册地址就不会更新
	b.RegisterAddress = ""
	return db.Table("COMPANY").Where("id=?", b.Id).Update(b).Error
}

func DealWHCompany() {
	data, err := GetCompanyList()
	if err != nil {
		fmt.Println(err)
	}
	DealPoint()
	for i := range data {
		lon := data[i].Lon
		lat := data[i].Lat
		fmt.Println(data[i].Name, data[i].Name, "南通市"+data[i].RegisterAddress+data[i].Name)
		da := egis.GetLonAndLat(data[i].Name)
		fmt.Println(data[i].Name, data[i].Name, "南通市"+data[i].RegisterAddress+data[i].Name)
		fmt.Println(da.Result.Level, da.Result.Location.Lng, da.Result.Location.Lat)
		fmt.Println(data[i].Name, data[i].RegisterAddress)
		if !CheckIsInGeo(lon, lat) || (lon == 120.89078395928195 && lat == 31.98307056920128) {
			/*ss := ""
			if data[i].Name != "" {
				add := []rune(data[i].Address)
				ss = string(add[0:3])
			}
			/*if !strings.Contains(data[i].Address, "南通") {
				data[i].Address = "南通市" + data[i].Address
			}*/
			/*data[i].Address = ss + data[i].Address
			fmt.Println(data[i].Code, data[i].EnterpriseName, data[i].Address)*/
			da := egis.GetLonAndLat(data[i].Name)
			fmt.Println(data[i].Name, data[i].Name, "南通市"+data[i].RegisterAddress+data[i].Name)
			fmt.Println(da.Result.Level, da.Result.Location.Lng, da.Result.Location.Lat)
			if da.Result.Location.Lng != 0 && da.Result.Location.Lat != 0 {
				data[i].Lon = da.Result.Location.Lng
				data[i].Lat = da.Result.Location.Lng
				data[i].GeoInfo = utils.GenGeostr(data[i].Lon, data[i].Lat)
				data[i].CenterPoint = data[i].GeoInfo
				data[i].UpdateFlag = "byaddr"
				err := UpdateCompanyList(data[i])

				if err != nil {
					fmt.Println(err.Error())
				}
				time.Sleep(time.Millisecond * 100)
			}

		}

	}
}

type TzqCompanyTownship struct {
	ID                       int     `gorm:"column:id" json:"id" form:"id"`
	Name                     string  `gorm:"column:name" json:"name" form:"name"`
	Representative           string  `gorm:"column:representative" json:"representative" form:"representative"`
	StaffMessage             string  `gorm:"column:staffMessage" json:"staffMessage" form:"staffMessage"`
	Employeesnum             string  `gorm:"column:employeesNum" json:"employees_num" form:"employees_num"`
	Fulltimesafetyproduction string  `gorm:"column:fullTimeSafetyProduction" json:"fullTsp" form:"fullTsp"`
	Ifsafetydirector         int     `gorm:"column:ifSafetyDirector" json:"if_safety_director" form:"if_safety_director"`
	Ifmanagementagency       int     `gorm:"column:ifManagementAgency" json:"if_management_agency" form:"if_management_agency"`
	Registeredaddress        string  `gorm:"column:registeredAddress" json:"registered_address" form:"registered_address"`
	Type                     string  `gorm:"column:type" json:"type" form:"type"`
	License                  string  `gorm:"column:license" json:"license" form:"license"`
	Contactperson            string  `gorm:"column:contactPerson" json:"contact_person" form:"contact_person"`
	Password                 string  `gorm:"column:password" json:"password" form:"password"`
	Initialpassword          string  `gorm:"column:initialPassword" json:"initial_password" form:"initial_password"`
	Phone                    string  `gorm:"column:phone" json:"phone" form:"phone"`
	Standardizationlevel     string  `gorm:"column:standardizationLevel" json:"standardization_level" form:"standardization_level"`
	Effectivedate            string  `gorm:"column:effectiveDate" json:"effective_date" form:"effective_date"`
	Safetyorganization       string  `gorm:"column:SafetyOrganization" json:"_safety_organization" form:"_safety_organization"`
	Lon                      float64 `gorm:"column:longitude" json:"longitude" form:"longitude"`
	Lat                      float64 `gorm:"column:latitude" json:"latitude" form:"latitude"`
	Zoomlevel                int     `gorm:"column:zoomLevel" json:"zoom_level" form:"zoom_level"`
	Administrationid         int     `gorm:"column:administrationID" json:"administrationid" form:"administrationid"`
	// CreatedAt                string `gorm:"column:created_at" json:"created_at" form:"created_at"`
	// UpdatedAt                string `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
	IsDelete    int    `gorm:"column:is_delete" json:"-" form:"-"`
	NodeCode    string `gorm:"column:NODE_CODE" json:"node_code" form:"node_code"`
	GeoInfo     string `gorm:"column:GEO_INFO" json:"geo_info" form:"geo_info"`
	CenterPoint string `gorm:"column:CENTER_POINT" json:"center_point" form:"center_point"`
	Level       string `gorm:"column:LEVEL" json:"level" form:"level"`
	XzqhID      int    `gorm:"column:XZQH_ID" json:"xzqh_id" form:"xzqh_id"`
	XzqhName    string `gorm:"column:XZQH_NAME" json:"xzqh_name" form:"xzqh_name"`
	// 好丑
	NodeCode1    string `gorm:"column:NODE_CODE1;default:null" json:"-" form:"-"`
	NodeCode2    string `gorm:"column:NODE_CODE2;default:null" json:"-" form:"-"`
	NodeCode3    string `gorm:"column:NODE_CODE3;default:null" json:"-" form:"-"`
	NodeCode4    string `gorm:"column:NODE_CODE4;default:null" json:"-" form:"-"`
	NodeCode5    string `gorm:"column:NODE_CODE5;default:null" json:"-" form:"-"`
	UpdateFlag   string `gorm:"column:update_flag" json:"update_flag" form:"update_flag"`
	DangersArray string `gorm:"column:dangers_array" json:"dangers_array" form:"dangers_array"`
	// RYXXV2       string `gorm:"column:ryxx_new" json:"ryxx_new" form:"ryxx_new"`
	// 新增的几个数量
	ZdwxyNum      int    `gorm:"column:zdwxy_num" json:"zdwxy_num" form:"zdwxy_num"`
	CyryNum       int    `gorm:"column:cyry_num" json:"cyry_num" form:"cyry_num"`
	ZzaqscglryNum int    `gorm:"column:zzaqscglry_num" json:"zzaqscglry_num" form:"zzaqscglry_num"`
	ZcaqgcsNum    int    `gorm:"column:zcaqgcs_num" json:"zcaqgcs_num" form:"zcaqgcs_num"`
	TzzydczrydNum int    `gorm:"column:tzzydczryd_num" json:"tzzydczryd_num" form:"tzzydczryd_num"`
	Aqscgljg      string `gorm:"column:aqscgljg" json:"aqscgljg" form:"aqscgljg"`
	AqscgljgJyfw  string `gorm:"column:aqscgljg_jyfw" json:"aqscgljg_jyfw" form:"aqscgljg_jyfw"`
	HangyeDl      string `gorm:"column:hangye_dl" json:"hangye_dl" form:"hangye_dl"`
	HangyeXl      string `gorm:"column:hangye_xl" json:"hangye_xl" form:"hangye_xl"`
	Townname      string `gorm:"column:townname" json:"townname" form:"townname"` // 新增字段

	// Certificate []TzqCompanyCertificate ` json:"Certificates" form:"Certificates"`
	PageNo   int    `json:"pageNo" form:"pageNo" gorm:"-"`
	PageSize int    `json:"pageSize" form:"pageSize" gorm:"-"`
	Ty       string `json:"ty" form:"ty" gorm:"-"`
	Value    string `json:"value" form:"value" gorm:"-"`
	// 查询条件新增 每个镇登陆只显示该镇下的企业
	Owner string `json:"owner" form:"owner" gorm:"-"`
}

// TableName 表名
func (o *TzqCompanyTownship) TableName() string {
	return "tzq_company_township"
}

func GetTzqCompanyList() ([]*TzqCompanyTownship, error) {
	db := dbaccess.OpenGorm()
	res := make([]*TzqCompanyTownship, 0)
	err := db.Table("tzq_company_township").Find(&res).Error
	if err != nil {
		logs.Error(err.Error())
	}
	return res, err
}

func UpdateTzqCompanyList(b *TzqCompanyTownship) error {
	db := dbaccess.OpenGorm()
	// 这样注册地址就不会更新
	b.Registeredaddress = ""
	return db.Table("tzq_company_township").Where("id=?", b.ID).Update(b).Error
}

func DealTzqCompany() {
	data, err := GetTzqCompanyList()
	if err != nil {
		fmt.Println(err)
	}
	DealPoint()
	for i := range data {
		lon := data[i].Lon
		lat := data[i].Lat

		fmt.Println(data[i].Name, data[i].Registeredaddress)
		if lon != 0 && !CheckIsInGeo(lon, lat) || (lon == 120.89078395928195 && lat == 31.98307056920128) {
			fmt.Println(data[i].Name, data[i].Registeredaddress, lon, lat)
			da := egis.GetLonAndLat(data[i].Name)
			fmt.Println(da.Result.Level, da.Result.Location.Lng, da.Result.Location.Lat)
			if da.Result.Location.Lng != 0 && da.Result.Location.Lat != 0 {
				data[i].Lon = da.Result.Location.Lng
				data[i].Lat = da.Result.Location.Lat
				data[i].GeoInfo = utils.GenGeostr(data[i].Lon, data[i].Lat)
				data[i].CenterPoint = data[i].GeoInfo
				data[i].UpdateFlag = "byaddr"
				err := UpdateTzqCompanyList(data[i])

				if err != nil {
					fmt.Println(err.Error())
				}
				// time.Sleep(time.Millisecond * 100)
			}

		}

	}
}

type CompanyTmp struct {
	Code int `json:"code"`
	Data struct {
		CountList []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			ParentID int    `json:"parentID"`
			Count    int    `json:"count"`
			Children []struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				ParentID int    `json:"parentID"`
				Count    int    `json:"count"`
			} `json:"children,omitempty"`
		} `json:"countList"`
		Count          int `json:"count"`
		SearchSiteList []struct {
			ID                  int         `json:"id"`
			RegulatorDepartment string      `json:"regulatorDepartment"`
			AdministrationID    int         `json:"administrationID"`
			CompanyID           int         `json:"companyID"`
			Longitude           float64     `json:"longitude"`
			Latitude            float64     `json:"latitude"`
			ZoomLevel           int         `json:"zoomLevel"`
			Type                string      `json:"type"`
			BusinessType        interface{} `json:"businessType,omitempty"`
			SiteHazard          []struct {
				ID            int    `json:"id"`
				AccidentGrade string `json:"accidentGrade"`
			} `json:"siteHazard"`
			Level []struct {
				Value string `json:"value"`
			} `json:"level"`
			Company struct {
				ID                   int         `json:"id"`
				Name                 string      `json:"name"`
				StandardizationLevel string      `json:"standardizationLevel"`
				StaffMessage         interface{} `json:"staffMessage"`
				Phone                string      `json:"phone"`
				Representative       string      `json:"representative"`
			} `json:"company"`
			Administration struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				ParentID int    `json:"parentID"`
			} `json:"administration"`
			SiteChemical []struct {
				ChemicalID    interface{} `json:"chemicalID"`
				SiteChemicals interface{} `json:"siteChemicals"`
			} `json:"siteChemical"`
		} `json:"searchSiteList"`
	} `json:"data"`
	Message string `json:"message"`
}

func getInfo(url string) *http.Request {
	// cookie1 := &http.Cookie{Name: "EGG_SESS", Value: "LwYhBKvCSoh9Bajtw3xWn2Udh8m1M3uKPbq_XwUzxEHpsDnPcBLWS-t1kev1M7F2", HttpOnly: true}
	// cookie2 := &http.Cookie{Name: "csrfToken", Value: "t5yJi5jem1Dd3jk6t25bRdLv", HttpOnly: true}
	req, err := http.NewRequest("GET", url, strings.NewReader("{}"))
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729)")
	return req
}

func GetAdminNameById(id int) string {
	var name string
	switch id {
	case 2:
		name = "南通市"
	case 3:
		name = "海门区"
	case 4:
		name = "如皋市"
	case 5:
		name = "海安市"
	case 6:
		name = "启东市"
	case 7:
		name = "如东县"
	case 8:
		name = "通州区"
	case 9:
		name = "崇川区"
	case 10:
		name = "港闸区"
	case 11:
		name = "开发区"
	case 12:
		name = "通州湾"
	case 13:
		name = "苏锡通"

	}

	return name

}

// 更新企业的红橙黄蓝等级和NODE_CODE      注意不是NODE_CODE2
func Company2MySQL() error {
	// Login()
	uri := "/api/web/statistics/company?administrationID=2"
	client := &http.Client{}
	fmt.Println(uri)
	req := getInfo("http://dev.codenai.com:4455" + uri)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	tmp := CompanyTmp{}
	err = json.Unmarshal(body, &tmp)
	// fmt.Println(string(body))
	// fmt.Println(err, tmp, req.URL)
	// SetDBAStr("root", "123456", "t.codenai.com", "3310", "ntemergency")
	db, _ := dbaccess.OpenDB()
	defer db.Close()
	sql := `UPDATE COMPANY set LEVEL = ?,XZQH_ID = ? , NODE_CODE = ?,XZQH_NAME =? where id = ?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	// fmt.Println(err)
	// count := 0
	for _, v := range tmp.Data.SearchSiteList {
		adnmid := 2
		if v.Administration.ParentID == 2 {
			adnmid = v.Administration.ID
		} else {
			adnmid = v.Administration.ParentID
		}
		adnmid = GetTrueXzqId(adnmid)
		adnm := GetAdminNameById(adnmid)
		node := ""
		fmt.Println(v.Company.Name, v.Company.ID, v.Level)
		if len(v.Level) > 0 {

			switch v.Level[0].Value {
			case "红":
				node = "T5901"
			case "橙":
				node = "T5902"
			case "黄":
				node = "T5903"
			case "蓝":
				node = "T5904"

			}
			_, err := stmt.Exec(v.Level[0].Value, adnmid, node, adnm, v.Company.ID)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			_, err := stmt.Exec("", adnmid, node, adnm, v.Company.ID)
			if err != nil {
				fmt.Println(err)
			}
		}

		/*fmt.Println(err)
		count++
		fmt.Println(count)
		fmt.Println(err)*/
		/*id, err := res.LastInsertId()
		fmt.Println(id)*/
	}
	return nil
}

func GetTrueXzqId(xzqid int) int {
	switch xzqid {
	case 10:
		return 9
	default:
		return xzqid

	}
}

var pointArray [][2]float64

func DealPoint() {

	points := []float64{
		32.60091000,
		120.20233000,
		32.59973000,
		120.20362000,
		32.59007000,
		120.20741000,
		32.58880000,
		120.21136000,
		32.58851000,
		120.21633000,
		32.58300000,
		120.22104000,
		32.58193000,
		120.22463000,
		32.58269000,
		120.22675000,
		32.58684000,
		120.23064000,
		32.58792000,
		120.23755000,
		32.59195000,
		120.24120000,
		32.59335000,
		120.24510000,
		32.59380000,
		120.24882000,
		32.59757000,
		120.25293000,
		32.59777000,
		120.25467000,
		32.59722000,
		120.25507000,
		32.58400000,
		120.25870000,
		32.57805000,
		120.26203000,
		32.57487000,
		120.26287000,
		32.57127000,
		120.26053000,
		32.56927000,
		120.26013000,
		32.56456000,
		120.26428000,
		32.56262000,
		120.26415000,
		32.56141000,
		120.26233000,
		32.54999000,
		120.26737000,
		32.54306000,
		120.26628000,
		32.52759000,
		120.26598000,
		32.51456000,
		120.26370000,
		32.51358000,
		120.26562000,
		32.50752000,
		120.26486000,
		32.50584000,
		120.26796000,
		32.50181000,
		120.27013000,
		32.50050000,
		120.26916000,
		32.49892000,
		120.25966000,
		32.49302000,
		120.25896000,
		32.48798000,
		120.26486000,
		32.48392000,
		120.26681000,
		32.47028000,
		120.27109000,
		32.46157000,
		120.27268000,
		32.44460000,
		120.27363000,
		32.43835000,
		120.27610000,
		32.43495000,
		120.27823000,
		32.42263000,
		120.27979000,
		32.41463000,
		120.28280000,
		32.40015000,
		120.28427000,
		32.39056000,
		120.28829000,
		32.38852000,
		120.28977000,
		32.38526000,
		120.28797000,
		32.38203000,
		120.28508000,
		32.37964000,
		120.28453000,
		32.37064000,
		120.28745000,
		32.36645000,
		120.29078000,
		32.36621000,
		120.29269000,
		32.37183000,
		120.29977000,
		32.37197000,
		120.30092000,
		32.37115000,
		120.30234000,
		32.36431000,
		120.30733000,
		32.36155000,
		120.31012000,
		32.36158000,
		120.31404000,
		32.36474000,
		120.31463000,
		32.36801000,
		120.32124000,
		32.36838000,
		120.32798000,
		32.37326000,
		120.33632000,
		32.37456000,
		120.34229000,
		32.37660000,
		120.34522000,
		32.37764000,
		120.34818000,
		32.37790000,
		120.35165000,
		32.36892000,
		120.35237000,
		32.36774000,
		120.35147000,
		32.36462000,
		120.35477000,
		32.35799000,
		120.35972000,
		32.35516000,
		120.35353000,
		32.35346000,
		120.35182000,
		32.35208000,
		120.35082000,
		32.34765000,
		120.35019000,
		32.34387000,
		120.34680000,
		32.34004000,
		120.34621000,
		32.33583000,
		120.35291000,
		32.33302000,
		120.35388000,
		32.32493000,
		120.35429000,
		32.31679000,
		120.36217000,
		32.30967000,
		120.36537000,
		32.30882000,
		120.36523000,
		32.30446000,
		120.35231000,
		32.28319000,
		120.35287000,
		32.27866000,
		120.35241000,
		32.27715000,
		120.35025000,
		32.26796000,
		120.34648000,
		32.25599000,
		120.35088000,
		32.25464000,
		120.34959000,
		32.25442000,
		120.34731000,
		32.24785000,
		120.34796000,
		32.24394000,
		120.35154000,
		32.23872000,
		120.35705000,
		32.23962000,
		120.36140000,
		32.23782000,
		120.36461000,
		32.22677000,
		120.37004000,
		32.21764000,
		120.37182000,
		32.21575000,
		120.37091000,
		32.21430000,
		120.36684000,
		32.21154000,
		120.36705000,
		32.20704000,
		120.36558000,
		32.20332000,
		120.36300000,
		32.19336000,
		120.35803000,
		32.19170000,
		120.35558000,
		32.18207000,
		120.35460000,
		32.18069000,
		120.35243000,
		32.17922000,
		120.34651000,
		32.17333000,
		120.34606000,
		32.17098000,
		120.34466000,
		32.17007000,
		120.34322000,
		32.15055000,
		120.34743000,
		32.15085000,
		120.34997000,
		32.14950000,
		120.35105000,
		32.13853000,
		120.35565000,
		32.13129000,
		120.35617000,
		32.13118000,
		120.36172000,
		32.12805000,
		120.37609000,
		32.12937000,
		120.37633000,
		32.12967000,
		120.38791000,
		32.12410000,
		120.42322000,
		32.12403000,
		120.43536000,
		32.11870000,
		120.44566000,
		32.10423000,
		120.48325000,
		32.09915000,
		120.48678000,
		32.09238000,
		120.51133000,
		32.09143000,
		120.51244000,
		32.08622000,
		120.51360000,
		32.08337000,
		120.51540000,
		32.06597000,
		120.51856000,
		32.06247000,
		120.52089000,
		32.05929000,
		120.52395000,
		32.05650000,
		120.52949000,
		32.05570000,
		120.54480000,
		32.05442000,
		120.54734000,
		32.05038000,
		120.55053000,
		32.04463000,
		120.55012000,
		32.04111000,
		120.54797000,
		32.03282000,
		120.55912000,
		32.02860000,
		120.55783000,
		32.01422000,
		120.54907000,
		32.00993000,
		120.55804000,
		32.00365000,
		120.57600000,
		32.00310000,
		120.59529000,
		32.00362000,
		120.63432000,
		32.00312000,
		120.64463000,
		32.00524000,
		120.67672000,
		32.00680000,
		120.71610000,
		32.00601000,
		120.76938000,
		32.00478000,
		120.77847000,
		32.00326000,
		120.78233000,
		31.99795000,
		120.78827000,
		31.88857000,
		120.85324000,
		31.87297000,
		120.86010000,
		31.83463000,
		120.88359000,
		31.81968000,
		120.89507000,
		31.80873000,
		120.91628000,
		31.78292000,
		121.00710000,
		31.78349000,
		121.01736000,
		31.78307000,
		121.06068000,
		31.77701000,
		121.07689000,
		31.76269000,
		121.10148000,
		31.77539000,
		121.12125000,
		31.79664000,
		121.15729000,
		31.83070000,
		121.20349000,
		31.84307000,
		121.22581000,
		31.85593000,
		121.25447000,
		31.86666000,
		121.28538000,
		31.86753000,
		121.29235000,
		31.86695000,
		121.31421000,
		31.86219000,
		121.33013000,
		31.85185000,
		121.35191000,
		31.84087000,
		121.36804000,
		31.81170000,
		121.40271000,
		31.79672000,
		121.41389000,
		31.78601000,
		121.41898000,
		31.77039000,
		121.43336000,
		31.76214000,
		121.44747000,
		31.75343000,
		121.47371000,
		31.75234000,
		121.48706000,
		31.74762000,
		121.51276000,
		31.70509000,
		121.59769000,
		31.70413000,
		121.59778000,
		31.70338000,
		121.59966000,
		31.70661000,
		121.60299000,
		31.70008000,
		121.62189000,
		31.69924000,
		121.62920000,
		31.68798000,
		121.67237000,
		31.66659000,
		121.72804000,
		31.66510000,
		121.73632000,
		31.65661000,
		121.75591000,
		31.63691000,
		121.81676000,
		31.61907000,
		121.86804000,
		31.61871000,
		121.88935000,
		31.61549000,
		121.90462000,
		31.61435000,
		121.99916000,
		31.68926000,
		121.99894000,
		31.81983000,
		121.99707000,
		31.88051000,
		121.99152000,
		31.95316000,
		121.98172000,
		31.99155000,
		121.95504000,
		32.04688000,
		121.88675000,
		32.07626000,
		121.84891000,
		32.10020000,
		121.81972000,
		32.13278000,
		121.71079000,
		32.13706000,
		121.66950000,
		32.14021000,
		121.62489000,
		32.15833000,
		121.56013000,
		32.16147000,
		121.55106000,
		32.16270000,
		121.51419000,
		32.16236000,
		121.50734000,
		32.16625000,
		121.49762000,
		32.18589000,
		121.49058000,
		32.21514000,
		121.48546000,
		32.21606000,
		121.50012000,
		32.20522000,
		121.51303000,
		32.19431000,
		121.52922000,
		32.18875000,
		121.54037000,
		32.18498000,
		121.55068000,
		32.18200000,
		121.55998000,
		32.18166000,
		121.56721000,
		32.18270000,
		121.57153000,
		32.18745000,
		121.57424000,
		32.21121000,
		121.57547000,
		32.22844000,
		121.57524000,
		32.23253000,
		121.57114000,
		32.23515000,
		121.56469000,
		32.24216000,
		121.53991000,
		32.24450000,
		121.52605000,
		32.24490000,
		121.50016000,
		32.24705000,
		121.48941000,
		32.25009000,
		121.48193000,
		32.35094000,
		121.47678000,
		32.41258000,
		121.46518000,
		32.45141000,
		121.45210000,
		32.52529000,
		121.44074000,
		32.53149000,
		121.43797000,
		32.53561000,
		121.43374000,
		32.53756000,
		121.42775000,
		32.53495000,
		121.25007000,
		32.53610000,
		121.18334000,
		32.56644000,
		121.11201000,
		32.59396000,
		121.06646000,
		32.63560000,
		121.01093000,
		32.64115000,
		121.00524000,
		32.64595000,
		121.00352000,
		32.65126000,
		120.99703000,
		32.65435000,
		120.99505000,
		32.66066000,
		120.99498000,
		32.63318000,
		120.89969000,
		32.63162000,
		120.89583000,
		32.63089000,
		120.86266000,
		32.61930000,
		120.86280000,
		32.61855000,
		120.85017000,
		32.60608000,
		120.85057000,
		32.60584000,
		120.85225000,
		32.61164000,
		120.85742000,
		32.60962000,
		120.86231000,
		32.60954000,
		120.86415000,
		32.60722000,
		120.86485000,
		32.60457000,
		120.85252000,
		32.60491000,
		120.84854000,
		32.60677000,
		120.84344000,
		32.60674000,
		120.83727000,
		32.59440000,
		120.83147000,
		32.59185000,
		120.82362000,
		32.58400000,
		120.80812000,
		32.57866000,
		120.79440000,
		32.57831000,
		120.79137000,
		32.58130000,
		120.79297000,
		32.59283000,
		120.79475000,
		32.59710000,
		120.79513000,
		32.59912000,
		120.79469000,
		32.60058000,
		120.79348000,
		32.60398000,
		120.75994000,
		32.60386000,
		120.75448000,
		32.59830000,
		120.74100000,
		32.59530000,
		120.73842000,
		32.58586000,
		120.73938000,
		32.58514000,
		120.73841000,
		32.58520000,
		120.73288000,
		32.59083000,
		120.73098000,
		32.59144000,
		120.73010000,
		32.59163000,
		120.72104000,
		32.59973000,
		120.72030000,
		32.60090000,
		120.71163000,
		32.60405000,
		120.70844000,
		32.60467000,
		120.70525000,
		32.60078000,
		120.69360000,
		32.59605000,
		120.68258000,
		32.59139000,
		120.67573000,
		32.58758000,
		120.67158000,
		32.58097000,
		120.66933000,
		32.57757000,
		120.66526000,
		32.57602000,
		120.66273000,
		32.57539000,
		120.66011000,
		32.57543000,
		120.65651000,
		32.57625000,
		120.65414000,
		32.57833000,
		120.65240000,
		32.58174000,
		120.65334000,
		32.58494000,
		120.64928000,
		32.58948000,
		120.64992000,
		32.59098000,
		120.65080000,
		32.59441000,
		120.65609000,
		32.59825000,
		120.65481000,
		32.60344000,
		120.65693000,
		32.61318000,
		120.64399000,
		32.61345000,
		120.64224000,
		32.61210000,
		120.64109000,
		32.61150000,
		120.63904000,
		32.61196000,
		120.63643000,
		32.61258000,
		120.63452000,
		32.61636000,
		120.63340000,
		32.61714000,
		120.63251000,
		32.61751000,
		120.62979000,
		32.61305000,
		120.62563000,
		32.61110000,
		120.62492000,
		32.60651000,
		120.61886000,
		32.60601000,
		120.61775000,
		32.60673000,
		120.61570000,
		32.60774000,
		120.61545000,
		32.62596000,
		120.62187000,
		32.62858000,
		120.62065000,
		32.63458000,
		120.61258000,
		32.63560000,
		120.61273000,
		32.63636000,
		120.61549000,
		32.63975000,
		120.61480000,
		32.64298000,
		120.60822000,
		32.64634000,
		120.60545000,
		32.65001000,
		120.60048000,
		32.64435000,
		120.59208000,
		32.63965000,
		120.59052000,
		32.63772000,
		120.58692000,
		32.63744000,
		120.58434000,
		32.63499000,
		120.57896000,
		32.63543000,
		120.57365000,
		32.63698000,
		120.56902000,
		32.63522000,
		120.56128000,
		32.63588000,
		120.55949000,
		32.63578000,
		120.54991000,
		32.63225000,
		120.54456000,
		32.63621000,
		120.52794000,
		32.63753000,
		120.51328000,
		32.64078000,
		120.50684000,
		32.64100000,
		120.50438000,
		32.63963000,
		120.49587000,
		32.62431000,
		120.49736000,
		32.62501000,
		120.49023000,
		32.63124000,
		120.48470000,
		32.63184000,
		120.48344000,
		32.63213000,
		120.47992000,
		32.63057000,
		120.47628000,
		32.62864000,
		120.47559000,
		32.63226000,
		120.46860000,
		32.63238000,
		120.46714000,
		32.63226000,
		120.46519000,
		32.63036000,
		120.46154000,
		32.63078000,
		120.45434000,
		32.62989000,
		120.44788000,
		32.62910000,
		120.44609000,
		32.62928000,
		120.44313000,
		32.63910000,
		120.44321000,
		32.64177000,
		120.44063000,
		32.64231000,
		120.44182000,
		32.64709000,
		120.43697000,
		32.67545000,
		120.43206000,
		32.67668000,
		120.43023000,
		32.67751000,
		120.41941000,
		32.67942000,
		120.41578000,
		32.68240000,
		120.41513000,
		32.68371000,
		120.41333000,
		32.68656000,
		120.41290000,
		32.68909000,
		120.41158000,
		32.68871000,
		120.40481000,
		32.68981000,
		120.40060000,
		32.69385000,
		120.39054000,
		32.70723000,
		120.39251000,
		32.70931000,
		120.38666000,
		32.70852000,
		120.38414000,
		32.71114000,
		120.37893000,
		32.71088000,
		120.37615000,
		32.70489000,
		120.36671000,
		32.70587000,
		120.35846000,
		32.70547000,
		120.35468000,
		32.70409000,
		120.35006000,
		32.70513000,
		120.34537000,
		32.70480000,
		120.34094000,
		32.70296000,
		120.33269000,
		32.70091000,
		120.33331000,
		32.69733000,
		120.33267000,
		32.69784000,
		120.32727000,
		32.69947000,
		120.32680000,
		32.70014000,
		120.32539000,
		32.69930000,
		120.31951000,
		32.70181000,
		120.31945000,
		32.70440000,
		120.31557000,
		32.70162000,
		120.31027000,
		32.68758000,
		120.29267000,
		32.68077000,
		120.28082000,
		32.67930000,
		120.27628000,
		32.67983000,
		120.25198000,
		32.67771000,
		120.23479000,
		32.67588000,
		120.22570000,
		32.65838000,
		120.22507000,
		32.65616000,
		120.22880000,
		32.65366000,
		120.22950000,
		32.64994000,
		120.22917000,
		32.64839000,
		120.22211000,
		32.64693000,
		120.22160000,
		32.64219000,
		120.22189000,
		32.63896000,
		120.22290000,
		32.63433000,
		120.22219000,
		32.63426000,
		120.21943000,
		32.62826000,
		120.21873000,
		32.62801000,
		120.21782000,
		32.62868000,
		120.21673000,
		32.62563000,
		120.21354000,
		32.60818000,
		120.20942000,
		32.60561000,
		120.21350000,
		32.60272000,
		120.21358000,
		32.60287000,
		120.21062000,
		32.60448000,
		120.20629000,
		32.60417000,
		120.20466000,
		32.60350000,
		120.20428000,
		32.60091000,
		120.20233000,
	}

	for i := range points {
		if i%2 == 0 {
			tmp := [2]float64{}
			tmp[0] = points[i+1]
			tmp[1] = points[i]

			pointArray = append(pointArray, tmp)
		}
	}
}

func CheckIsInGeo(lon, lat float64) bool {
	var iSum int
	var iCount int
	var dLon1, dLon2, dLat1, dLat2, dLon float64
	// fmt.Println(len(pointArray))
	if len(pointArray) < 3 {

		return false
	}
	iCount = len(pointArray)
	for i := 0; i < iCount; i++ {
		if i == iCount-1 {
			dLon1 = pointArray[i][0]
			dLat1 = pointArray[i][1]
			dLon2 = pointArray[0][0]
			dLat2 = pointArray[0][1]
		} else {
			dLon1 = pointArray[i][0]
			dLat1 = pointArray[i][1]
			dLon2 = pointArray[i+1][0]
			dLat2 = pointArray[i+1][1]
		}
		// 以下语句判断A点是否在边的两端点的水平平行线之间，在则可能有交点，开始判断交点是否在左射线上
		if ((lat >= dLat1) && (lat < dLat2)) || ((lat >= dLat2) && (lat < dLat1)) {
			if math.Abs(dLat1-dLat2) > 0 {
				// 得到 A点向左射线与边的交点的x坐标：
				dLon = dLon1 - ((dLon1-dLon2)*(dLat1-lat))/(dLat1-dLat2)
				if dLon < lon {
					iSum++
				}

			}
		}
	}
	if iSum%2 != 0 {
		return true
	} else {
		return false
	}

}

func StartTask() {
	DealPoint()
	list, err := ship.GetShipTrans()
	if err != nil {
		fmt.Println(err)
	}
	for i := range list {
		ll := list[i]
		lon := ll.Lon
		lat := ll.Lat
		if !CheckIsInGeo(lon, lat) || (lon == 120.89078395928195 && lat == 31.98307056920128) {
			fmt.Println(ll.WaitPort, ll.WaitLoc, lon, lat)
			da := egis.GetLonAndLat(ll.WaitPort + ll.WaitLoc)
			fmt.Println(da.Result.Level, da.Result.Location.Lng, da.Result.Location.Lat)
			if da.Result.Location.Lng != 0 && da.Result.Location.Lat != 0 {
				ll.Lon = da.Result.Location.Lng
				ll.Lat = da.Result.Location.Lat
				ll.GeoInfo = utils.GenGeostr(ll.Lon, ll.Lat)

				err := ship.UpdateShip(ll)

				if err != nil {
					fmt.Println(err.Error())
				}
				// time.Sleep(time.Millisecond * 100)
			}

		}

	}
}
