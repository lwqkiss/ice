package addressBook

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"path"
	"time"
)

//通州湾通讯录

func TzwTask() {
	fmt.Println("开始处理 " + path.Join("./addressBook/通讯录-通州湾.xlsx"))
	f, err := excelize.OpenFile(path.Join("./addressBook/通讯录-通州湾.xlsx"))
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(f.SheetCount)
	//fmt.Println(f.GetRows("Sheet1"))

	//处理得到的数据
	handleData(f.GetRows("Sheet1"))
}

func handleData(data [][]string) {
	//0 姓名 1 办公电话 2 办公短号 3 手机号码 4 集团短号 5 标记
	dept := 0
	job := 0
	for i := range data {
		per := data[i]
		if per[5] == "1" {
			job = 0
			dept = 0
			//新建部门
			tmp := &Department{
				DeptName:   per[0],
				DeptType:   "dept",
				ParentId:   0,
				IsDelete:   0,
				UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
			}
			AddDepartment(tmp)
			dept = tmp.Id
			continue
		}
		if per[5] == "2" {
			job = 0
			//新建职务
			tmp := &Department{
				DeptName:   per[0],
				DeptType:   "job",
				ParentId:   dept,
				IsDelete:   0,
				UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
			}
			AddDepartment(tmp)
			job = tmp.Id
			continue
		}

		//处理正常的数据
		trueJobId := 0
		if job != 0 {
			trueJobId = job
		} else {
			if dept == 0 {
				continue
			}
			trueJobId = dept
		}
		if per[0] == "" && per[1] == "" && per[2] == "" && per[3] == "" && per[4] == "" {
			continue
		}
		ur := &User{
			Id:          0,
			UserName:    per[0],
			Mobilephone: per[3],
			Telphone:    per[1],
			Address:     "",
			IsDelete:    0,
			UpdateTime:  time.Now().Format("2006-01-02 15:04:05"),
			IsEmPerson:  "否",
			WorkNum:     per[2],
			CropNum:     per[4],

			DpList: []int{trueJobId},
		}
		ii, err := AddUser(ur)
		fmt.Println(ii, err)
	}
}
