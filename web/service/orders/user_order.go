package orders

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"

	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	//"reflect"

	"web/component/cfgutils"
	"web/component/errcode"
	"web/component/idutils"
	"web/component/wxutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/component/aliutils"
	"web/component/orderutils"
	"web/component/randutils"

	//"web/component/rongcloud"
	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/clientmodel"
	"web/models/coursemodel"
	"web/models/ordermodel"
	"web/models/rendermodel"
	"web/models/reqparamodel"
	"web/models/usermodel"
	"web/service/getter"
	"web/service/immsgs"

	"web/service/orderpays"
	"web/service/orderups"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userOrderRouterBuilder()
}

func userOrderRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/user/mine/customer/singleorder", CreateSingleOrder)
	m.Post("/user/mine/customer/dingjinorder", CreateDingjinOrder)
	m.Post("/user/mine/customer/courseorder", CreateCourseOrder)

	m.Get("/user/mine/order/:OrderId", GetOrderById)

	m.Get("/user/mine/customer/order", GetCustomerOrder)
	m.Get("/user/mine/provider/order", GetProviderOrder)

	m.Get("/user/mine/customer/order/:OrderId", GetCustomerOrder)
	m.Get("/user/mine/provider/order/:OrderId", GetProviderOrder)
}

func NewUserOrderControllerObject(ctrl *UserOrderController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForCreate: ctrl.exInit4Create,
		ExpendInitFuncForGet:    ctrl.exInit4Get,
		ExpendInitFuncForUpdate: nil,
		MoreProcessForCreate:    ctrl.more4Create,

		CheckParamFuncForCreate: ctrl.checkUserId,
		CheckParamFuncForGet:    ctrl.checkUserId,
		CheckParamFuncForUpdate: nil,

		WhereCondFuncForGet:    ctrl.compCond4Get,
		WhereCondFuncForUpdate: nil,

		AppendMoreResultFunc: ctrl.appendMoreInfo4Order,
	}

	return obj
}

func NewUserOrdercontroller() *UserOrderController {
	ctrl := &UserOrderController{
		tableName: "web_orders",
		genIdFlag: "order_id",
		isOnlyId:  false,
		userId:    0,
		orderType: 0,
	}

	ctrl.initDB()

	return ctrl
}

type UserOrderController struct {
	tableName string
	genIdFlag string
	db        *sql.DB
	tx        *sql.Tx

	isOnlyId  bool
	userId    int
	orderType int
}

func (this *UserOrderController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}

	if this.tx == nil {
		var err error
		this.tx, err = this.db.Begin()
		if err != nil {
			panic(err)
		}
	}
}

func (this *UserOrderController) closeDB() {
	if this.tx != nil {
		this.tx.Rollback()
	}

	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}
func (this *UserOrderController) getDB() *sql.DB {
	return this.db
}
func (this *UserOrderController) getTX() *sql.Tx {
	return this.tx
}

func (this *UserOrderController) getTableName() string {
	return this.tableName
}

func (this *UserOrderController) getGenIdFlag() string {
	return this.genIdFlag
}

func (this *UserOrderController) compCond4Get(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {

	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	if !this.isOnlyId {
		if uid, err := strconv.ParseInt(params.RouterParams["UserId"], 10, 32); err == nil {
			if strings.Contains(params.ShortUrl, "customer/order") {
				idWhere["customer_id"] = int(uid)
				idRlue["customer_id"] = "="
			} else if strings.Contains(params.ShortUrl, "provider/order") {
				idWhere["provider_id"] = int(uid)
				idRlue["provider_id"] = "="
			} else {
				idWhere["customer_id"] = int(0)
				idRlue["customer_id"] = "="
			}
		} else {
			idWhere["customer_id"] = int(0)
			idRlue["customer_id"] = "="
		}
	}

	if orderid, err := strconv.ParseInt(params.RouterParams["OrderId"], 10, 32); err == nil {
		idWhere["order_id"] = int(orderid)
		idRlue["order_id"] = "="
	}

	log.Println("compWhereCondition : ", idWhere)

	return idWhere, idRlue
}

func (this *UserOrderController) checkUserId(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	log.Println("userid:" + headParams.RouterParams["UserId"])

	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	return true
}

func (this *UserOrderController) exInit4Create(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	info, ok := reqInfo.(*ordermodel.OrderInfo)
	if !ok {
		return errors.New("req info type wrong, not order info ")
	}

	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)

	if this.orderType == orderutils.Order_Type_Course {
		return this.exInit4CourseOrderCreate(info, int(userid))
	}

	if len(info.CustomerExpect) < 10 {
		return errors.New(" customer expect too little ")
	}

	if len(info.CustomerIntroduce) < 10 {
		return errors.New(" customer introduce too little ")
	}

	if len(info.CustomerDateAdvice) < 5 && len(info.ExpiredDate) < 5 {
		return errors.New(" customer date advice too little, or expired date is empty ")
	}

	var db *sql.DB = this.getDB()

	//info := this.getObjInfo()
	log.Println("user order exInit4Create info : ", info)

	serveinfo, err := getter.GetModelInfoGetter().GetServiceByServiceId(db, info.ServiceId, nil, nil)
	if err != nil {
		log.Println(info.ServiceId, "can not get service info")
		return err
	}

	if serveinfo.UserId == int(userid) {
		log.Println("can order with self service")
		return errors.New("can order with self service")
	}

	dt := time.Now().Format("2006-01-02 15:04:05")

	info.OrderType = this.orderType

	info.ServiceName = serveinfo.ServiceName
	info.ProviderId = serveinfo.UserId
	info.OrderPrice = serveinfo.Price
	info.OrderStatus = orderutils.Order_Status_Wait_Provider_Accept

	// info.PrepayType = 1
	// info.PrepayStatus = orderutils.Pay_Status_Init
	// info.PrepayMoney =
	info.PrepayTime = "1970-01-01 00:00:00"

	// info.PayType =
	// info.PayStatus =
	// info.PayMoney =
	info.PayTime = "1970-01-01 00:00:00"

	info.CreateTime = dt
	info.UpdateTime = dt
	info.OverTime = time.Now().Add(time.Hour * 24 * 90).Format("2006-01-02 15:04:05")

	if info.ExpiredDate == "" {
		info.ExpiredDate = time.Now().Add(time.Hour * 24 * 60).Format("2006-01-02 15:04:05")
	}

	expiredate, err := time.Parse(("2006-01-02 15:04:05"), info.ExpiredDate)
	if err != nil {
		log.Println("order expiredate error : ", info)
		return errors.New("expiredate is wrong!")
	}

	if expiredate.Unix() < time.Now().Add(time.Hour*24*20).Unix() {
		log.Println("order expiredate too short : ", info)
		return errors.New("expiredate is short!")
	}

	info.ArragedDateOp = orderutils.Order_DateArrange_Op_UnSelect
	//info.ArragedDate = "1970-01-01 00:00:00"

	info.ArragedDateOption1 = time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
	info.ArragedDateOption2 = time.Now().Add(time.Hour * 48).Format("2006-01-02 15:04:05")

	info.Feedback = ""
	info.OverReason = ""

	info.PayCouponId01 = 0
	info.PayCouponMoney01 = 0

	// if info.UserCouponId01 > 0 {
	// 	couponinfo01, err := getter.GetModelInfoGetter().GetCouponByUserCouponId(db, info.UserCouponId01)
	// 	if err != nil || couponinfo01 == nil || couponinfo01.UserId != int(userid) {
	// 		return errors.New("userId in coupon is not equal with current userId")
	// 	}

	// 	if couponinfo01.CouponStatus != orderutils.User_Coupon_Status_Ready {
	// 		return errors.New("coupun is not effect")
	// 	}

	// 	info.PayCouponId01 = couponinfo01.CouponId
	// 	info.PayCouponMoney01 = couponinfo01.Money
	// } else {
	// 	info.PayCouponId01 = 0
	// 	info.PayCouponMoney01 = 0
	// }

	info.CustomerId = int(userid)
	info.OrderId = idutils.GetId(this.getGenIdFlag())

	//md5 := randutils.BuildMd5PWPhoneStringV64(fmt.Sprint(info.OrderId)+dt+fmt.Sprint(userid), "")
	//md5 = strings.ToUpper(md5)

	//info.OrderOutId = md5[0:4] + "-" + md5[4:8] + "-" + md5[8:12] + "-" + md5[12:16]

	info.OrderOutId = randutils.BuildMD5OrderOutId(fmt.Sprint(info.OrderId)+dt+fmt.Sprint(userid), info.OrderId)

	return nil
}

func (this *UserOrderController) exInit4CourseOrderCreate(info *ordermodel.OrderInfo, userid int) error {

	courseinfo, err := getter.GetModelInfoGetter().GetCourseMainByCourseId(this.getDB(), info.ServiceId)
	if err != nil {
		log.Println(info.ServiceId, "can not get course info")
		return err
	}

	if courseinfo.UserId == int(userid) {
		log.Println("can order with self course")
		return errors.New("can order with self course")
	}

	dt := time.Now().Format("2006-01-02 15:04:05")

	info.OrderType = this.orderType

	info.ServiceName = courseinfo.Title
	info.ProviderId = courseinfo.UserId
	info.OrderPrice = courseinfo.NowPrice
	info.OrderStatus = orderutils.Order_Status_Wait_Customer_PrePay

	// info.PrepayType = 1
	// info.PrepayStatus = orderutils.Pay_Status_Init
	// info.PrepayMoney =
	info.PrepayTime = "1970-01-01 00:00:00"

	// info.PayType =
	// info.PayStatus =
	// info.PayMoney =
	info.PayTime = "1970-01-01 00:00:00"

	info.CreateTime = dt
	info.UpdateTime = dt
	info.OverTime = time.Now().Add(time.Hour * 24 * 90).Format("2006-01-02 15:04:05")
	info.ExpiredDate = time.Now().Add(time.Hour * 24 * 90).Format("2006-01-02 15:04:05")

	info.ArragedDateOp = orderutils.Order_DateArrange_Op_UnSelect
	//info.ArragedDate = "1970-01-01 00:00:00"

	info.ArragedDateOption1 = dt
	info.ArragedDateOption2 = dt

	info.Feedback = ""
	info.OverReason = ""

	info.PayCouponId01 = 0
	info.PayCouponMoney01 = 0

	info.CustomerId = (userid)
	info.OrderId = idutils.GetId(this.getGenIdFlag())

	//md5 := randutils.BuildMd5PWPhoneStringV64(fmt.Sprint(info.OrderId)+dt+fmt.Sprint(userid), "")
	//md5 = strings.ToUpper(md5)

	//info.OrderOutId = md5[0:4] + "-" + md5[4:8] + "-" + md5[8:12] + "-" + md5[12:16]
	info.OrderOutId = randutils.BuildMD5OrderOutId(fmt.Sprint(info.OrderId)+dt+fmt.Sprint(userid), info.OrderId)
	return nil
}

func (this *UserOrderController) more4Create(reqInfo basemodel.ObjectUtilBaseIf) error {
	info, ok := reqInfo.(*ordermodel.OrderInfo)
	if !ok {
		log.Println("order more for create, req info type wrong, not order info ")
		return errors.New("req info type wrong, not order info ")
	}

	err := immsgs.SendPlatformNotifyMsg(info.ProviderId, "Order", "您有新订单...", info.OrderId, info.OrderType)
	log.Println("order, more for create ", err)

	user, err1 := getter.GetModelInfoGetter().GetUserByUserId(this.getDB(), info.ProviderId, nil, nil)
	if err1 == nil {
		err := aliutils.AliSmsMsgSend(user.Phone, user.UserName)
		if err != nil {
			log.Println("ali sms code send failed, ", err)
		}
	} else {
		log.Println(err1)
	}

	return nil
}

func (this *UserOrderController) exInit4Get(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	this.userId = int(userid)

	return nil
}

func (this *UserOrderController) getNextOrderStatus(orderInfo ordermodel.OrderInfo) int {
	ordertype := orderInfo.OrderType
	svrid := orderInfo.ServiceId

	switch ordertype {
	case orderutils.Order_Type_Single_Business:
		return orderutils.Order_Status_Wait_Arrange_Date

	case orderutils.Order_Type_With_Front_Money_Business:
		return orderutils.Order_Status_Wait_Arrange_Date

	case orderutils.Order_Type_Course:
		cm, err := getter.GetModelInfoGetter().GetCourseMainByCourseId(this.getDB(), svrid)
		if err != nil {
			return orderutils.Order_Status_Wait_Provider_Feedback
		}

		switch cm.CourseType {
		case coursemodel.Const_Course_Type_Online:
			return orderutils.Order_Status_Wait_Comment
		}

		return orderutils.Order_Status_Wait_Provider_Feedback
	}

	return orderutils.Order_Status_Wait_Arrange_Date
}

func (this *UserOrderController) appendMoreInfo4Order(result *[]interface{}) *[]interface{} {
	fakeren := &rendermodel.FakeMrtiniRender{}

	totals := len(*result)

	orderInfos := []interface{}{}

	found := false
	ui := &usermodel.UserInfo{}
	for _, odIf := range *result {
		if orderInfo, ok := odIf.(ordermodel.OrderInfo); ok {
			clientInfo := &clientmodel.ClientOrderInfo{}
			if this.userId != orderInfo.CustomerId {
				orderInfo.OrderOutId = ""
			}

			clientInfo.OrderInfo = &orderInfo

			expiredate, err := time.Parse(("2006-01-02 15:04:05"), orderInfo.ExpiredDate)
			if err != nil {
				log.Println("order expiredate error : ", orderInfo)
				clientInfo.IsExpired = true
				clientInfo.IsForbidden = true
			}

			if expiredate.Unix() < time.Now().Unix() {
				clientInfo.IsExpired = true
			}

			if expiredate.Add(time.Hour*24*60).Unix() < time.Now().Unix() {
				clientInfo.IsForbidden = true
			}

			usInfo, err := getter.GetModelInfoGetter().GetUserByUserId(this.getDB(), orderInfo.CustomerId, ui.GetSkipFieldsForOrderQuery(), nil)
			if err == nil {
				clientInfo.CustomerInfo = usInfo
			} else {
				log.Println(err)
			}

			upInfo, err := getter.GetModelInfoGetter().GetUserByUserId(this.getDB(), orderInfo.ProviderId, ui.GetSkipFieldsForOrderQuery(), nil)
			if err == nil {
				clientInfo.ProviderInfo = upInfo
			} else {
				log.Println(err)
			}

			// if orderInfo.CustomerId == this.userId || orderInfo.ProviderId == this.userId {
			// 	found = true
			// }

			if (totals == 1) && orderInfo.OrderType == orderutils.Order_Type_With_Front_Money_Business {
				clientInfo.OrderSubInfos, err = getter.GetModelInfoGetter().GetMultiOrderSubsByOrderIds(this.getDB(), orderInfo.OrderId)
				if err != nil {
					log.Print(err)
				}
			}

			isComplete := false
			isOk := true
			if (totals == 1) && orderInfo.OrderStatus == orderutils.Order_Status_Wait_Customer_PrePay && orderInfo.PrepayStatus == orderutils.Pay_Status_Wait_Notify {
				err, rsp := wxutils.WXPayQueryOrder(wxutils.CompWxPayTradeId(orderInfo.OrderId, orderInfo.OrderType, orderInfo.OrderStatus))
				//err = nil //rsp = wxutils.GetWXQueryRspForDebug()
				if err != nil {
					log.Println(err)
				} else {
					upfields := []string{"PrepayStatus"}
					paystatus := wxutils.WxpayParseTradeStatus(rsp.TradeState)
					nextOrderStatus := orderInfo.OrderStatus
					if paystatus == orderutils.Pay_Status_Ok {
						nextOrderStatus = this.getNextOrderStatus(orderInfo)
						if nextOrderStatus == orderutils.Order_Status_Wait_Comment {
							isComplete = true
						}
						upfields = append(upfields, "OrderStatus", "PrepayMoney")
					} else if paystatus == orderutils.Pay_Status_Wait_Notify {
						//
					} else if paystatus == orderutils.Pay_Status_Refund {
						nextOrderStatus = orderutils.Order_Status_Pay_Closed
						upfields = append(upfields, "OrderStatus")
					} else {
						//
					}

					money := int(rsp.TotalFee)
					if money == 0 {
						money = orderInfo.PrepayMoney
					}

					ok := orderups.UpdateOrderPayType(this.getDB(), this.getTX(),
						orderInfo.OrderId,
						paystatus, nextOrderStatus, money,
						upfields,
						nil, nil, fakeren)
					if !ok {
						isOk = false
						log.Println(" update order prepay status failed ")
					}
				}
			}

			if (totals == 1) && orderInfo.OrderStatus == orderutils.Order_Status_Wait_Customer_Complete && orderInfo.PayStatus == orderutils.Pay_Status_Wait_Notify {
				err, rsp := wxutils.WXPayQueryOrder(wxutils.CompWxPayTradeId(orderInfo.OrderId, orderInfo.OrderType, orderInfo.OrderStatus))
				if err != nil {
					log.Println(err)
				} else {
					upfields := []string{"PayStatus"}
					paystatus := wxutils.WxpayParseTradeStatus(rsp.TradeState)
					nextOrderStatus := orderInfo.OrderStatus
					if paystatus == orderutils.Pay_Status_Ok {
						nextOrderStatus = orderutils.Order_Status_Wait_Comment
						isComplete = true
						upfields = append(upfields, "OrderStatus", "PayMoney")
					} else if paystatus == orderutils.Pay_Status_Wait_Notify {
						//
					} else if paystatus == orderutils.Pay_Status_Refund {
						nextOrderStatus = orderutils.Order_Status_Pay_Closed
						upfields = append(upfields, "OrderStatus")
					} else {
						//
					}

					money := int(rsp.TotalFee)
					if money == 0 {
						money = orderInfo.PayMoney
					}

					ok := orderups.UpdateOrderPayType(this.getDB(), this.getTX(),
						orderInfo.OrderId,
						paystatus, nextOrderStatus, money,
						upfields,
						nil, nil, fakeren)
					if !ok {
						isOk = false
						log.Println(" update order prepay status failed ")
					}
				}
			}

			if isComplete {
				err = orderpays.Order_PayCompleteToProvider(this.getDB(), this.getTX(), orderInfo.OrderPrice, &orderInfo, &orderInfo, nil, nil, fakeren)
				if err != nil {
					isOk = false
					log.Println("pay complete to provider failed : ", err)
				}
			}

			if isComplete && isOk {
				this.getTX().Commit()
			}

			orderInfos = append(orderInfos, clientInfo)
		}
	}

	if found {
		return result
	}

	return &orderInfos
}

func createOrder(ordertype int, headParams *reqparamodel.HttpReqParams, req *http.Request, r render.Render) {
	ctrl := NewUserOrdercontroller()
	ctrl.orderType = ordertype

	defer ctrl.closeDB()

	obj := NewUserOrderControllerObject(ctrl)

	orderInfo := ordermodel.NewOrderInfo()

	obj.Util_CreateObjectWithId(orderInfo, headParams, req, r)
}

func CreateSingleOrder(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	createOrder(orderutils.Order_Type_Single_Business, headParams, req, r)
}

func CreateDingjinOrder(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	createOrder(orderutils.Order_Type_With_Front_Money_Business, headParams, req, r)
}

func CreateCourseOrder(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	createOrder(orderutils.Order_Type_Course, headParams, req, r)
}

func GetCustomerOrder(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserOrdercontroller()
	defer ctrl.closeDB()

	obj := NewUserOrderControllerObject(ctrl)

	orderInfo := ordermodel.NewOrderInfo()

	obj.Util_GetObjectWithId(orderInfo, headParams, req, r, orderInfo.GetSkipFieldsForCustomerQuery(), nil)
}

func GetProviderOrder(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserOrdercontroller()
	defer ctrl.closeDB()
	obj := NewUserOrderControllerObject(ctrl)

	orderInfo := ordermodel.NewOrderInfo()

	obj.Util_GetObjectWithId(orderInfo, headParams, req, r, orderInfo.GetSkipFieldsForProviderQuery(), nil)
}

func GetOrderById(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserOrdercontroller()
	defer ctrl.closeDB()
	ctrl.isOnlyId = true

	obj := NewUserOrderControllerObject(ctrl)

	orderInfo := ordermodel.NewOrderInfo()

	obj.Util_GetObjectWithId(orderInfo, headParams, req, r, orderInfo.GetSkipFieldsForProviderQuery(), nil)
}
