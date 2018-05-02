package ordermodel

import (
	"web/component/objutils"
)

type PayRecordInfo struct {
	Record     string
	SourceId   int
	DestId     int
	CreateTime string
}

func (this *PayRecordInfo) GetUniqId() int {
	return -1
}

func (this *PayRecordInfo) GetUniqIdName() string {
	return "RecordId"
}
func (this *PayRecordInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *PayRecordInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *PayRecordInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
