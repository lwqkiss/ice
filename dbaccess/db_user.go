package dbaccess

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

/**
 * @author @codenai.com
 * @date 2020/2/12
 */
// User 救援队伍用户类
type User struct {
	Id                  int     `json:"id" form:"id" gorm:"column:ID"`
	UserName            string  `json:"username" form:"username" gorm:"column:USERNAME"` //用户名称
	Password            string  `json:"password" form:"password" gorm:"column:PASSWORD"` //密码
	TEAMNAME            string  `json:"teamname" form:"teamname" gorm:"column:TEAM_NAME"`
	MUS                 string  `json:"mus" form:"mus" gorm:"column:MUS"`                                                   // MUS 主管单位
	ABILITYTYPE         string  `json:"abilitytype" form:"abilitytype" gorm:"column:ABILITY_TYPE"`                          // ABILITYTYPE 能力类型
	TEAMTYPE            string  `json:"teamtype" form:"teamtype" gorm:"column:TEAM_TYPE"`                                   // TEAMTYPE 队伍性质
	PERSONNUMBER        int     `json:"personnumber" form:"personnumber" gorm:"column:PERSON_NUMBER"`                       // PERSONNUMBER 队伍人数
	SOF                 string  `json:"sof" form:"sof" gorm:"column:SOF"`                                                   // SOF 经费来源
	TEAMADDRESS         string  `json:"teamaddress" form:"teamaddress" gorm:"column:TEAM_ADDRESS"`                          // TEAMADDRESS 办公地址
	ISORGANIZED         string  `json:"isorganized" form:"isorganized" gorm:"column:IS_ORGANIZED"`                          // ISORGANIZED 是否有组织架构
	TEAMEQUIPMENTUNIT   string  `json:"teamequipmentunit" form:"teamequipmentunit" gorm:"column:TEAM_EQUIPMENT_UNIT"`       // TEAMEQUIPMENTUNIT 装备情况 单位
	TEAMEQUIPMENTPERSON string  `json:"teamequipmentperson" form:"teamequipmentperson" gorm:"column:TEAM_EQUIPMENT_PERSON"` // TEAMEQUIPMENTPERSON 装备情况 个人
	GOODS               string  `json:"goods" form:"goods" gorm:"column:GOODS"`                                             // GOODS 物资情况
	TEAMLEADERNAME      string  `json:"teamleadername" form:"teamleadername" gorm:"column:TEAM_LEADER_NAME"`                // TEAMLEADERNAME 救援队长姓名
	TEAMPHONE           string  `json:"teamphone" form:"teamphone" gorm:"column:TEAM_PHONE"`                                // TEAMPHONE 救援队值班电话
	XZQHID              int     `json:"xzqhid" form:"xzqhid" gorm:"column:XZQH_ID"`                                         // XZQHID 行政区划id
	XZQHNAME            string  `json:"xzqhname" form:"xzqhname" gorm:"column:XZQH_NAME"`
	LON                 float64 `json:"lon" form:"lon" gorm:"column:LON"` // LON 经度
	LAT                 float64 `json:"lat" form:"lat" gorm:"column:LAT"` // LAT 纬度
	IsDelete            int     `json:"isdelete" form:"isdelete" gorm:"cloumn:IS_DELETE"`
	HasXfc              int     `json:"hasxfc" form:"hasxfc" gorm:"column:HAS_XFC"`
	TOTAL_ASSETS        string  `json:"totalassets" form:"totalassets" gorm:"column:TOTAL_ASSETS"`
	ESTABLISH_TIME      string  `json:"establishtime" from:"establishtime" gorm:"column:ESTABLISH_TIME"`
	BACK_COL            string  `json:"backcol" form:"backcol" gorm:"column:BACK_COL"`
	LastWarningTime     string  `json:"lastwarningtime" form:"lastwarningtime" gorm:"column:LAST_WARNING_TIME"`
	GeoInfo             string  `json:"geo_info" form:"geo_info" gorm:"column:GEO_INFO"`
}

func (b *User) TableName() string {
	return "EMERGENCY_TEAM"
}

func CheckUser(u User) (User, error) {
	db := OpenGorm()
	user := &User{}
	sql1 := `select * from EMERGENCY_TEAM where USERNAME = ? and IS_DELETE = 0`
	err := db.Raw(sql1, u.UserName).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return *user, fmt.Errorf("用户名错误")
	}
	//fmt.Println(*user.Password,*u.Password)
	if u.Password == "" {
		return *user, fmt.Errorf("密码不能为空")
	}
	if u.Password == user.Password {
		return *user, nil
	} else {
		return *user, fmt.Errorf("密码错误")
	}

}
func AddTeam(u *User) error {
	db := OpenGorm()
	u.XZQHID, _ = GetAdminIdByName(strings.TrimSpace(u.XZQHNAME))
	if u.UserName != "" {
		if !CheckUserName(u.UserName, u.Id) {
			return errors.New("用户名已存在 " + u.UserName)
		}
	}

	return db.Create(u).Error
}
func UpdateUser(u User) error {
	db := OpenGorm()
	if u.UserName != "" {
		//这个username重复
		if !CheckUserName(u.UserName, u.Id) {
			return errors.New("用户名已存在 " + u.UserName)
		}
	}
	return db.Model(User{}).Where("ID=?", u.Id).Update(u).Error
}

func GetUserById(id int) (User, error) {
	db := OpenGorm()
	user := User{}
	err := db.Where("ID = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func ListUser(pageNo, pageSize int, key string) ([]*User, error) {
	var offset = (pageNo - 1) * pageSize
	db := OpenGorm()
	userlist := []*User{}
	qryStr := ""
	if key != "" {
		qryStr += " and CONCAT_WS(',',TEAM_NAME,GOODS,ABILITY_TYPE,TEAM_ADDRESS,XZQH_NAME,TEAM_EQUIPMENT_UNIT,TEAM_EQUIPMENT_PERSON) REGEXP '" + key + "'"
	}
	err := db.Where("IS_DELETE = 0 " + qryStr).Offset(offset).Limit(pageSize).Find(&userlist).Error
	return userlist, err
}

func ListUserAll() ([]*User, error) {

	db := OpenGorm()
	userlist := []*User{}
	qryStr := ""

	err := db.Where("IS_DELETE = 0 " + qryStr).Find(&userlist).Error
	return userlist, err
}

func CountUser(b *User, key string) (int64, error) {
	db := OpenGorm()
	var count int64
	qryStr := ""
	if key != "" {
		qryStr += " and CONCAT_WS(',',TEAM_NAME,GOODS,ABILITY_TYPE,TEAM_ADDRESS,XZQH_NAME,TEAM_EQUIPMENT_UNIT,TEAM_EQUIPMENT_PERSON) REGEXP '" + key + "'"
	}
	db.Table("EMERGENCY_TEAM").Where("IS_DELETE=0" + qryStr).Count(&count)
	return count, nil
}

func DelTeam(id int) error {
	db := OpenGorm()
	err := db.Table("EMERGENCY_TEAM").Where("ID=?", id).Update("IS_DELETE", 1).Error
	return err
}

func CheckUserName(username string, id int) bool {
	db := OpenGorm()
	var count int64
	err := db.Table("EMERGENCY_TEAM").Where("IS_DELETE=0 and USERNAME = ? and ID != ?", username, id).Count(&count).Error
	if err != nil {
		return false
	}

	if count > 0 {
		return false
	} else {
		return true
	}
}

//导入队伍
func ImportTeam(u []*User) error {
	db := OpenGorm()
	for i, _ := range u {
		//得到行政区id
		u[i].XZQHID, _ = GetAdminIdByName(strings.TrimSpace(u[i].XZQHNAME))
		tmp, err := GetTeamByName(u[i].TEAMNAME)
		fmt.Println(tmp)
		//找不到队伍名称一样的，那么新增
		if err != nil && err == gorm.ErrRecordNotFound {
			err = db.Table("EMERGENCY_TEAM").Create(u[i]).Error

			if err != nil {

				return err
			}
			continue
		} else if err != nil {

			return err
		}

		//找到名称一样的，而且没有错误，那么更新
		err = db.Table("EMERGENCY_TEAM").Where("ID = ? ", tmp.Id).Update(u[i]).Error
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

	}

	return nil
}

func GetTeamByName(teamname string) (*User, error) {
	db := OpenGorm()
	user := &User{}
	err := db.Where("TEAM_NAME = ? and IS_DELETE = 0", teamname).First(user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetTeamByNameOrPhone(teamName, phone string) (User, error) {
	db := OpenGorm()
	user := User{}
	if phone != "" {
		err := db.Where("(TEAM_NAME like ? or TEAM_PHONE = ?) and IS_DELETE = 0", "%"+teamName+"%", phone).First(&user).Error
		if err != nil {
			return user, err
		}
	}

	err := db.Where("(TEAM_NAME like ? ) and IS_DELETE = 0", "%"+teamName+"%").First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
