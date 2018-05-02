package usermodel

import (
	//"log"
	//"reflect"
	"web/component/objutils"
	"web/models/basemodel"
)

type UserRegisterInfo struct {
	Phone    string `schema:"Phone" form:"Phone" binding:"required"`
	Password string `schema:"Password" form:"Password"`
	Code     string `schema:"Code"`
}

func (this *UserRegisterInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

type UserPasswordInfo struct {
	Phone    string `schema:"Phone"`
	OldPw    string `schema:"OldPw" form:"OldPw"`
	Password string `schema:"Password" form:"Password"`
	Code     string `schema:"Code"`
}

type UserInfo struct {
	UserId int

	LoginName string `schema:"LoginName" form:"LoginName"`
	LoginPw   string

	UserName string `schema:"UserName" form:"UserName"`

	QqId          string `schema:"QqId" form:"QqId"`
	WeixinOpenId  string `schema:"WeixinOpenId" form:"WeixinOpenId"`
	WeixinUnionId string `schema:"WeixinUnionId" form:"WeixinUnionId"`

	WeiboId string `schema:"WeiboId" form:"WeiboId"`

	NickName string `schema:"NickName" form:"NickName"`
	Sex      int    `schema:"Sex" form:"Sex"`
	City     int    `schema:"City" form:"City"`
	Age      int    `schema:"Age" form:"Age"`
	Phone    string `schema:"Phone"  form:"Phone"`
	Email    string `schema:"Email" form:"Email"`

	OpenContact string `schema:"OpenContact"  form:"OpenContact"`

	UserType int `schema:"UserType"  form:"UserType"`

	EmailWork  string `schema:"EmailWork" form:"EmailWork"`
	CreateTime string
	UpdateTime string

	Industry string `schema:"Industry" form:"Industry"`

	Firm       string `schema:"Firm" form:"Firm"`
	SimpleFirm string

	JobTitle  string `schema:"JobTitle" form:"JobTitle"`
	Introduce string `schema:"Introduce" form:"Introduce"`

	Portrait string

	ServedNum  int
	JobsNum    int
	CommentNum int
	AvgStar    float32

	AccountBalance int

	IsTuijian   int
	TuijianInfo string
	TuijianImg  string
	TuijianUid  int

	IsRenzheng   int
	RenzhengTime string
	RenzhengInfo string
	RenzhengUser int

	IsZhiying int
	IsVip     int
}

func (this *UserInfo) IsUserIdValid() bool {
	if this.UserId < 10000 {
		return false
	}

	return true
}

func (this *UserInfo) GetUserId() int {
	return this.UserId
}

func (this *UserInfo) GetUniqId() int {
	return this.UserId
}

func (this *UserInfo) GetUniqIdName() string {
	return "UserId"
}

func (this *UserInfo) GetSkipFieldsForSelfQuery() []string {
	return []string{"LoginPw"}
}

func (this *UserInfo) GetSkipFieldsForOrderQuery() []string {
	return []string{"LoginName", "LoginPw", "CreateTime", "AccountBalance"}
}

func (this *UserInfo) GetSkipFieldsForOpenQuery() []string {
	return []string{"LoginName", "LoginPw", "Phone", "Email", "CreateTime", "AccountBalance"}
}

func (this *UserInfo) GetSkipFieldsForUpdate() []string {
	return []string{"UserId", "LoginName", "LoginPw", "CreateTime", "Portrait", "AccountBalance", "ServedNum", "CommentNum", "AvgStar", "JobsNum"}
}

func (this *UserInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *UserInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *UserInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

func NewUserInfo() *UserInfo {
	info := NewUserInfoEx()
	return &info
}

func NewUserInfoEx() UserInfo {
	return UserInfo{
		UserId:    -1,
		LoginName: basemodel.Default_Value_For_String_Field,
		LoginPw:   basemodel.Default_Value_For_String_Field,

		UserName:      basemodel.Default_Value_For_String_Field,
		QqId:          basemodel.Default_Value_For_String_Field,
		WeixinOpenId:  basemodel.Default_Value_For_String_Field,
		WeixinUnionId: basemodel.Default_Value_For_String_Field,
		WeiboId:       basemodel.Default_Value_For_String_Field,

		NickName: basemodel.Default_Value_For_String_Field,
		Sex:      -1,
		City:     -1,
		Age:      -1,

		Phone: basemodel.Default_Value_For_String_Field,
		Email: basemodel.Default_Value_For_String_Field,

		EmailWork:  basemodel.Default_Value_For_String_Field,
		CreateTime: basemodel.Default_Value_For_String_Field,
		UpdateTime: basemodel.Default_Value_For_String_Field,

		Industry:  basemodel.Default_Value_For_String_Field,
		Firm:      basemodel.Default_Value_For_String_Field,
		JobTitle:  basemodel.Default_Value_For_String_Field,
		Introduce: basemodel.Default_Value_For_String_Field,

		Portrait: basemodel.Default_Value_For_String_Field,
	}
}
