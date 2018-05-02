package usermodel

import (
	"web/component/objutils"
)

const (
	Const_Customer_Servier_Phone_Main = "15022589186"
)

type SmsCodeInfo struct {
	Phone      string `schema:"Phone"`
	Code       string
	CreateTime string
}

func (this *SmsCodeInfo) GetUniqId() int {
	return 1
}

func (this *SmsCodeInfo) GetUniqIdName() string {
	return "Phone"
}

func (this *SmsCodeInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *SmsCodeInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *SmsCodeInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
