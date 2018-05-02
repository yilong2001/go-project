package others

import (
	//"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"web/component/aliutils"

	"errors"
	"net/http"
	"strconv"
	//"strings"
	"time"
	//"reflect"

	//"web/component/cfgutils"
	//"web/component/errcode"
	"web/component/idutils"
	//"web/component/wxutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/component/orderutils"
	//"web/component/rongcloud"
	//"web/dal/sqldrv"
	"web/models/basemodel"
	//"web/models/clientmodel"
	//"web/models/ordermodel"
	//"web/models/rendermodel"
	"web/models/reqparamodel"
	"web/models/usermodel"
	//"web/service/getter"
	//"web/service/immsgs"
	//"web/service/orderups"
	"web/service/ctrlbase"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userRequireRouterBuilder()
}

func userRequireRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/require", CreateUserRequire)
	m.Get("/require", GetRequires)
}

func NewUserRequireAdminQueryControllerObject(ctrl *UserRequireController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		ExpendInitFuncForGet: ctrl.exInit4Get,
		CheckParamFuncForGet: ctrl.check,
		WhereCondFuncForGet:  ctrl.adminCond,
		AppendMoreResultFunc: nil,
	}

	return obj
}

func NewUserRequireQueryControllerObject(ctrl *UserRequireController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		ExpendInitFuncForGet: ctrl.exInit4Get,
		CheckParamFuncForGet: ctrl.check,
		WhereCondFuncForGet:  ctrl.compCond4GetSub,
		AppendMoreResultFunc: nil,
	}

	return obj
}

func NewUserRequireControllerObject(ctrl *UserRequireController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		ExpendInitFuncForCreate: ctrl.exInit4Create,
		CheckParamFuncForCreate: ctrl.check,
	}

	return obj
}

func NewUserRequireController() *UserRequireController {
	ctrl := new(UserRequireController)
	ctrl.TableName = "web_requires"
	ctrl.GenIdFlag = "user_require_id"

	ctrl.InitDB()

	return ctrl
}

type UserRequireController struct {
	ctrlbase.CtrlBaseController
}

func (this *UserRequireController) exInit4Create(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	info, ok := reqInfo.(*usermodel.UserRequireInfo)
	if !ok {
		return errors.New("req info type is not user require info")
	}

	info.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	info.RequireId = idutils.GetId(this.GetGenIdFlag())

	if userid, err := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32); err == nil && userid > 0 {
		info.RequireUserId = int(userid)
	}

	err := aliutils.AliSmsMsgSend(usermodel.Const_Customer_Servier_Phone_Main, "新需求")
	if err != nil {
		log.Println("ali sms send failed, ", err)
	}

	return nil
}

func (this *UserRequireController) check(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func (this *UserRequireController) exInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	return nil
}

func (this *UserRequireController) compCond4GetSub(headParams *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	dt := time.Now().Add(-90 * 24 * time.Hour).Format("2006-01-02 15:04:05")
	idWhere["create_time"] = dt
	idRlue["create_time"] = " > "

	found := false

	mail := headParams.URLParams.Get("Email")
	if mail != "" && len(mail) > 5 {
		idWhere["email"] = "%" + mail + "%"
		idRlue["email"] = " like "
		found = true
	} else {

	}

	phone := headParams.URLParams.Get("Phone")
	if mail != "" && len(mail) > 5 {
		idWhere["phone"] = "%" + phone + "%"
		idRlue["phone"] = " like "
		found = true
	}

	if !found {
		if userid, err := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32); err == nil && userid > 10000 {
			idWhere["require_user_id"] = int(userid)
			idRlue["require_user_id"] = " = "
		}
	}

	if len(idWhere) == 1 {
		idWhere["require_user_id"] = int(1)
		idRlue["require_user_id"] = " = "
	}

	//log.Println("compWhereCondition : ", idWhere)

	return idWhere, idRlue
}

func (this *UserRequireController) adminCond(headParams *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	dt := time.Now().Add(-90 * 24 * time.Hour).Format("2006-01-02 15:04:05")
	idWhere["create_time"] = dt
	idRlue["create_time"] = " > "

	mail := headParams.URLParams.Get("Email")
	if mail != "" && len(mail) > 5 {
		idWhere["email"] = "%" + mail + "%"
		idRlue["email"] = " like "
	}

	phone := headParams.URLParams.Get("Phone")
	if mail != "" && len(mail) > 5 {
		idWhere["phone"] = "%" + phone + "%"
		idRlue["phone"] = " like "
	}

	//log.Println("compWhereCondition : ", idWhere)

	return idWhere, idRlue
}

func CreateUserRequire(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserRequireController()
	defer ctrl.CloseDB()

	obj := NewUserRequireControllerObject(ctrl)

	info := &usermodel.UserRequireInfo{}

	res := obj.Util_CreateObjectWithId(info, headParams, req, ren)
	if res {
		ctrl.GetTX().Commit()
	}
}

func GetRequires(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserRequireController()
	defer ctrl.CloseDB()

	obj := NewUserRequireQueryControllerObject(ctrl)

	info := &usermodel.UserRequireInfo{}

	res := obj.Util_GetObjectWithId(info, headParams, req, r, nil, nil)
	if res {
		ctrl.GetTX().Commit()
	}
}

func GetRequiresForAdmin(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserRequireController()
	defer ctrl.CloseDB()

	obj := NewUserRequireAdminQueryControllerObject(ctrl)

	info := &usermodel.UserRequireInfo{}

	res := obj.Util_GetObjectWithId(info, headParams, req, r, nil, nil)
	if res {
		ctrl.GetTX().Commit()
	}
}
