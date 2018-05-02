package adminmodel

import (
	"web/component/objutils"
)

type AdminReviewerResultInfo struct {
	Ids  string `schema:"ids"`
	Info string `schema:"info"`
}

func (this *AdminReviewerResultInfo) GetUniqId() int {
	return -1
}

func (this *AdminReviewerResultInfo) GetUniqIdName() string {
	return ""
}
func (this *AdminReviewerResultInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *AdminReviewerResultInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *AdminReviewerResultInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}
