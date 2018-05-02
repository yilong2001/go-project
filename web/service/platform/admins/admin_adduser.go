package admins

import (
	//"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"

	"errors"
	"net/http"
	"strconv"
	//"strings"
	"fmt"
	"time"
	//"reflect"

	//"web/component/cfgutils"
	"web/component/errcode"
	"web/component/idutils"
	//"web/component/wxutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/component/filterutils"
	//"web/component/orderutils"
	//"web/component/randutils"
	//"web/component/rongcloud"
	//"web/dal/sqldrv"
	"web/models/basemodel"
	//"web/models/clientmodel"
	//"web/models/condmodel"
	//"web/models/coursemodel"
	"web/models/platform/adminmodel"
	"web/models/rendermodel"
	"web/models/reqparamodel"
	"web/models/usermodel"
	//"web/service/getter"
	//"web/service/immsgs"
	//"web/service/orderups"
	"web/service/ctrlbase"
	"web/service/routers"
	"web/service/users"
	"web/service/utils"
)

func init() {
	adminAdduserRouterBuilder()
}

func adminAdduserRouterBuilder() {
	m := routers.GetRouterHandler()
	m.Get("/admin/reviewer/adduser", GetAdminAddusers)
	m.Post("/admin/reviewer/adduser", AdminAddUser)
}

type AdminAddusersController struct {
	ctrlbase.CtrlBaseController
	adminId int
}

func NewAdminAddusersControllerObject(ctrl *AdminAddusersController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		ExpendInitFuncForCreate: ctrl.exInit4Create,
		CheckParamFuncForCreate: ctrl.checkUserId,
		MoreProcessForCreate:    nil,

		ExpendInitFuncForGet: ctrl.exInit4Get,
		CheckParamFuncForGet: ctrl.checkUserId,
		WhereCondFuncForGet:  ctrl.compCond4Get,
		//WhereCondComposerForGet: ctrl.compComposer4OpenGet,
		AppendMoreResultFunc: nil,
	}

	return obj
}

func NewAdminAddusersController() *AdminAddusersController {
	ctrl := new(AdminAddusersController)
	ctrl.TableName = "web_admin_addusers"
	ctrl.GenIdFlag = "admin_adduser_id"

	ctrl.InitDB()

	return ctrl
}

func (this *AdminAddusersController) exInit4Create(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	info, ok := reqInfo.(*adminmodel.AdminAdduserInfo)
	if !ok {
		return errors.New("req info type wrong, not AdminAdduserInfo")
	}

	dt := time.Now().Format("2006-01-02 15:04:05")

	info.AddId = idutils.GetId(this.GetGenIdFlag())
	info.CreateTime = dt

	tmp := fmt.Sprint(1000+this.adminId) + fmt.Sprint(info.AddId)
	//_, err := strconv.ParseInt(tmp, 10, 32)
	//if err != nil {
	//	return err
	//}

	info.UserLoginId = tmp

	info.AdminUserId = this.adminId

	log.Println("new AdminAdduserInfo : ", info)

	return nil
}

func (this *AdminAddusersController) checkUserId(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !isAdminToken(headParams) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not admin user!"))
		return false
	}

	aid, err := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not wrong!"))
		return false
	}

	this.adminId = int(aid)
	return true
}

func (this *AdminAddusersController) exInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *AdminAddusersController) compCond4Get(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	if this.adminId > 1 {
		idWhere["admin_user_id"] = this.adminId
		idRlue["admin_user_id"] = " = "
	}

	return idWhere, idRlue
}

func AdminAddUser(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewAdminAddusersController()
	defer ctrl.CloseDB()

	obj := NewAdminAddusersControllerObject(ctrl)

	info := &adminmodel.AdminAdduserInfo{}
	fakeren := &rendermodel.FakeMrtiniRender{}

	res := obj.Util_CreateObjectWithId(info, headParams, req, fakeren)
	if res {
		reginfo := &usermodel.UserRegisterInfo{}
		reginfo.Phone = info.UserLoginId
		reginfo.Password = "4280d89a5a03f812751f504cc10ee8a5" //randutils.BuildMd5PWPhoneStringV2("123456", "")

		err, tokendb := users.DoRegister(ctrl.GetDB(), true, reginfo, headParams, req, false)
		if err != nil {
			r.JSON(200, err)
			return
		}

		log.Println(tokendb)

		ctrl.GetTX().Commit()

		r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": info.UserLoginId})
	} else {
		r.JSON(fakeren.GetStatus(), fakeren.GetVal())
	}
}

func GetAdminAddusers(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewAdminAddusersController()
	defer ctrl.CloseDB()

	obj := NewAdminAddusersControllerObject(ctrl)

	info := &adminmodel.AdminAdduserInfo{}

	res := obj.Util_GetObjectWithId(info, headParams, req, r, nil, nil)
	if res {
		ctrl.GetTX().Commit()
	}
}
