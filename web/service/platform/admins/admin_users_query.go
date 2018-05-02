package admins

import (
	"database/sql"
	"github.com/go-martini/martini"
	//"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	//"errors"
	"net/http"
	"strconv"
	//"strings"
	//"time"
	//"reflect"
	"net/url"

	"web/component/cfgutils"
	"web/component/errcode"
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
	openUserWithinAdminRouterBuilder()
}

func openUserWithinAdminRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Get("/admin/reviewer/users", GetUserInfo4AdminReview)
}

func NewUserWithinAdminControllerObject(ctrl *UserWithinAdminController) *utils.ObjectWithIdUtil {
	return &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInit4AdminGet,
		CheckParamFuncForGet: ctrl.check4AdminGet,
		WhereCondFuncForGet:  ctrl.compWhereCond,
	}
}

func NewUserWithinAdminController() *UserWithinAdminController {
	ctrl := &UserWithinAdminController{
		tableName: "web_users",
	}

	ctrl.initDB()
	return ctrl
}

type UserWithinAdminController struct {
	tableName string
	db        *sql.DB
}

func (this *UserWithinAdminController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UserWithinAdminController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UserWithinAdminController) getDB() *sql.DB {
	return this.db
}

func (this *UserWithinAdminController) getTableName() string {
	return this.tableName
}

func (this *UserWithinAdminController) exInit4AdminGet(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *UserWithinAdminController) check4AdminGet(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !isAdminToken(headParams) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "not admin user"))
		return false
	}

	return true
}

func (this *UserWithinAdminController) compWhereCond(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

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
	if tuijian, err := strconv.ParseInt(IsTuijian, 10, 32); err == nil && IsTuijian != "" {
		idWhere["is_tuijian"] = tuijian
		idRlue["is_tuijian"] = " = "
	}

	return idWhere, idRlue
}

func GetUserInfo4AdminReview(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	userInfo := &usermodel.UserInfo{}

	ctrl := NewUserWithinAdminController()
	defer ctrl.closeDB()
	obj := NewUserWithinAdminControllerObject(ctrl)

	obj.Util_GetObjectWithId(userInfo, headParams, req, r, userInfo.GetSkipFieldsForOpenQuery(), nil)
}
