package adminmodel

import (
	"web/component/objutils"
)

type AdminAdduserInfo struct {
	AddId       int
	UserLoginId string
	UserSource  string
	SourceName  string
	CreateTime  string
	AdminUserId int
}

func (this *AdminAdduserInfo) GetUniqId() int {
	return this.AddId
}

func (this *AdminAdduserInfo) GetUniqIdName() string {
	return "AddId"
}

func (this *AdminAdduserInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *AdminAdduserInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *AdminAdduserInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
