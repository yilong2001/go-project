package userups

import (
	"database/sql"
	//"log"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	//"fmt"
	"log"
	//"net/http"
	//"strconv"
	//"strings"
	"encoding/json"
	//"errors"
	//"time"
	//"reflect"

	"web/component/cfgutils"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/condmodel"
	"web/models/reqparamodel"

	//"web/models/servemodel"
	"web/models/usermodel"
	//"web/service/routers"
	"web/service/utils"
)

type UserWXIdUpController struct {
	tableName string
	db        *sql.DB
	tx        *sql.Tx
	userId    int
	openid    string
	unionid   string
}

func (this *UserWXIdUpController) getTX() *sql.Tx {
	return this.tx
}

func (this *UserWXIdUpController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UserWXIdUpController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UserWXIdUpController) getDB() *sql.DB {
	return this.db
}

func (this *UserWXIdUpController) getTableName() string {
	return this.tableName
}

func (this *UserWXIdUpController) checkParams(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func (this *UserWXIdUpController) exInfoInit(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *UserWXIdUpController) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {
	compser := condmodel.NewCondComposerLinker(" or ")

	where := map[string]interface{}{"user_id": this.userId}
	rule := map[string]string{"user_id": " = "}
	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	compser.SetItem(compSub)
	compser.SetNext(nil)

	ji, _ := json.Marshal(compser)
	log.Println("update service num", string(ji))

	return compser
}

func (this *UserWXIdUpController) calcUpFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {
	fds := map[string]string{}

	fds["weixin_open_id"] = this.openid

	fds["weixin_union_id"] = this.unionid

	return fds
}

func UpdateUserWXId(dbin *sql.DB,
	txin *sql.Tx,
	wxOpenId, wxUnionId string,
	uid int,
	ren render.Render) bool {

	upfields := []string{}
	if wxOpenId != "" {
		upfields = append(upfields, "WeixinOpenId")
	}

	if wxUnionId != "" {
		upfields = append(upfields, "WeixinUnionId")
	}

	if len(upfields) == 0 || uid == 0 {
		return true
	}

	ctrl := &UserWXIdUpController{
		tableName: "web_users",
		db:        dbin,
		tx:        txin,
		userId:    uid,
		openid:    wxOpenId,
		unionid:   wxUnionId,
	}

	if dbin == nil && txin == nil {
		ctrl.initDB()
		defer ctrl.closeDB()
	}

	obj := &utils.UpdateObjectWithIdUtil{
		TableName:       ctrl.getTableName(),
		Db:              ctrl.getDB(),
		Tx:              ctrl.getTX(),
		FormUnParseFlag: 1,

		ExInitFunc:             ctrl.exInfoInit,
		CheckParamFunc:         ctrl.checkParams,
		CondCompserFunc:        ctrl.condCompser,
		CalcedUpdateFieldsFunc: ctrl.calcUpFields,
		MoreProcessFunc:        nil,
	}

	info := usermodel.NewUserInfo()
	info.UserId = uid
	info.WeixinOpenId = wxOpenId
	info.WeixinUnionId = wxUnionId

	return obj.Update_With_MultiInObject(info, info, nil, nil, ren, nil, upfields)
}
