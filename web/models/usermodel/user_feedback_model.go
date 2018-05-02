package usermodel

import (
	"web/component/objutils"
)

type UserFeedbackInfo struct {
	UserFeedbackId int
	Info           string `schema:"Info"`
	Contact        string `schema:"Contact"`
	CreateTime     string
}

func (this *UserFeedbackInfo) GetUniqId() int {
	return this.UserFeedbackId
}

func (this *UserFeedbackInfo) GetUniqIdName() string {
	return "UserFeedbackId"
}

func (this *UserFeedbackInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *UserFeedbackInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *UserFeedbackInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
