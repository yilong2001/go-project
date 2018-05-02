package serveups

import (
	"database/sql"
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

	"web/models/servemodel"
	"web/models/usermodel"
	//"web/service/routers"
	"web/service/utils"
)

type UserServiceFavouriteController struct {
	tableName string
	db        *sql.DB
	tx        *sql.Tx
}

func (this *UserServiceFavouriteController) getTX() *sql.Tx {
	return this.tx
}

func (this *UserServiceFavouriteController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UserServiceFavouriteController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UserServiceFavouriteController) getDB() *sql.DB {
	return this.db
}

func (this *UserServiceFavouriteController) getTableName() string {
	return this.tableName
}

func (this *UserServiceFavouriteController) checkParams(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func (this *UserServiceFavouriteController) exInfoInit(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *UserServiceFavouriteController) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {
	compser := condmodel.NewCondComposerLinker(" or ")
	di, ok := destinfo.(*servemodel.ServeInfo)
	if !ok {
		return nil
	}

	where := map[string]interface{}{"service_id": di.ServiceId, "favourite_num": 0}
	rule := map[string]string{"service_id": " = ", "favourite_num": " > "}
	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	compser.SetItem(compSub)
	compser.SetNext(nil)

	// compser := condmodel.NewCondComposer(" or ")

	// di, ok := destinfo.(*servemodel.ServeInfo)
	// if !ok {
	// 	return compser
	// }

	// where := map[string]interface{}{"service_id": di.ServiceId}
	// rule := map[string]string{"service_id": " = "}
	// compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	// compser.AddItem(compSub)

	ji, _ := json.Marshal(compser)
	log.Println("update service num", string(ji))

	return compser
}

func (this *UserServiceFavouriteController) calcAddFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	fds := map[string]string{"favourite_num": " favourite_num = favourite_num + 1"}

	return fds
}

func (this *UserServiceFavouriteController) calcDecFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	fds := map[string]string{"favourite_num": " favourite_num = favourite_num - 1"}

	return fds
}

func AddNumOfServiceFavourite(dbin *sql.DB,
	txin *sql.Tx,
	ufi *usermodel.UserFavouriteInfo,
	headParams *reqparamodel.HttpReqParams,
	req *http.Request,
	ren render.Render) bool {
	ctrl := &UserServiceFavouriteController{
		tableName: "web_services",
		db:        dbin,
		tx:        txin,
	}

	obj := &utils.UpdateObjectWithIdUtil{
		TableName:       ctrl.getTableName(),
		Db:              ctrl.getDB(),
		Tx:              ctrl.getTX(),
		FormUnParseFlag: 1,

		ExInitFunc:             ctrl.exInfoInit,
		CheckParamFunc:         ctrl.checkParams,
		CondCompserFunc:        ctrl.condCompser,
		CalcedUpdateFieldsFunc: ctrl.calcAddFields,
		MoreProcessFunc:        nil,
	}

	info := servemodel.NewServeInfo()
	info.ServiceId = ufi.DestId

	return obj.Update_With_MultiInObject(ufi, info, headParams, req, ren, nil, []string{"FavouriteNum"})
}

func DecNumOfServiceFavourite(dbin *sql.DB,
	txin *sql.Tx,
	ufi *usermodel.UserFavouriteInfo,
	headParams *reqparamodel.HttpReqParams,
	req *http.Request,
	ren render.Render) bool {
	ctrl := &UserServiceFavouriteController{
		tableName: "web_services",
		db:        dbin,
		tx:        txin,
	}

	obj := &utils.UpdateObjectWithIdUtil{
		TableName:       ctrl.getTableName(),
		Db:              ctrl.getDB(),
		Tx:              ctrl.getTX(),
		FormUnParseFlag: 1,

		ExInitFunc:             ctrl.exInfoInit,
		CheckParamFunc:         ctrl.checkParams,
		CondCompserFunc:        ctrl.condCompser,
		CalcedUpdateFieldsFunc: ctrl.calcDecFields,
		MoreProcessFunc:        nil,
	}

	info := servemodel.NewServeInfo()
	info.ServiceId = ufi.DestId

	return obj.Update_With_MultiInObject(ufi, info, headParams, req, ren, nil, []string{"FavouriteNum"})
}
