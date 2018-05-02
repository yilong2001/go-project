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
	"web/models/condmodel"
	"web/models/reqparamodel"
	//"web/service/routers"
)

func (this *UpdateObjectWithIdUtil) Update_With_MultiInObject(orginfo basemodel.ObjectUtilBaseIf,
	destinfo basemodel.ObjectUtilBaseIf,
	params *reqparamodel.HttpReqParams,
	req *http.Request,
	ren render.Render,
	skipFields []string,
	specFields []string) bool {
	log.Println("Util_UpdateObjectInfoWithId start...")

	if !this.CheckParamFunc.(func(*reqparamodel.HttpReqParams, render.Render) bool)(params, ren) {
		return false
	}

	//info := this.getObjInfo()
	if this.FormUnParseFlag == 0 {
		err := objutils.ParseObjectWithForm(orginfo, req)
		if err != nil {
			ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, err.Error()))
			return false
		}

		log.Println(orginfo)
	}

	err := this.ExInitFunc.(func(basemodel.ObjectUtilBaseIf, basemodel.ObjectUtilBaseIf, *reqparamodel.HttpReqParams) error)(orginfo, destinfo, params)
	if err != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, err.Error()))
		return false
	}

	//fieldIfArrs, _ := info.GetFieldsWithNotDefaultValue()
	postFields := []string{}
	if params != nil {
		postFields = params.PostFields
	}

	destFields := []string{}
	var fieldIfArrs map[string]interface{}
	if skipFields == nil && specFields == nil {
		fieldIfArrs, _ = destinfo.GetFieldsWithSpecs(postFields)
	} else if specFields != nil {
		fieldIfArrs, _ = destinfo.GetFieldsWithSpecs(specFields)
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

		fieldIfArrs, _ = destinfo.GetFieldsWithSpecs(destFields)
	}

	res, fns := IsFieldsValueOk(fieldIfArrs)
	if !res {
		detail := "the info (" + fns + ") is not correct"
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, detail))
		return false
	}

	condCompser := this.CondCompserFunc.(func(basemodel.ObjectUtilBaseIf, basemodel.ObjectUtilBaseIf, *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker)(orginfo, destinfo, params)

	if condCompser.Item == nil {
		detail := "query condition is null!"
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Query_Where_Error, detail))
		return false
	}

	var calcUpdateFields map[string]string = nil

	if this.CalcedUpdateFieldsFunc != nil {
		calcUpdateFields = this.CalcedUpdateFieldsFunc.(func(basemodel.ObjectUtilBaseIf, basemodel.ObjectUtilBaseIf, *reqparamodel.HttpReqParams) map[string]string)(orginfo, destinfo, params)
	}

	msqls, args := sqlutils.Sqls_CompUpdateWithComposerLinker(this.getTableName(), fieldIfArrs, condCompser, calcUpdateFields)

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

	if this.MoreProcessFunc != nil {
		morers := this.MoreProcessFunc.(func(basemodel.ObjectUtilBaseIf) *errcode.ErrRsp)(destinfo)

		if morers != nil {
			ren.JSON(200, morers)
			return false
		}
	}

	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": destinfo.GetUniqId()})
	return true
}
