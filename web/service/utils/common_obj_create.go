package utils

import (
	//"database/sql"
	"errors"
	"log"
	//"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	//"strconv"
	//"time"

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

func (this *ObjectWithIdUtil) Util_CreateObjectWithId(info basemodel.ObjectUtilBaseIf,
	params *reqparamodel.HttpReqParams,
	req *http.Request,
	r render.Render) bool {

	if !this.getCheckParamFuncForCreate().(func(*reqparamodel.HttpReqParams, render.Render) bool)(params, r) {
		return false
	}

	//info := this.getObjInfo()

	if this.FormUnParseFlag == 0 {
		err := objutils.ParseObjectWithForm(info, req)
		if err != nil {
			log.Println("Util_CreateObjectWithId parse form faield : ", err)
			r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, err.Error()))
			return false
		}
	}

	err := this.getExpendObjInitFuncForCreate().(func(basemodel.ObjectUtilBaseIf, *reqparamodel.HttpReqParams) error)(info, params)
	if err != nil {
		log.Println("Util_CreateObjectWithId ex create failed : ", err)
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, err.Error()))
		return false
	}

	log.Println("obj info : ", info)

	fieldIfArrs, _ := info.GetWholeFields()
	res, fns := IsFieldsValueOk(fieldIfArrs)
	if !res {
		detail := "the info (" + fns + ") is not correct"
		log.Println("Util_CreateObjectWithId get fields failed : ", detail)
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, detail))
		return false
	}

	msqls, args := sqlutils.Sqls_CompInsert(this.getTableName(), fieldIfArrs)
	log.Println("Sqls_CompInsert sqls : ", msqls)

	if this.getTX() == nil {
		err = sqlutils.Sqls_Do_PrepareAndExec(this.getDB(), msqls, args)
		if err != nil {
			log.Println("Util_CreateObjectWithId sql exec failed : ", err)
			r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
			return false
		}
	} else {
		err = sqlutils.Sqls_Do_PrepareAndExec_Tx(this.getTX(), msqls, args)
		if err != nil {
			log.Println("Util_CreateObjectWithId sql exec tx failed : ", err)
			r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
			return false
		}
	}

	if this.getMoreProcessForCreate() != nil {
		err := this.getMoreProcessForCreate().(func(basemodel.ObjectUtilBaseIf) error)(info)
		if err != nil {
			log.Println("Util_CreateObjectWithId more create failed : ", err)
			r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
			return false
		}
	}

	r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": info.GetUniqId()})
	return true
}

func (this *ObjectWithIdUtil) CreateObjectWithInfo(info basemodel.ObjectUtilBaseIf) error {
	fieldIfArrs, _ := info.GetWholeFields()
	res, fns := IsFieldsValueOk(fieldIfArrs)
	if !res {
		detail := "the info (" + fns + ") is not correct"
		return errors.New(detail)
		//return errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, detail)
	}

	msqls, args := sqlutils.Sqls_CompInsert(this.getTableName(), fieldIfArrs)

	if this.getTX() != nil {
		err := sqlutils.Sqls_Do_PrepareAndExec_Tx(this.getTX(), msqls, args)
		if err != nil {
			return err
			//return errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error())
		}
	} else {
		err := sqlutils.Sqls_Do_PrepareAndExec(this.getDB(), msqls, args)
		if err != nil {
			return err
			//return errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error())
		}
	}

	return nil
}
