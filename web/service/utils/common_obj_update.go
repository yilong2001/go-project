package utils

import (
	//"database/sql"
	"log"
	//"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	//"strconv"
	//"time"
	"strings"

	//"web/component/cfgutils"
	//"web/component/idutils"
	"web/component/errcode"
	"web/component/objutils"
	//"web/component/pageutils"
	"web/component/sqlutils"

	//"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/reqparamodel"
	//"web/service/routers"
)

func (this *ObjectWithIdUtil) Util_UpdateObjectInfoWithId(info basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams, req *http.Request, ren render.Render, skipFields []string, specFields []string) bool {
	log.Println("Util_UpdateObjectInfoWithId start...")

	if !this.getCheckParamFuncForUpdate().(func(*reqparamodel.HttpReqParams, render.Render) bool)(params, ren) {
		return false
	}

	//info := this.getObjInfo()
	if this.FormUnParseFlag == 0 {
		err := objutils.ParseObjectWithForm(info, req)
		if err != nil {
			log.Println(err)
			ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, err.Error()))
			return false
		}
	}

	log.Println(info)

	err := this.getExpendObjInitFuncForUpdate().(func(basemodel.ObjectUtilBaseIf, *reqparamodel.HttpReqParams) error)(info, params)
	if err != nil {
		log.Println(err)
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, err.Error()))
		return false
	}

	//fieldIfArrs, _ := info.GetFieldsWithNotDefaultValue()
	postFields := params.PostFields
	destFields := []string{}
	var fieldIfArrs map[string]interface{}
	if skipFields == nil && specFields == nil {
		fieldIfArrs, _ = info.GetFieldsWithSpecs(postFields)
	} else if specFields != nil {
		// for _, pf := range postFields {
		// 	for _, sf := range specFields {
		// 		if strings.ToLower(pf) == strings.ToLower(sf) {
		// 			destFields = append(destFields, sf)
		// 		}
		// 	}
		// }
		fieldIfArrs, _ = info.GetFieldsWithSpecs(specFields)
	} else {
		for _, pf := range postFields {
			isNotFound := true
			for _, sf := range skipFields {
				if strings.ToLower(pf) == strings.ToLower(sf) {
					isNotFound = false
				}
			}

			if isNotFound {
				destFields = append(destFields, pf)
			}
		}

		fieldIfArrs, _ = info.GetFieldsWithSpecs(destFields)
	}

	res, fns := IsFieldsValueOk(fieldIfArrs)
	if !res {
		detail := "the info (" + fns + ") is not correct"
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, detail))
		return false
	}

	whereCond, ruleCond := this.getWhereCondFuncForUpdate().(func(*reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string))(params)

	if len(whereCond) < 1 {
		detail := "query condition is null!"
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Query_Where_Error, detail))
		return false
	}

	msqls, args := sqlutils.Sqls_CompUpdate(this.getTableName(), fieldIfArrs, whereCond, ruleCond)

	log.Println(msqls)
	log.Println(args)

	if this.getTX() == nil {
		err = sqlutils.Sqls_Do_PrepareAndExec(this.getDB(), msqls, args)
		if err != nil {
			ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
			return false
		}
	} else {
		err = sqlutils.Sqls_Do_PrepareAndExec_Tx(this.getTX(), msqls, args)
		if err != nil {
			ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
			return false
		}
	}

	if this.getMoreProcessForUpdate() != nil {
		err := this.getMoreProcessForUpdate().(func(basemodel.ObjectUtilBaseIf, *reqparamodel.HttpReqParams) error)(info, params)
		if err != nil {
			ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
			return false
		}
	}

	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": info.GetUniqId()})
	return true
}
