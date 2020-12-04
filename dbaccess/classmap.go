package dbaccess

type ClassMap struct {
	ID       string `gorm:"column:id" json:"id" form:"id"`
	Type     string `gorm:"column:type" json:"type" form:"type"`
	Big      string `gorm:"column:big" json:"big" form:"big"`
	Mid      string `gorm:"column:mid" json:"mid" form:"mid"`
	Small    string `gorm:"column:small" json:"small" form:"small"`
	Bigg     string `gorm:"column:bigg" json:"bigg" form:"bigg"`
	Midd     string `gorm:"column:midd" json:"midd" form:"midd"`
	Smalll   string `gorm:"column:smalll" json:"smalll" form:"smalll"`
	Nodecode string `gorm:"column:nodecode" json:"nodecode" form:"nodecode"`
}

func (b *ClassMap) TableName() string {
	return "CLASS_MAP"
}

func GetClassMapInfo(typeid string) (*ClassMap, error) {
	db := OpenGorm()
	res := ClassMap{}
	err := db.Model(&ClassMap{}).Where("TYPE = ?", typeid).First(&res).Error
	if err != nil {
		return &res, err
	}
	return &res, err
}
