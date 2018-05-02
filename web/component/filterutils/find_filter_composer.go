package filterutils

import (
	"strconv"
	"web/component/objutils"
	"web/models/condmodel"
	"web/models/reqparamodel"
)

func GetComposerOfFilters(params *reqparamodel.HttpReqParams, uniqIdName, uniqTypeName string) *condmodel.CondComposerLinker {

	composer1 := condmodel.NewCondComposerLinker("and")

	idWhere1 := make(map[string]interface{})
	idRlue1 := make(map[string]string)

	idWhere1["del_status"] = int(0)
	idRlue1["del_status"] = " = "
	if uid, err := strconv.ParseInt(params.RouterParams[uniqIdName], 10, 32); err == nil {
		tmp, _ := objutils.CamelToUnderLine(uniqIdName)
		idWhere1[tmp] = int(uid)
		idRlue1[tmp] = " = "
	}

	//idWhere1["audit_status"] = int(1)
	//idRlue1["audit_status"] = " = "

	utype := params.URLParams.Get(uniqTypeName)
	if utype != "" {
		utypecode, err := strconv.ParseInt(utype, 10, 32)
		if err == nil && utypecode > 0 {
			tmp, _ := objutils.CamelToUnderLine(uniqTypeName)
			idWhere1[tmp] = utypecode
			idRlue1[tmp] = " = "
		}
	}

	city := params.URLParams.Get("City")
	if city != "" {
		cityCode, err := strconv.ParseInt(city, 10, 32)
		if err == nil && cityCode > 0 {
			idWhere1["city"] = int(cityCode)
			idRlue1["city"] = " = "
		}
	}

	directionType := params.URLParams.Get("DirectionType")
	if directionType != "" {
		dt, err := strconv.ParseInt(directionType, 10, 32)
		if err == nil && dt > 0 {
			idWhere1["direction_type"] = int(dt)
			idRlue1["direction_type"] = " = "
		}
	}

	cur := composer1

	industry := params.URLParams.Get("Industry")
	if industry != "" {
		industryCode, err := strconv.ParseInt(industry, 10, 32)
		if err == nil && industryCode > 0 {
			if industryCode%1000 == 0 {

				compsubt1 := condmodel.NewCondComposerItem(map[string]interface{}{"industry": industryCode}, map[string]string{"industry": ">="}, " and ")
				composert1 := condmodel.NewCondComposerLinker("and")
				composert1.SetItem(compsubt1)
				cur.SetNext(composert1)
				cur = composert1

				compsubt2 := condmodel.NewCondComposerItem(map[string]interface{}{"industry": industryCode + 99}, map[string]string{"industry": "<="}, " and ")
				composert2 := condmodel.NewCondComposerLinker("and")
				composert2.SetItem(compsubt2)
				cur.SetNext(composert2)
				cur = composert2
			} else {
				idWhere1["industry"] = industryCode
				idRlue1["industry"] = " = "
			}
		}
	}

	compsub1 := condmodel.NewCondComposerItem(idWhere1, idRlue1, " and ")
	composer1.SetItem(compsub1)

	return composer1
}
