package newsmodel

import (
	//"log"
	//"strings"
	"web/component/objutils"
)

type VideoInfo struct {
	Id        int
	Title     string
	Url       string
	RefUrl    string
	VideoType int
	//OrgDate    string
	IsSub  int
	SubNum int

	CreateTime string
	Descrip    string
	//Content    string
	//ImgSrc     string
	OrgWeb     string
	NewsClass  string
	CommentNum int
	Tags       string
}

func (this *VideoInfo) GetUniqId() int {
	return this.Id
}

func (this *VideoInfo) GetUniqIdName() string {
	return "Id"
}

func (this *VideoInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *VideoInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}

func (this *VideoInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
