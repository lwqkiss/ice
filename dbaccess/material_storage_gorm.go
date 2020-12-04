package dbaccess

//物资存储情况

type MaterialStorage struct {
	Id            int     `json:"id" form:"id" gorm:"column:ID;primary_key"`
	WarehouseId   int     `json:"warehouseid" form:"warehouseid" gorm:"column:WAREHOUSE_ID"`
	Supplies      string  `json:"supplies" form:"supplies" gorm:"column:SUPPLIES"`
	Quantity      float64 `json:"quantity" form:"quantity" gorm:"column:QUANTITY"`
	SuppliesModel string  `json:"suppliesmodel" form:"suppliesmodel" gorm:"column:SUPPLIES_MODEL"`
	ExpireDate    string  `json:"expiredate" form:"expiredate" gorm:"column:EXPIRE_DATE"`
	IsDelete      int     `json:"isdelete" form:"isdelete" gorm:"column:IS_DELETE"`
	WarehouseName string  `json:"warehousename" form:"warehousename" gorm:"column:WAREHOUSE_NAME"`
	RfidNo        string  `json:"rfidno" form:"rfidno" gorm:"column:RFID_NO"`
	CheckStatus   string  `json:"checkstatus" form:"checkstatus" gorm:"column:CHECK_STATUS"`
	PageNo        int     `json:"pageNo" form:"pageNo" gorm:"-"`
	PageSize      int     `json:"pageSize" form:"pageSize" gorm:"-"`
	Key           string  `json:"key" form:"key" gorm:"-"`
	Danwei        string  `json:"danwei" form:"danwei" gorm:"column:DANWEI"`
}

func (b *MaterialStorage) TableName() string {
	return "MATERIAL_STORAGE"
}

func AddMaterialStorage(b *MaterialStorage) error {
	db := OpenGorm()
	err := db.Create(b).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateMaterialStorage(b *MaterialStorage) error {
	db := OpenGorm()
	return db.Model(MaterialStorage{}).Where("ID=?", b.Id).Update(b).Error
}

func GetMaterialStorageById(id int) (*MaterialStorage, error) {
	db := OpenGorm()
	b := &MaterialStorage{}
	err := db.Where("ID=?  and IS_DELETE = 0", id).First(b).Error
	return b, err
}

func GetMaterialStorageByIdAndExpire(id int, expire string) (*MaterialStorage, error) {
	db := OpenGorm()
	b := &MaterialStorage{}
	err := db.Where("ID=? and EXPIRE_DATE = ? and  IS_DELETE = 0", id, expire).First(b).Error
	return b, err
}

func GetMaterialStorageByNameAndExpireByWarehouseID(name string, expire string, warehouseid int) (*MaterialStorage, error) {
	db := OpenGorm()
	b := &MaterialStorage{}
	err := db.Where("SUPPLIES=? and EXPIRE_DATE = ? and IS_DELETE = 0  and WAREHOUSE_ID = ?", name, expire, warehouseid).First(b).Error
	return b, err
}

func ListMaterialStorage(b *MaterialStorage) ([]*MaterialStorage, error) {
	db := OpenGorm()
	bis := make([]*MaterialStorage, 0)

	if b.Key == "" {
		err := db.Where("IS_DELETE = 0 and WAREHOUSE_ID = ?", b.WarehouseId).Offset((b.PageNo - 1) * b.PageSize).Limit(b.PageSize).Find(&bis).Error
		return bis, err
	}
	err := db.Where("IS_DELETE = 0 and WAREHOUSE_ID = ? and CONCAT_WS(',',SUPPLIES,SUPPLIES_MODEL,EXPIRE_DATE) REGEXP ?", b.WarehouseId, b.Key).Offset((b.PageNo - 1) * b.PageSize).Limit(b.PageSize).Find(&bis).Error
	return bis, err
}

func CountMaterialStorage(b *MaterialStorage) (int64, error) {
	db := OpenGorm()
	var count int64

	if b.Key == "" {
		err := db.Model(MaterialStorage{}).Where("IS_DELETE = 0 and WAREHOUSE_ID = ?", b.WarehouseId).Count(&count).Error
		return count, err
	}
	err := db.Model(MaterialStorage{}).Where("IS_DELETE = 0 and WAREHOUSE_ID = ? and CONCAT_WS(',',SUPPLIES,SUPPLIES_MODEL,EXPIRE_DATE) REGEXP ?", b.WarehouseId, b.Key).Count(&count).Error
	return count, err
}

func DeleteMaterialStorage(id int) error {
	db := OpenGorm()
	return db.Model(&MaterialStorage{}).Where("ID=?", id).Update("IS_DELETE", 1).Error
}

func DeleteAll(warehouseid int) error {
	db := OpenGorm()
	return db.Model(&MaterialStorage{}).Where("WAREHOUSE_ID=?", warehouseid).Update("IS_DELETE", 1).Error
}
