package pays

import (
	"database/sql"
	//"github.com/go-martini/martini"
	//"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	//"errors"
	//"net/http"
	//"strconv"
	//"strings"
	"time"
	//"reflect"

	//"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/dal/sqldrv"
	//"web/models/basemodel"
	//"web/models/clientmodel"
	//"web/models/jobmodel"
	"web/models/ordermodel"
	//"web/models/rendermodel"
	"web/models/reqparamodel"
	//"web/models/servemodel"
	//"web/models/usermodel"
	//"web/service/getter"
	//"web/service/routers"
	//"web/service/serves"
	"web/service/utils"
)

func NewPayRecordController(dbin *sql.DB, txin *sql.Tx, srcid, dstid int, infoin string) *PayRecordController {
	ctrl := &PayRecordController{
		tableName: "web_pay_records",
		db:        dbin,
		tx:        txin,
		srcId:     srcid,
		destId:    dstid,
		info:      infoin,
	}

	return ctrl
}

type PayRecordController struct {
	tableName string
	info      string
	srcId     int
	destId    int
	db        *sql.DB
	tx        *sql.Tx
}

func (this *PayRecordController) commitDB() {
	this.tx.Commit()
}

func (this *PayRecordController) getDB() *sql.DB {
	return this.db
}
func (this *PayRecordController) getTX() *sql.Tx {
	return this.tx
}
func (this *PayRecordController) getTableName() string {
	return this.tableName
}

func (this *PayRecordController) getRecrod() *ordermodel.PayRecordInfo {

	record := &ordermodel.PayRecordInfo{}

	record.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	record.Record = this.info
	record.SourceId = this.srcId
	record.DestId = this.destId

	return record
}

func (this *PayRecordController) check(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func AddPayRecord(dbin *sql.DB, txin *sql.Tx, srcId, dstId int, info string) error {
	ctrl := NewPayRecordController(dbin, txin, srcId, dstId, info)

	object := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.db,
		Tx:        ctrl.tx,
	}

	return object.CreateObjectWithInfo(ctrl.getRecrod())
}
