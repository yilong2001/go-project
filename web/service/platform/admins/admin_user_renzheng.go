package admins

import (
	//"database/sql"
	//"encoding/json"
	//"fmt"
	"github.com/go-martini/martini"
	//"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"errors"
	"net/http"
	"strconv"
	//"strings"
	"time"
	//"reflect"

	//"web/component/cfgutils"
	"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/reqparamodel"
	//"web/models/tokenmodel"

	"web/models/condmodel"
	//"web/models/platform/adminmodel"
	"web/models/usermodel"
	//"web/models/rendermodel"
	//"web/models/servemodel"

	"web/service/ctrlbase"
	//"web/service/getter"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	openUserRenzhengWithinAdminRouterBuilder()
}

func openUserRenzhengWithinAdminRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/admin/reviewer/users/renzheng", UpdateRenzhengUser4AdminReview)
}

func NewUserRenzhengWithinAdminControllerObject(ctrl *UserRenzhengWithinAdminController) *utils.UpdateObjectWithIdUtil {
	obj := &utils.UpdateObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		ExInitFunc:      ctrl.exInfoInitUpdate,
		CheckParamFunc:  ctrl.check4Update,
		CondCompserFunc: ctrl.condCompser,
	}

	return obj
}

func NewUserRenzhengWithinAdminController() *UserRenzhengWithinAdminController {
	ctrl := new(UserRenzhengWithinAdminController)
	ctrl.TableName = "web_users"

	ctrl.InitDB()
	return ctrl
}

type UserRenzhengWithinAdminController struct {
	ctrlbase.CtrlBaseController
}

func (this *UserRenzhengWithinAdminController) exInfoInitUpdate(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) error {

	urinfo, ok := orginfo.(*usermodel.UserInfo)
	if !ok {
		return errors.New("orginfo type is not AdminReviewerResultInfo")
	}

	if urinfo.UserId == 0 {
		return errors.New("to be tuijian userid can not be 0")
	}

	//dtinfo, ok := destinfo.(*usermodel.UserInfo)
	//if !ok {
	//  return errors.New("orginfo type is not AdminReviewerResultInfo")
	//}

	//dtinfo.IsRenzheng = 1
	//dtinfo.TuijianInfo = urinfo.TuijianInfo
	//dtinfo.TuijianImg = urinfo.TuijianImg

	uid, err := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)
	if err != nil {

		return err
	}
	dt := time.Now().Format("2006-01-02 15:04:05")

	urinfo.RenzhengTime = dt
	urinfo.RenzhengUser = int(uid)

	return nil
}

func (this *UserRenzhengWithinAdminController) check4Update(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !isAdminToken(headParams) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "not admin user"))
		return false
	}

	return true
}

func (this *UserRenzhengWithinAdminController) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {

	compser := condmodel.NewCondComposerLinker("and")
	root := compser

	where := map[string]interface{}{}
	rule := map[string]string{}

	urinfo, ok := orginfo.(*usermodel.UserInfo)
	if !ok {
		where["user_id"] = 0
		rule["user_id"] = " = "
	} else {
		where["user_id"] = urinfo.UserId
		rule["user_id"] = " = "
	}

	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	compser.SetItem(compSub)

	return root
}

func UpdateRenzhengUser4AdminReview(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	info := &usermodel.UserInfo{}

	ctrl := NewUserRenzhengWithinAdminController()
	defer ctrl.CloseDB()

	obj := NewUserRenzhengWithinAdminControllerObject(ctrl)
	res := obj.Update_With_MultiInObject(info, info, headParams, req, r, nil, []string{"IsRenzheng", "RenzhengTime", "RenzhengUser", "RenzhengInfo"})
	if res {
		ctrl.GetTX().Commit()
	}
}
