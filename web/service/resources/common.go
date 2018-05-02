package resources

import (
	//"database/sql"
	//"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"

	"errors"
	//"net/http"
	"strconv"
	"strings"
	"time"
	//"reflect"

	//"web/component/cfgutils"
	"web/component/errcode"
	"web/component/idutils"
	//"web/component/wxutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/component/orderutils"
	//"web/component/rongcloud"
	//"web/dal/sqldrv"
	"web/models/basemodel"
	//"web/models/clientmodel"
	"web/models/resourcemodel"
	//"web/models/rendermodel"
	"web/models/reqparamodel"
	//"web/models/usermodel"
	//"web/service/getter"
	//"web/service/immsgs"
	//"web/service/orderups"
	"web/service/ctrlbase"
	//"web/service/routers"
	"web/service/utils"
)

func NewUserResourceControllerQueryObject(ctrl *UserResourceController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		ExpendInitFuncForGet: ctrl.exInit4Get,
		CheckParamFuncForGet: ctrl.check,
		WhereCondFuncForGet:  ctrl.compCond4Get,
		AppendMoreResultFunc: nil,
	}

	return obj
}

func NewUserResourceControllerObject(ctrl *UserResourceController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		ExpendInitFuncForCreate: ctrl.exInit4Create,
		CheckParamFuncForCreate: ctrl.checkUserId,
		MoreProcessForCreate:    nil,
	}

	return obj
}

type UserResourceController struct {
	ctrlbase.CtrlBaseController
	userId       int
	courseId     int
	serviceId    int
	resourceType int
	referType    int
}

func NewUserResourceController(resourceTy, referTy int) *UserResourceController {
	ctrl := new(UserResourceController)
	ctrl.TableName = "web_user_resources"
	ctrl.GenIdFlag = "user_resource_id"

	ctrl.InitDB()

	ctrl.resourceType = resourceTy
	ctrl.referType = referTy

	return ctrl
}

func getResourceType(headParams *reqparamodel.HttpReqParams) int {
	if strings.Contains(headParams.ShortUrl, "pic") {
		return resourcemodel.Const_User_Resource_Pic
	}

	if strings.Contains(headParams.ShortUrl, "article") {
		return resourcemodel.Const_User_Resource_Article
	}

	if strings.Contains(headParams.ShortUrl, "video") {
		return resourcemodel.Const_User_Resource_Video
	}

	return resourcemodel.Const_User_Resource_Article
}

func getReferType(headParams *reqparamodel.HttpReqParams) int {
	if strings.Contains(headParams.ShortUrl, "course") {
		return resourcemodel.Const_User_Resource_Refer_Course
	}

	if strings.Contains(headParams.ShortUrl, "service") {
		return resourcemodel.Const_User_Resource_Refer_Service
	}

	return resourcemodel.Const_User_Resource_Refer_Course
}

func (this *UserResourceController) exInit4Create(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	info, ok := reqInfo.(*resourcemodel.UserResourceInfo)
	if !ok {
		return errors.New("UserResourceInfo type error")
	}

	dt := time.Now().Format("2006-01-02 15:04:05")

	info.CreateTime = dt
	info.ResourceId = idutils.GetId(this.GetGenIdFlag())
	info.UserId = this.userId
	info.ReferId = this.getReferId(headParams)
	info.ReferType = this.referType
	info.ResourceType = this.resourceType

	log.Println("new couser main info : ", info)

	return nil
}

func (this *UserResourceController) getReferId(headParams *reqparamodel.HttpReqParams) int {
	if strings.Contains(headParams.ShortUrl, "course") {
		return this.courseId
	}

	if strings.Contains(headParams.ShortUrl, "service") {
		return this.serviceId
	}

	return 0
}

func (this *UserResourceController) checkCourseId(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !strings.Contains(headParams.ShortUrl, "course") {
		return true
	}

	courseid, err := strconv.ParseInt(headParams.RouterParams["CourseId"], 10, 32)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "CourseId is not correct!"))
		return false
	}

	this.courseId = int(courseid)
	return true
}

func (this *UserResourceController) checkServiceId(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !strings.Contains(headParams.ShortUrl, "service") {
		return true
	}

	serviceid, err := strconv.ParseInt(headParams.RouterParams["ServiceId"], 10, 32)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "ServiceId is not correct!"))
		return false
	}

	this.serviceId = int(serviceid)
	return true
}

func (this *UserResourceController) checkUserId(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	if !this.checkCourseId(headParams, r) {
		return false
	}

	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	this.userId = int(userid)

	return true
}

func (this *UserResourceController) check(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !this.checkCourseId(headParams, r) {
		return false
	}

	if !this.checkServiceId(headParams, r) {
		return false
	}

	return true
}

func (this *UserResourceController) exInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	return nil
}

func (this *UserResourceController) compCond4Get(headParams *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	idWhere["del_status"] = int(0)
	idRlue["del_status"] = " = "

	idWhere["resource_type"] = this.resourceType
	idRlue["resource_type"] = " = "

	idWhere["refer_type"] = this.referType
	idRlue["refer_type"] = " = "

	idWhere["refer_id"] = this.getReferId(headParams)
	idRlue["refer_id"] = " = "

	log.Println("UserResourceController get condition : ", idWhere)

	return idWhere, idRlue
}
