package resourcemodel

import (
	//"log"
	//"strings"
	"web/component/objutils"
)

const (
	Const_User_Resource_Article = 0
	Const_User_Resource_Pic     = 1
	Const_User_Resource_Video   = 2
)

const (
	Const_User_Resource_Refer_Service = 0
	Const_User_Resource_Refer_Job     = 1
	Const_User_Resource_Refer_Course  = 2
)

type UserResourceInfo struct {
	ResourceId   int
	ResourceType int
	UserId       int
	ReferId      int
	ReferType    int
	Title        string
	Url          string
	CreateTime   string
	ClickNum     int
	Content      string
}

func (this *UserResourceInfo) GetUniqId() int {
	return this.ResourceId
}

func (this *UserResourceInfo) GetUniqIdName() string {
	return "ResourceId"
}

func (this *UserResourceInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *UserResourceInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}

func (this *UserResourceInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
