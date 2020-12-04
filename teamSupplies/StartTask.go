package teamSupplies

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	"lwq/dbaccess"
	"lwq/utils"
	"path"
	"strconv"
	"strings"
)

func StartTask() {

	fileList := utils.GetFileListByPath("./南通市应急救援队伍装备汇总/崇川")
	//fmt.Println(fileList)
	for i, _ := range fileList {
		fmt.Println("开始处理 " + path.Join("./南通市应急救援队伍装备汇总/崇川", fileList[i]))
		f, err := excelize.OpenFile(path.Join("./南通市应急救援队伍装备汇总/崇川", fileList[i]))
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(f.SheetCount)
		//fmt.Println(f.GetRows("Sheet1"))
		if len(f.GetRows("Sheet1")) == 0 {
			utils.WriteHistoryFile(fileList[i], "./", "无数据的表.txt")
		}

		//处理得到的数据
		HandleData(f.GetRows("Sheet1"))
	}

}

type HandleSupplies struct {
	Id            int
	Supplies      string
	Count         float64
	MeasureUnit   string
	Model         string
	Dept          string
	Address       string
	ContactPerson string
	Phone         string
}

func HandleData(a [][]string) {
	/**
	0:   序号
	1:   物资名称
	2:   数量
	3:   单位
	4:   规格型号
	5:   所属单位（部门）
	6:   存储位置
	7:   联系人
	8:   联系人手机
	9:   十四五 计划增加数量
	10:  十四五达标数量

	*/
	lastFlag := ""
	for i, _ := range a {
		//从数据行开始
		if i >= 2 {

			vo := a[i]
			//如果 数量一行没有值，那么跳过
			if vo[2] == "" && vo[3] == "" && vo[1] == "" && vo[4] == "" && vo[5] == "" && vo[6] == "" && vo[7] == "" && vo[8] == "" {
				continue
			} else {
				if vo[1] == "" {
					vo[1] = lastFlag
				} else {
					lastFlag = vo[1]
				}
				if vo[5] == "" {
					vo[5] = vo[6]
				}
				if vo[5] == "" {
					continue
				}
				/*if vo[5] == "消防支队" || vo[5] == "卫健委" || vo[5] == "交通运输局" || vo[5] == "市政和园林局" || vo[5] == "发改委" || vo[5] == "水利局" {
					vo[5], vo[6] = vo[6], vo[5]
				}
				if vo[5] == "市应急物资储备库" {
					vo[5] = "南通市应急救灾物资储备库"

				}*/

				//有值的处理，说明是有数据的行
				/*if vo[5] == "启东市消防救援大队" {
					vo[5] = vo[6]
				}*/
				id, _ := strconv.Atoi(vo[0])
				count, _ := strconv.ParseFloat(vo[2], 64)
				tmp := HandleSupplies{
					Id:            id,
					Supplies:      strings.TrimSpace(vo[1]),
					Count:         count,
					MeasureUnit:   strings.TrimSpace(vo[3]),
					Model:         strings.TrimSpace(vo[4]),
					Dept:          strings.TrimSpace(vo[5]),
					Address:       strings.TrimSpace(vo[6]),
					ContactPerson: strings.TrimSpace(vo[7]),
					Phone:         strings.TrimSpace(vo[8]),
				}
				tmp.Dept = HandleDeptName(tmp.Dept)

				//先找到仓库，没有就新建
				warehouse := &dbaccess.WarehouseInformation{}
				//消防队的处理，启东是这样处理的，其他地方不一样  每个地方其实都不一样，所以要分地区处理

				team, err := dbaccess.GetTeamByNameOrPhone(tmp.Dept, tmp.Phone)
				if err != nil {
					//fmt.Println(err.Error())
					if err == gorm.ErrRecordNotFound {
						//先在仓库表找,队伍没找到的情况
						tt, err := dbaccess.GetWarehouseInformationByName(tmp.Dept)
						if err != nil {
							if err == gorm.ErrRecordNotFound {
								tmpWarehouse := dbaccess.WarehouseInformation{

									Name:          tmp.Dept,
									Location:      tmp.Address,
									Longtitude:    0,
									Latitude:      0,
									ContactPerson: tmp.ContactPerson,
									Phone:         tmp.Phone,
									XzqhName:      "崇川区",
									XzqhId:        9,
									NodeCode:      "",
									GeoInfo:       "",
									CenterPoint:   "",
									ManageUnit:    "崇川区物资管理",
									ManageUnitID:  34,
									LoginCode:     "123456",
									IsDelete:      0,
								}
								if strings.Contains(tmp.Dept, "水务") {
									tmpWarehouse.NodeCode = "T140102"
								} else if strings.Contains(tmp.Dept, "公司") {
									tmpWarehouse.NodeCode = "T140101"
								} else {
									tmpWarehouse.NodeCode = "T140103"
								}
								fmt.Println("新增仓库 ", tmpWarehouse.Name)
								err := dbaccess.AddWarehouseInformationByTeam(&tmpWarehouse)
								if err != nil {
									fmt.Println(err.Error())
								}
								tt = &tmpWarehouse
							}
						}
						warehouse = tt
					}

				} else if team.Id != 0 {

					tt, err := dbaccess.GetWarehouseInformationById(team.Id)
					if err != nil {
						if err == gorm.ErrRecordNotFound {

							tmpWarehouse := dbaccess.WarehouseInformation{
								Id:            team.Id,
								Name:          team.TEAMNAME,
								Location:      team.TEAMADDRESS,
								Longtitude:    team.LON,
								Latitude:      team.LAT,
								ContactPerson: team.TEAMLEADERNAME,
								Phone:         team.TEAMPHONE,
								XzqhName:      team.XZQHNAME,
								XzqhId:        team.XZQHID,
								NodeCode:      "T140105",
								GeoInfo:       team.GeoInfo,
								CenterPoint:   team.GeoInfo,
								ManageUnit:    "救援队伍物资搜集",
								ManageUnitID:  23,
								LoginCode:     "123456",
								IsDelete:      0,
							}
							fmt.Println("新增仓库 ", tmpWarehouse.Name)
							err := dbaccess.AddWarehouseInformationByTeam(&tmpWarehouse)
							if err != nil {
								fmt.Println(err.Error())
							}
							tt = &tmpWarehouse
						}
					}
					warehouse = tt
				}

				//到这里，有了仓库的id了

				//现在有仓库了
				//那么检查物资是否在标准化名称里面，在的话，就更新单位，不在，那么新增
				tmpName, err := dbaccess.GetMaterialNameStandardizationByName(tmp.Supplies)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						ttt := dbaccess.MaterialNameStandardization{

							StandardName: tmp.Supplies,
							Alias:        "",
							Demand:       0,
							MeasureUnit:  tmp.MeasureUnit,
							IsDelete:     0,
							PageNo:       0,
							PageSize:     0,
							CategoryName: "消防应急救援装备",
							CategoryId:   7,
							StandardNo:   "",
							Key:          "",
						}
						fmt.Println("新增标准化名称 ", ttt.StandardName)
						dbaccess.AddMaterialNameStandardization(&ttt)
					}
				} else {
					tmpName.MeasureUnit = tmp.MeasureUnit
					fmt.Println("更新标准化单位", tmpName)
					dbaccess.UpdateMaterialNameStandardization(tmpName)
				}

				//新增物资

				t1 := dbaccess.MaterialStorage{

					WarehouseId:   warehouse.Id,
					Supplies:      tmp.Supplies,
					Quantity:      tmp.Count,
					SuppliesModel: tmp.Model,
					ExpireDate:    "2099-01-01",
					IsDelete:      0,
					WarehouseName: warehouse.Name,
					RfidNo:        "",
					CheckStatus:   "",
					PageNo:        0,
					PageSize:      0,
					Key:           "",
					Danwei:        tmp.MeasureUnit,
				}

				err = dbaccess.AddMaterialStorage(&t1)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

}

func HandleDeptName(dept string) string {
	switch dept {
	case "南纤消防队":
		return "南通醋酸纤维有限公司专职消防队"
	case "如东县消防大队、如东县沿海经济开发区":
		return "如东县洋口港经济开发区专职消防队"

	case "南通消防支队崇川区大队钟秀中队":
		return "钟秀政府专职队"
	case "城管":
		return "崇川区城管"
	case "二院":
		return "通州区第二人民医院"
	case "三院":
		return "通州区第三人民医院"
	case "海门高新区（海门市秀山西路581号）":
		return "海门高新区防汛应急救援小组"
	case "江苏斯德雷特通光光纤有限公司（南通市海门市北海西路219号）":
		return "江苏斯德雷特光纤救援小组"
	case "开发区消防":
		return "通州区经济技术开发区专职消防队"
	case "工业园区":
		return "通州区工业园区"
	case "四院":
		return "通州区第四人民医院"
	case "五院":
		return "通州区第五人民医院"
	case "六院":
		return "通州区第六人民医院"
	case "一院":
		return "通州区第一人民医院"
	case "海安大队":
		return "海安中队"

	case "气象局":
		return "海安气象局"
	case "执法局":
		return "海安执法局"
	case "生态环境监测站":
		return "海安生态环境监测站"
	case "城建集团":
		return "通州城建集团"
	case "区交运局":
		return "通州区交运局"
	case "消防大队":
		return "通州湾消防大队"
	case "通州区消防大队":
		return "通州消防综合救援中队"
	case "公安局":
		return "通州区公安局"
	case "应急管理局":
		return "通州区应急管理局"
	case "经发局":
		return "通州湾经发局"
	case "区人武部":
		return "通州区人武部"
	case "人武部":
		return "通州区人武部"
	case "水利局":
		return "南通市水利局"
	case "交通运输局":
		return "南通市交通运输局"
	case "市政和园林局":
		return "南通市市政和园林局"
	default:
		return dept

	}
}
