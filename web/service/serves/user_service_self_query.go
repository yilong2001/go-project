package serves

import (
	"database/sql"
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

	"web/component/cfgutils"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/reqparamodel"
	"web/models/servemodel"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userSelfServiceRouterBuilderEx()
}

func userSelfServiceRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Get("/user/mine/service", GetSelfServeInfoWithSelfId)
	m.Get("/user/mine/service/:ServiceId", GetSelfServeInfoWithSelfId)
}

func NewUseSelfServiceControllerExObject(ctrl *UseSelfServiceControllerEx) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInfoInit4Get,
		CheckParamFuncForGet: checkUserIdWithParamsForQuery,
		WhereCondFuncForGet:  compWhereCondition,
	}
	return obj
}

func NewUseSelfServiceControllerEx() *UseSelfServiceControllerEx {
	ctrl := &UseSelfServiceControllerEx{
		tableName: "web_services",
		genIdFlag: "service_id",
	}
	ctrl.initDB()
	return ctrl
}

type UseSelfServiceControllerEx struct {
	tableName string
	genIdFlag string
	db        *sql.DB
}

func (this *UseSelfServiceControllerEx) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UseSelfServiceControllerEx) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UseSelfServiceControllerEx) getDB() *sql.DB {
	return this.db
}

func (this *UseSelfServiceControllerEx) getTableName() string {
	return this.tableName
}

func (this *UseSelfServiceControllerEx) getGenIdFlag() string {
	return this.genIdFlag
}

func (this *UseSelfServiceControllerEx) exInfoInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
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

func GetSelfServeInfoWithSelfId(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUseSelfServiceControllerEx()
	defer ctrl.closeDB()
	obj := NewUseSelfServiceControllerExObject(ctrl)

	info := &servemodel.ServeInfo{}

	obj.Util_GetObjectWithId(info, headParams, req, r, info.GetSkipFieldsForOpenQuery(), nil)
}
