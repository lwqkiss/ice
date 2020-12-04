package qianzhi

import (
	"fmt"
	"lwq/dbaccess"
	"strings"
)

func DealTriger() {

	db, _ := dbaccess.OpenDB()
	table := "T06024_SAFETY_EVENT_COLLECT"
	tmpTable := "T06024_SAFET2003201731149_TMP"
	triger := "2003201731149"
	t1 := "trigger_insert_" + triger
	t2 := "trigger_modify_" + triger
	t3 := "trigger_delete_" + triger
	col := getColumns(table)
	qrcol := "new." + strings.Replace(col, ",", ", new.", -1)
	oldcol := "old." + strings.Replace(col, ",", ", old.", -1)
	//fmt.Println(col)

	sql1 := `CREATE TRIGGER ` + t1 + ` AFTER INSERT ON ` + table + ` FOR EACH ROW insert into ` + tmpTable + `(` + col + `, KTL_FLG, KTL_TRI_FIELD)
values(` + qrcol + `, 'N', null);`
	sql2 := `CREATE TRIGGER ` + t2 + ` AFTER UPDATE ON ` + table + ` FOR EACH ROW insert into ` + tmpTable + `(` + col + `, KTL_FLG, KTL_TRI_FIELD)
	values(` + qrcol + `, 'M', null);`
	sql3 := `CREATE TRIGGER ` + t3 + ` AFTER DELETE ON ` + table + ` FOR EACH ROW insert into ` + tmpTable + `(` + col + `, KTL_FLG, KTL_TRI_FIELD)
	values(` + oldcol + `, 'D', null);`

	//fmt.Println(sql1)
	//fmt.Println(sql2)
	//fmt.Println(sql3)
	res, err := db.Exec(sql1)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res.RowsAffected())
	res, err = db.Exec(sql2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res.RowsAffected())
	res, err = db.Exec(sql3)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res.RowsAffected())
}

type RES struct {
	Col string `json:"col" gorm:"column:col"`
}

func getColumns(tbname string) string {
	db := dbaccess.OpenGorm()
	res := RES{}
	err := db.Raw("SELECT GROUP_CONCAT(COLUMN_NAME) col FROM information_schema.`COLUMNS` where TABLE_SCHEMA = 'T06024_yjgl' and TABLE_NAME = ?", tbname).First(&res).Error
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	fmt.Println(res.Col)
	return res.Col
}
