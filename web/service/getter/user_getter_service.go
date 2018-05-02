package getter

import (
	"database/sql"
	"log"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	//"github.com/martini-contrib/render"

	"errors"
	//"net/http"
	//"strconv"
	//"strings"
	//"time"
	//"reflect"

	"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/servemodel"
	//"web/service/utils"
)

func transferIf2ServiceInfo(result *[]interface{}) *[]servemodel.ServeInfo {
	outs := []servemodel.ServeInfo{}
	for _, rst := range *result {
		info, ok := rst.(servemodel.ServeInfo)
		if !ok {
			log.Println("userinfo type error ", rst)
			continue
		}

		outs = append(outs, info)
	}

	return &outs
}

func (this *ModelInfoGetter) GetServiceByServiceId(db *sql.DB, serviceid int, skip, specs []string) (*servemodel.ServeInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	serveinfo := servemodel.NewServeInfo()

	var fieldAddrIfArrs map[string]interface{}
	if skip == nil && specs == nil {
		_, fieldAddrIfArrs = serveinfo.GetWholeFields()
	} else if skip == nil {
		_, fieldAddrIfArrs = serveinfo.GetFieldsWithSpecs(specs)
	} else {
		_, fieldAddrIfArrs = serveinfo.GetFieldsWithSkip(skip)
	}

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["service_id"] = serviceid
	ruleCond["service_id"] = " = "
	whereCond["del_status"] = 0
	ruleCond["del_status"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetServiceTableName(), serveinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(servemodel.ServeInfo)
	if !ok {
		return nil, errors.New("ServeInfo output type is wrong")
	}

	log.Println("GetServiceByServiceId", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetServicesByServiceIds(db *sql.DB, serviceid []int) (*[]servemodel.ServeInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	serveinfo := servemodel.NewServeInfo()

	_, fieldAddrIfArrs := serveinfo.GetWholeFields()

	ruleCond := make(map[string]string)
	ruleCond["service_id"] = " = "
	ruleCond["del_status"] = " = "

	allWhereConds := make([]map[string]interface{}, 0)
	for _, si := range serviceid {
		whereCond := map[string]interface{}{}
		whereCond["service_id"] = si
		whereCond["del_status"] = 0
		allWhereConds = append(allWhereConds, whereCond)
	}

	result, err := sqlutils.Sqls_GetMultiInfoWithMultiConds(db2, this.GetServiceTableName(), serveinfo, fieldAddrIfArrs, allWhereConds, ruleCond)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	outs := transferIf2ServiceInfo(result)

	return outs, nil
}
