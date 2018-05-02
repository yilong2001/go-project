package tokens

import (
	"database/sql"
	"log"
	"strings"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	//"strconv"
	//"strings"
	"reflect"
	"time"
	"web/component/cfgutils"

	"web/component/errcode"
	"web/component/keyutils"
	//"web/component/randutils"

	// "web/component/idutils"
	// "web/component/objutils"
	// "web/component/sqlutils"
	"web/dal/sqldrv"
	// "web/models/usermodel"
	//"crypto/md5"
	//"crypto/rand"
	//"encoding/base64"
	//"encoding/hex"

	"web/models/reqparamodel"
	"web/models/tokenmodel"

	"web/service/routers"
	// "web/service/utils"

	//"github.com/pborman/uuid"

	jwt "github.com/dgrijalva/jwt-go"
)

func init() {
	userOpenTokenRouterBuilder()
}

func userOpenTokenRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Get("/token/open", GetOpenToken)
}

func NewOpenTokenController() *OpenTokenController {
	ctrl := &OpenTokenController{}
	ctrl.initDB()
	return ctrl
}

type OpenTokenController struct {
	db *sql.DB
}

func (this *OpenTokenController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *OpenTokenController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *OpenTokenController) getDB() *sql.DB {
	return this.db
}

func GetOpenToken(reqparams *reqparamodel.HttpReqParams, ren render.Render, req *http.Request) {
	log.Println(reqparams)

	ctrl := NewOpenTokenController()
	defer ctrl.closeDB()

	tokenDb, err := GetAndSaveNewToken(ctrl.getDB(), -1, time.Now().Add(time.Hour*240).Unix(), tokenmodel.Const_Token_Type_Public, "")
	if err != nil {
		ren.JSON(200, errcode.NewErrRsp2(errcode.Err_Token_SaveStore_Error, err.Error()))
		return
	}

	out := map[string]string{"token": tokenDb.Token, "uuid": tokenDb.Uuid,
		"uid": "-1"}

	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": []interface{}{out}})
}

func createOpenToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenClaimsModel := tokenmodel.NewDefaultOpenTokenClaims()

	object := reflect.ValueOf(tokenClaimsModel)
	myref := object.Elem()
	typeOfType := myref.Type()

	for i := 0; i < myref.NumField(); i++ {
		field := myref.Field(i)

		token.Claims[strings.ToLower(typeOfType.Field(i).Name)] = field.Interface()
	}

	// Sign and get the complete encoded token as a string
	//log.Println("key : ", GetHS256Key())
	tokenString, err := token.SignedString(keyutils.GetHS256Key())

	return tokenString, err
}
