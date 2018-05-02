package utils

import (
	//"database/sql"
	"log"
	//"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	//"strconv"
	//"time"
	//"strings"

	//"web/component/cfgutils"
	//"web/component/idutils"
	"web/component/errcode"
	//"web/component/objutils"
	//"web/component/pageutils"
	"web/component/sqlutils"

	//"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/reqparamodel"
	//"web/service/routers"
)

func (this *ObjectWithIdUtil) Util_UpdateDelStatusWithId(info basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams, req *http.Request, ren render.Render) bool {
	log.Println("Util_UpdateObjectInfoWithId start...")

	if !this.getCheckParamFuncForUpdateDelStatus().(func(*reqparamodel.HttpReqParams, render.Render) bool)(params, ren) {
		return false
	}

	var delStatus int = 0
	log.Println("req.Method : ", req.Method)
	if req.Method == "DELETE" {
		delStatus = 1
	}

	whereCond, ruleCond := this.getWhereCondFuncForUpdateDelStatus().(func(*reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string))(params)

	if this.getTX() == nil {
		err := sqlutils.Sqls_UpdateDelStatus(this.getDB(), this.getTableName(), delStatus, whereCond, ruleCond)
		if err != nil {
			ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
			return false
		}
	} else {
		err := sqlutils.Sqls_UpdateDelStatus_Tx(this.getTX(), this.getTableName(), delStatus, whereCond, ruleCond)
		if err != nil {
			ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
			return false
		}
	}

	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": ""})

	return true
}
