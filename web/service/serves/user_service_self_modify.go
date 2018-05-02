package serves

import (
	"database/sql"
	"github.com/go-martini/martini"
	"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"net/http"
	"strconv"
	//"strings"
	"errors"
	"time"
	//"reflect"

	"web/component/cfgutils"
	"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/reqparamodel"
	"web/models/servemodel"
	//"web/models/usermodel"

	"web/service/getter"
	"web/service/routers"
	"web/service/tags"
	"web/service/utils"
)

func init() {
	userServiceModifyRouterBuilderEx()
}

func userServiceModifyRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Post("/user/mine/service", CreateSelfServe)
	m.Post("/user/mine/service/:ServiceId", UpdateSelfServeInfoWithUserId)

	m.Delete("/user/mine/service/:ServiceId", RemoveSelfServeInfoWithUserId)
}

func NewUseServiceModifyControllerExObject(ctrl *UseServiceModifyControllerEx) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForCreate: ctrl.exInit4Create,
		CheckParamFuncForCreate: checkUserIdWithParamsForModify,

		ExpendInitFuncForUpdate: ctrl.exInit4Update,
		CheckParamFuncForUpdate: checkServiceIdWithParamsForModify,
		WhereCondFuncForUpdate:  compWhereCondition,

		CheckParamFuncForUpdateDelStatus: checkServiceIdWithParamsForModify,
		ExpendInitFuncForUpdateDelStatus: ctrl.exInit4Del,
		WhereCondFuncForUpdateDelStatus:  ctrl.compDelCond,
	}

	return obj
}

func NewUseServiceModifyControllerEx() *UseServiceModifyControllerEx {
	ctrl := &UseServiceModifyControllerEx{
		tableName: "web_services",
		genIdFlag: "service_id",
	}
	ctrl.initDB()
	return ctrl
}

type UseServiceModifyControllerEx struct {
	tableName string
	genIdFlag string
	db        *sql.DB
}

func (this *UseServiceModifyControllerEx) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UseServiceModifyControllerEx) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UseServiceModifyControllerEx) getDB() *sql.DB {
	return this.db
}

func (this *UseServiceModifyControllerEx) getTableName() string {
	return this.tableName
}

func (this *UseServiceModifyControllerEx) getGenIdFlag() string {
	return this.genIdFlag
}

func (this *UseServiceModifyControllerEx) exInit4Create(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)

	info, ok := reqInfo.(*servemodel.ServeInfo)
	if !ok {
		log.Println("svr exInit4Create", "req info type is not service info", userid)
		return errors.New("req info type is not service info")
	}
	//info := this.getServeInfo()

	info.UserId = int(userid)
	info.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	info.UpdateTime = info.CreateTime
	info.AuditTime = "2006-01-02 15:04:05"

	//ui := &usermodel.UserInfo{}
	userInfo, err := getter.GetModelInfoGetter().GetUserByUserId(this.getDB(), info.UserId, nil, nil)
	if err != nil {
		log.Println("svr exInit4Create", err)
		return err
	}

	info.City = userInfo.City
	if info.City < 110100 {
		info.City = 110100
	}
	info.IsZhiying = userInfo.IsZhiying
	//info.Industry

	info.ServiceId = idutils.GetId(this.getGenIdFlag())

	return nil
}

func (this *UseServiceModifyControllerEx) exInit4Update(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	serviceid, _ := strconv.ParseInt(headParams.RouterParams["ServiceId"], 10, 32)

	info, ok := reqInfo.(*servemodel.ServeInfo)
	if !ok {
		return errors.New("req info type is not service info")
	}
	//info := this.getServeInfo()

	info.UserId = int(userid)

	info.AuditTime = "2006-01-02 15:04:05"
	info.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	info.ServiceId = int(serviceid)

	return nil
}

func (this *UseServiceModifyControllerEx) exInit4Del(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	return nil
}

func (this *UseServiceModifyControllerEx) compDelCond(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	svrid, _ := strconv.ParseInt(params.RouterParams["ServiceId"], 10, 32)
	userid, _ := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)

	idWhere["service_id"] = int(svrid)
	idRlue["service_id"] = "="

	idWhere["user_id"] = int(userid)
	idRlue["user_id"] = "="

	return idWhere, idRlue
}

func CreateSelfServe(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUseServiceModifyControllerEx()
	defer ctrl.closeDB()
	obj := NewUseServiceModifyControllerExObject(ctrl)

	info := &servemodel.ServeInfo{}

	res := obj.Util_CreateObjectWithId(info, headParams, req, r)
	if res {
		tags.AddServiceTags(ctrl.getDB(), nil, info, headParams, req, r)
	}
}

func UpdateSelfServeInfoWithUserId(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUseServiceModifyControllerEx()
	defer ctrl.closeDB()
	obj := NewUseServiceModifyControllerExObject(ctrl)

	info := servemodel.NewServeInfo()

	obj.Util_UpdateObjectInfoWithId(info, headParams, req, ren, info.GetSkipFieldsForUpdate(), nil)
}

func RemoveSelfServeInfoWithUserId(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUseServiceModifyControllerEx()
	defer ctrl.closeDB()
	obj := NewUseServiceModifyControllerExObject(ctrl)

	info := servemodel.NewServeInfo()

	obj.Util_UpdateDelStatusWithId(info, headParams, req, ren)

}
