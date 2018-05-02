package jobmodel

import (
	//"log"
	//"reflect"
	"web/component/objutils"
	//"web/models/basemodel"
)

type JobInfo struct {
	JobId       int    `schema:"JobId"`
	UserId      int    `schema:"UserId"`
	Title       string `schema:"Title"`
	ProjectName string `schema:"ProjectName"`
	WorkYears   int    `schema:"WorkYears"`
	Address     string `schema:"Address"`
	City        int    `schema:"City"`
	Industry    int    `schema:"Industry"`
	Education   int    `schema:"Education"`
	Attraction  string `schema:"Attraction"`

	Available   int8   `schema:"Available"`
	Description string `schema:"Description"`
	SalaryMin   int    `schema:"SalaryMin"`
	SalaryMax   int    `schema:"SalaryMax"`

	Tag1 string `schema:"Tag1"`
	Tag2 string `schema:"Tag2"`
	Tag3 string `schema:"Tag3"`
	Tag4 string `schema:"Tag4"`
	Tag5 string `schema:"Tag5"`

	CreateTime string
	UpdateTime string

	FavouriteNum int

	AuditName string
	AuditTime string
}

func (this *JobInfo) IsJobIdValid() bool {
	if this.JobId < 10000 {
		return false
	}

	return true
}

func (this *JobInfo) IsUserIdValid() bool {
	if this.UserId < 10000 {
		return false
	}

	return true
}

func (this *JobInfo) GetUserId() int {
	return this.UserId
}

func (this *JobInfo) GetUniqId() int {
	return this.JobId
}

func (this *JobInfo) GetUniqIdName() string {
	return "JobId"
}

func (this *JobInfo) GetSkipFieldsForOpenQuery() []string {
	return []string{"AuditTime", "AuditName", "AuditStatus", "AuditInfo"}
}

func (this *JobInfo) GetSkipFieldsForUpdate() []string {
	return []string{"JobId", "UserId", "CreateTime", "AuditTime", "AuditName"}
}

func (this *JobInfo) GetSkipFieldsForAdmin() []string {
	return []string{}
}
func (this *JobInfo) GetSpecFieldsForAdminUpdate() []string {
	return []string{"AuditTime", "AuditName", "AuditStatus", "AuditInfo"}
}

func (this *JobInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *JobInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *JobInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

func NewJobInfo() *JobInfo {
	return &JobInfo{
	// JobId:  -1,
	// UserId: -1,

	// JobName:     basemodel.Default_Value_For_String_Field,
	// ProjectName: basemodel.Default_Value_For_String_Field,

	// WorkYears:   -1,
	// Address:     basemodel.Default_Value_For_String_Field,
	// Available:   -1,
	// Description: basemodel.Default_Value_For_String_Field,
	// SalaryMin:   -1,
	// SalaryMax:   -1,

	// Tag1: basemodel.Default_Value_For_String_Field,
	// Tag2: basemodel.Default_Value_For_String_Field,
	// Tag3: basemodel.Default_Value_For_String_Field,
	// Tag4: basemodel.Default_Value_For_String_Field,
	// Tag5: basemodel.Default_Value_For_String_Field,

	// CreateTime: basemodel.Default_Value_For_String_Field,
	// UpdateTime: basemodel.Default_Value_For_String_Field,
	// AuditName:  basemodel.Default_Value_For_String_Field,
	// AuditTime:  basemodel.Default_Value_For_String_Field,
	}
}
