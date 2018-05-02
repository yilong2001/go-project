package userups

import (
	"database/sql"
	"fmt"
	//"log"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"log"
	"net/http"
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

type UserUpdatedFieldNumController struct {
	tableName   string
	db          *sql.DB
	tx          *sql.Tx
	deta        int
	upServedNum bool
}

func (this *UserUpdatedFieldNumController) getTX() *sql.Tx {
	return this.tx
}

func (this *UserUpdatedFieldNumController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UserUpdatedFieldNumController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UserUpdatedFieldNumController) getDB() *sql.DB {
	return this.db
}

func (this *UserUpdatedFieldNumController) getTableName() string {
	return this.tableName
}

func (this *UserUpdatedFieldNumController) checkParams(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func (this *UserUpdatedFieldNumController) exInfoInit(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *UserUpdatedFieldNumController) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {
	compser := condmodel.NewCondComposerLinker(" or ")

	ui, ok := destinfo.(*usermodel.UserInfo)
	if !ok {
		return compser
	}

	where := map[string]interface{}{"user_id": ui.UserId}
	rule := map[string]string{"user_id": " = "}

	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	compser.SetItem(compSub)
	compser.SetNext(nil)

	//compser.AddItem(compSub)

	ji, _ := json.Marshal(compser)
	log.Println("update service num", string(ji))

	return compser
}

func (this *UserUpdatedFieldNumController) condCompserDecAccount(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {
	compser := condmodel.NewCondComposerLinker(" or ")

	ui, ok := destinfo.(*usermodel.UserInfo)
	if !ok {
		return compser
	}

	where := map[string]interface{}{"user_id": ui.UserId, "account_balance": this.deta}
	rule := map[string]string{"user_id": " = ", "account_balance": " >= "}
	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	compser.SetItem(compSub)
	compser.SetNext(nil)

	//compser.AddItem(compSub)

	ji, _ := json.Marshal(compser)
	log.Println("update account_balance ", string(ji))

	return compser
}

func (this *UserUpdatedFieldNumController) calcServedAddFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	fds := map[string]string{"served_num": " served_num = served_num + 1"}

	return fds
}

func (this *UserUpdatedFieldNumController) calcServedDecFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	fds := map[string]string{"served_num": " served_num = served_num - 1"}

	return fds
}

func (this *UserUpdatedFieldNumController) calcBalanceAddFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	fds := map[string]string{"account_balance": " account_balance = account_balance + " + fmt.Sprint(this.deta)}

	if this.upServedNum {
		fds["served_num"] = " served_num = served_num + 1 "
	}

	return fds
}

func (this *UserUpdatedFieldNumController) calcBalanceDecFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	fds := map[string]string{"account_balance": " account_balance = account_balance - " + fmt.Sprint(this.deta)}

	return fds
}

func AddServedNumOfUserSelf(dbin *sql.DB,
	txin *sql.Tx,
	userid int,
	headParams *reqparamodel.HttpReqParams,
	req *http.Request,
	ren render.Render) bool {
	ctrl := &UserUpdatedFieldNumController{
		tableName: "web_users",
		db:        dbin,
		tx:        txin,

		upServedNum: false,
	}

	obj := &utils.UpdateObjectWithIdUtil{
		TableName:       ctrl.getTableName(),
		Db:              ctrl.getDB(),
		Tx:              ctrl.getTX(),
		FormUnParseFlag: 1,

		ExInitFunc:             ctrl.exInfoInit,
		CheckParamFunc:         ctrl.checkParams,
		CondCompserFunc:        ctrl.condCompser,
		CalcedUpdateFieldsFunc: ctrl.calcServedAddFields,
		MoreProcessFunc:        nil,
	}

	info := &usermodel.UserInfo{}
	info.UserId = userid

	return obj.Update_With_MultiInObject(nil, info, headParams, req, ren, nil, []string{"ServedNum"})
}

func AddAccountBanlancOfUser(dbin *sql.DB,
	txin *sql.Tx,
	userid int,
	addnum int,
	isUpServedNum bool,
	headParams *reqparamodel.HttpReqParams,
	req *http.Request,
	ren render.Render) bool {
	ctrl := &UserUpdatedFieldNumController{
		tableName: "web_users",
		db:        dbin,
		tx:        txin,
		deta:      addnum,

		upServedNum: isUpServedNum,
	}

	obj := &utils.UpdateObjectWithIdUtil{
		TableName:       ctrl.getTableName(),
		Db:              ctrl.getDB(),
		Tx:              ctrl.getTX(),
		FormUnParseFlag: 1,

		ExInitFunc:             ctrl.exInfoInit,
		CheckParamFunc:         ctrl.checkParams,
		CondCompserFunc:        ctrl.condCompser,
		CalcedUpdateFieldsFunc: ctrl.calcBalanceAddFields,
		MoreProcessFunc:        nil,
	}

	info := &usermodel.UserInfo{}
	info.UserId = userid

	specFields := []string{"AccountBalance"}
	if isUpServedNum {
		specFields = append(specFields, "ServedNum")
	}

	return obj.Update_With_MultiInObject(nil, info, headParams, req, ren, nil, specFields)
}

func DecAccountBanlancOfUser(dbin *sql.DB,
	txin *sql.Tx,
	userid int,
	decnum int,
	headParams *reqparamodel.HttpReqParams,
	req *http.Request,
	ren render.Render) bool {
	ctrl := &UserUpdatedFieldNumController{
		tableName: "web_users",
		db:        dbin,
		tx:        txin,
		deta:      decnum,

		upServedNum: false,
	}

	obj := &utils.UpdateObjectWithIdUtil{
		TableName:       ctrl.getTableName(),
		Db:              ctrl.getDB(),
		Tx:              ctrl.getTX(),
		FormUnParseFlag: 1,

		ExInitFunc:             ctrl.exInfoInit,
		CheckParamFunc:         ctrl.checkParams,
		CondCompserFunc:        ctrl.condCompser,
		CalcedUpdateFieldsFunc: ctrl.calcBalanceDecFields,
		MoreProcessFunc:        nil,
	}

	info := &usermodel.UserInfo{}
	info.UserId = userid

	return obj.Update_With_MultiInObject(nil, info, headParams, req, ren, nil, []string{"AccountBalance"})
}
