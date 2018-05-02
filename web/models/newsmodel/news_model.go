package newsmodel

import (
	//"log"
	//"strings"
	"web/component/objutils"
)

type NewsInfo struct {
	Id         int
	Title      string
	Url        string
	OrgDate    string
	CreateTime string
	Descrip    string
	Content    string
	ImgSrc     string
	OrgWeb     string
	NewsClass  string
	CommentNum int
	Tags       string
}

func (this *NewsInfo) GetUniqId() int {
	return this.Id
}

func (this *NewsInfo) GetUniqIdName() string {
	return "Id"
}

func (this *NewsInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *NewsInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}

func (this *NewsInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

type NewsCommentInfo struct {
	CommentId   int
	CommentType int
	UserId      int
	SrcId       int
	Comment     string `schema:"Comment"`
	Star        int    `schema:"Star"`
	ReplyNum    int
	UpNum       int
	DownNum     int
	CreateTime  string
}

func (this *NewsCommentInfo) GetUniqId() int {
	return this.CommentId
}

func (this *NewsCommentInfo) GetUniqIdName() string {
	return "CommentId"
}

func (this *NewsCommentInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *NewsCommentInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}

func (this *NewsCommentInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
