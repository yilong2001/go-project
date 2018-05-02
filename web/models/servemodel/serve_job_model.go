package servemodel

import (
	"web/component/objutils"
)

type ServeJobInfo struct {
	ServiceJobId int

	ServiceId     int
	ServiceUserId int

	JobId     int `schema:"JobId"`
	JobUserId int

	CreateTime string
}

func (this *ServeJobInfo) GetUniqId() int {
	return this.ServiceJobId
}

func (this *ServeJobInfo) GetUniqIdName() string {
	return "ServiceJobId"
}

func (this *ServeJobInfo) GetSkipFieldsForOpenQuery() []string {
	return []string{}
}

func (this *ServeJobInfo) GetSkipFieldsForUpdate() []string {
	return []string{"ServiceJobId", "ServiceId", "ServiceUserId", "JobId", "JobUserId", "CreateTime"}
}

func (this *ServeJobInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *ServeJobInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *ServeJobInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
