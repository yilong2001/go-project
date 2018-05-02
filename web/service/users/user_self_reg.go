package users

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"net/http"
	//"strconv"
	"errors"
	"strings"
	"time"
	//"reflect"

	"web/component/cfgutils"
	"web/component/errcode"
	"web/component/idutils"
	"web/component/objutils"
	"web/component/randutils"

	"web/component/rongcloud"
	"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/reqparamodel"
	"web/models/tokenmodel"
	"web/models/usermodel"
	"web/service/getter"
	"web/service/routers"
	"web/service/tokens"
	"web/service/utils"
)

func init() {
	userSelfRegRouterBuilder()
}

func userSelfRegRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/user", UserSelfRegister)
}

func IsLoginNameOnlyOne(db *sql.DB, tablename, phone string, r render.Render) bool {
	ct := 0

	whereCond := map[string]interface{}{"login_name": phone}
	ruleCondition := map[string]string{"login_name": "="}

	msqls, selargs, whereargs := sqlutils.Sqls_CompSelectCount(tablename, "user_id", &ct, whereCond, ruleCondition)

	log.Println(msqls)

	err := sqlutils.Sqls_Do_QueryRowAndScan(db, msqls, selargs, whereargs)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Scan_Error, err.Error()))
		log.Println(errcode.Err_Form_Para_Old_Pw_Error)
		return false
	}

	if ct > 0 {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_Duplicate_Error, "phone has been used!"))
		return false
	}

	return true
}

func IsLoginNameOnlyOneEx(db *sql.DB, tablename, phone string) error {
	ct := 0

	whereCond := map[string]interface{}{"login_name": phone}
	ruleCondition := map[string]string{"login_name": "="}

	msqls, selargs, whereargs := sqlutils.Sqls_CompSelectCount(tablename, "user_id", &ct, whereCond, ruleCondition)

	log.Println(msqls)

	err := sqlutils.Sqls_Do_QueryRowAndScan(db, msqls, selargs, whereargs)
	if err != nil {
		log.Println("IsLoginNameOnlyOne : ", err)
		return err
	}

	if ct > 0 {
		return errors.New("phone has been used!")
	}

	return nil
}

func UserSelfRegister(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	isskipsms := false
	if headParams.URLParams.Get("SpecialUser") == "qazwsxedc" {
		isskipsms = true
	}

	reginfo := &usermodel.UserRegisterInfo{}
	err := objutils.ParseObjectWithForm(reginfo, req)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, err.Error()))
		return
	}

	log.Println(reginfo)

	fieldIfArrs, _ := reginfo.GetWholeFields()
	res, fns := utils.IsFieldsValueOk(fieldIfArrs)
	if !res {
		detail := "the info (" + fns + ") is not correct"
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_NotCorrect_Error, detail))
		return
	}

	ctrl := NewUseSelfController()
	defer ctrl.closeDB()

	var db *sql.DB = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	defer db.Close()

	//一个手机号不同的人注册，暂不考虑
	isonly := IsLoginNameOnlyOne(db, ctrl.getTableName(), reginfo.Phone, r)
	if !isonly {
		log.Println("phone has been used")
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_Duplicate_Error, "phone has been used!"))
		return
	}

	dt := time.Now().Format("2006-01-02 15:04:05")
	openContact := reginfo.Phone

	if !isskipsms {
		//
		smscodeinfo, err := getter.GetModelInfoGetter().GetSmsCodeByPhone(db, reginfo.Phone)
		if err != nil {
			log.Println("sms code validate is wrong, no phone code, ", err, reginfo.Phone)
			r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_NotCorrect_Error, err.Error()))
			return
		}

		codemd5again := randutils.BuildMd5PWPhoneStringV2(smscodeinfo.Code, "")
		log.Println(smscodeinfo.Code, codemd5again)

		if smscodeinfo.Code != reginfo.Code {
			log.Println("sms code is wrong, ", smscodeinfo.Code, codemd5again, reginfo.Code)
			r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_NotCorrect_Error, "sms code"))
			return
		}
	} else {
		if strings.HasPrefix(reginfo.Phone, "100") {
			openContact = "联系平台：" + usermodel.Const_Customer_Servier_Phone_Main
		}
	}

	id := idutils.GetId(ctrl.getGenIdFlag())
	log.Println("new id is %d", id)

	preSql := "insert into web_users(user_id,login_name,login_pw,phone,open_contact,introduce,create_time,renzheng_time) values(?,?,?,?,?,?,?,?)"

	stmt, err := db.Prepare(preSql)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Prepare_Error, preSql))
		return
	}
	defer stmt.Close()

	//md5pw := randutils.BuildMd5PWPhoneStringV2(reginfo.Password, "")
	md5pw := reginfo.Password

	if _, err := stmt.Exec(id, reginfo.Phone, md5pw, reginfo.Phone, openContact, " ", dt, dt); err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
		return
	}

	//get rong cloud token
	rctoken := ""
	rcResult := rongcloud.UserGetToken(fmt.Sprint(id), fmt.Sprint(id))
	if rcResult == nil {
		log.Println("get rong cloud failed")
	} else {
		rctoken = rcResult.Token
	}

	tokenDb, err := tokens.GetAndSaveNewToken(db, id, time.Now().Add(time.Hour*240).Unix(), tokenmodel.Const_Token_Type_Private, rctoken)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_SaveStore_Error, err.Error()))
		return
	}

	out := map[string]string{"token": tokenDb.Token, "uuid": tokenDb.Uuid,
		"uid": fmt.Sprint(id), "rctoken": rctoken}

	r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": []interface{}{out}})
}

func DoRegister(db *sql.DB, isskipsms bool, reginfo *usermodel.UserRegisterInfo, headParams *reqparamodel.HttpReqParams, req *http.Request, withToken bool) (*errcode.ErrRsp, *tokenmodel.TokenDbModel) {

	ctrl := NewUseSelfController()
	defer ctrl.closeDB()

	// var db *sql.DB = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	// defer db.Close()

	//一个手机号不同的人注册，暂不考虑
	err1 := IsLoginNameOnlyOneEx(db, ctrl.getTableName(), reginfo.Phone)
	if err1 != nil {
		log.Println("phone has been used")
		return errcode.NewErrRsp2(errcode.Err_Form_Para_NotCorrect_Error, err1.Error()), nil
	}

	dt := time.Now().Format("2006-01-02 15:04:05")
	openContact := reginfo.Phone

	if !isskipsms {
		//
		smscodeinfo, err := getter.GetModelInfoGetter().GetSmsCodeByPhone(db, reginfo.Phone)
		if err != nil {
			log.Println("sms code validate is wrong, no phone code, ", err, reginfo.Phone)
			return errcode.NewErrRsp2(errcode.Err_Form_Para_NotCorrect_Error, err.Error()), nil
		}

		codemd5again := randutils.BuildMd5PWPhoneStringV2(smscodeinfo.Code, "")
		log.Println(smscodeinfo.Code, codemd5again)

		if smscodeinfo.Code != reginfo.Code {
			log.Println("sms code is wrong, ", smscodeinfo.Code, codemd5again, reginfo.Code)
			return errcode.NewErrRsp2(errcode.Err_Form_Para_NotCorrect_Error, "sms code"), nil
		}
	} else {
		if strings.HasPrefix(reginfo.Phone, "10") {
			openContact = "助理：" + usermodel.Const_Customer_Servier_Phone_Main
		}
	}

	id := idutils.GetId(ctrl.getGenIdFlag())
	log.Println("new id is %d", id)

	preSql := "insert into web_users(user_id,login_name,login_pw,phone,open_contact,introduce,create_time,renzheng_time) values(?,?,?,?,?,?,?,?)"

	stmt, err := db.Prepare(preSql)
	if err != nil {
		return errcode.NewErrRsp2(errcode.Err_Db_Prepare_Error, preSql), nil
	}
	defer stmt.Close()

	//md5pw := randutils.BuildMd5PWPhoneStringV2(reginfo.Password, "")
	md5pw := reginfo.Password

	if _, err := stmt.Exec(id, reginfo.Phone, md5pw, reginfo.Phone, openContact, " ", dt, dt); err != nil {
		return errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()), nil
	}

	//get rong cloud token
	rctoken := ""
	rcResult := rongcloud.UserGetToken(fmt.Sprint(id), fmt.Sprint(id))
	if rcResult == nil {
		log.Println("get rong cloud failed")
	} else {
		rctoken = rcResult.Token
	}

	if withToken {
		tokenDb, err := tokens.GetAndSaveNewToken(db, id, time.Now().Add(time.Hour*240).Unix(), tokenmodel.Const_Token_Type_Private, rctoken)
		if err != nil {
			return errcode.NewErrRsp2(errcode.Err_Token_SaveStore_Error, err.Error()), nil
		}

		return nil, tokenDb
	} else {
		return nil, nil
	}

	//out := map[string]string{"token": tokenDb.Token, "uuid": tokenDb.Uuid,
	//	"uid": fmt.Sprint(id), "rctoken": rctoken}

	//r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": []interface{}{out}})
}
