package users

import (
	"database/sql"
	//"errors"
	"github.com/go-martini/martini" //
	"log"
	//"io"
	"fmt"
	"net/http"
	"strconv"
	"time"
	//"os"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"web/component/cfgutils"
	"web/component/errcode"
	//"web/component/fileutils"
	"web/component/objutils"
	"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/reqparamodel"
	"web/models/tokenmodel"
	"web/models/usermodel"
	"web/service/getter"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userPwRouterBuilder()
}

func userPwRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/user/mine/password", UpdateUserPassword)
}

var globalUserPasswordController *UsePasswordController = &UsePasswordController{
	tableName: "web_users",
}

func getGlobalUserPwController() *UsePasswordController {
	return globalUserPasswordController
}

type UsePasswordController struct {
	tableName string
}

func (this *UsePasswordController) getTableName() string {
	return this.tableName
}

func (this *UsePasswordController) isOldPasswordValid(db *sql.DB, pwinfo *usermodel.UserPasswordInfo,
	userid int, r render.Render) bool {

	whereCond := map[string]interface{}{"user_id": userid, "phone": pwinfo.Phone}
	ruleCondition := map[string]string{"user_id": "=", "phone": " = "}

	ct := 0
	msqls, selargs, whereargs := sqlutils.Sqls_CompSelectCount(this.getTableName(), "user_id", &ct, whereCond, ruleCondition)

	log.Println(msqls)

	err := sqlutils.Sqls_Do_QueryRowAndScan(db, msqls, selargs, whereargs)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Scan_Error, err.Error()))
		log.Println(err)
		return false
	}

	if ct == 0 {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_Old_Pw_Error, "old passowrd is not correct!"))
		log.Println("old passowrd is not correct!")
		return false
	}

	return true
}

func (this *UsePasswordController) checkParams(db *sql.DB, headParams *reqparamodel.HttpReqParams, req *http.Request, r render.Render) (bool, *usermodel.UserPasswordInfo) {
	if headParams.TokenParams["TokenType"] != fmt.Sprint(tokenmodel.Const_Token_Type_Private) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "User Token is not correct!"))
		return false, nil
	}

	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false, nil
	}

	pwinfo := &usermodel.UserPasswordInfo{}
	err := objutils.ParseObjectWithForm(pwinfo, req)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, "password is wrong!"))
		return false, nil
	}

	if !utils.IsFieldCorrectWithRule("password", pwinfo.Password) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_Pw_Error, "new password is not correct!"))
		return false, nil
	}

	//
	smscodeinfo, err := getter.GetModelInfoGetter().GetSmsCodeByPhone(db, pwinfo.Phone)
	if err != nil {
		log.Println("sms code validate is wrong, no phone code, ", err, pwinfo.Phone)
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_NotCorrect_Error, err.Error()))
		return false, nil
	}

	if smscodeinfo.Code != pwinfo.Code {
		log.Println("sms code is wrong, ", smscodeinfo.Code, pwinfo.Code)
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_NotCorrect_Error, err.Error()))
		return false, nil
	}

	return true, pwinfo
}

func (this *UsePasswordController) changePassword(db *sql.DB, pwinfo *usermodel.UserPasswordInfo, userid int, r render.Render) {

	dt := time.Now().Format("2006-01-02 15:04:05")

	upCondition := map[string]interface{}{"login_pw": pwinfo.Password,
		"update_time": dt}

	whereCondition := map[string]interface{}{"user_id": userid}
	ruleCondition := map[string]string{"user_id": "="}

	msqls, args := sqlutils.Sqls_CompUpdate(this.getTableName(), upCondition, whereCondition, ruleCondition)

	log.Println(msqls)
	log.Println(args)

	err := sqlutils.Sqls_Do_PrepareAndExec(db, msqls, args)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
		return
	}

	r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": ""})
}

func UpdateUserPassword(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)
	log.Println(headParams.RouterParams["UserId"], " change password!")

	ctrl := &UsePasswordController{
		tableName: "web_users",
	}

	var db *sql.DB = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	defer db.Close()

	res, pwinfo := ctrl.checkParams(db, headParams, req, r)
	if !res {
		return
	}

	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)

	if !ctrl.isOldPasswordValid(db, pwinfo, int(userid), r) {
		return
	}

	ctrl.changePassword(db, pwinfo, int(userid), r)
}
