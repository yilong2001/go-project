package thirdpay

import (
	//"database/sql"
	//"errors"
	"fmt"
	"log"
	"net/http"
	//"strconv"
	//"web/component/cfgutils"
	"web/component/aliutils"
	"web/component/orderutils"
	//"web/component/wxutils"
	//"web/dal/sqldrv"
	//"web/models/ordermodel"
	//"web/models/rendermodel"
	//"web/service/getter"
	//"web/service/orderups"
	//"web/service/pays"
	"web/service/routers"
)

func init() {
	orderStatusAliPayRouterBuilder()
}

func orderStatusAliPayRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/third/alipay/notify", ProcessAlipayCallback)
	m.Get("/third/alipay/notify", ProcessAlipayCallbackGet)
}

func GetAliPayNotifyUrl() string {
	return "/third/alipay/notify"
}

func ProcessAlipayCallbackGet(w http.ResponseWriter, r *http.Request) {
	w.(http.ResponseWriter).WriteHeader(http.StatusOK)
	fmt.Fprint(w.(http.ResponseWriter), "test ok!")
	return
}

func alipayParseTradeStatus(tradestatus string) int {
	if tradestatus == "WAIT_BUYER_PAY" {
		return orderutils.Pay_Status_Wait_Notify
	}

	if tradestatus == "TRADE_CLOSED" {
		return orderutils.Pay_Status_Fail
	}

	if tradestatus == "TRADE_SUCCESS" {
		return orderutils.Pay_Status_Ok
	}

	if tradestatus == "TRADE_FINISHED" {
		return orderutils.Pay_Status_Ok
	}

	return orderutils.Pay_Status_Wait_Notify
}

func ProcessAlipayCallback(w http.ResponseWriter, r *http.Request) {
	err, result := aliutils.AlipayNotify(r)
	if err != nil {
		w.(http.ResponseWriter).WriteHeader(http.StatusOK)
		fmt.Fprint(w.(http.ResponseWriter), "SUCCESS")
		return
	}

	log.Println("alipay notify: ", result)

	isuccess := false
	if alipayParseTradeStatus(result.TradeStatus) == orderutils.Pay_Status_Ok {
		isuccess = true
	} else if alipayParseTradeStatus(result.TradeStatus) == orderutils.Pay_Status_Wait_Notify {
		w.(http.ResponseWriter).WriteHeader(http.StatusOK)
		fmt.Fprint(w.(http.ResponseWriter), "SUCCESS")
		return
	} else {
		isuccess = false
	}

	orderId := result.OutTradeNo
	totalFee := result.TotalAmount

	myCtrl := NewOrderStatusWxPayController()
	myCtrl.initDB()
	defer myCtrl.closeDB()

	err = myCtrl.initOrder(orderId)
	if err != nil {
		log.Println("alipay cb init order failed:", orderId)
		w.(http.ResponseWriter).WriteHeader(http.StatusOK)
		fmt.Fprint(w.(http.ResponseWriter), "SUCCESS")
		return
	}

	if isuccess { //付款成功，处理订单
		myCtrl.isWxpaySuccess = true
	} else {
		myCtrl.isWxpaySuccess = false
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
		log.Println("alipay cb process order failed: ", orderId, localFee)
	} else {
	}

	w.(http.ResponseWriter).WriteHeader(http.StatusOK)
	fmt.Fprint(w.(http.ResponseWriter), "SUCCESS")
}
