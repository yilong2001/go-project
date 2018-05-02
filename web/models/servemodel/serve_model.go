package servemodel

import (
	"web/component/objutils"
	//"web/models/basemodel"
)

type ServeInfo struct {
	ServiceId        int
	ServiceType      int    `schema:"ServiceType"`
	ServiceName      string `schema:"ServiceName"`
	ServiceItem      string `schema:"ServiceItem"`
	ServiceItemOther string `schema:"ServiceItemOther"`
	UserId           int

	City     int
	Industry int `schema:"Industry"`

	DirectionType int `schema:"DirectionType"`

	Description string `schema:"Description"`
	CostDesc    string `schema:"CostDesc"`
	Price       int    `schema:"Price"`

	Duration string `schema:"Duration"`

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

	IsZhiying int
}

func NewServeInfo() *ServeInfo {
	return &ServeInfo{}
	// return &ServeInfo{}
	// 	ServiceId:        -1,
	// 	ServiceType:      0,
	// 	ServiceName:      basemodel.Default_Value_For_String_Field,
	// 	ServiceItem:      "",
	// 	ServiceItemOther: "",

	// 	UserId:      -1,
	// 	Description: basemodel.Default_Value_For_String_Field,
	// 	Price:       -1,
	// 	Duration:    "",

	// 	CreateTime:  basemodel.Default_Value_For_String_Field,
	// 	UpdateTime:  basemodel.Default_Value_For_String_Field,
	// 	AuditStatus: 0,
	// 	AuditTime:   basemodel.Default_Value_For_String_Field,
	// 	AuditName:   basemodel.Default_Value_For_String_Field,
	// }
}

func (this *ServeInfo) IsServeIdValid() bool {
	if this.ServiceId < 10000 {
		return false
	}

	return true
}

func (this *ServeInfo) IsUserIdValid() bool {
	if this.UserId < 10000 {
		return false
	}

	return true
}

func (this *ServeInfo) GetUserId() int {
	return this.UserId
}

func (this *ServeInfo) GetUniqId() int {
	return this.ServiceId
}

func (this *ServeInfo) GetUniqIdName() string {
	return "ServiceId"
}

func (this *ServeInfo) GetSkipFieldsForOpenQuery() []string {
	return []string{"AuditTime", "AuditName", "AuditStatus", "AuditInfo"}
}

func (this *ServeInfo) GetSkipFieldsForUpdate() []string {
	return []string{"ServiceId", "ServiceType", "UserId", "City", "CreateTime", "AuditTime", "AuditName"}
}

func (this *ServeInfo) GetSkipFieldsForAdmin() []string {
	return []string{}
}

func (this *ServeInfo) GetSpecFieldsForAdminUpdate() []string {
	return []string{"AuditTime", "AuditName", "AuditStatus", "AuditInfo"}
}

func (this *ServeInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *ServeInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *ServeInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
