package tokens

import (
	"fmt"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"strconv"
	//"strings"
	//"reflect"
	"database/sql"
	"time"
	"web/component/cfgutils"
	"web/dal/sqldrv"

	"web/component/errcode"
	"web/component/randutils"

	"web/models/reqparamodel"

	"web/service/getter"
	"web/service/routers"
	//"github.com/pborman/uuid"
)

func init() {
	userRefreshTokenRouterBuilder()
}

func userRefreshTokenRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Get("/token/refresh", RefreshToken)
}

func NewRefreshTokenController() *RefreshTokenController {
	ctrl := &RefreshTokenController{}
	ctrl.initDB()
	return ctrl
}

type RefreshTokenController struct {
	db *sql.DB
}

func (this *RefreshTokenController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *RefreshTokenController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *RefreshTokenController) getDB() *sql.DB {
	return this.db
}

func (this *RefreshTokenController) isPwiValid(pwi string, uid string, uuid string) bool {
	destPwi := randutils.BuildRawMd5String(uid, uuid)
	if pwi != destPwi {
		return false
	}

	return true
}

func RefreshToken(reqparams *reqparamodel.HttpReqParams, ren render.Render, req *http.Request) {
	log.Println("RefreshToken req", reqparams)

	ctrl := NewRefreshTokenController()
	defer ctrl.closeDB()

	tokenDb, err := getter.GetModelInfoGetter().GetTokenModelByTId(ctrl.getDB(), reqparams.TokenParams["jti"])
	if err != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetNormalToken_Error, err.Error()))
		return
	}

	log.Println("tokenDb", tokenDb)

	if !ctrl.isPwiValid(reqparams.TokenParams["pwi"], fmt.Sprint(tokenDb.Uid), tokenDb.Uuid) {
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

	if exp > time.Now().Add(time.Hour*2400).Unix() {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_GetPrivateToken_Error, "expire time is too large"))
		return
	}

	tokenDb2, err := GetAndSaveNewToken(ctrl.getDB(), tokenDb.Uid, exp, tokenDb.TokenType, tokenDb.RongCloudToken)
	if err != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_SaveStore_Error, err.Error()))
		return
	}

	out := map[string]string{"token": tokenDb2.Token, "uuid": tokenDb2.Uuid, "rctoken": tokenDb.RongCloudToken}

	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": []interface{}{out}})
}
