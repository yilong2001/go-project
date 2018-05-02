package servemodel

import (
	"web/component/objutils"
)

type ServeTagInfo struct {
	ServiceTagId int
	ServiceId    int
	ServiceType  int
	Tag          string
	CreateTime   string
}

func (this *ServeTagInfo) GetUniqId() int {
	return this.ServiceTagId
}

func (this *ServeTagInfo) GetUniqIdName() string {
	return "ServiceTagId"
}

func (this *ServeTagInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *ServeTagInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *ServeTagInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

type ServeTagCountsInfo struct {
	ServiceType int
	Tag         string
	Count       int
}

func (this *ServeTagCountsInfo) GetSelectPartsOfGroupSql() (string, []interface{}) {
	selFieldIfs := make([]interface{}, 0)
	sqls := "select count(service_tag_id) as count, service_type, tag"

	_, fieldAddrIfArr := objutils.GetWholeFields(this)

	selFieldIfs = append(selFieldIfs, fieldAddrIfArr["count"])
	selFieldIfs = append(selFieldIfs, fieldAddrIfArr["service_type"])
	selFieldIfs = append(selFieldIfs, fieldAddrIfArr["tag"])

	return sqls, selFieldIfs
}

func (this *ServeTagCountsInfo) GetGroupPartsOfGroupSql() string {
	return " service_type,tag "
}

func GetAllServiceTypes() []int {
	return []int{0, 1, 2, 3}
}
