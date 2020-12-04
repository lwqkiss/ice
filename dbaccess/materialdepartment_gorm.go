package dbaccess

type MaterialDepartment struct {
	Id                int    `json:"id" form:"id" gorm:"column:ID"`
	DepartmentName    string `json:"departmentname" form:"departmentname" gorm:"column:DEPARTMENT_NAME"`
	DepartmentAccount string `json:"departmentaccount" form:"departmentaccount" gorm:"column:DEPARTMENT_ACCOUNT"`
	ContactPerson     string `json:"contactperson" form:"contactperson" gorm:"column:CONTACT_PERSON"`
	CellPhone         string `json:"contactnumber" form:"contactnumber" gorm:"column:CELL_PHONE"`
	Explanation       string `json:"explanation" form:"explanation" gorm:"column:EXPLANATION"`
	IsDelete          int    `json:"isdelete" form:"isdelete" gorm:"column:IS_DELETE"`
	LoginCode         string `json:"logincode" form:"logincode" gorm:"column:LOGIN_CODE"`
	LastTime          string `json:"lasttime" form:"lasttime" gorm:"column:LAST_TIME"`
	PageNo            int    `json:"pageNo" form:"pageNo" gorm:"-"`
	PageSize          int    `json:"pageSize" form:"pageSize" gorm:"-"`
	LoginType
	Key string `json:"key" form:"key" gorm:"-"`
}

func (b *MaterialDepartment) TableName() string {
	return "MATERIAL_DEPARTMENT"
}
func AddMaterialDepartment(b *MaterialDepartment) error {
	db := OpenGorm()
	return db.Create(b).Error
}

func UpdateMaterialDepartment(b *MaterialDepartment) error {
	db := OpenGorm()
	return db.Model(MaterialDepartment{}).Where("ID=?", b.Id).Update(b).Error
}

func GetMaterialDepartmentById(id int) (*MaterialDepartment, error) {
	db := OpenGorm()
	b := &MaterialDepartment{}
	err := db.Where("ID=? and IS_DELETE = 0", id).First(b).Error
	return b, err
}
func ListMaterialDepartment(b *MaterialDepartment) ([]*MaterialDepartment, error) {
	db := OpenGorm()
	bis := make([]*MaterialDepartment, 0)

	err := db.Where("IS_DELETE =0 and LAST_TIME like ? and DEPARTMENT_NAME like ?", "%"+b.LastTime+"%", "%"+b.DepartmentName+"%").Offset((b.PageNo - 1) * b.PageSize).Limit(b.PageSize).Find(&bis).Error
	return bis, err

}

func CountMaterialDepartment(b *MaterialDepartment) (int64, error) {
	db := OpenGorm()
	var count int64
	err := db.Model(MaterialDepartment{}).Where("IS_DELETE =0 and LAST_TIME like ? and DEPARTMENT_NAME like ?", "%"+b.LastTime+"%", "%"+b.DepartmentName+"%").Count(&count).Error
	return count, err
}

func DeleteMaterialDepartment(id int64) error {
	db := OpenGorm()
	return db.Model(&MaterialDepartment{}).Where("ID=?", id).Update("IS_DELETE", 1).Error
}

func GetMaterialDepartmentByPhone(phone string) (*MaterialDepartment, error) {
	db := OpenGorm()
	b := &MaterialDepartment{}
	err := db.Where("CELL_PHONE=?", phone).First(b).Error
	return b, err
}
