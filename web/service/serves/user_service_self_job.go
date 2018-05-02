package serves

import (
	"database/sql"
	//"fmt"
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
	"web/component/sqlutils"

	"web/component/errcode"
	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/reqparamodel"
	"web/models/servemodel"
	"web/service/getter"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userServiceJobRouterBuilderEx()
}

func userServiceJobRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Post("/user/mine/service/:ServiceId/job/:JobId", AddJobOnUserService)
	m.Delete("/user/mine/service/:ServiceId/job/:JobId", RemoveJobFromUserService)
	m.Get("/user/mine/service/:ServiceId/job", GetJobsOnUserService)
}

func NewServiceJobJobControllerObject(ctrl *ServiceJobJobController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		CheckParamFuncForCreate: ctrl.checkPara,
		ExpendInitFuncForCreate: ctrl.exInit4Create,

		CheckParamFuncForGet: ctrl.checkGetPara,
		ExpendInitFuncForGet: ctrl.exInit4Get,
		WhereCondFuncForGet:  ctrl.compGetCond,

		CheckParamFuncForUpdateDelStatus: ctrl.checkPara,
		ExpendInitFuncForUpdateDelStatus: ctrl.exInit4Del,
		WhereCondFuncForUpdateDelStatus:  ctrl.compUpdateCond,
	}
	return obj
}

func NewServiceJobJobController() *ServiceJobJobController {
	ctrl := &ServiceJobJobController{
		tableName: "web_service_jobs",
		genIdFlag: "service_job_id",
	}
	ctrl.initDB()
	return ctrl
}

type ServiceJobJobController struct {
	tableName string
	genIdFlag string
	db        *sql.DB
}

func (this *ServiceJobJobController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *ServiceJobJobController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *ServiceJobJobController) getDB() *sql.DB {
	return this.db
}

func (this *ServiceJobJobController) getTableName() string {
	return this.tableName
}

func (this *ServiceJobJobController) getGenIdFlag() string {
	return this.genIdFlag
}

func (this *ServiceJobJobController) isUserServiceJobExist(headParams *reqparamodel.HttpReqParams) bool {
	uid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	sid, _ := strconv.ParseInt(headParams.RouterParams["ServiceId"], 10, 32)
	jid, _ := strconv.ParseInt(headParams.RouterParams["JobId"], 10, 32)

	whereCond := map[string]interface{}{"service_id": int(sid), "job_id": int(jid), "service_user_id": int(uid)}
	ruleCondition := map[string]string{"service_id": " = ", "job_id": " = ", "service_user_id": " = "}

	ct, err := sqlutils.Sqls_GetCounts(this.getDB(), this.getTableName(), this.getGenIdFlag(), whereCond, ruleCondition)
	if err != nil {
		log.Println(err)
		return false
	}

	if ct == 0 {
		return false
	}

	return true
}

func (this *ServiceJobJobController) exInit4Create(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	serviceid, _ := strconv.ParseInt(headParams.RouterParams["ServiceId"], 10, 32)
	jobid, _ := strconv.ParseInt(headParams.RouterParams["JobId"], 10, 32)

	info, ok := reqInfo.(*servemodel.ServeJobInfo)
	if !ok {
		return errors.New("req info type is not service job info")
	}

	db := sqldrv.GetDb(cfgutils.GetWebApiConfig())
	defer db.Close()

	svrInfo, err := getter.GetModelInfoGetter().GetServiceByServiceId(db, int(serviceid), nil, nil)
	if err != nil {
		return err
	}

	if svrInfo.UserId != int(userid) {
		return errors.New("user id in post form is wrong!")
	}

	info.ServiceId = int(serviceid)
	info.ServiceUserId = int(userid)
	info.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	jobInfo, err := getter.GetModelInfoGetter().GetJobByJobId(db, int(jobid), nil, nil)
	if err != nil {
		return err
	}

	info.JobId = int(jobid)
	info.JobUserId = jobInfo.UserId

	info.ServiceJobId = idutils.GetId(this.getGenIdFlag())

	return nil
}

func (this *ServiceJobJobController) compUpdateCond(headParams *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	serviceid, _ := strconv.ParseInt(headParams.RouterParams["ServiceId"], 10, 32)
	jobid, _ := strconv.ParseInt(headParams.RouterParams["JobId"], 10, 32)

	idWhere := map[string]interface{}{"service_id": int(serviceid), "job_id": int(jobid), "service_user_id": int(userid)}
	ruleCond := map[string]string{"service_id": " = ", "job_id": " = ", "service_user_id": " = "}

	return idWhere, ruleCond
}

func (this *ServiceJobJobController) updateDelStatus(status int, headParams *reqparamodel.HttpReqParams) error {
	upateField := make(map[string]interface{})
	upateField["del_status"] = status

	where, rule := this.compUpdateCond(headParams)

	msqls, args := sqlutils.Sqls_CompUpdate(this.getTableName(), upateField, where, rule)
	err := sqlutils.Sqls_Do_PrepareAndExec(this.getDB(), msqls, args)

	if err != nil {
		return err
	}

	return nil
}

func (this *ServiceJobJobController) checkGetPara(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	if !utils.IsFieldCorrectWithRule("service_id", headParams.RouterParams["ServiceId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "ServiceId is not correct!"))
		return false
	}

	return true
}

func (this *ServiceJobJobController) checkPara(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !this.checkGetPara(headParams, r) {
		return false
	}

	if !utils.IsFieldCorrectWithRule("job_id", headParams.RouterParams["JobId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "JobId is not correct!"))
		return false
	}

	return true
}

func (this *ServiceJobJobController) exInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	return nil
}

func (this *ServiceJobJobController) exInit4Del(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	return nil
}

func (this *ServiceJobJobController) compGetCond(headParams *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	serviceid, _ := strconv.ParseInt(headParams.RouterParams["ServiceId"], 10, 32)
	//jobid, _ := strconv.ParseInt(headParams.RouterParams["JobId"], 10, 32)

	idWhere := map[string]interface{}{"service_id": int(serviceid), "del_status": 0, "service_user_id": int(userid)}
	ruleCond := map[string]string{"service_id": " = ", "del_status": " = ", "service_user_id": " = "}

	return idWhere, ruleCond
}

func AddJobOnUserService(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewServiceJobJobController()
	defer ctrl.closeDB()

	obj := NewServiceJobJobControllerObject(ctrl)

	if !ctrl.checkPara(headParams, ren) {
		return
	}

	if ctrl.isUserServiceJobExist(headParams) {
		err := ctrl.updateDelStatus(0, headParams)
		if err != nil {
			ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
		} else {
			ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": ""})
		}
		return
	}

	info := &servemodel.ServeJobInfo{}

	obj.Util_CreateObjectWithId(info, headParams, req, ren)
}

func RemoveJobFromUserService(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewServiceJobJobController()
	defer ctrl.closeDB()
	obj := NewServiceJobJobControllerObject(ctrl)

	info := &servemodel.ServeJobInfo{}

	obj.Util_UpdateDelStatusWithId(info, headParams, req, ren)
}

func GetJobsOnUserService(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewServiceJobJobController()
	defer ctrl.closeDB()
	obj := NewServiceJobJobControllerObject(ctrl)

	info := &servemodel.ServeJobInfo{}

	obj.Util_GetObjectWithId(info, headParams, req, ren, nil, nil)
}
