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
	"time"
	//"reflect"

	"web/component/cfgutils"
	"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/component/orderutils"
	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/ordermodel"
	"web/models/rendermodel"
	"web/models/reqparamodel"
	"web/service/getter"
	"web/service/routers"
	//"web/service/users"
	"web/service/utils"
)

func init() {
	orderStatusRouterBuilder()
}

func orderStatusRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/user/mine/provider/order/:OrderId/accepter", ProcessOrderStatus)

	m.Post("/user/mine/provider/order/:OrderId/rejecter", ProcessOrderStatus)

	m.Post("/user/mine/customer/order/:OrderId/prepayer", ProcessOrderStatus)

	m.Post("/user/mine/provider/order/:OrderId/cost/dingjin", ProcessOrderStatus)

	m.Post("/user/mine/customer/order/:OrderId/date/accepter", ProcessOrderStatus)

	m.Post("/user/mine/provider/order/:OrderId/feedback", ProcessOrderStatus)

	m.Post("/user/mine/customer/order/:OrderId/complete", ProcessOrderStatus)
	m.Post("/user/mine/customer/order/:OrderId/uncomplete", ProcessOrderStatus)

	m.Post("/user/mine/customer/order/:OrderId/comment", ProcessOrderStatus)

	m.Post("/user/mine/customer/order/:OrderId/cancel", ProcessOrderStatus)

	m.Post("/user/mine/provider/order/:OrderId/cancel", ProcessOrderStatus)

	m.Post("/user/mine/provider/order/:OrderId/date/option", ProcessOrderStatus)

	m.Post("/user/mine/provider/order/:OrderId/feedback/plus", ProcessOrderStatus)

	m.Post("/user/mine/customer/order/:OrderId/comment/plus", ProcessOrderStatus)

	m.Post("/user/mine/provider/order/:OrderId/comment/plus", ProcessOrderStatus)
}

func NewOrderStatusController() *OrderStatusController {
	ctrl := &OrderStatusController{
		tableName:      "web_orders",
		oldOrderInfo:   nil,
		statusTransfer: nil,
	}

	ctrl.initDB()
	return ctrl
}

func NewOrderStatusControllerWithDb(dbin *sql.DB, txin *sql.Tx, orderinfo *ordermodel.OrderInfo, sf *ordermodel.OrderStatusTransfer) *OrderStatusController {
	ctrl := &OrderStatusController{
		tableName:      "web_orders",
		oldOrderInfo:   orderinfo,
		statusTransfer: sf,
		db:             dbin,
		tx:             txin,
	}

	return ctrl
}

type OrderStatusController struct {
	tableName      string
	oldOrderInfo   *ordermodel.OrderInfo
	statusTransfer *ordermodel.OrderStatusTransfer
	db             *sql.DB
	tx             *sql.Tx
}

func (this *OrderStatusController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}

	if this.tx == nil {
		var err error = nil
		this.tx, err = this.db.Begin()
		if err != nil {
			panic(err)
		}
	}
}

func (this *OrderStatusController) closeDB() {
	if this.tx != nil {
		this.tx.Rollback()
		this.tx = nil
	}

	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *OrderStatusController) getDB() *sql.DB {
	return this.db
}

func (this *OrderStatusController) getTX() *sql.Tx {
	return this.tx
}

func (this *OrderStatusController) commitTX() {
	this.tx.Commit()
}

func (this *OrderStatusController) getTableName() string {
	return this.tableName
}

func (this *OrderStatusController) getOldOrder() *ordermodel.OrderInfo {
	return this.oldOrderInfo
}

func (this *OrderStatusController) exInit4Single(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	info, ok := reqInfo.(*ordermodel.OrderInfo)
	if !ok {
		return errors.New("req info type wrong, not order info ")
	}

	//info.ExpiredDate = time.Now().Add(time.Hour * 24 * 60).Format("2006-01-02 15:04:05")

	err := this.exInit(info)
	if err != nil {
		return err
	}

	err = this.exInit4Rollback(info)
	if err != nil {
		return err
	}

	err = this.exInit4Prepay(info, this.oldOrderInfo.OrderPrice)

	return err
}

func (this *OrderStatusController) exInit(info *ordermodel.OrderInfo) error {
	info.OrderId = this.oldOrderInfo.OrderId
	info.OrderStatus = this.statusTransfer.NextStatus

	dt := time.Now().Format("2006-01-02 15:04:05")

	info.PrepayTime = dt
	info.PayTime = dt

	if info.Star > 5 {
		info.Star = 5
	}

	info.UpdateTime = dt
	info.OverTime = dt

	if this.oldOrderInfo.OrderStatus > orderutils.Order_Status_Wait_Customer_PrePay && this.oldOrderInfo.OrderStatus < orderutils.Order_Status_Wait_Comment {
		info.AccountLocked = 0
	}

	return nil
}

func (this *OrderStatusController) exInit4Rollback(info *ordermodel.OrderInfo) error {
	if !strings.HasSuffix(this.statusTransfer.Router, "cancel") {
		return nil
	}

	if this.oldOrderInfo.OrderStatus <= orderutils.Order_Status_Wait_Customer_PrePay {
		return nil
	}

	if this.oldOrderInfo.OrderStatus >= orderutils.Order_Status_Wait_Comment {
		return nil
	}

	isRollback := false
	if this.oldOrderInfo.OrderType == orderutils.Order_Type_Single_Business && this.statusTransfer.Role == orderutils.Order_Role_Customer {
		isRollback = orderSingle_isRollbackForCustomerCancel(this.oldOrderInfo)
	} else if this.oldOrderInfo.OrderType == orderutils.Order_Type_With_Front_Money_Business {
		isRollback = true
	} else {
		isRollback = true
	}

	if isRollback {
		info.PrepayStatus = orderutils.Pay_Status_Refund
	} else {
		info.PrepayStatus = this.oldOrderInfo.PrepayStatus
	}

	return nil
}

func (this *OrderStatusController) exInit4PrepayCheck(info *ordermodel.OrderInfo) error {
	if this.oldOrderInfo.OrderType == orderutils.Order_Type_Single_Business || this.oldOrderInfo.OrderType == orderutils.Order_Type_With_Front_Money_Business {
		if (this.statusTransfer.CurrentStatus ==
			orderutils.Order_Status_Wait_Customer_PrePay) &&
			(this.statusTransfer.NextStatus !=
				orderutils.Order_Status_Wait_Arrange_Date) {
			log.Println("order status exInit4Prepay(single), but status is wrong, cur, next", this.statusTransfer.CurrentStatus, this.statusTransfer.NextStatus)
			return errors.New("order status exInit4Prepay(single), but status is wrong, cur, next")
		}
	}

	if this.oldOrderInfo.OrderType == orderutils.Order_Type_Course {
		if (this.statusTransfer.CurrentStatus ==
			orderutils.Order_Status_Wait_Customer_PrePay) &&
			(this.statusTransfer.NextStatus !=
				orderutils.Order_Status_Wait_Provider_Feedback && this.statusTransfer.NextStatus !=
				orderutils.Order_Status_Wait_Comment) {
			log.Println("order status exInit4Prepay(course), but status is wrong, cur, next", this.statusTransfer)
			return errors.New("order status exInit4Prepay(course), but status is wrong, cur, next")
		}
	}

	return nil
}

func (this *OrderStatusController) exInit4Prepay(info *ordermodel.OrderInfo, cost int) error {
	// if (this.statusTransfer.CurrentStatus !=
	// 	orderutils.Order_Status_Wait_Customer_PrePay) ||
	// 	(this.statusTransfer.NextStatus !=
	// 		orderutils.Order_Status_Wait_Arrange_Date) {
	// 	log.Println("order status exInit4Prepay, but status is wrong, cur, next", this.statusTransfer.CurrentStatus, this.statusTransfer.NextStatus)
	// 	return nil
	// }

	if strings.HasSuffix(this.statusTransfer.Router, "cancel") {
		return nil
	}

	//if this.oldOrderInfo.OrderPrice == 0 {
	//	return nil
	//}

	//should be prepay status now
	if this.statusTransfer.CurrentStatus !=
		orderutils.Order_Status_Wait_Customer_PrePay {
		return nil
	}

	if err0 := this.exInit4PrepayCheck(info); err0 != nil {
		return err0
	}

	log.Println("this.oldOrderInfo.PrepayStatus : ", this.oldOrderInfo.PrepayStatus)

	if this.oldOrderInfo.PrepayStatus != orderutils.Pay_Status_Init && this.oldOrderInfo.PrepayStatus != orderutils.Pay_Status_Fail {
		log.Println("Error: the order is paying : ", this.oldOrderInfo.OrderId, this.oldOrderInfo.PrepayStatus)
		return errors.New("Error: the order is paying")
	}

	couponMoney := 0

	if cost > 0 && info.UserCouponId01 > 0 {
		couponinfo01, err := getter.GetModelInfoGetter().GetCouponByUserCouponId(this.getDB(), info.UserCouponId01)
		if err != nil || couponinfo01 == nil || couponinfo01.UserId != this.oldOrderInfo.CustomerId {
			return errors.New("userId in coupon is not equal with current userId")
		}

		if couponinfo01.CouponStatus != orderutils.User_Coupon_Status_Availabe {
			return errors.New("coupun is not effect")
		}

		et, err := time.Parse("2006-01-02 15:04:05", couponinfo01.ExpireTime)
		if err != nil || et.Unix() <= time.Now().Unix() {
			return errors.New("coupon is expired")
		}

		couponMoney = couponinfo01.Money

		info.PayCouponMoney01 = couponinfo01.Money
	}

	log.Println("exInit4Prepay, cost : ", cost, couponMoney)

	if cost > couponMoney {
		if info.PrepayType == orderutils.Pay_Type_Account_Balance {
			info.AccountLocked = cost - couponMoney
			info.PrepayMoney = cost - couponMoney
			info.PrepayStatus = orderutils.Pay_Status_Ok
		} else if info.PrepayType == orderutils.Pay_Type_WeiXin {
			info.PrepayMoney = cost - couponMoney
			info.PrepayStatus = orderutils.Pay_Status_Wait_Notify
			//keep order status unchange, until wxpay notify us
			info.OrderStatus = orderutils.Order_Status_Wait_Customer_PrePay
		} else if info.PrepayType == orderutils.Pay_Type_ZhiFuBao {
			info.PrepayMoney = cost - couponMoney
			info.PrepayStatus = orderutils.Pay_Status_Init
			//keep order status unchange, until wxpay notify us
			info.OrderStatus = orderutils.Order_Status_Wait_Customer_PrePay
		} else {
			return errors.New("pay type is error")
		}
	} else {
		info.PrepayStatus = orderutils.Pay_Status_Ok
	}

	return nil
}

func (this *OrderStatusController) exInit4Dingjin(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {

	info, ok := reqInfo.(*ordermodel.OrderInfo)
	if !ok {
		return errors.New("req info type wrong, not order info ")
	}

	//dur := time.Duration(this.oldOrderInfo.ContractDuration)
	if info.ExpiredDate != "" {
		expiredate, err := time.Parse(("2006-01-02 15:04:05"), info.ExpiredDate)
		if err != nil {
			log.Println("order expiredate error : ", info)
			return errors.New("expiredate is wrong!")
		}

		oldexpiredate, err := time.Parse(("2006-01-02 15:04:05"), this.oldOrderInfo.ExpiredDate)
		if err != nil {
			log.Println("oldexpiredate error : ", this.oldOrderInfo)
		} else {
			if expiredate.Unix() < oldexpiredate.Unix() {
				return errors.New("can not modify expire date more littler")
			}
		}
	} else {
		info.ExpiredDate = this.oldOrderInfo.ExpiredDate
	}

	err := this.exInit(info)
	if err != nil {
		return nil
	}

	err = this.exInit4Rollback(info)
	if err != nil {
		return err
	}

	err = this.exInit4Prepay(info, this.oldOrderInfo.FrontMoney)
	if err != nil {
		return err
	}

	err = this.exInit4DingjinPay(info)

	return err
}

func (this *OrderStatusController) exInit4DingjinPay(info *ordermodel.OrderInfo) error {
	if (this.statusTransfer.CurrentStatus !=
		orderutils.Order_Status_Wait_Customer_Complete) ||
		(this.statusTransfer.NextStatus !=
			orderutils.Order_Status_Wait_Comment) {
		return nil
	}

	if this.oldOrderInfo.PayStatus != orderutils.Pay_Status_Init && this.oldOrderInfo.PayStatus != orderutils.Pay_Status_Fail {
		log.Println("Error: the order is paying ", this.oldOrderInfo.OrderId, this.oldOrderInfo.PayStatus)
		return errors.New("Error: the order is paying or payed, refund")
	}

	if strings.HasSuffix(this.statusTransfer.Router, "uncomplete") {
		return nil
	}

	orderSubInfos, err := getter.GetModelInfoGetter().GetMultiOrderSubsByOrderIds(this.getDB(), this.oldOrderInfo.OrderId)
	if err != nil {
		log.Print(err)
		return err
	}

	totalCost := 0

	for _, subInfo := range *orderSubInfos {
		if subInfo.OverStatus == 0 {
			totalCost = totalCost + subInfo.ActualNum*subInfo.UnitCost
		}
	}

	info.PayedTotal = totalCost
	extraCost := orderFrontMoney_CalcExtraCost(totalCost, this.oldOrderInfo.FrontMoney, this.oldOrderInfo.PayCouponMoney01)

	if info.PayType == orderutils.Pay_Type_Account_Balance {
		info.PayStatus = orderutils.Pay_Status_Ok
		info.PayMoney = extraCost
	} else if info.PayType == orderutils.Pay_Type_WeiXin {
		info.PayStatus = orderutils.Pay_Status_Wait_Notify
		info.PayMoney = extraCost
		//keep order status unchange, until wxpay notify us
		info.OrderStatus = orderutils.Order_Status_Wait_Customer_Complete
	} else if info.PayType == orderutils.Pay_Type_ZhiFuBao {
		info.PayStatus = orderutils.Pay_Status_Init
		info.PayMoney = extraCost
		//keep order status unchange, until wxpay notify us
		info.OrderStatus = orderutils.Order_Status_Wait_Customer_Complete
	} else {
		return errors.New("pay type is error")
	}

	return nil
}

func (this *OrderStatusController) checkReqInfo(params *reqparamodel.HttpReqParams, r render.Render) bool {
	log.Println("userid:" + params.RouterParams["UserId"])

	if !utils.IsFieldCorrectWithRule("user_id", params.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	if !utils.IsFieldCorrectWithRule("order_id", params.RouterParams["OrderId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_OrderId_Error, "OrderId is not correct!"))
		return false
	}

	//userid, _ := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)

	// orderid, _ := strconv.ParseInt(params.RouterParams["OrderId"], 10, 32)

	// orderinfo, err := getter.GetModelInfoGetter().GetOrderByOrderId(this.db, int(orderid))
	// if err != nil {
	// 	r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_OrderId_Error, err.Error()))
	// 	return false
	// }

	// if strings.Contains(params.ShortUrl, "customer/order") {
	// 	if orderinfo.CustomerId != int(userid) {
	// 		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
	// 		return false
	// 	}
	// } else {
	// 	if orderinfo.ProviderId != int(userid) {
	// 		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
	// 		return false
	// 	}
	// }

	//this.oldOrderInfo = orderinfo
	return true
}

func (this *OrderStatusController) checkParams(params *reqparamodel.HttpReqParams, r render.Render) bool {

	if !this.checkReqInfo(params, r) {
		return false
	}

	// if this.oldOrderInfo.OrderStatus > orderutils.Order_Status_Over {
	// 	log.Println("cur status : ", this.oldOrderInfo.OrderStatus)
	// 	r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_OrderOver_Error, "order status error, has been over"))
	// 	return false
	// }

	// if this.oldOrderInfo.OrderStatus != this.reqStatus && this.reqStatus <= orderutils.Order_Status_Over {
	// 	log.Println("cur status : ", this.oldOrderInfo.OrderStatus)
	// 	r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_OrderStatus_Error, "order status error"))
	// 	return false
	// }

	return true
}

func (this *OrderStatusController) compWhereCond(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	userid, _ := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)

	orderid, _ := strconv.ParseInt(params.RouterParams["OrderId"], 10, 32)

	idWhere["order_id"] = int(orderid)
	idRlue["order_id"] = "="

	if strings.Contains(params.ShortUrl, "customer/order") {
		idWhere["customer_id"] = int(userid)
		idRlue["customer_id"] = "="
	} else {
		idWhere["provider_id"] = int(userid)
		idRlue["provider_id"] = "="
	}
	return idWhere, idRlue
}

func (this *OrderStatusController) getOldOrderInfoWithReqParam(params *reqparamodel.HttpReqParams) (*ordermodel.OrderInfo, error) {
	orderid, err := strconv.ParseInt(params.RouterParams["OrderId"], 10, 32)
	if err != nil {
		return nil, err
	}

	orderinfo, err := getter.GetModelInfoGetter().GetOrderByOrderId(this.getDB(), int(orderid))
	if err != nil {
		log.Println("GetOrderByOrderId", err)
		return nil, err
	}

	userid, _ := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)
	if strings.Contains(params.ShortUrl, "customer/order") {
		if orderinfo.CustomerId != int(userid) {
			return nil, errors.New("orderinfo.CustomerId != int(userid)")
		}
	} else {
		if orderinfo.ProviderId != int(userid) {
			return nil, errors.New("orderinfo.ProviderId != int(userid)")
		}
	}

	return orderinfo, nil
}

func (this *OrderStatusController) getOrderRoleWithReqParams(params *reqparamodel.HttpReqParams) int {
	if strings.Contains(params.ShortUrl, "mine/provider") {
		return orderutils.Order_Role_Provider
	}

	if strings.Contains(params.ShortUrl, "mine/customer") {
		return orderutils.Order_Role_Customer
	}

	panic("url router wrong")
	return orderutils.Order_Role_Customer
}

func ProcessOrderStatusV0(params *reqparamodel.HttpReqParams, mrtiniparams martini.Params, req *http.Request, ren render.Render) {
	params.MergeMartiniParams(mrtiniparams)

	ctrl := NewOrderStatusControllerWithParam(params, ren)
	if ctrl == nil {
		return
	}

	defer ctrl.closeDB()

	object := NewUserOrderStatusControllerObject4Single(ctrl)
	if ctrl.oldOrderInfo.OrderType == orderutils.Order_Type_With_Front_Money_Business {
		object = NewUserOrderStatusControllerObject4Dingjin(ctrl)
	}

	neworderInfo := ordermodel.NewOrderInfo()
	fakeren := &rendermodel.FakeMrtiniRender{}

	res := object.Util_UpdateObjectInfoWithId(neworderInfo, params, req, fakeren, nil, ctrl.statusTransfer.ChangedFields)

	appendout := map[string]interface{}{}

	next := false
	if res {
		if ctrl.statusTransfer.PayFunction != nil {
			err := ctrl.statusTransfer.PayFunction.(func(*sql.DB, *sql.Tx, *ordermodel.OrderInfo, *ordermodel.OrderInfo, *reqparamodel.HttpReqParams, *http.Request, *map[string]interface{}) error)(ctrl.getDB(), ctrl.getTX(), ctrl.oldOrderInfo, neworderInfo, params, req, &appendout)
			if err == nil {
				//good
				next = true
			} else {
				log.Println("ctrl.statusTransfer.PayFunction", err)
				ren.JSON(200,
					errcode.NewErrRsp2(errcode.Err_Form_Para_OrderStatus_Error, err.Error()))
				return
			}
		} else {
			//good
			next = true
		}
	} else {
		log.Println("Util_UpdateObjectInfoWithId failed : ", fakeren.GetVal())
		ren.JSON(200, fakeren.GetVal())
		return
	}

	if next {
		next = false
		if ctrl.statusTransfer.MoreFunction != nil {
			err := ctrl.statusTransfer.MoreFunction.(func(*sql.DB, *sql.Tx, *ordermodel.OrderInfo, *ordermodel.OrderInfo, *reqparamodel.HttpReqParams, *http.Request, *map[string]interface{}) error)(ctrl.getDB(), ctrl.getTX(), ctrl.oldOrderInfo, neworderInfo, params, req, &appendout)
			if err == nil {
				//good
				next = true
			} else {
				log.Println("ctrl.statusTransfer.MoreFunction", err)
				ren.JSON(200,
					errcode.NewErrRsp2(errcode.Err_Form_Para_OrderStatus_Error, err.Error()))
				return
			}
		} else {
			//good
			next = true
		}
	}

	ctrl.commitTX()

	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": []interface{}{appendout}})
}

func ProcessOrderStatus(params *reqparamodel.HttpReqParams, mrtiniparams martini.Params, req *http.Request, ren render.Render) {
	params.MergeMartiniParams(mrtiniparams)

	ctrl := NewOrderStatusControllerWithParam(params, ren)
	if ctrl == nil {
		return
	}

	defer ctrl.closeDB()

	object := NewUserOrderStatusControllerObject4Single(ctrl)
	if ctrl.oldOrderInfo.OrderType == orderutils.Order_Type_With_Front_Money_Business {
		object = NewUserOrderStatusControllerObject4Dingjin(ctrl)
	}

	ProcessOrderStatusWithCtrl(ctrl, object, params, req, ren)
}

func ProcessOrderStatusWithCtrl(ctrl *OrderStatusController, object *utils.ObjectWithIdUtil, params *reqparamodel.HttpReqParams, req *http.Request, ren render.Render) {

	neworderInfo := ordermodel.NewOrderInfo()
	fakeren := &rendermodel.FakeMrtiniRender{}

	res := object.Util_UpdateObjectInfoWithId(neworderInfo, params, req, fakeren, nil, ctrl.statusTransfer.ChangedFields)

	appendout := map[string]interface{}{}

	next := false
	if res {
		if ctrl.statusTransfer.PayFunction != nil {
			err := ctrl.statusTransfer.PayFunction.(func(*sql.DB, *sql.Tx, *ordermodel.OrderInfo, *ordermodel.OrderInfo, *reqparamodel.HttpReqParams, *http.Request, *map[string]interface{}) error)(ctrl.getDB(), ctrl.getTX(), ctrl.oldOrderInfo, neworderInfo, params, req, &appendout)
			if err == nil {
				//good
				next = true
			} else {
				log.Println("ctrl.statusTransfer.PayFunction", err)
				ren.JSON(200,
					errcode.NewErrRsp2(errcode.Err_Form_Para_OrderStatus_Error, err.Error()))
				return
			}
		} else {
			//good
			next = true
		}
	} else {
		log.Println("Util_UpdateObjectInfoWithId failed : ", fakeren.GetVal())
		ren.JSON(200, fakeren.GetVal())
		return
	}

	if next {
		next = false
		if ctrl.statusTransfer.MoreFunction != nil {
			err := ctrl.statusTransfer.MoreFunction.(func(*sql.DB, *sql.Tx, *ordermodel.OrderInfo, *ordermodel.OrderInfo, *reqparamodel.HttpReqParams, *http.Request, *map[string]interface{}) error)(ctrl.getDB(), ctrl.getTX(), ctrl.oldOrderInfo, neworderInfo, params, req, &appendout)
			if err == nil {
				//good
				next = true
			} else {
				log.Println("ctrl.statusTransfer.MoreFunction", err)
				ren.JSON(200,
					errcode.NewErrRsp2(errcode.Err_Form_Para_OrderStatus_Error, err.Error()))
				return
			}
		} else {
			//good
			next = true
		}
	}

	ctrl.commitTX()

	ren.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": []interface{}{appendout}})
}

func NewOrderStatusControllerWithParam(params *reqparamodel.HttpReqParams, ren render.Render) *OrderStatusController {
	ctrl := NewOrderStatusController()

	orderinfo, err := ctrl.getOldOrderInfoWithReqParam(params)
	if err != nil {
		ren.JSON(200,
			errcode.NewErrRsp2(errcode.Err_Form_Para_OrderStatus_Error, err.Error()))
		return nil
	}

	role := ctrl.getOrderRoleWithReqParams(params)

	destMachine := GetStatusTransferMachine(orderinfo.OrderType, orderinfo.OrderStatus, role, params.ShortUrl, orderinfo)
	if destMachine == nil {
		ren.JSON(200,
			errcode.NewErrRsp2(errcode.Err_Form_Para_OrderStatus_Error, "url router is wrong."))
		return nil
	}

	ctrl.oldOrderInfo = orderinfo
	ctrl.statusTransfer = destMachine

	return ctrl
}

func NewUserOrderStatusControllerObject4Single(ctrl *OrderStatusController) *utils.ObjectWithIdUtil {
	object := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),
		Tx:        ctrl.getTX(),

		CheckParamFuncForUpdate: ctrl.checkParams,
		ExpendInitFuncForUpdate: ctrl.exInit4Single,
		WhereCondFuncForUpdate:  ctrl.compWhereCond,
	}

	return object
}

func NewUserOrderStatusControllerObject4Dingjin(ctrl *OrderStatusController) *utils.ObjectWithIdUtil {
	object := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),
		Tx:        ctrl.getTX(),

		CheckParamFuncForUpdate: ctrl.checkParams,
		ExpendInitFuncForUpdate: ctrl.exInit4Dingjin,
		WhereCondFuncForUpdate:  ctrl.compWhereCond,
	}

	return object
}

func NewUserOrderStatusControllerObject4Course(ctrl *OrderStatusController) *utils.ObjectWithIdUtil {
	object := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),
		Tx:        ctrl.getTX(),

		CheckParamFuncForUpdate: ctrl.checkParams,
		ExpendInitFuncForUpdate: ctrl.exInit4Single,
		WhereCondFuncForUpdate:  ctrl.compWhereCond,
	}

	return object
}
