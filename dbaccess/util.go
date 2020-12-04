package dbaccess

import "fmt"

func GetCodeByName(str string) (string, error) {
	var nodecode string
	db, err := OpenDB()
	defer db.Close()
	if err != nil {
		return "", fmt.Errorf("数据库连接异常%w", err)
	}
	er := db.QueryRow("select NODE_CODE from INFO_CATALOG_TREE where NODE_NAME =?", str).Scan(&nodecode)
	return nodecode, er

}
func GetIdByAdministrationName(adnm string) (int, string, error) {

	db, err := OpenDB()
	defer db.Close()
	if err != nil {
		return 0, "", err.Error
	}
	id := 0
	str := ""
	sql := `select id,name from administration where name = ?`
	row := db.QueryRow(sql, adnm)
	er := row.Scan(&id, &str)
	if er != nil {
		return 0, "", er
	}
	return id, str, nil

}
func GetNewCodeByName(str string) (string, error) {
	var nodecode string
	db, err := OpenDB()
	defer db.Close()
	if err != nil {
		return "", fmt.Errorf("数据库连接异常%w", err)
	}
	er := db.QueryRow("select NODE_CODE from INFO_CATALOG_TREE where NODE_NAME =?", str).Scan(&nodecode)
	return nodecode, er

}
