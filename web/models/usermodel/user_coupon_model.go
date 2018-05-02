package usermodel

import (
	"web/component/objutils"
)

type UserCouponInfo struct {
	UserCouponId   int
	UserCouponSign string `schema"UserCouponSign"`
	UserId         int
	CouponId       int
	CouponName     string
	CouponStatus   int
	Money          int
	CreateTime     string
	ExpireTime     string
}

func (this *UserCouponInfo) IsUserIdValid() bool {
	if this.UserId < 10000 || this.UserId < 10000 {
		return false
	}

	return true
}

func (this *UserCouponInfo) GetUserId() int {
	return this.UserId
}

func (this *UserCouponInfo) GetUniqId() int {
	return this.UserCouponId
}

func (this *UserCouponInfo) GetUniqIdName() string {
	return "UserCouponId"
}

func (this *UserCouponInfo) GetSkipFieldsForSelfQuery() []string {
	return []string{"UserCouponSign"}
}

func (this *UserCouponInfo) GetSkipFieldsForUpdate() []string {
	return []string{"UserCouponId", "UserCouponSign", "UserId", "CouponName", "CouponId", "Money", "CreateTime"}
}

func (this *UserCouponInfo) GetSpecFieldsForUpdateStatus() []string {
	return []string{"CouponStatus"}
}

func (this *UserCouponInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *UserCouponInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *UserCouponInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

func NewUserCouponInfo() *UserCouponInfo {
	return &UserCouponInfo{
		UserCouponId:   -1,
		UserCouponSign: "",
		UserId:         -1,
		CouponId:       -1,
		CouponName:     "",
		CouponStatus:   0,
		Money:          0,
		CreateTime:     "",
		ExpireTime:     "",
	}
}
