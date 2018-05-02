package admins

import (
	//"database/sql"
	"encoding/json"
	//"fmt"
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
	"web/models/platform/adminmodel"
	//"web/models/rendermodel"
	"web/models/coursemodel"

	"web/service/ctrlbase"
	"web/service/getter"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	adminCourseRouterBuilderEx()
}

func adminCourseRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Get("/admin/reviewer/course/all", GetCourseInfo4AdminReview)
	m.Get("/admin/reviewer/course/new", GetCourseInfo4AdminReview)
	m.Get("/admin/reviewer/course/accept", GetCourseInfo4AdminReview)
	m.Get("/admin/reviewer/course/reject", GetCourseInfo4AdminReview)

	m.Post("/admin/reviewer/service/accept", AcceptCourseByAdmin)
	m.Post("/admin/reviewer/service/reject", RejectCourseByAdmin)
}

func NewAdminCourseControllerExObject(ctrl *AdminCourseControllerEx) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.GetTableName(),
		Db:        ctrl.GetDB(),
		Tx:        ctrl.GetTX(),

		ExpendInitFuncForGet: ctrl.exInfoInit4Get,
		CheckParamFuncForGet: ctrl.check4Get,
		WhereCondFuncForGet:  ctrl.compWhereCond4Get,
	}

	return obj
}

func NewAdminCourseControllerEx() *AdminCourseControllerEx {
	ctrl := new(AdminCourseControllerEx)
	ctrl.TableName = "web_courses"
	ctrl.userId = -1
	ctrl.reviewType = adminmodel.Const_Admin_Review_Accept

	ctrl.InitDB()
	return ctrl
}

type AdminCourseControllerEx struct {
	ctrlbase.CtrlBaseController
	getType    string
	adminName  string
	userId     int
	reviewType int
	upIds      string
	upInfo     string
}

func (this *AdminCourseControllerEx) exInfoInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)

	_, ok := reqInfo.(*coursemodel.CourseMainInfo)
	if !ok {
		return errors.New("req info type is not service info")
	}

	this.getType = getAdminQueryTypeByURL(headParams.ShortUrl)

	log.Println(userid)

	//info := this.getServeInfo()

	//info.UserId = int(userid)
	//info.ServiceId = int(serviceid)

	return nil
}

func (this *AdminCourseControllerEx) compWhereCond4Get(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	if svrid, err := strconv.ParseInt(params.RouterParams["CourseId"], 10, 32); err == nil {
		idWhere["course_id"] = int(svrid)
		idRlue["course_id"] = " = "
	}

	idWhere["del_status"] = int(0)
	idRlue["del_status"] = " = "

	//log.Println("getAdminStatus", getAdminStatus(this.getType))
	adminStatus := getAdminStatus(this.getType)
	if adminStatus >= 0 {
		idWhere["audit_status"] = adminStatus
		idRlue["audit_status"] = " = "
	}

	log.Println("compWhereCondition : ", idWhere)

	return idWhere, idRlue
}

func (this *AdminCourseControllerEx) check4Get(headParams *reqparamodel.HttpReqParams, r render.Render) bool {

	if !isAdminToken(headParams) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "not admin user"))
		return false
	}

	return true
}

func (this *AdminCourseControllerEx) check4Update(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !isAdminToken(headParams) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "not admin user"))
		return false
	}

	userid, err := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "there is no userid?"))
		return false
	}

	this.userId = int(userid)

	return true
}

func GetCourseInfo4AdminReview(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewAdminCourseControllerEx()
	defer ctrl.CloseDB()
	//ctrl.getType = "new"

	obj := NewAdminCourseControllerExObject(ctrl)

	info := &coursemodel.CourseMainInfo{}

	res := obj.Util_GetObjectWithId(info, headParams, req, r, nil, []string{"CourseId", "Title", "TitleMore", "UserId", "AuditTime", "AuditName", "AuditStatus", "AuditInfo"})
	if res {
		ctrl.GetTX().Commit()
	}

}

func (this *AdminCourseControllerEx) exInfoInitUpdate(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) error {

	oi, ok := orginfo.(*adminmodel.AdminReviewerResultInfo)
	if !ok {
		return errors.New("orginfo type is not AdminReviewerResultInfo")
	}

	this.upIds = oi.Ids

	di, ok := destinfo.(*coursemodel.CourseMainInfo)
	if !ok {
		return errors.New("destinfo type is not ServeInfo")
	}

	adminuserinfo, err := getter.GetModelInfoGetter().GetAdminUserByUserId(this.GetDB(), this.userId, nil, nil)
	if err != nil {
		return err
	}

	di.AuditName = adminuserinfo.UserName
	di.AuditTime = time.Now().Format("2006-01-02 15:04:05")
	di.AuditStatus = this.reviewType
	di.AuditInfo = oi.Info

	return nil
}

func (this *AdminCourseControllerEx) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {

	compser := condmodel.NewCondComposerLinker("or")
	root := compser

	count := 0
	ids := strings.Split(this.upIds, ",")
	for _, id := range ids {
		nid := strings.TrimSpace(id)
		outid, err := strconv.ParseInt(nid, 10, 32)
		if err == nil {
			count = count + 1
			where := map[string]interface{}{"course_id": outid}
			rule := map[string]string{"course_id": " = "}
			compSub := condmodel.NewCondComposerItem(where, rule, " and ")
			if compser.Item == nil {
				compser.SetItem(compSub)
			} else {
				nextcompser := condmodel.NewCondComposerLinker("or")
				nextcompser.SetItem(compSub)
				compser.SetNext(nextcompser)
				compser = nextcompser
			}
		}
	}

	if count == 0 {
		where := map[string]interface{}{"course_id": 0}
		rule := map[string]string{"course_id": " = "}
		compSub := condmodel.NewCondComposerItem(where, rule, " and ")

		compser.SetItem(compSub)
		compser.SetNext(nil)
	}

	ji, _ := json.Marshal(root)

	log.Println("audit update condition", this.upIds)
	log.Println("audit update condition", string(ji))

	return root
}

func NewReviewAdminCourseControllerExObject(ctrl *AdminCourseControllerEx) *utils.UpdateObjectWithIdUtil {
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

func AcceptCourseByAdmin(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	orginfo := &adminmodel.AdminReviewerResultInfo{}
	destinfo := &coursemodel.CourseMainInfo{}

	ctrl := NewAdminCourseControllerEx()
	defer ctrl.CloseDB()

	ctrl.reviewType = adminmodel.Const_Admin_Review_Accept

	obj := NewReviewAdminCourseControllerExObject(ctrl)
	res := obj.Update_With_MultiInObject(orginfo, destinfo, headParams, req, r, nil, []string{"AuditTime", "AuditName", "AuditStatus", "AuditInfo"})
	if res {
		ctrl.GetTX().Commit()
	}
}

func RejectCourseByAdmin(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	orginfo := &adminmodel.AdminReviewerResultInfo{}
	destinfo := &coursemodel.CourseMainInfo{}

	ctrl := NewAdminCourseControllerEx()
	defer ctrl.CloseDB()

	ctrl.reviewType = adminmodel.Const_Admin_Review_Reject

	obj := NewReviewAdminCourseControllerExObject(ctrl)
	res := obj.Update_With_MultiInObject(orginfo, destinfo, headParams, req, r, nil, []string{"AuditTime", "AuditName", "AuditStatus", "AuditInfo"})
	if res {
		ctrl.GetTX().Commit()
	}
}
