package excelToCreateSQL

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"lwq/utils"
	"strings"
)

/**

CREATE TABLE `ACCESS_REPORT` (
  `ID` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `TYPES_OF` varchar(128) DEFAULT NULL COMMENT '类型',
  `ALARM_TIME` varchar(64) DEFAULT NULL COMMENT '接警时间',
  `STATUS` varchar(64) DEFAULT NULL COMMENT '状态',
  `RESPONSE_LEVEL` varchar(64) DEFAULT NULL COMMENT '响应级别',
  `OBJECT_PRIMARY_KEY_ID` int(11) DEFAULT NULL COMMENT '对象主键id',
  `STARTING_TIME` varchar(64) DEFAULT NULL COMMENT '开始时间',
  `END_TIME` varchar(64) DEFAULT NULL COMMENT '结束时间',
  `CANCELLATION_TIME` varchar(64) DEFAULT NULL COMMENT '取消时间',
  `UPDATE_TIME` varchar(32) DEFAULT NULL COMMENT '更新时间',
  `IS_DELETE` int(11) DEFAULT NULL COMMENT '是否删除',
  `SOURCE` varchar(200) DEFAULT NULL COMMENT '来源',
  `IS_TEST` varchar(5) DEFAULT NULL COMMENT '是否测试',
  `ACCESS_TIME` varchar(64) DEFAULT NULL COMMENT '填写时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=100029 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='接报表'

*/

func Start() {
	f, err := excelize.OpenFile("./数据文件/二期已发布共享市级部门资源目录-对接.xlsx")

	if err != nil {
		log.Fatal(err)
	}

	dataRow := f.GetRows("数据项信息")

	sqlDrop := "drop table  IF EXISTS "
	sqlCreate := ""
	count := 0
	pri := ""
	// tableName := ""
	tableComment := ""
	for i := range dataRow {

		l := dataRow[i]
		if strings.Contains(l[0], "zhxx") || strings.Contains(l[0], "finish") {
			if count > 0 {
				if pri != "" {
					sqlCreate += " `zhxx_insert_time` timestamp(6) NOT NULL DEFAULT  CURRENT_TIMESTAMP(6) COMMENT '创建时间',\n"
					sqlCreate += "PRIMARY KEY (`" + pri + "`) USING BTREE\n"
				} else {
					sqlCreate += " `zhxx_insert_time` timestamp(6) NOT NULL DEFAULT  CURRENT_TIMESTAMP(6) COMMENT '创建时间'\n"
				}

				sqlCreate += ") ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='" + tableComment + "';"
				// 执行建表  直接打印建表语句即可
				fmt.Println(sqlCreate)
				utils.WriteFile(sqlCreate)
				// 初始化sql
				sqlCreate = ""
			}
			// tableName = l[0]
			pri = ""
			tableComment = l[1]
			sqlDrop += "`" + l[0] + "`;\n"
			sqlCreate += "DROP TABLE IF EXISTS `" + l[0] + "`;\n"
			sqlCreate += "CREATE TABLE `" + l[0] + "`(\n"

			count++
		} else {

			dataType := ""
			switch l[5] {
			case "数字":
				dataType = "double"
			case "整型":
				dataType = "int(11)"
			case "大字段":
				dataType = "text"
			case "文本":
				dataType = "varchar(100)"
			case "日期":
				dataType = "datetime"
			case "时间":
				dataType = "date"

			}
			if l[6] == "主键id" {
				pri = l[3]
				// if dataType == "varchar(2048)" {
				// 	dataType="varchar(255)"
				// }
				sqlCreate += "`" + l[3] + "` " + dataType + " NOT NULL COMMENT '" + l[2] + "',\n"
			} else {
				sqlCreate += "`" + l[3] + "` " + dataType + " DEFAULT NULL COMMENT '" + l[2] + "',\n"
			}

		}

	}

}

func Start1() {
	f, err := excelize.OpenFile("./数据文件/二期已发布共享市级部门资源目录-对接.xlsx.xlsx")

	if err != nil {
		log.Fatal(err)
	}

	dataRow := f.GetRows("数据项信息")

	st := "{"
	col := ""
	col2 := ""
	cname := ""
	count := 0
	lasttable := ""
	for i := range dataRow {

		l := dataRow[i]
		if i == 0 {
			lasttable = l[0]
			cname = l[1]
		}
		if strings.Contains(l[0], "zhxx") || strings.Contains(l[0], "finish") {

			if count > 0 {
				st = strings.TrimRight(st, ",")
				st = st + "}"
				col = strings.TrimRight(col, ",")
				fmt.Println(lasttable, cname)
				fmt.Println(st)
				fmt.Println(col)
				fmt.Println(col2)
				st = "{"
				col = ""
				col2 = ""
			}
			lasttable = l[0]
			cname = l[1]
			count++
		} else {
			if l[6] == "1" {
				col += "`" + l[3] + "`,"
				col2 += "`" + l[3] + "`,"
				st += "\"" + strings.ToLower(l[3]) + "\"" + ":" + "\"" + l[2] + "\"" + ","

			}

		}

	}

}
