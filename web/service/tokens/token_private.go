package tokens

import (
	"database/sql"
	"fmt"
	"log"
	//"strings"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	//"strings"
	//"reflect"
	"time"
	"web/component/cfgutils"

	"web/component/errcode"
	"web/component/randutils"
	"web/component/rongcloud"

	//"web/component/idutils"
	// "web/component/objutils"
	// "web/component/sqlutils"
	"web/dal/sqldrv"
	//"web/models/usermodel"
	//"crypto/md5"
	//"crypto/rand"
	//"encoding/base64"
	//"encoding/hex"

	"web/models/reqparamodel"
	"web/models/tokenmodel"

	"web/service/getter"
	"web/service/routers"
	// "web/service/utils"
	//"github.com/pborman/uuid"
	//jwt "github.com/dgrijalva/jwt-go"
)

func init() {
	userPrivateTokenRouterBuilder()
}

func userPrivateTokenRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Get("/token/private", GetPrivateToken)
}

func NewPrivateTokenController() *PrivateTokenController {
	ctrl := &PrivateTokenController{}
	ctrl.initDB()
	return ctrl
}

type PrivateTokenController struct {
	db *sql.DB
}

func (this *PrivateTokenController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *PrivateTokenController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *PrivateTokenController) getDB() *sql.DB {
	return this.db
}

func (this *PrivateTokenController) isPwiValid(pwi string, pwmd5 string, phone string) bool {
	destPwi := randutils.BuildRawMd5String(pwmd5, phone)
	// destPwi := ""
	// if len(pwmd5) != 32 {
	// 	destPwi = randutils.BuildRawMd5String(pwmd5, phone)
	// } else {
	// 	destPwi = randutils.BuildMd5PWPhoneStringV2(pwmd5, phone)
	// }

	if pwi != destPwi {
		return false
	}

	return true
}

func GetPrivateToken(reqparams *reqparamodel.HttpReqParams, ren render.Render, req *http.Request) {
	log.Println("GetToken req", reqparams)

	ctrl := NewPrivateTokenController()
	defer ctrl.closeDB()

	userinfo, err := getter.GetModelInfoGetter().GetUserByLoginName(ctrl.getDB(), reqparams.TokenParams["uid"])
	if err != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetPrivateToken_Error, err.Error()))
		return
	}

	log.Println("GetToken userinfo", userinfo)

	if !ctrl.isPwiValid(reqparams.TokenParams["pwi"], userinfo.LoginPw, userinfo.LoginName) {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetToken_PW_Error, "pwi is not correct"))
		return
	}

	exp, err1 := strconv.ParseInt(reqparams.TokenParams["exp"], 10, 64)
	if err1 != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetPrivateToken_Error, err1.Error()))
		return
	}

	if exp < time.Now().Add(time.Second*300).Unix() {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetPrivateToken_Error, "expire time is too short"))
		return
	}

	//get rong cloud token
	rctoken := ""
	rcResult := rongcloud.UserGetToken(fmt.Sprint(userinfo.UserId), fmt.Sprint(userinfo.UserId))
	if rcResult == nil {
		log.Println("get rong cloud failed")
	} else {
		rctoken = rcResult.Token
	}

	tokenDb, err := GetAndSaveNewToken(ctrl.getDB(), userinfo.UserId, exp, tokenmodel.Const_Token_Type_Private, rctoken)
	if err != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_SaveStore_Error, err.Error()))
		return
	}

	out := map[string]string{"token": tokenDb.Token, "uuid": tokenDb.Uuid,
		"uid": fmt.Sprint(userinfo.UserId), "rctoken": rctoken}

	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": []interface{}{out}})
}
