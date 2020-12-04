package ship

import (
	"lwq/dbaccess"
)

type Dgcarriageliquid struct {
	ID              int     `gorm:"column:id" json:"id" form:"id"`
	ApplyNo         string  `gorm:"column:申报单编号" json:"申报单编号" form:"申报单编号"`
	ShipNo          string  `gorm:"column:船舶识别号" json:"船舶识别号" form:"船舶识别号"`
	IMONo           string  `gorm:"column:IMO编号" json:"IMO编号" form:"IMO编号"`
	InitNO          string  `gorm:"column:初始登记号" json:"初始登记号" form:"初始登记号"`
	LisenseNo       string  `gorm:"column:牌舶号" json:"牌舶号" form:"牌舶号"`
	ShipCallNo      string  `gorm:"column:船舶呼号" json:"船舶呼号" form:"船舶呼号"`
	Cname           string  `gorm:"column:中文船名" json:"中文船名" form:"中文船名"`
	Ename           string  `gorm:"column:英文船名" json:"英文船名" form:"英文船名"`
	Country         string  `gorm:"column:国籍" json:"国籍" form:"国籍"`
	ShipPort        string  `gorm:"column:船籍港" json:"船籍港" form:"船籍港"`
	ShipType        string  `gorm:"column:船舶种类" json:"船舶种类" form:"船舶种类"`
	TotalTon        string  `gorm:"column:总吨" json:"总吨" form:"总吨"`
	CleanTon        string  `gorm:"column:净吨" json:"净吨" form:"净吨"`
	TotalReferTon   string  `gorm:"column:参考总载重吨" json:"参考总载重吨" form:"参考总载重吨"`
	MMSI            string  `gorm:"column:MMSI" json:"MMSI" form:"MMSI"`
	Vessel          string  `gorm:"column:航次" json:"航次" form:"航次"`
	InOutPort       string  `gorm:"column:进出港" json:"进出港" form:"进出港"`
	Trade           string  `gorm:"column:内外贸" json:"内外贸" form:"内外贸"`
	WaitPort        string  `gorm:"column:靠泊港口" json:"靠泊港口" form:"靠泊港口"`
	AwayTime        string  `gorm:"column:抵离港时间" json:"抵离港时间" form:"抵离港时间"`
	WaitLoc         string  `gorm:"column:靠泊泊位" json:"靠泊泊位" form:"靠泊泊位"`
	StartWorkTime   string  `gorm:"column:开始作业时间" json:"开始作业时间" form:"开始作业时间"`
	EndWorkTime     string  `gorm:"column:结束作业时间" json:"结束作业时间" form:"结束作业时间"`
	ContactPerson   string  `gorm:"column:联系人姓名" json:"联系人姓名" form:"联系人姓名"`
	ContactPhone    string  `gorm:"column:联系方式" json:"联系方式" form:"联系方式"`
	UnitName        string  `gorm:"column:申报单位名称" json:"申报单位名称" form:"申报单位名称"`
	ApplyPersonName string  `gorm:"column:申报员姓名" json:"申报员姓名" form:"申报员姓名"`
	ApplyCert       string  `gorm:"column:申报员证书编号" json:"申报员证书编号" form:"申报员证书编号"`
	CaptainNo       string  `gorm:"column:船长编号" json:"船长编号" form:"船长编号"`
	CaptainName     string  `gorm:"column:船长姓名" json:"船长姓名" form:"船长姓名"`
	CaptainCert     string  `gorm:"column:船长证书编号" json:"船长证书编号" form:"船长证书编号"`
	OrderNo         string  `gorm:"column:序号" json:"序号" form:"序号"`
	GoodsApplyNo    string  `gorm:"column:货报申报单编号" json:"货报申报单编号" form:"货报申报单编号"`
	GoodsAway       string  `gorm:"column:货物流向" json:"货物流向" form:"货物流向"`
	GoodsKind       string  `gorm:"column:货物种类" json:"货物种类" form:"货物种类"`
	UNNO            string  `gorm:"column:UNNO" json:"UNNO" form:"UNNO"`
	TrueCname       string  `gorm:"column:正确运输名称中文" json:"正确运输名称中文" form:"正确运输名称中文"`
	TrueEname       string  `gorm:"column:正确运输名称英文" json:"正确运输名称英文" form:"正确运输名称英文"`
	TrueCnameRemark string  `gorm:"column:正确运输名称中文补充说明" json:"正确运输名称中文补充说明" form:"正确运输名称中文补充说明"`
	TrueEnameRemark string  `gorm:"column:正确运输名称英文补充说明" json:"正确运输名称英文补充说明" form:"正确运输名称英文补充说明"`
	DangerType      string  `gorm:"column:危险类别" json:"危险类别" form:"危险类别"`
	PolluteType     string  `gorm:"column:污染类别" json:"污染类别" form:"污染类别"`
	OilType         string  `gorm:"column:油类性质" json:"油类性质" form:"油类性质"`
	FlashPoint      string  `gorm:"column:闪点" json:"闪点" form:"闪点"`
	WorkType        string  `gorm:"column:作业方式" json:"作业方式" form:"作业方式"`
	OilNo           string  `gorm:"column:液货舱编号" json:"液货舱编号" form:"液货舱编号"`
	StartPort       string  `gorm:"column:起运港" json:"起运港" form:"起运港"`
	GoodPort        string  `gorm:"column:装货港" json:"装货港" form:"装货港"`
	UnloadPort      string  `gorm:"column:卸货港" json:"卸货港" form:"卸货港"`
	TargetPort      string  `gorm:"column:目的港" json:"目的港" form:"目的港"`
	SendPerson      string  `gorm:"column:发货人" json:"发货人" form:"发货人"`
	CheckPerson     string  `gorm:"column:托运人" json:"托运人" form:"托运人"`
	TransPerson     string  `gorm:"column:承运人" json:"承运人" form:"承运人"`
	RecPerson       string  `gorm:"column:收货人" json:"收货人" form:"收货人"`
	GoodsOwner      string  `gorm:"column:货主" json:"货主" form:"货主"`
	GetNo           string  `gorm:"column:提货单号" json:"提货单号" form:"提货单号"`
	ApproveTime     string  `gorm:"column:审批时间" json:"审批时间" form:"审批时间"`
	Lon             float64 `gorm:"column:lon" json:"lon" form:"lon"`
	Lat             float64 `gorm:"column:lat" json:"lat" form:"lat"`
	GeoInfo         string  `gorm:"column:GEO_INFO" json:"GEO_INFO" form:"GEO_INFO"`
	NodeCode        string  `gorm:"column:NODE_CODE" json:"NODE_CODE" form:"NODE_CODE"`
	St              string  `gorm:"-" json:"st" form:"st"`
	Et              string  `gorm:"-" json:"et" form:"et"`
}

func (a *Dgcarriageliquid) TableName() string {
	return "dgCarriageLiquid"
}

func GetShipTrans() ([]*Dgcarriageliquid, error) {
	db := dbaccess.OpenGorm()
	res := make([]*Dgcarriageliquid, 0)

	err := db.Table("dgCarriageLiquid").Find(&res).Error

	return res, err

}

func UpdateShip(b *Dgcarriageliquid) error {
	db := dbaccess.OpenGorm()
	// 这样注册地址就不会更新

	return db.Table("dgCarriageLiquid").Where("id=?", b.ID).Update(b).Error
}
