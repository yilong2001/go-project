package serves

import (
	//"database/sql"
	//"log"
	"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"net/http"
	"strconv"
	//"strings"
	"errors"
	//"time"
	//"reflect"

	//"web/component/cfgutils"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/dal/sqldrv"
	"web/component/errcode"

	"web/models/basemodel"
	"web/models/reqparamodel"
	"web/models/servemodel"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userServiceRouterBuilderEx()
}

func userServiceRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Get("/user/open/:DestId/service", GetOpenServeInfoWithUserId)
	m.Get("/user/open/:DestId/service/:ServiceId", GetOpenServeInfoWithUserId)
}

func NewUseServiceControllerExObject(ctrl *ServeDefaultController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInit4OpenGet,
		CheckParamFuncForGet: ctrl.check4OpenGet,
		WhereCondFuncForGet:  ctrl.compWhereCond4OpenUserGet,
		AppendMoreResultFunc: ctrl.appendUserInfo4Result,
	}

	return obj
}

func (this *ServeDefaultController) check4OpenGet(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["DestId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "DestId is not correct!"))
		return false
	}

	return true
}

func (this *ServeDefaultController) exInit4OpenGet(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	serviceid, _ := strconv.ParseInt(headParams.RouterParams["ServiceId"], 10, 32)

	info, ok := reqInfo.(*servemodel.ServeInfo)
	if !ok {
		return errors.New("req info type is not service info")
	}
	//info := this.getServeInfo()

	info.UserId = int(userid)
	info.ServiceId = int(serviceid)

	return nil
}

func (this *ServeDefaultController) compWhereCond4OpenUserGet(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	if svrid, err := strconv.ParseInt(params.RouterParams["ServiceId"], 10, 32); err == nil {
		idWhere["service_id"] = int(svrid)
		idRlue["service_id"] = "="
	}

	userid, _ := strconv.ParseInt(params.RouterParams["DestId"], 10, 32)
	idWhere["user_id"] = int(userid)
	idRlue["user_id"] = " = "

	idWhere["del_status"] = int(0)
	idRlue["del_status"] = " = "

	//idWhere["audit_status"] = int(1)
	//idRlue["audit_status"] = " = "

	return idWhere, idRlue
}

func GetOpenServeInfoWithUserId(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewServeDefaultController()
	defer ctrl.closeDB()
	obj := NewUseServiceControllerExObject(ctrl)

	info := &servemodel.ServeInfo{}

	obj.Util_GetObjectWithId(info, headParams, req, r, info.GetSkipFieldsForOpenQuery(), nil)
}
