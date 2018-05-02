package tokens

import (
	//"database/sql"
	//"log"
	"github.com/martini-contrib/render"
	"net/http"
	//"strings"
	//"strconv"
	//"strings"
	//"reflect"
	//"time"

	"web/component/keyutils"
	//"web/component/errcode"
	//"web/component/randutils"

	"web/models/reqparamodel"
	//"web/models/tokenmodel"

	"web/service/routers"
	// "web/service/utils"
	//"github.com/pborman/uuid"
	//jwt "github.com/dgrijalva/jwt-go"
)

func init() {
	userTokenRsaKeyRouterBuilder()
}

func userTokenRsaKeyRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Get("/token/rsa/key", GetRSAPublicKey)
}

type RsaKeyController struct {
}

func GetRSAPublicKey(reqparams *reqparamodel.HttpReqParams, ren render.Render, req *http.Request) {

	keys := keyutils.GetAllRSAKeys()

	out := map[string]string{"ver": string(keys[0].Ver), "pub": string(keys[0].PublicKey)}
	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": []interface{}{out}})
	return
}
