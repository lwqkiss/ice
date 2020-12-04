package video

import (
	"fmt"
	"log"
	"lwq/dbaccess"
	"math"
	"strconv"
	"strings"
)

type Device struct {
	NodeType        int    `json:"nodeType" gorm:"-"`
	CameraFunctions string `json:"cameraFunctions,omitempty" gorm:"column:cameraFunctions"`
	CameraType      string `json:"cameraType,omitempty"  gorm:"column:cameraType"`
	//Category        string `json:"category" gorm:"column:category"`
	ChannelID   string `json:"channelId" gorm:"column:channelId"`
	ChannelSeq  int64  `json:"channelSeq,omitempty" gorm:"column:channelSeq"`
	ChannelType string `json:"channelType" gorm:"column:channelType"`
	DeviceCode  string `json:"deviceCode,omitempty" gorm:"column:deviceCode"`
	//DeviceID        string `json:"deviceId" gorm:"column:deviceId"`
	GpsX string `json:"gpsX" gorm:"column:gpsX"`
	GpsY string `json:"gpsY" gorm:"column:gpsY"`
	Icon string `json:"icon,omitempty" gorm:"column:icon"`
	ID   string `json:"id" gorm:"column:id"`
	//	IntelliFlag     int64  `json:"intelliFlag" gorm:"column:intelliFlag"`
	IntelliState int64 `json:"intelliState,omitempty" gorm:"column:intelliState"`
	//	IP              string `json:"ip" gorm:"column:ip"`
	IsParent bool `json:"isParent,omitempty" gorm:"column:isParent"`
	//	Manufacturer    string `json:"manufacturer" gorm:"column:manufacturer"`
	Name string `json:"name" gorm:"column:name"`
	//NodeType        int64  `json:"nodeType" gorm:"column:nodeType"`
	Online   string `json:"online,omitempty" gorm:"column:online"`
	OrgCode  string `json:"orgCode,omitempty" gorm:"column:orgCode"`
	OrgType  string `json:"orgType,omitempty" gorm:"column:orgType"`
	PaasID   string `json:"paasId,omitempty" gorm:"column:paasId"`
	ParentID string `json:"parentId,omitempty" gorm:"column:parentId"`
	Sort     int64  `json:"sort,omitempty" gorm:"column:sort"`
	Status   int64  `json:"status,omitempty" gorm:"column:status"`
	//SubType         string `json:"subType" gorm:"column:subType"`
	UnitType string `json:"unitType,omitempty" gorm:"column:unitType"`
	NodeCode string `json:"nodecode" gorm:"column:NODE_CODE"`
	GeoInfo  string `json:"geoinfo" gorm:"column:GEO_INFO"`
	//OrgName  string `json:"orgname" gorm:"column:ORG_NAME"`
	Address string `json:"address" gorm:"column:Address"`
}

func (d Device) TableName() string {
	return "POLICE_VIDEO"
}

func ListVideo() ([]*Device, error) {
	res := make([]*Device, 0)
	db := dbaccess.OpenGorm()
	err := db.Table("POLICE_VIDEO").Find(&res).Error
	return res, err

}

func UpdteVideo(b *Device) {
	db := dbaccess.OpenGorm()
	err := db.Table("POLICE_VIDEO").Where("id = ?", b.ID).Update(b).Error
	if err != nil {
		fmt.Print(err.Error())
	}
}

func StartTask() {
	list, err := ListVideo()
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	DealPoint()
	for i, _ := range list {
		ll := list[i]
		ff := GeoToPoint(ll.GeoInfo)
		/*da := egis.GetLonAndLat("南通市" + ll.Address)
		fmt.Println(ll.ID, ll.Name, ll.Address)
		fmt.Println(da.Result.Level, da.Result.Location.Lng, da.Result.Location.Lat)*/
		if CheckIsInTzq(ff[0], ff[1]) {
			ll.NodeCode = "T7004"
			fmt.Println("更新", ll.Name)
			UpdteVideo(ll)
		}
	}
}

func GeoToPoint(str string) []float64 {
	defer func() {
		if err := recover(); err != nil {
			log.Println(str)
		}
	}()
	// str = "POINT(1 1)"
	arr := []float64{}
	left := strings.Index(str, "(")
	right := strings.Index(str, ")")
	numbers := str[left+1 : right]
	strarr := strings.Split(numbers, " ")

	for _, e := range strarr {
		num, e := strconv.ParseFloat(strings.Trim(e, " "), 64)
		if e != nil {
			log.Println(e)
		}
		arr = append(arr, num)
	}
	return arr
}

var pointArray [][2]float64

func DealPoint() {

	points := []float64{
		31.95388063,
		120.86368067,
		31.92411291,
		120.87402857,
		31.93361713,
		120.91410608,
		31.95861980,
		120.89948719,
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

//检查是否在通州区
func CheckIsInTzq(lon, lat float64) bool {
	var iSum int
	var iCount int
	var dLon1, dLon2, dLat1, dLat2, dLon float64
	//fmt.Println(len(pointArray))
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
		//以下语句判断A点是否在边的两端点的水平平行线之间，在则可能有交点，开始判断交点是否在左射线上
		if ((lat >= dLat1) && (lat < dLat2)) || ((lat >= dLat2) && (lat < dLat1)) {
			if math.Abs(dLat1-dLat2) > 0 {
				//得到 A点向左射线与边的交点的x坐标：
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
