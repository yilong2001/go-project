package newsmodel

import (
	//"log"
	//"strings"
	"web/component/objutils"
)

type UserNewsTagInfo struct {
	TagId      int
	UserId     int
	Tag        string `schema:"Tag"`
	CreateTime string
}

func (this *UserNewsTagInfo) GetUniqId() int {
	return this.TagId
}

func (this *UserNewsTagInfo) GetUniqIdName() string {
	return "TagId"
}

func (this *UserNewsTagInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *UserNewsTagInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}

func (this *UserNewsTagInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
