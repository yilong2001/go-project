package couponmodel

import (
	"web/component/objutils"
	//"web/models/basemodel"
)

type CouponInfo struct {
	CouponId   int
	CouponSign string
	CouponName string `schema:"CouponName"`
	Money      int    `schema:"Money"`
	TotalNum   int    `schema:"TotalNum"`
	ApplyNum   int
	UsedNum    int
	CreateTime string
	ExpireTime string `schema:"ExpireTime"`
}

func (this *CouponInfo) GetUniqId() int {
	return this.CouponId
}

func (this *CouponInfo) GetUniqIdName() string {
	return "CouponId"
}

func (this *CouponInfo) GetSkipFieldsForOpenQuery() []string {
	return []string{"CoupSign"}
}

func (this *CouponInfo) GetSkipFieldsForUpdate() []string {
	return []string{"CouponId", "CouponSign", "CreateTime", "ApplyNum", "UsedNum", "CreateTime"}
}

func (this *CouponInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *CouponInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *CouponInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

func NewCouponInfo() *CouponInfo {
	return &CouponInfo{
		CouponId:   -1,
		CouponSign: "",
		CouponName: "",
		Money:      0,
		TotalNum:   0,
		ApplyNum:   0,
		UsedNum:    0,
		CreateTime: "",
		ExpireTime: "",
	}
}
