package servicers

import (
	"database/sql"
	"github.com/go-martini/martini"
	"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	//"errors"
	"net/http"
	"net/url"
	"strconv"
	//"strings"
	//"time"
	//"reflect"

	"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/firmmodel"
	"web/models/reqparamodel"
	"web/models/usermodel"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	openServicersRouterBuilder()
}

func openServicersRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Get("/servicer", GetOpenServicers)
	m.Get("/servicer/tuijian", GetOpenTuijianServicers)
}

func NewOpenServicerControllerObject(ctrl *OpenServicerController) *utils.ObjectWithIdUtil {
	return &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInit4Get,
		CheckParamFuncForGet: ctrl.check4Open,
		WhereCondFuncForGet:  ctrl.compWhereCond,
	}
}

func NewOpenServicerController() *OpenServicerController {
	ctrl := &OpenServicerController{
		tableName: "web_users",
	}

	ctrl.initDB()
	return ctrl
}

type OpenServicerController struct {
	tableName string
	db        *sql.DB
}

func (this *OpenServicerController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *OpenServicerController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *OpenServicerController) getDB() *sql.DB {
	return this.db
}

func (this *OpenServicerController) getTableName() string {
	return this.tableName
}

func (this *OpenServicerController) exInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *OpenServicerController) check4Open(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	//firm := headParams.URLParams.Get("Firm")
	// if firm == "" || len(firm) < 3 {
	// 	log.Println("no firm info with req")
	// 	r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, err.Error()))
	// 	return false
	// }

	return true
}

func (this *OpenServicerController) compWhereCond(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	log.Println("compWhereCondition HttpReqParams", params)

	firm, err := url.QueryUnescape(params.URLParams.Get("Firm"))
	if err == nil && len(firm) > 1 {
		simplefirm := firmmodel.MakeCompareName(firm)
		idWhere["simple_firm"] = simplefirm
		idRlue["simple_firm"] = " = "
	}

	Name := params.URLParams.Get("Name")
	if len(Name) > 1 {
		idWhere["user_name"] = Name
		idRlue["user_name"] = " = "
	}

	Phone := params.URLParams.Get("Phone")
	if len(Phone) > 1 {
		idWhere["phone"] = Phone
		idRlue["phone"] = " = "
	}

	City := params.URLParams.Get("City")
	if city, err := strconv.ParseInt(City, 10, 32); err == nil && len(City) > 0 && city > 0 {
		idWhere["city"] = city
		idRlue["city"] = " = "
	}

	Industry := params.URLParams.Get("Industry")
	if industry, err := strconv.ParseInt(Industry, 10, 32); err == nil && len(Industry) > 0 && industry > 0 {
		idWhere["industry"] = industry
		idRlue["industry"] = " = "
	}

	IsTuijian := params.URLParams.Get("IsTuijian")
	if _, err := strconv.ParseInt(IsTuijian, 10, 32); err == nil && IsTuijian != "" {
		idWhere["is_tuijian"] = 1
		idRlue["is_tuijian"] = " = "
	}

	if len(idWhere) == 0 {
		idWhere["is_tuijian"] = 1
		idRlue["is_tuijian"] = " = "
	}

	log.Println("OpenServicerController : ", idWhere)

	return idWhere, idRlue
}

func GetOpenServicers(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	userInfo := &usermodel.UserInfo{}

	ctrl := NewOpenServicerController()
	defer ctrl.closeDB()
	obj := NewOpenServicerControllerObject(ctrl)

	obj.Util_GetObjectWithId(userInfo, headParams, req, r, userInfo.GetSkipFieldsForOpenQuery(), nil)
}

func GetOpenTuijianServicers(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	if headParams.URLParams.Get("IsTuijian") == "" {
		headParams.URLParams.Set("IsTuijian", "1")
	}

	userInfo := &usermodel.UserInfo{}

	ctrl := NewOpenServicerController()
	defer ctrl.closeDB()
	obj := NewOpenServicerControllerObject(ctrl)

	obj.Util_GetObjectWithId(userInfo, headParams, req, r, userInfo.GetSkipFieldsForOpenQuery(), nil)
}
