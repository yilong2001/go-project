package users

import (
	"database/sql"
	"github.com/go-martini/martini"
	"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	//"reflect"

	"web/component/cfgutils"
	"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/firmmodel"
	"web/models/reqparamodel"
	"web/models/usermodel"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userSelfRouterBuilder()
}

func userSelfRouterBuilder() {
	m := routers.GetRouterHandler()

	//Get("user/mine")
	//Post("user/mine")
	//Get("user/open")
	//Get("user/open/:userid")

	m.Post("/user/mine", UserSelfUpdateUserInfo)
	m.Get("/user/mine", GetUserSelfInfoWithId)

	m.Get("/user/open/:DestId", GetUserOpenInfoWithId)
}

func NewUseSelfControllerObject(ctrl *UseSelfController) *utils.ObjectWithIdUtil {
	return &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForCreate: ctrl.exInfoInit4Create,
		CheckParamFuncForCreate: checkUserIdWithParamsForModify,

		ExpendInitFuncForGet: ctrl.exInfoInit4Get,
		CheckParamFuncForGet: checkUserIdWithParamsForQuery,
		WhereCondFuncForGet:  ctrl.compWhereCond,

		ExpendInitFuncForUpdate: ctrl.exInfoInit4Update,
		CheckParamFuncForUpdate: checkUserIdWithParamsForModify,
		WhereCondFuncForUpdate:  ctrl.compWhereCond,
	}
}

func NewUseSelfControllerObjectForOpen(ctrl *UseSelfController) *utils.ObjectWithIdUtil {
	return &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInit4Open,
		CheckParamFuncForGet: ctrl.check4Open,
		WhereCondFuncForGet:  ctrl.compWhereCond,
	}
}

func NewUseSelfController() *UseSelfController {
	ctrl := &UseSelfController{
		tableName: "web_users",
		genIdFlag: "user_id",
	}

	ctrl.initDB()
	return ctrl
}

type UseSelfController struct {
	tableName string
	db        *sql.DB
	genIdFlag string
}

func (this *UseSelfController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UseSelfController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UseSelfController) getDB() *sql.DB {
	return this.db
}

func (this *UseSelfController) getTableName() string {
	return this.tableName
}

func (this *UseSelfController) getGenIdFlag() string {
	return this.genIdFlag
}

func (this *UseSelfController) exInfoInit4Create(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	info, ok := reqInfo.(*usermodel.UserInfo)
	if !ok {
		return errors.New("user info type error")
	}

	info.SimpleFirm = firmmodel.MakeCompareName(info.Firm)

	return nil
	//TODO
}

func (this *UseSelfController) exInfoInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)

	info, ok := reqInfo.(*usermodel.UserInfo)
	if !ok {
		return errors.New("req info type is not user info")
	}

	info.UserId = int(userid)

	return nil
}

func (this *UseSelfController) exInit4Open(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["DestId"], 10, 32)

	info, ok := reqInfo.(*usermodel.UserInfo)
	if !ok {
		return errors.New("req info type is not user info")
	}

	info.UserId = int(userid)

	return nil
}

func (this *UseSelfController) check4Open(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["DestId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "Dest UserId is not correct!"))
		return false
	}

	return true
}

func (this *UseSelfController) exInfoInit4Update(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)

	info, ok := reqInfo.(*usermodel.UserInfo)
	if !ok {
		return errors.New("req info type is not user info")
	}

	info.SimpleFirm = firmmodel.MakeCompareName(info.Firm)
	info.UserId = int(userid)

	info.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (this *UseSelfController) compWhereCond(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	log.Println("compWhereCondition HttpReqParams", params)

	uid := params.RouterParams["UserId"]
	if params.RouterParams["DestId"] != "" {
		uid = params.RouterParams["DestId"]
	}

	id, err := strconv.ParseInt(uid, 10, 32)
	if err == nil {
		idWhere["user_id"] = int(id)
		idRlue["user_id"] = "="
	} else {
		idWhere["user_id"] = int(0)
		idRlue["user_id"] = "="
	}

	log.Println("compWhereCondition : ", idWhere)

	return idWhere, idRlue
}

func GetUserSelfInfoWithId(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	userInfo := &usermodel.UserInfo{}

	ctrl := NewUseSelfController()
	defer ctrl.closeDB()
	obj := NewUseSelfControllerObject(ctrl)

	obj.Util_GetObjectWithId(userInfo, headParams, req, r, userInfo.GetSkipFieldsForSelfQuery(), nil)
}

func GetUserOpenInfoWithId(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	userInfo := &usermodel.UserInfo{}

	ctrl := NewUseSelfController()
	defer ctrl.closeDB()
	obj := NewUseSelfControllerObjectForOpen(ctrl)

	obj.Util_GetObjectWithId(userInfo, headParams, req, r, userInfo.GetSkipFieldsForOpenQuery(), nil)
}

func UserSelfUpdateUserInfo(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)
	for _, pf := range headParams.PostFields {
		if strings.ToLower(pf) == "firm" {
			headParams.PostFields = append(headParams.PostFields, "SimpleFirm")
			break
		}
	}

	userInfo := usermodel.NewUserInfo()
	ctrl := NewUseSelfController()
	defer ctrl.closeDB()
	obj := NewUseSelfControllerObject(ctrl)

	skips := userInfo.GetSkipFieldsForUpdate()

	obj.Util_UpdateObjectInfoWithId(userInfo, headParams, req, ren, skips, nil)
}
