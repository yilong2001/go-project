package coursemodel

import (
	//"log"
	//"strings"
	"web/component/objutils"
)

const (
	Const_Course_Type_Open   = 1
	Const_Course_Type_Inner  = 2
	Const_Course_Type_Online = 3
)

type CourseMainInfo struct {
	CourseId int
	UserId   int

	Title     string `schema:"Title"`
	TitleMore string `schema:"TitleMore"`

	CourseType int `schema:"CourseType"`

	DirectionType int `schema:"DirectionType"`

	OrgPrice     int    `schema:"OrgPrice"`
	NowPrice     int    `schema:"NowPrice"`
	PriceDescrip string `schema:"PriceDescrip"`

	Notice string `schema:"Notice"`

	Aganda    string `schema:"Aganda"`
	Advantage string `schema:"Advantage"`

	TargetPeople string `schema:"TargetPeople"`
	Effect       string `schema:"Effect"`
	Teachers     string `schema:"Teachers"`

	OnlineType int    `schema:"OnlineType"`
	City       int    `schema:"City"`
	Industry   int    `schema:"Industry"`
	Address    string `schema:"Address"`

	IsPaused int `schema:"IsPaused"`

	CoverPic    string
	CourseClass int
	IsZhiying   int
	ClassPeriod int

	CreateTime string
	UpdateTime string

	FavouriteNum int
	CommentNum   int
	ServedNum    int
	AvgStar      float32

	AuditStatus int
	AuditInfo   string
	AuditName   string
	AuditTime   string
}

func (this *CourseMainInfo) GetUniqId() int {
	return this.CourseId
}

func (this *CourseMainInfo) GetUniqIdName() string {
	return "CourseId"
}

func (this *CourseMainInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *CourseMainInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}

func (this *CourseMainInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

type CourseCatalogFirstInfo struct {
	CatalogFirstId int
	CourseId       int
	Title          string
	Descrip        string

	VideoPath string

	CreateTime string
	UpdateTime string
}

func (this *CourseCatalogFirstInfo) GetUniqId() int {
	return this.CatalogFirstId
}

func (this *CourseCatalogFirstInfo) GetUniqIdName() string {
	return "CatalogFirstId"
}

func (this *CourseCatalogFirstInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *CourseCatalogFirstInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)
}

func (this *CourseCatalogFirstInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

type CourseCatalogSecondInfo struct {
	CatalogSecondId int
	CatalogFirstId  int
	Title           string
	Descrip         string

	VideoPath string

	CreateTime string
	UpdateTime string
}

func (this *CourseCatalogSecondInfo) GetUniqId() int {
	return this.CatalogSecondId
}

func (this *CourseCatalogSecondInfo) GetUniqIdName() string {
	return "CatalogSecondId"
}

func (this *CourseCatalogSecondInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *CourseCatalogSecondInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)
}

func (this *CourseCatalogSecondInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
