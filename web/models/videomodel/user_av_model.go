package videomodel

import (
	"web/component/objutils"
)

type UserAVAbbrs struct {
	UserAvId    int
	UserId      int
	UserPathMd5 string
	PathMd5     string
	Path        string
	PlayNum     int
	ExpiredTime string
	CreatedTime string
}

func (this *UserAVAbbrs) GetUniqId() int {
	return this.UserAvId
}

func (this *UserAVAbbrs) GetUniqIdName() string {
	return "UserAvId"
}

func (this *UserAVAbbrs) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *UserAVAbbrs) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *UserAVAbbrs) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
