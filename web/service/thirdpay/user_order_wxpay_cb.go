package thirdpay

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	//"strconv"
	"web/component/cfgutils"
	"web/component/orderutils"
	"web/component/wxutils"
	"web/dal/sqldrv"
	"web/models/coursemodel"
	"web/models/ordermodel"
	"web/models/rendermodel"
	"web/service/getter"
	"web/service/orderpays"
	"web/service/orderups"
	"web/service/pays"
	"web/service/routers"
)

func init() {
	orderStatusWxPayRouterBuilder()
}

func orderStatusWxPayRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/third/wxpay/notify", ProcessWxpayCallback)
	m.Get("/third/wxpay/notify", ProcessWxpayCallbackGet)
}

func ProcessWxpayCallbackGet(w http.ResponseWriter, r *http.Request) {
	//w.(http.ResponseWriter).WriteHeader(http.StatusOK)
	//fmt.Fprint(w.(http.ResponseWriter), "test ok!")
	orderId := "24001A0A3A0"
	totalFee := 5

	myCtrl := NewOrderStatusWxPayController()
	myCtrl.initDB()
	defer myCtrl.closeDB()

	err := myCtrl.initOrder(orderId)
	if err != nil {
		rsp := wxutils.CompXmlStrForRspMsg("FAIL", err.Error())
		w.(http.ResponseWriter).WriteHeader(http.StatusOK)
		log.Println("wxpay cb init order failed:", orderId)
		fmt.Fprint(w.(http.ResponseWriter), rsp)
		return
	}

	myCtrl.isWxpaySuccess = true

	if totalFee == 0 {
		if myCtrl.oldOrderInfo.OrderStatus == orderutils.Order_Status_Wait_Customer_PrePay {
			totalFee = myCtrl.oldOrderInfo.PrepayMoney
		} else {
			totalFee = myCtrl.oldOrderInfo.PayMoney
		}
	}

	localFee := totalFee

	err = myCtrl.process(localFee)
	if err != nil {
		log.Println("wxpay cb process order failed: ", orderId, localFee)
		rsp := wxutils.CompXmlStrForRspMsg("FAIL", err.Error())
		w.(http.ResponseWriter).WriteHeader(http.StatusOK)
		fmt.Fprint(w.(http.ResponseWriter), rsp)
	} else {
		rsp := wxutils.CompXmlStrForRspMsg("SUCCESS", "OK")
		w.(http.ResponseWriter).WriteHeader(http.StatusOK)
		fmt.Fprint(w.(http.ResponseWriter), rsp)
	}

	return
}

func ProcessWxpayCallback(w http.ResponseWriter, r *http.Request) {
	returnErr, resultErr, notifyInfo := wxutils.WxpayNotifyCallback(w, r)
	if returnErr != nil {
		log.Println("wxpay cb  returnErr")
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if notifyInfo == nil {
		log.Println("wxpay cb  notifyInfo")
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orderId := notifyInfo.OutTradeNo
	totalFee := int(notifyInfo.TotalFee)

	myCtrl := NewOrderStatusWxPayController()
	myCtrl.initDB()
	defer myCtrl.closeDB()

	err := myCtrl.initOrder(orderId)
	if err != nil {
		rsp := wxutils.CompXmlStrForRspMsg("FAIL", err.Error())
		w.(http.ResponseWriter).WriteHeader(http.StatusOK)
		log.Println("wxpay cb init order failed:", orderId)
		fmt.Fprint(w.(http.ResponseWriter), rsp)
		return
	}

	if resultErr != nil {
		myCtrl.isWxpaySuccess = false
	} else {
		myCtrl.isWxpaySuccess = true
	}

	if totalFee == 0 {
		if myCtrl.oldOrderInfo.OrderStatus == orderutils.Order_Status_Wait_Customer_PrePay {
			totalFee = myCtrl.oldOrderInfo.PrepayMoney
		} else {
			totalFee = myCtrl.oldOrderInfo.PayMoney
		}
	}

	localFee := totalFee

	err = myCtrl.process(localFee)
	if err != nil {
		log.Println("wxpay cb process order failed: ", orderId, localFee)
		rsp := wxutils.CompXmlStrForRspMsg("FAIL", err.Error())
		w.(http.ResponseWriter).WriteHeader(http.StatusOK)
		fmt.Fprint(w.(http.ResponseWriter), rsp)
	} else {
		rsp := wxutils.CompXmlStrForRspMsg("SUCCESS", "OK")
		w.(http.ResponseWriter).WriteHeader(http.StatusOK)
		fmt.Fprint(w.(http.ResponseWriter), rsp)
	}
}

type OrderStatusWxPayController struct {
	tableName      string
	oldOrderInfo   *ordermodel.OrderInfo
	isWxpaySuccess bool
	db             *sql.DB
	tx             *sql.Tx

	isPayComplete bool
}

func NewOrderStatusWxPayController() *OrderStatusWxPayController {
	return &OrderStatusWxPayController{
		tableName:      "web_orders",
		oldOrderInfo:   nil,
		isWxpaySuccess: false,
		db:             nil,
		tx:             nil,

		isPayComplete: false,
	}
}

func (this *OrderStatusWxPayController) initDB() {
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

func (this *OrderStatusWxPayController) closeDB() {
	if this.tx != nil {
		this.tx.Rollback()
		this.tx = nil
	}

	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *OrderStatusWxPayController) getDB() *sql.DB {
	return this.db
}

func (this *OrderStatusWxPayController) getTX() *sql.Tx {
	return this.tx
}

func (this *OrderStatusWxPayController) commitTX() {
	this.tx.Commit()
}

func (this *OrderStatusWxPayController) getTableName() string {
	return this.tableName
}

func (this *OrderStatusWxPayController) getNextOrderStatus() int {
	ordertype := this.oldOrderInfo.OrderType
	svrid := this.oldOrderInfo.ServiceId

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
			this.isPayComplete = true
			return orderutils.Order_Status_Wait_Comment
		}

		return orderutils.Order_Status_Wait_Provider_Feedback
	}

	return orderutils.Order_Status_Wait_Arrange_Date
}

func (this *OrderStatusWxPayController) processPrepay(totalFee int) error {
	if this.oldOrderInfo.OrderStatus != orderutils.Order_Status_Wait_Customer_PrePay {
		log.Println("processPrepay order status wrong, ", this.oldOrderInfo.OrderStatus)
		return nil
	}

	if this.oldOrderInfo.PrepayType == orderutils.Pay_Type_ZhiFuBao {
		if this.oldOrderInfo.PrepayStatus == orderutils.Pay_Status_Ok {
			log.Println("has pay success!")
			return nil
		}
	} else if this.oldOrderInfo.PrepayStatus != orderutils.Pay_Status_Wait_Notify {
		log.Println("has been processed")
		return nil
	}

	upfields := []string{"PrepayMoney", "PrepayStatus", "OrderStatus"}

	nextOrderStatus := this.getNextOrderStatus()

	paystatus := orderutils.Pay_Status_Ok
	payMoney := totalFee
	if !this.isWxpaySuccess {
		paystatus = orderutils.Pay_Status_Fail
		nextOrderStatus = this.oldOrderInfo.OrderStatus
		payMoney = 0
	}

	fakeren := &rendermodel.FakeMrtiniRender{}

	ok := orderups.UpdateOrderPayType(this.getDB(), this.getTX(),
		this.oldOrderInfo.OrderId,
		paystatus, nextOrderStatus, payMoney,
		upfields,
		nil, nil, fakeren)
	if !ok {
		log.Println(" update order prepay status failed ")
		return errors.New(" update order prepay status failed ")
	}

	if this.isWxpaySuccess {
		err := pays.AddPayRecord(this.getDB(), this.getTX(), this.oldOrderInfo.CustomerId, orderutils.GetSystmePayerId(), fmt.Sprintf("order id %d, user %d prepay with weixinpay %d (fen)", this.oldOrderInfo.OrderId, this.oldOrderInfo.CustomerId, totalFee))
		if err != nil {
			log.Println("add pay record failed", err)
			return err
		}

		if this.isPayComplete {
			err = orderpays.Order_PayCompleteToProvider(this.getDB(), this.getTX(), this.oldOrderInfo.OrderPrice, this.oldOrderInfo, this.oldOrderInfo, nil, nil, fakeren)
			if err != nil {
				log.Println("pay complete to provider failed : ", err)
				return err
			}
		}
	}

	this.commitTX()

	return nil
}

func (this *OrderStatusWxPayController) processPay(totalFee int) error {
	if this.oldOrderInfo.OrderType != orderutils.Order_Type_With_Front_Money_Business || this.oldOrderInfo.OrderStatus != orderutils.Order_Status_Wait_Customer_Complete {
		log.Println("processPay order type&status wrong, ", this.oldOrderInfo.OrderType, this.oldOrderInfo.OrderStatus)
		return nil
	}

	if this.oldOrderInfo.PayType == orderutils.Pay_Type_ZhiFuBao {
		if this.oldOrderInfo.PayStatus == orderutils.Pay_Status_Ok {
			log.Println("has pay success!")
			return nil
		}
	} else if this.oldOrderInfo.PayStatus != orderutils.Pay_Status_Wait_Notify {
		log.Println("has been processed")
		return nil
	}

	upfields := []string{"PayMoney", "PayStatus", "OrderStatus"}

	nextOrderStatus := orderutils.Order_Status_Wait_Comment

	paystatus := orderutils.Pay_Status_Ok
	payMoney := totalFee

	if !this.isWxpaySuccess {
		paystatus = orderutils.Pay_Status_Fail
		nextOrderStatus = this.oldOrderInfo.OrderStatus
		payMoney = 0
	}

	fakeren := &rendermodel.FakeMrtiniRender{}

	ok := orderups.UpdateOrderPayType(this.getDB(), this.getTX(),
		this.oldOrderInfo.OrderId,
		paystatus, nextOrderStatus, payMoney,
		upfields,
		nil, nil, fakeren)
	if !ok {
		return errors.New(" update order pay status failed ")
	}

	if this.isWxpaySuccess {
		err := pays.AddPayRecord(this.getDB(), this.getTX(), this.oldOrderInfo.CustomerId, orderutils.GetSystmePayerId(), fmt.Sprintf("order id %d, user %d pay with weixinpay %d (fen)", this.oldOrderInfo.OrderId, this.oldOrderInfo.CustomerId, totalFee))
		if err != nil {
			log.Println("add pay record failed : ", err)
			return err
		}

		lastCost := payMoney + this.oldOrderInfo.FrontMoney
		//this.oldOrderInfo.PayedTotal
		//if lastCost < this.oldOrderInfo.FrontMoney {
		//	lastCost = this.oldOrderInfo.FrontMoney
		//}

		err = orderpays.Order_PayCompleteToProvider(this.getDB(), this.getTX(), lastCost, this.oldOrderInfo, this.oldOrderInfo, nil, nil, fakeren)
		if err != nil {
			log.Println("pay complete to provider failed : ", err)
			return err
		}
	}

	this.commitTX()

	return nil
}

func (this *OrderStatusWxPayController) process(totalFee int) error {
	if this.oldOrderInfo.OrderStatus == orderutils.Order_Status_Wait_Customer_PrePay {
		err := this.processPrepay(totalFee)
		return err
	}

	if this.oldOrderInfo.OrderStatus == orderutils.Order_Status_Wait_Customer_Complete {
		err := this.processPay(totalFee)
		return err
	}

	return nil
}

func (this *OrderStatusWxPayController) initOrder(reqorderid string) error {
	orderid, ordertype, orderstatus, _ := wxutils.ParseWxPayTradeId(reqorderid)

	if (orderid + ordertype + orderstatus) == 0 {
		log.Println("wxpay order id wrong : ", reqorderid)
		return errors.New("wxpay order id wrong : " + reqorderid)
	}

	orderinfo, err := getter.GetModelInfoGetter().GetOrderByOrderId(this.getDB(), orderid)
	if err != nil {
		log.Println("wxpay  GetOrderByOrderId ", err)
		return err
	}

	if orderstatus != orderinfo.OrderStatus || ordertype != orderinfo.OrderType {
		log.Println("wxpay order id wrong : ", reqorderid)
		return errors.New("wxpay order id wrong : " + reqorderid)
	}

	this.oldOrderInfo = orderinfo
	return nil
}
