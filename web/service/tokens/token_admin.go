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

	"web/dal/sqldrv"

	//"web/models/platform/adminmodel"
	"web/models/reqparamodel"
	"web/models/tokenmodel"

	"web/service/getter"
	"web/service/routers"
)

func init() {
	userAdminTokenRouterBuilder()
}

func userAdminTokenRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Get("/token/admin", GetAdminToken)
}

func NewAdminTokenController() *AdminTokenController {
	ctrl := &AdminTokenController{}
	ctrl.initDB()
	return ctrl
}

type AdminTokenController struct {
	db *sql.DB
}

func (this *AdminTokenController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *AdminTokenController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *AdminTokenController) getDB() *sql.DB {
	return this.db
}

func (this *AdminTokenController) isPwiValid(pwi string, pwmd5 string, phone string) bool {
	destPwi := randutils.BuildRawMd5String(pwmd5, phone)
	if pwi != destPwi {
		return false
	}

	return true
}

func GetAdminToken(reqparams *reqparamodel.HttpReqParams, ren render.Render, req *http.Request) {
	log.Println("GetToken req", reqparams)

	ctrl := NewAdminTokenController()
	defer ctrl.closeDB()

	userinfo, err := getter.GetModelInfoGetter().GetAdminUserByLoginName(ctrl.getDB(), reqparams.TokenParams["uid"])
	if err != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetAdminToken_Error, err.Error()))
		return
	}

	log.Println("GetToken userinfo", userinfo)

	if !ctrl.isPwiValid(reqparams.TokenParams["pwi"], userinfo.LoginPw, userinfo.LoginName) {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetToken_PW_Error, "pwi is not correct"))
		return
	}

	exp, err1 := strconv.ParseInt(reqparams.TokenParams["exp"], 10, 64)
	if err1 != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetAdminToken_Error, err1.Error()))
		return
	}

	if exp < time.Now().Add(time.Second*300).Unix() {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetAdminToken_Error, "expire time is too short"))
		return
	}

	tokenDb, err := GetAndSaveNewToken(ctrl.getDB(), userinfo.UserId, exp, tokenmodel.Const_Token_Type_Admin, "")
	if err != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_SaveStore_Error, err.Error()))
		return
	}

	out := map[string]string{"token": tokenDb.Token, "uuid": tokenDb.Uuid,
		"uid": fmt.Sprint(userinfo.UserId)}

	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": []interface{}{out}})
}
