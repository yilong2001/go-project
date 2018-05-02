package serveups

import (
	"database/sql"
	//"log"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"fmt"
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

	"web/models/servemodel"
	//"web/models/usermodel"
	//"web/service/routers"
	"web/service/utils"
)

type UserServiceCommentController struct {
	tableName    string
	db           *sql.DB
	tx           *sql.Tx
	commentLevel int
}

func (this *UserServiceCommentController) getTX() *sql.Tx {
	return this.tx
}

func (this *UserServiceCommentController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UserServiceCommentController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UserServiceCommentController) getDB() *sql.DB {
	return this.db
}

func (this *UserServiceCommentController) getTableName() string {
	return this.tableName
}

func (this *UserServiceCommentController) checkParams(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func (this *UserServiceCommentController) exInfoInit(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *UserServiceCommentController) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {
	compser := condmodel.NewCondComposerLinker(" or ")
	di, ok := destinfo.(*servemodel.ServeInfo)
	if !ok {
		return nil
	}

	where := map[string]interface{}{"service_id": di.ServiceId}
	rule := map[string]string{"service_id": " = "}
	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	compser.SetItem(compSub)
	compser.SetNext(nil)

	ji, _ := json.Marshal(compser)
	log.Println("update service num", string(ji))

	return compser
}

func (this *UserServiceCommentController) calcUpFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {
	fds := map[string]string{}

	fds["avg_star"] = " avg_star = (comment_num * avg_star + " + fmt.Sprint(this.commentLevel) + ")/(comment_num+1)"

	fds["comment_num"] = " comment_num = comment_num + 1 "

	return fds
}

func (this *UserServiceCommentController) calcServedNumFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	fds := map[string]string{"served_num": " served_num = served_num + 1 "}

	return fds
}

func UpdateUserServiceCommentNum(dbin *sql.DB,
	txin *sql.Tx,
	level int,
	svrid int,
	ren render.Render) bool {
	ctrl := &UserServiceCommentController{
		tableName:    "web_services",
		db:           dbin,
		tx:           txin,
		commentLevel: level,
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

	info := servemodel.NewServeInfo()
	info.ServiceId = svrid

	return obj.Update_With_MultiInObject(info, info, nil, nil, ren, nil, []string{"CommentNum", "AvgStar"})
}

func AddNumOfServicedNum(dbin *sql.DB,
	txin *sql.Tx,
	svid int,
	ren render.Render) bool {
	ctrl := &UserServiceCommentController{
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
		CalcedUpdateFieldsFunc: ctrl.calcServedNumFields,
		MoreProcessFunc:        nil,
	}

	info := servemodel.NewServeInfo()
	info.ServiceId = svid

	return obj.Update_With_MultiInObject(info, info, nil, nil, ren, nil, []string{"ServedNum"})
}
