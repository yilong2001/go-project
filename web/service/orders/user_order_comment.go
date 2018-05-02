package orders

import (
	"database/sql"
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/render"
	//"log"

	//"errors"
	//"net/http"
	//"strconv"
	//"strings"
	//"time"
	//"reflect"

	//"web/component/cfgutils"
	//"web/component/errcode"
	"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/component/orderutils"
	//"web/dal/sqldrv"
	//"web/models/basemodel"
	"web/models/ordermodel"
	//"web/models/reqparamodel"
	//"web/service/getter"
	//"web/service/routers"
	"web/service/utils"
)

type UserOrderCommentController struct {
	tableName     string
	genIdFlag     string
	commentIdFlag string
	db            *sql.DB
	tx            *sql.Tx
}

func NewUserOrderCommentController(dbin *sql.DB, txin *sql.Tx) *UserOrderCommentController {
	return &UserOrderCommentController{
		tableName:     "web_order_comments",
		genIdFlag:     "order_comment_id",
		commentIdFlag: "comment_id",
		db:            dbin,
		tx:            txin,
	}
}

func (this *UserOrderCommentController) getTableName() string {
	return this.tableName
}
func (this *UserOrderCommentController) getGenIdFlag() string {
	return this.genIdFlag
}
func (this *UserOrderCommentController) getCommentIdFlag() string {
	return this.commentIdFlag
}
func (this *UserOrderCommentController) addComment(info *ordermodel.OrderCommentInfo) error {
	object := &utils.ObjectWithIdUtil{
		TableName: this.getTableName(),
		Db:        this.db,
		Tx:        this.tx,
	}

	info.CommentId = idutils.GetId(this.getGenIdFlag())
	return object.CreateObjectWithInfo(info)
}
