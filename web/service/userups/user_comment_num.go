package userups

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

	//"web/models/servemodel"
	"web/models/usermodel"
	//"web/service/routers"
	"web/service/utils"
)

type UserCommentController struct {
	tableName    string
	db           *sql.DB
	tx           *sql.Tx
	commentLevel int
	userId       int
}

func (this *UserCommentController) getTX() *sql.Tx {
	return this.tx
}

func (this *UserCommentController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UserCommentController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UserCommentController) getDB() *sql.DB {
	return this.db
}

func (this *UserCommentController) getTableName() string {
	return this.tableName
}

func (this *UserCommentController) checkParams(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func (this *UserCommentController) exInfoInit(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *UserCommentController) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {
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

func (this *UserCommentController) calcUpFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {
	fds := map[string]string{}

	fds["avg_star"] = " avg_star = (comment_num * avg_star + " + fmt.Sprint(this.commentLevel) + ")/(comment_num+1)"

	fds["comment_num"] = " comment_num = comment_num + 1"

	return fds
}

func UpdateUserCommentNum(dbin *sql.DB,
	txin *sql.Tx,
	level int,
	uid int,
	ren render.Render) bool {
	ctrl := &UserCommentController{
		tableName:    "web_users",
		db:           dbin,
		tx:           txin,
		commentLevel: level,
		userId:       uid,
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

	return obj.Update_With_MultiInObject(info, info, nil, nil, ren, nil, []string{"CommentNum", "AvgStar"})
}
