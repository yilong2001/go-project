package firmmodel

import (
	"log"
	"strings"
	"web/component/objutils"
)

type FirmInfo struct {
	FirmId      int
	Name        string `schema:"Name"`
	SimpleName  string `schema:"SimpleName"`
	CompareName string
	MainWeb     string `schema:"MainWeb"`
	Address     string `schema:"Address"`
	KefuPhone   string `schema:"KefuPhone"`
	Introduce   string `schema:"Introduce"`
	Honor       string `schema:"Honor"`
	Logo        string `schema:"logo"`

	AdminUserId int

	CreateTime string
	UpdateTime string

	AuditStatus int
	AuditInfo   string
	AuditName   string
	AuditTime   string

	IsTuijian   int
	TuijianInfo string
	TuijianImg  string
	TuijianUid  int
}

func (this *FirmInfo) GetUniqId() int {
	return this.FirmId
}

func (this *FirmInfo) GetUniqIdName() string {
	return "FirmId"
}

func (this *FirmInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *FirmInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}

func (this *FirmInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

var (
	FirmNameFixedSuffixs []string = []string{"股份有限责任公司", "有限责任公司", "股份有限公司", "集团有限责任公司", "集团有限公司", "集团股份公司", "有限公司", "股份公司", "公司"}

	ChangYongCi []string = []string{"科技", "技术"}
)

func MakeCompareName(name string) string {
	out := name
	for _, suf := range FirmNameFixedSuffixs {
		if strings.Contains(name, suf) {
			log.Print(name, suf)
			out = strings.Replace(name, suf, "", -1)
			break
		}
	}

	tmps := strings.Split(out, "(")
	if len(tmps) == 1 {
		tmps = strings.Split(out, "（")
	}

	if len(tmps) == 1 {
		return out
	}

	tmps2 := strings.Split(tmps[1], ")")
	if len(tmps2) == 1 {
		tmps2 = strings.Split(out, "）")
	}

	result := ""
	if len(tmps2) > 1 {
		result = tmps[0] + tmps2[1]
	} else {
		result = tmps[0] + tmps2[0]
	}

	for _, cyc := range ChangYongCi {
		result = strings.Replace(result, cyc, "", -1)
	}

	return result
}
