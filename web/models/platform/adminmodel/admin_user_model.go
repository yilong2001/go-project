package adminmodel

import (
	"web/component/objutils"
)

type AdminUserInfo struct {
	UserId    int
	LoginName string
	LoginPw   string
	UserName  string
	Role      int
	Email     string
	Introduce string
	Portrait  string
}

func (this *AdminUserInfo) GetUniqId() int {
	return this.UserId
}

func (this *AdminUserInfo) GetUniqIdName() string {
	return "CouponId"
}

func (this *AdminUserInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *AdminUserInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *AdminUserInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
