package orders

import (
	"database/sql"
	//"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"

	"errors"
	"net/http"
	"strconv"
	"strings"
	//"time"
	//"reflect"

	"web/component/cfgutils"
	"web/component/errcode"
	"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/component/orderutils"
	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/ordermodel"
	//"web/models/rendermodel"
	"web/models/reqparamodel"
	"web/service/getter"
	"web/service/routers"
	//"web/service/users"
	"web/service/utils"
)

func init() {
	orderSubRouterBuilder()
}

func orderSubRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/user/mine/customer/orderjob", CreateOrderSub)
	m.Post("/user/mine/customer/orderjob/:OrderSubId", UpdateOrderSubJob)
	m.Post("/user/mine/provider/orderjob/:OrderSubId", UpdateOrderSubJob)

	m.Delete("/user/mine/customer/orderjob/:OrderSubId", RemoveOrderSubJob)
}

func NewOrderSubController(dbin *sql.DB, txin *sql.Tx) *OrderSubController {
	ctrl := &OrderSubController{
		tableName: "web_order_subs",
		genIdFlag: "order_sub_id",
		db:        dbin,
		tx:        txin,
	}

	return ctrl
}

func NewOrderSubControllerObject(ctrl *OrderSubController) *utils.ObjectWithIdUtil {
	object := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),
		Tx:        ctrl.getTX(),

		ExpendInitFuncForCreate: ctrl.exInit4Create,
		CheckParamFuncForCreate: ctrl.checkParams,

		CheckParamFuncForUpdate: ctrl.checkUpParams,
		ExpendInitFuncForUpdate: ctrl.exInit4Up,
		WhereCondFuncForUpdate:  ctrl.compWhereCond,
		MoreProcessForUpdate:    ctrl.moreForUp,

		CheckParamFuncForUpdateDelStatus: ctrl.checkUpParams,
		ExpendInitFuncForUpdateDelStatus: ctrl.exInit4Del,
		WhereCondFuncForUpdateDelStatus:  ctrl.compWhereCond,
	}

	return object
}

type OrderSubController struct {
	tableName string
	genIdFlag string
	db        *sql.DB
	tx        *sql.Tx
	subId     int
	userId    int

	orderInfo *ordermodel.OrderInfo
}

func (this *OrderSubController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *OrderSubController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *OrderSubController) getDB() *sql.DB {
	return this.db
}

func (this *OrderSubController) getTX() *sql.Tx {
	return this.tx
}

func (this *OrderSubController) getTableName() string {
	return this.tableName
}

func (this *OrderSubController) getGenIdFlag() string {
	return this.genIdFlag
}

func (this *OrderSubController) exInit4Create(reqInfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) error {

	info, ok := reqInfo.(*ordermodel.OrderSubInfo)
	if !ok {
		return errors.New("req info type wrong, not order info ")
	}

	orderinfo, err := getter.GetModelInfoGetter().GetOrderByOrderId(this.getDB(), int(info.OrderId))
	if err != nil {
		return err
	}

	if orderinfo.OrderStatus != orderutils.Order_Status_Wait_Provider_Accept {
		return errors.New(" can not add sub order for current status ")
	}

	userid, _ := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)
	if orderinfo.CustomerId != int(userid) {
		return errors.New("userid is error")
	}

	info.OrderSubId = idutils.GetId(this.getGenIdFlag())

	info.ReferType = orderutils.Order_Sub_Refer_Type_Job
	info.ActualNum = 0
	info.UnitCost = 0
	info.OverTime = "1970-01-01 00:00:00"

	return nil
}

func (this *OrderSubController) checkParams(params *reqparamodel.HttpReqParams, r render.Render) bool {
	if !utils.IsFieldCorrectWithRule("user_id", params.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	return true
}

func (this *OrderSubController) exInit4Up(reqInfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) error {

	info, ok := reqInfo.(*ordermodel.OrderSubInfo)
	if !ok {
		return errors.New("req info type wrong, not order info ")
	}

	log.Println(info)

	ordersubinfo, err := getter.GetModelInfoGetter().GetOrderSubByOrderSubId(this.getDB(), int(this.subId))
	if err != nil {
		log.Println("OrderSubController:exInit4Up:GetOrderSubByOrderSubId", err)
		return err
	}

	orderinfo, err := getter.GetModelInfoGetter().GetOrderByOrderId(this.getDB(), int(ordersubinfo.OrderId))
	if err != nil {
		log.Println("OrderSubController:exInit4Up:GetOrderByOrderId", err)
		return err
	}

	if orderinfo.OrderStatus > orderutils.Order_Status_Wait_Customer_Complete {
		return errors.New(" can not add sub order for current status ")
	}

	if orderinfo.PayStatus > orderutils.Pay_Status_Init {
		return errors.New(" paying order can not update sub info ")
	}

	// if orderinfo.CustomerId != this.userId {
	// 	return errors.New(" current user is not customer user ")
	// }

	//info.ReferType = orderutils.Order_Sub_Refer_Type_Job
	//info.OverTime = "1970-01-01 00:00:00"

	this.orderInfo = orderinfo

	return nil
}

func (this *OrderSubController) exInit4Del(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	return nil
}

func (this *OrderSubController) checkUpParams(params *reqparamodel.HttpReqParams, r render.Render) bool {
	if !utils.IsFieldCorrectWithRule("user_id", params.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	if !utils.IsFieldCorrectWithRule("order_sub_id", params.RouterParams["OrderSubId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	usrid, _ := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)
	subid, _ := strconv.ParseInt(params.RouterParams["OrderSubId"], 10, 32)

	this.userId = int(usrid)
	this.subId = int(subid)

	return true
}

func (this *OrderSubController) compWhereCond(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	idWhere["order_sub_id"] = int(this.subId)
	idRlue["order_sub_id"] = " = "

	// if strings.Contains(params.ShortUrl, "customer/order") {
	// 	idWhere["customer_id"] = int(userid)
	// 	idRlue["customer_id"] = " = "
	// } else {
	// 	idWhere["provider_id"] = int(userid)
	// 	idRlue["provider_id"] = " = "
	// }
	return idWhere, idRlue
}

func (this *OrderSubController) moreForUp(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	order_More_NotifyMsg(this.getDB(), nil, this.orderInfo, nil, headParams, nil, nil)
	return nil
}

func CreateOrderSub(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewOrderSubController(nil, nil)
	ctrl.initDB()
	defer ctrl.closeDB()

	obj := NewOrderSubControllerObject(ctrl)

	info := &ordermodel.OrderSubInfo{}

	obj.Util_CreateObjectWithId(info, headParams, req, r)
}

func getSpecFieldsForUpdate(headParams *reqparamodel.HttpReqParams) []string {
	if strings.Contains(headParams.ShortUrl, "provider") {
		return []string{"UnitCost"}
	}

	outf := []string{}

	postfields := headParams.PostFields
	for _, pf := range postfields {
		if strings.ToLower(pf) == strings.ToLower("ReferName") {
			outf = append(outf, "ReferName")
		}
		if strings.ToLower(pf) == strings.ToLower("ExpectNum") {
			outf = append(outf, "ExpectNum")
		}
		if strings.ToLower(pf) == strings.ToLower("ActualNum") {
			outf = append(outf, "ActualNum")
		}
	}

	return outf
}

func UpdateOrderSubJob(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewOrderSubController(nil, nil)
	ctrl.initDB()
	defer ctrl.closeDB()

	obj := NewOrderSubControllerObject(ctrl)

	info := &ordermodel.OrderSubInfo{}

	obj.Util_UpdateObjectInfoWithId(info, headParams, req, r, nil, getSpecFieldsForUpdate(headParams))
}

func RemoveOrderSubJob(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewOrderSubController(nil, nil)
	defer ctrl.closeDB()
	obj := NewOrderSubControllerObject(ctrl)

	info := &ordermodel.OrderSubInfo{}

	obj.Util_UpdateDelStatusWithId(info, headParams, req, ren)
}
