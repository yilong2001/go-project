package orderups

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
	//"web/models/usermodel"
	"web/models/ordermodel"
	//"web/service/routers"
	"web/service/utils"
)

type OrderPayTypeUpController struct {
	tableName       string
	db              *sql.DB
	tx              *sql.Tx
	orderId         int
	destPayStatus   int
	destOrderStatus int
	payMoney        int
}

func (this *OrderPayTypeUpController) getTX() *sql.Tx {
	return this.tx
}

func (this *OrderPayTypeUpController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *OrderPayTypeUpController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *OrderPayTypeUpController) getDB() *sql.DB {
	return this.db
}

func (this *OrderPayTypeUpController) getTableName() string {
	return this.tableName
}

func (this *OrderPayTypeUpController) checkParams(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func (this *OrderPayTypeUpController) exInfoInit(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	return nil
}

func (this *OrderPayTypeUpController) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {
	compser := condmodel.NewCondComposerLinker(" or ")

	where := map[string]interface{}{"order_id": this.orderId}
	rule := map[string]string{"order_id": " = "}

	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	compser.SetItem(compSub)
	compser.SetNext(nil)

	ji, _ := json.Marshal(compser)
	log.Println("update : ", string(ji))

	return compser
}

func (this *OrderPayTypeUpController) calcPaytypeUpFields(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) map[string]string {

	fds := map[string]string{"prepay_status": fmt.Sprint(this.destPayStatus)}
	fds["prepay_money"] = fmt.Sprint(this.payMoney)
	fds["pay_status"] = fmt.Sprint(this.destPayStatus)
	fds["pay_money"] = fmt.Sprint(this.payMoney)
	fds["order_status"] = fmt.Sprint(this.destOrderStatus)

	return fds
}

func UpdateOrderPayType(dbin *sql.DB,
	txin *sql.Tx,
	orderid int,
	payStatus int,
	orderStatus int,
	paymoney int,
	upFields []string,
	headParams *reqparamodel.HttpReqParams,
	req *http.Request,
	ren render.Render) bool {
	ctrl := &OrderPayTypeUpController{
		tableName:       "web_orders",
		db:              dbin,
		tx:              txin,
		orderId:         orderid,
		destPayStatus:   payStatus,
		destOrderStatus: orderStatus,
		payMoney:        paymoney,
	}

	obj := &utils.UpdateObjectWithIdUtil{
		TableName:       ctrl.getTableName(),
		Db:              ctrl.getDB(),
		Tx:              ctrl.getTX(),
		FormUnParseFlag: 1,

		ExInitFunc:             ctrl.exInfoInit,
		CheckParamFunc:         ctrl.checkParams,
		CondCompserFunc:        ctrl.condCompser,
		CalcedUpdateFieldsFunc: ctrl.calcPaytypeUpFields,
		MoreProcessFunc:        nil,
	}

	info := &ordermodel.OrderInfo{}
	info.OrderId = orderid

	return obj.Update_With_MultiInObject(nil, info, headParams, req, ren, nil, upFields)
}
