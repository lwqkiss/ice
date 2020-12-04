package main

import (
	"crypto/sha256"
	"fmt"
	"lwq/company"
	"lwq/dbaccess"
	"lwq/egis"
	"lwq/excelToCreateSQL"
	"lwq/utils"
	"net/http"
	"strings"
)

/**
 * @author @codenai.com
 * @date 2020/2/25
 */

/*type AAA struct {
	A string `json:"a"`
	B string `json:"b"`
}*/

type Total2Data struct {
	Tablecode string     `json:"tablecode"`
	Count     int        `json:"count"`
	IsLast    int        `json:"islast"`
	Datas     [][]string `json:"datas"`
}

func main() {
	// 创建 UUID v4
	/*u1:= uuid.Must(uuid.NewV4())
	println(`生成的UUID v4：`)
	println(u1.String())*/

	// 创建可以进行错误处理的 UUID v4
	/*u2 := uuid.NewV4()

		println(`生成的UUID v4：`)
		println(u2.String())

		// 解析 字符串 到 UUID
		u2, err2 := uuid.FromString(`6ba7b810-9dad-11d1-80b4-00c04fd430c8`)
		if err2 != nil {
			println(`解析 字符串 到 UUID 时出错`)
			panic(err2)
		}
		println(`解析 字符串 到 UUID 成功！解析到的 UUID 如下：`)
		println(u2.String())
		//ss := "bbb"
		cc := `""324234""`
		Ab := AAA{
			A: `21412
	  232232`,
			B: cc,
		}

		str, err := json.Marshal(Ab)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(str))
	*/
	// dbaccess.SetDBAStr("root", "123456", "101.133.168.208", "3306", "tzq")
	// db, er := dbaccess.OpenDB()
	// if er != nil {
	// 	fmt.Println(er.Error)
	// }
	// sql2 := `select ID from EMERGENCY_TEAM where ID <=4427`
	//
	// rows, err := db.Query(sql2)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// list := make([]int, 0)
	// for rows.Next() {
	// 	tmp := 0
	// 	rows.Scan(&tmp)
	// 	list = append(list, tmp)
	// }
	// fmt.Println(list)
	//
	// sql1 := `INSERT INTO TEAM_RELATION(TEAM_ID, TYPE_ID) VALUES (?, ?)`
	//
	// stmt, err := db.Prepare(sql1)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// for _, i := range list {
	// 	fmt.Println(i)
	// 	for j := 1; j <= 31; j++ {
	// 		_, err := stmt.Exec(i, j)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	}
	// }
	//
	// fmt.Println("23333")
	/*ss,_:=	hex.DecodeString("1a")
	  fmt.Println(string(ss))*/
	/*fi := &FireEnginesInfo{
		Id:          "45455",
		Name:        "",
		TeamName:    "",
		Desc:        "",
		Direct:      0,
		GpsStatus:   0,
		Speed:       0,
		TerminalNo:  "",
		State:       0,
		Type:        0,
		GpsTime:     "",
		Lon:         0,
		Lat:         0,
		PLon:        0,
		PLat:        0,
		Height:      0,
		Address:     "",
		Sim:         "",
		IsBus:       0,
		XzqhId:      0,
		XzqhName:    "",
		GeoInfo:     "",
		CenterPoint: "",
		NodeCode:    "",
		IsDelete:    0,
	}

	f2 := &FireEnginesInfo{
		Id:          "4324444444444",
		Name:        "",
		TeamName:    "",
		Desc:        "",
		Direct:      0,
		GpsStatus:   0,
		Speed:       0,
		TerminalNo:  "",
		State:       0,
		Type:        0,
		GpsTime:     "",
		Lon:         0,
		Lat:         0,
		PLon:        0,
		PLat:        0,
		Height:      0,
		Address:     "",
		Sim:         "",
		IsBus:       0,
		XzqhId:      0,
		XzqhName:    "",
		GeoInfo:     "",
		CenterPoint: "",
		NodeCode:    "",
		IsDelete:    0,
	}
	f2.Id = "fsdffsd"
	fmt.Println(fi.Id,f2.Id)
	logs.Error("232423",f2.Id)*/

	/*var total2data = Total2Data{
			Tablecode: "tab.Tablecode",
			Count:     0,
			IsLast:    0,
			Datas:     [][]string{{"eee","aaa"},{"bbb","ccc"},
			},
		}
	//fmt.Println(total2data)
		total2data.Datas = [][]string{}
		//fmt.Println(total2data)

		by ,err := ioutil.ReadFile("./ill_act_guide_20200409_total2_01.data.ready")
		if err != nil {
			fmt.Println(err)

		}
		//fmt.Println(string(by))
		//b, err := simplifiedchinese.GBK.NewDecoder().Bytes(by)

		fmt.Println("valid=", utf8.ValidString(string(by)))

		utf8Encoder := mahonia.NewEncoder("UTF-8")
		newStr:= utf8Encoder.ConvertString(string(by))
		fmt.Println(newStr)
		WriteFile(newStr,"ss.txt")*/
	/*s := make([]string, 0)
		fmt.Println(s[:len(s)])

		s1 := "0030"
	ss,_:= 	strconv.Atoi(s1)
	fmt.Println(ss)*/

	/*sql3 := `select SUPPLIES from MATERIAL_STORAGE where SUPPLIES like '%%' and SUPPLIES not in
	(select STANDARD_NAME from MATERIAL_NAME_STANDARDIZATION) group by SUPPLIES ORDER BY SUPPLIES desc;`
		rows, err := db.Query(sql3)
		if err != nil {
			fmt.Println(err)
		}
		count := 130
		for rows.Next() {
			supplies := ""
			rows.Scan(&supplies)
			if supplies != "" {
				b := &MaterialNameStandardization{

					StandardName: supplies,
					Alias:        "",

					CategoryName: "消防应急救援装备",
					CategoryId:   7,
					StandardNo:   fmt.Sprintf("201%05d", count),
				}
				err := AddMaterialNameStandardization(b)
				if err != nil {
					fmt.Println(err)
				}
			}
			count++
		}*/

	/*

	 */
	/*list, err := dbaccess.ListWarehouseInformationAll()
	if err != nil {
		fmt.Println(err)
	}
	bc := gobaidumap.NewBaiduMapClient("eGYLzw6HoOIEKLTB2iwmzFGRTVhGvUZj")

	for i, _ := range list {
		cc := list[i].Name
		if string(cc[0]) == "市" {
			cc = "南通" + cc

		}
		if list[i].Longtitude <= 119 || list[i].Latitude <= 30 {
			addressToGEO, err := bc.GetGeoViaAddress(cc)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(cc, addressToGEO.Result.Location.Lng, addressToGEO.Result.Location.Lat)
				list[i].Longtitude = addressToGEO.Result.Location.Lng
				list[i].Latitude = addressToGEO.Result.Location.Lat
				list[i].GeoInfo = utils.GenGeostr(addressToGEO.Result.Location.Lng, addressToGEO.Result.Location.Lat)
				list[i].CenterPoint = list[i].GeoInfo
				dbaccess.UpdateWarehouseInformation(list[i])
			}
			if addressToGEO.Result.Location.Lng < 119 || addressToGEO.Result.Location.Lat <= 31 {
				fmt.Println(cc, "不是南通的")
			}
		}
	*/
	// time.Sleep(time.Second)
	// time.Hour
	// }
	// dbaccess.SetDBAStr("root", "123456", "101.133.168.208", "3306", "tzwyj")
	// company.DealCompanyLonAndLat()
	// data2mysql.Deal0808()
	// company.DealDangerCompany()
	// company.DealWHCompany()
	// company.Company2MySQL()
	// dbaccess.SetDBAStr("ntsyjglzhxxxt", "dfBjk63Kcs1", "127.0.0.1", "33061", "T06024_yjgl")
	// qianzhi.DealTriger()
	// addressBook.TzwTask()

	// const (
	//	i=7
	//	j
	//	k
	// )
	// // i j k分别等于多少
	//
	// fmt.Println(i,j,k)

	/*endTime := time.Now().Format("2006-01-02T15:04:05")
	fmt.Println(url.QueryEscape(endTime))*/
	// 请问输出结果

	// DealWarehouse()
	dbaccess.SetDBAStr("root", "123456", "101.133.168.208", "3306", "codenai")
	// company.DealTzqCompany()

	// type ZhxxSwapOutLastTime struct {
	// 	ID          int       `gorm:"column:id" json:"id" form:"id"`
	// 	SwapId      int       `gorm:"column:swap_id" json:"swap_id" form:"swap_id"`
	// 	SwapTable   string    `gorm:"column:swap_table" json:"swap_table" form:"swap_table"`
	// 	SwapTableId int       `gorm:"column:swap_table_id" json:"swap_table_id" form:"swap_table_id"`
	// 	TimeField   string    `gorm:"column:time_field" json:"time_field" form:"time_field"`
	// 	LastTime    time.Time `gorm:"column:last_time" json:"last_time" form:"last_time"`
	// }
	//
	// aa := ZhxxSwapOutLastTime{
	// 	ID:          0,
	// 	SwapId:      0,
	// 	SwapTable:   "",
	// 	SwapTableId: 0,
	// 	TimeField:   "",
	// 	LastTime:    time.Now(),
	// }
	// rv := reflect.ValueOf(aa)
	//
	// fmt.Println(rv.Kind())

	// l1 := []int{1, 2, 3, 4, 5}
	// var l2 []*int
	// for _, i := range l1 {
	// 	l2 = append(l2, &i)
	// }
	//
	// fmt.Println()

	// str := "hawk.1.header\n1353832234\nGET\n/resource/1?b=1&a=2\n127.0.0.1\n8000\nsome-app-data"
	// s1 := GenSha256(str)
	// by := []byte(s1)
	// ss := base64.StdEncoding.EncodeToString(by)
	//
	// fmt.Println(ss)
	//
	// // 对sha256算法进行hash加密,key随便设置
	// hash := hmac.New(sha256.New, []byte("werxhqb98rpaxn39848xrunpaw3489ruxnpa98w4rxn")) // 创建对应的sha256哈希加密算法
	// hash.Write(by)                                                                      // 写入加密数据
	// // c10a04b78bcbcc1c4cba37f6afe0fa60cbf08f6e0a1d93b09387f7069be1aeff
	// ss2 := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	// // /uYWR6W5vTbY3WKUAN6fa+7p1t+1Yl6hFxKeMLfR6kk=
	// fmt.Println(ss2)
	//
	// aut := &Auth{
	// 	Credentials: Credentials{
	// 		ID:       "wqe",
	// 		Key:      "21",
	// 		Hash:     sha256.New,
	// 		Data:     nil,
	// 		App:      "1111",
	// 		Delegate: "132432",
	// 	},
	// 	Method:          "GET",
	// 	RequestURI:      "/auth/hawk",
	// 	Host:            "postman-echo.com",
	// 	Port:            "443",
	// 	MAC:             nil,
	// 	Nonce:           "1111",
	// 	Ext:             "111",
	// 	Hash:            nil,
	// 	ReqHash:         false,
	// 	IsBewit:         false,
	// 	Timestamp:       time.Unix(1604907962, 0),
	// 	ActualTimestamp: time.Time{},
	// }
	// fmt.Println(base64.StdEncoding.EncodeToString(aut.mac(AuthHeader)))
	//
	// r, err := http.NewRequest("GET", "https://postman-echo.com/auth/hawk/documents and settings", strings.NewReader("{}"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// path1 := r.URL.RequestURI()
	// if r.URL.RawQuery != "" {
	// 	path1 = path1[:len(path1)-len(r.URL.RawQuery)-1]
	// }
	// fmt.Println(path1)
	// slash := strings.HasSuffix(path1, "/")
	// fmt.Println(slash)
	// path1 = path.Clean(path1)
	// fmt.Println(path1)
	// if path1 != "/" && slash {
	//
	// 	path1 += "/"
	// }
	// fmt.Println("111", path1)
	// // w.Write([]byte(path))
	//
	// fmt.Println(path.Clean("C:/a/b/../c"))
	// fmt.Println(path.Clean("./1.txt"))
	// prepareRequestV4(r)
	// fmt.Println(normuri(normuri(r.URL.Path)))
	//
	// fmt.Println()
	// fmt.Println(os.Getwd())
	excelToCreateSQL.Start1()
	//company.StartTask()
	//video.StartTask()
}

func prepareRequestV4(request *http.Request) *http.Request {
	necessaryDefaults := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
		// "X-Amz-Date":   timestampV4(),
	}

	for header, value := range necessaryDefaults {
		if request.Header.Get(header) == "" {
			request.Header.Set(header, value)
		}
	}

	if request.URL.Path == "" {
		request.URL.Path += "/"
	}

	return request
}

func normuri(uri string) string {
	parts := strings.Split(uri, "/")
	for i := range parts {
		parts[i] = encodePathFrag(parts[i])
	}
	return strings.Join(parts, "/")
}

func encodePathFrag(s string) string {
	hexCount := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}
	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		} else {
			t[j] = c
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte) bool {
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
		return false
	}
	if '0' <= c && c <= '9' {
		return false
	}
	if c == '-' || c == '_' || c == '.' || c == '~' {
		return false
	}
	return true
}
func GenSha256(passwd string) string {
	h := sha256.New()

	h.Write([]byte(passwd))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func test1() (x int) {
	defer fmt.Printf("in defer: x = %d\n", x)
	x = 7
	return 9
}

func test2() (x int) {
	x = 7
	defer fmt.Printf("in defer: x = %d\n", x)
	return 9
}

func test3() (x int) {
	defer func() {
		fmt.Printf("in defer: x = %d\n", x)
	}()

	x = 7
	return 9
}

func test4() (x int) {
	defer func(n int) {
		fmt.Printf("in defer x as parameter: x = %d\n", n)
		fmt.Printf("in defer x after return: x = %d\n", x)
	}(x)

	x = 7
	return 9
}

type FireEnginesInfo struct {
	Id          string  `json:"id" form:"id" gorm:"column:ID"`
	Name        string  `json:"name" form:"name" gorm:"column:NAME"`
	TeamName    string  `json:"teamname" form:"teamname" gorm:"column:TEAM_NAME"`
	Desc        string  `json:"desc" form:"desc" gorm:"column:XFC_DESC"`
	Direct      float64 `json:"direct" form:"direct" gorm:"column:DIRECT"`
	GpsStatus   int     `json:"gpsstatus" form:"gpsstatus" gorm:"column:GPS_STATUS"`
	Speed       float64 `json:"speed" form:"speed" gorm:"column:SPEED"`
	TerminalNo  string  `json:"terminalno" form:"terminalno" gorm:"column:TERMINAL_NO"`
	State       int     `json:"state" form:"state" gorm:"column:STATE"`
	Type        int     `json:"type" form:"type" gorm:"column:TYPE"`
	GpsTime     string  `json:"gpstime" form:"gpstime" gorm:"column:GPS_TIME"`
	Lon         float64 `json:"lon" form:"lon" gorm:"column:LON"`
	Lat         float64 `json:"lat" form:"lat" gorm:"column:LAT"`
	PLon        float64 `json:"plon" form:"plon" gorm:"column:P_LON"`
	PLat        float64 `json:"plat" form:"plat" gorm:"column:P_LAT"`
	Height      float64 `json:"height" form:"height" gorm:"column:HEIGHT"`
	Address     string  `json:"address" form:"address" gorm:"column:ADDRESS"`
	Sim         string  `json:"sim" form:"sim" gorm:"column:SIM"`
	IsBus       int     `json:"isbus" form:"isbus" gorm:"column:IS_BUS"`
	XzqhId      int     `json:"xzqhid" form:"xzqhid" gorm:"column:XZQH_ID"`
	XzqhName    string  `json:"xzqhname" form:"xzqhname" gorm:"column:XZQH_NAME"`
	GeoInfo     string  `json:"geoinfo" form:"geoinfo" gorm:"column:GEO_INFO"`
	CenterPoint string  `json:"centerpoint" form:"centerpoint" gorm:"column:CENTER_POINT"`
	NodeCode    string  `json:"nodecode" form:"nodecode" gorm:"coolumn:"`
}

func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	// 8bit中首个0bit前有多少个1bits
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}
func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				// 判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			// 其他情况说明不是utf-8
			return false
		}
	}
	return true
}

const (
	GBK     string = "GBK"
	UTF8    string = "UTF8"
	UNKNOWN string = "UNKNOWN"
)

func GetStrCoding(data []byte) string {
	if isUtf8(data) == true {
		return UTF8
	} else if isGBK(data) == true {
		return GBK
	} else {
		return UNKNOWN
	}
}

func isGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			// 编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			// 大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func validUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { // 与操作之后不为0，说明首位为1
				for (buf[i] & 0x80) != 0 {
					buf[i] <<= 1 // 左移一位
					nBytes++     // 记录字符共占几个字节
				}

				if nBytes < 2 || nBytes > 6 { // 因为UTF8编码单字符最多不超过6个字节
					return false
				}

				nBytes-- // 减掉首字节的一个计数
			}
		} else { // 处理多字节字符
			if buf[i]&0xc0 != 0x80 { // 判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

// 得到学校类型
func getSchoolType(str string) string {
	ii := strings.LastIndex(str, ";")
	if ii != -1 && len(str) > 0 {
		return str[ii:+1]
	}
	return ""
}

func DealWarehouse() {
	list, err := dbaccess.ListWarehouseInformationAll()
	if err != nil {
		fmt.Println(err)
	}

	company.DealPoint()

	for i, _ := range list {
		ll := list[i]
		i := i

		if !company.CheckIsInGeo(ll.Longtitude, ll.Latitude) || (ll.Longtitude == 120.89078395928195 && ll.Latitude == 31.98307056920128) {
			fmt.Println(ll.Name, ll.Location)
			da := egis.GetLonAndLat(ll.Location)
			fmt.Println(da.Result.Level, da.Result.Location.Lng, da.Result.Location.Lat)
			if err != nil {
				fmt.Println(err)
			} else {

				list[i].Longtitude = da.Result.Location.Lng
				list[i].Latitude = da.Result.Location.Lat
				list[i].GeoInfo = utils.GenGeostr(ll.Longtitude, ll.Latitude)
				list[i].CenterPoint = list[i].GeoInfo
				dbaccess.UpdateWarehouseInformation(list[i])
			}

		}

	}

}
