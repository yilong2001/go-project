package admins

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	"web/models/reqparamodel"
	"web/models/tokenmodel"

	"web/models/condmodel"
	"web/models/platform/adminmodel"
	//"web/models/rendermodel"
	"web/models/servemodel"

	"web/service/getter"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	adminServiceRouterBuilderEx()
}

func adminServiceRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Get("/admin/reviewer/service/all", GetServeInfo4AdminReview)
	m.Get("/admin/reviewer/service/new", GetServeInfo4AdminReview)
	m.Get("/admin/reviewer/service/accept", GetServeInfo4AdminReview)
	m.Get("/admin/reviewer/service/reject", GetServeInfo4AdminReview)

	m.Post("/admin/reviewer/service/accept", AcceptServiceByAdmin)
	m.Post("/admin/reviewer/service/reject", RejectServiceByAdmin)
}

func isAdminToken(para *reqparamodel.HttpReqParams) bool {
	if para.TokenParams["TokenType"] == fmt.Sprint(tokenmodel.Const_Token_Type_Admin) {
		return true
	}

	return false
}

func getAdminQueryTypeByURL(url string) string {
	getType := ""
	if strings.Contains(url, "/all") {
		getType = "all"
	} else if strings.Contains(url, "/new") {
		getType = "new"
	} else if strings.Contains(url, "/accept") {
		getType = "accept"
	} else {
		getType = "reject"
	}

	return getType
}

func getAdminStatus(gettype string) int {
	if gettype == "new" {
		return int(adminmodel.Const_Admin_Review_New)
	}

	if gettype == "accept" {
		return int(adminmodel.Const_Admin_Review_Accept)
	}

	if gettype == "reject" {
		return int(adminmodel.Const_Admin_Review_Reject)
	}

	return -1
}

func NewAdminServiceControllerExObject(ctrl *AdminServiceControllerEx) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInfoInit4Get,
		CheckParamFuncForGet: ctrl.check4Get,
		WhereCondFuncForGet:  ctrl.compWhereCond4Get,
	}

	return obj
}

func NewAdminServiceControllerEx() *AdminServiceControllerEx {
	ctrl := &AdminServiceControllerEx{
		tableName:  "web_services",
		userId:     -1,
		reviewType: adminmodel.Const_Admin_Review_Accept,
	}
	ctrl.initDB()
	return ctrl
}

type AdminServiceControllerEx struct {
	tableName  string
	db         *sql.DB
	getType    string
	adminName  string
	userId     int
	reviewType int
	upIds      string
	upInfo     string
}

func (this *AdminServiceControllerEx) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *AdminServiceControllerEx) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *AdminServiceControllerEx) getDB() *sql.DB {
	return this.db
}
func (this *AdminServiceControllerEx) getTableName() string {
	return this.tableName
}

func (this *AdminServiceControllerEx) exInfoInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)

	_, ok := reqInfo.(*servemodel.ServeInfo)
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

func (this *AdminServiceControllerEx) compWhereCond4Get(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	if svrid, err := strconv.ParseInt(params.RouterParams["ServiceId"], 10, 32); err == nil {
		idWhere["service_id"] = int(svrid)
		idRlue["service_id"] = " = "
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

func (this *AdminServiceControllerEx) check4Get(headParams *reqparamodel.HttpReqParams, r render.Render) bool {

	if !isAdminToken(headParams) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "not admin user"))
		return false
	}

	return true
}

func (this *AdminServiceControllerEx) check4Update(headParams *reqparamodel.HttpReqParams, r render.Render) bool {

	userid, err := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "there is no userid?"))
		return false
	}

	if userid <= 0 || userid > 10000 {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "userid range is error"))
		return false
	}

	this.userId = int(userid)

	return true
}

func GetServeInfo4AdminReview(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewAdminServiceControllerEx()
	defer ctrl.closeDB()
	//ctrl.getType = "new"

	obj := NewAdminServiceControllerExObject(ctrl)

	info := &servemodel.ServeInfo{}

	obj.Util_GetObjectWithId(info, headParams, req, r, info.GetSkipFieldsForAdmin(), nil)

}

func (this *AdminServiceControllerEx) exInfoInitUpdate(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) error {

	oi, ok := orginfo.(*adminmodel.AdminReviewerResultInfo)
	if !ok {
		return errors.New("orginfo type is not AdminReviewerResultInfo")
	}

	this.upIds = oi.Ids

	di, ok := destinfo.(*servemodel.ServeInfo)
	if !ok {
		return errors.New("destinfo type is not ServeInfo")
	}

	adminuserinfo, err := getter.GetModelInfoGetter().GetAdminUserByUserId(this.getDB(), this.userId, nil, nil)
	if err != nil {
		return err
	}

	di.AuditName = adminuserinfo.UserName
	di.AuditTime = time.Now().Format("2006-01-02 15:04:05")
	di.AuditStatus = this.reviewType
	di.AuditInfo = oi.Info

	return nil
}

func (this *AdminServiceControllerEx) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {

	compser := condmodel.NewCondComposerLinker("or")
	root := compser

	count := 0
	ids := strings.Split(this.upIds, ",")
	for _, id := range ids {
		nid := strings.TrimSpace(id)
		outid, err := strconv.ParseInt(nid, 10, 32)
		if err == nil {
			count = count + 1
			where := map[string]interface{}{"service_id": outid}
			rule := map[string]string{"service_id": " = "}
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
		where := map[string]interface{}{"service_id": 0}
		rule := map[string]string{"service_id": " = "}
		compSub := condmodel.NewCondComposerItem(where, rule, " and ")

		compser.SetItem(compSub)
		compser.SetNext(nil)
	}

	// compser := condmodel.NewCondComposer(" or ")

	// count := 0
	// ids := strings.Split(this.upIds, ",")
	// for _, id := range ids {
	// 	nid := strings.TrimSpace(id)
	// 	outid, err := strconv.ParseInt(nid, 10, 32)
	// 	if err == nil {
	// 		count = count + 1
	// 		where := map[string]interface{}{"service_id": outid}
	// 		rule := map[string]string{"service_id": " = "}
	// 		compSub := condmodel.NewCondComposerItem(where, rule, " and ")
	// 		compser.AddItem(compSub)
	// 	}
	// }

	// if count == 0 {
	// 	where := map[string]interface{}{"service_id": 0}
	// 	rule := map[string]string{"service_id": " = "}
	// 	compSub := condmodel.NewCondComposerItem(where, rule, " and ")
	// 	compser.AddItem(compSub)
	// }

	ji, _ := json.Marshal(root)

	log.Println("audit update condition", this.upIds)
	log.Println("audit update condition", string(ji))

	return root
}

func NewReviewAdminServiceControllerExObject(ctrl *AdminServiceControllerEx) *utils.UpdateObjectWithIdUtil {
	obj := &utils.UpdateObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExInitFunc:      ctrl.exInfoInitUpdate,
		CheckParamFunc:  ctrl.check4Update,
		CondCompserFunc: ctrl.condCompser,
	}

	return obj
}

func AcceptServiceByAdmin(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	orginfo := &adminmodel.AdminReviewerResultInfo{}
	destinfo := &servemodel.ServeInfo{}

	ctrl := NewAdminServiceControllerEx()
	defer ctrl.closeDB()

	ctrl.reviewType = adminmodel.Const_Admin_Review_Accept

	obj := NewReviewAdminServiceControllerExObject(ctrl)
	obj.Update_With_MultiInObject(orginfo, destinfo, headParams, req, r, nil, destinfo.GetSpecFieldsForAdminUpdate())
}

func RejectServiceByAdmin(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	orginfo := &adminmodel.AdminReviewerResultInfo{}
	destinfo := &servemodel.ServeInfo{}

	ctrl := NewAdminServiceControllerEx()
	defer ctrl.closeDB()

	ctrl.reviewType = adminmodel.Const_Admin_Review_Reject

	obj := NewReviewAdminServiceControllerExObject(ctrl)
	obj.Update_With_MultiInObject(orginfo, destinfo, headParams, req, r, nil, destinfo.GetSpecFieldsForAdminUpdate())
}
