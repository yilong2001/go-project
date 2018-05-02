package orders

import (
	//"database/sql"
	"github.com/go-martini/martini"
	"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"net/http"
	"strconv"
	//"strings"
	//"errors"
	//"time"
	//"reflect"

	//"web/component/cfgutils"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"
	"web/component/orderutils"

	"web/component/errcode"
	//"web/dal/sqldrv"
	"web/models/ordermodel"

	"web/models/clientmodel"
	//"web/models/condmodel"

	"web/models/basemodel"
	"web/models/reqparamodel"
	//"web/models/servemodel"
	"web/models/usermodel"
	"web/service/ctrlbase"
	"web/service/getter"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userServiceCommentRouterBuilderEx()
}

func userServiceCommentRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Get("/service/:ServiceId/comment", GetServiceCommentInfo)
	m.Get("/course/:CourseId/comment", GetCourseCommentInfo)
}

func NewUserServiceCommentQueryObject(ctrl *UserServiceCommentController) *utils.ObjectWithIdUtil {
	return &utils.ObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		CheckParamFuncForGet: ctrl.check4Service,
		ExpendInitFuncForGet: ctrl.exInit4Get,
		WhereCondFuncForGet:  ctrl.compServiceCond,
		AppendMoreResultFunc: ctrl.appendUserInfo4Result,
	}
}

func NewUserCourseCommentQueryObject(ctrl *UserServiceCommentController) *utils.ObjectWithIdUtil {
	return &utils.ObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		CheckParamFuncForGet: ctrl.check4Course,
		ExpendInitFuncForGet: ctrl.exInit4Get,
		WhereCondFuncForGet:  ctrl.compCourseCond,
		AppendMoreResultFunc: ctrl.appendUserInfo4Result,
	}
}

func NewUserServiceCommentController() *UserServiceCommentController {
	ctrl := new(UserServiceCommentController)
	ctrl.TableName = "web_order_comments"

	ctrl.InitDB()

	return ctrl
}

type UserServiceCommentController struct {
	ctrlbase.CtrlBaseController
	serviceId int
}

func (this *UserServiceCommentController) check4Service(headParams *reqparamodel.HttpReqParams, r render.Render) bool {

	if !utils.IsFieldCorrectWithRule("service_id", headParams.RouterParams["ServiceId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "ServiceId is not correct!"))
		return false
	}

	destid, _ := strconv.ParseInt(headParams.RouterParams["ServiceId"], 10, 32)

	this.serviceId = int(destid)

	return true
}

func (this *UserServiceCommentController) check4Course(headParams *reqparamodel.HttpReqParams, r render.Render) bool {

	if !utils.IsFieldCorrectWithRule("course_id", headParams.RouterParams["CourseId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "CourseId is not correct!"))
		return false
	}

	destid, _ := strconv.ParseInt(headParams.RouterParams["CourseId"], 10, 32)

	this.serviceId = int(destid)
	return true
}

func (this *UserServiceCommentController) exInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	return nil
}

func (this *UserServiceCommentController) compServiceCond(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	idWhere["service_id"] = this.serviceId
	idRlue["service_id"] = " = "

	idWhere["order_type"] = orderutils.Order_Type_Course
	idRlue["order_type"] = " < "

	return idWhere, idRlue
}

func (this *UserServiceCommentController) compCourseCond(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	idWhere["service_id"] = this.serviceId
	idRlue["service_id"] = " = "

	idWhere["order_type"] = orderutils.Order_Type_Course
	idRlue["order_type"] = " = "

	return idWhere, idRlue
}

func (this *UserServiceCommentController) appendUserInfo4Result(result *[]interface{}) *[]interface{} {
	if len(*result) < 1 {
		return result
	}

	//userids := []int{}
	clientInfos := []interface{}{}

	for _, cmtif := range *result {
		if cmtInfo, ok := cmtif.(ordermodel.OrderCommentInfo); ok {
			clientInfo := &clientmodel.ClientCommentInfo{}
			clientInfo.CommentInfo = &cmtInfo

			ui := &usermodel.UserInfo{}
			userInfo, err := getter.GetModelInfoGetter().GetUserByUserId(this.GetDB(), cmtInfo.CustomerId, ui.GetSkipFieldsForOpenQuery(), nil)
			if err == nil {
				clientInfo.UserInfo = userInfo
			} else {
				log.Print("user query wrong: *** ", err)
			}

			clientInfos = append(clientInfos, clientInfo)
		}
	}

	return &clientInfos
}

func GetServiceCommentInfo(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserServiceCommentController()
	defer ctrl.CloseDB()

	obj := NewUserServiceCommentQueryObject(ctrl)

	info := &ordermodel.OrderCommentInfo{}

	res := obj.Util_GetObjectWithId(info, headParams, req, ren, nil, nil)
	if res {
		ctrl.GetTX().Commit()
	}
}

func GetCourseCommentInfo(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserServiceCommentController()
	defer ctrl.CloseDB()

	obj := NewUserCourseCommentQueryObject(ctrl)

	info := &ordermodel.OrderCommentInfo{}

	res := obj.Util_GetObjectWithId(info, headParams, req, ren, nil, nil)
	if res {
		ctrl.GetTX().Commit()
	}
}
