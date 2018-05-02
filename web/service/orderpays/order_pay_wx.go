package orderpays

import (
	"database/sql"
	//"errors"
	"fmt"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	//"strings"
	//"web/component/aliutils"
	"web/component/orderutils"
	"web/component/wxutils"
	// "web/dal/sqldrv"
	// "web/models/basemodel"
	"web/models/ordermodel"
	//"web/models/rendermodel"
	"web/models/reqparamodel"
	//"web/service/coupons"
	"web/service/getter"
	//"web/service/immsgs"
	"web/service/pays"
	//"web/service/serveups"
	//"web/service/userups"
)

func Order_PayWithWxPay(db *sql.DB, tx *sql.Tx, cost int, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render, out *map[string]interface{}) error {
	if cost == 0 {
		log.Println("wxpay, but pay is zero")
		return nil
	}

	userinfo, _err := getter.GetModelInfoGetter().GetUserByUserId(db, oldOrder.CustomerId, nil, nil)
	if _err != nil {
		log.Println("Order_PayWithWxPay GetUserByUserId", _err)
		return _err
	}

	//for test
	err, outrsp := wxutils.WXPayUnifiedOrder(oldOrder.ServiceName,
		headParams.TokenParams["addr"],
		wxutils.CompWxPayTradeId(oldOrder.OrderId, oldOrder.OrderType, oldOrder.OrderStatus),
		wxutils.CompNotifyUrl(),
		cost, userinfo.WeixinOpenId)

	if err != nil {
		log.Println("Order_PayWithWxPay WXPayUnifiedOrder", err)
		return err
	}

	log.Println("Order_PayWithWxPay.WXPayUnifiedOrder ok ,", outrsp)

	(*out)["response"] = outrsp

	return nil
}

func Order_RollbackWithWxPay(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render) error {
	log.Println("Order_RollbackWithWxPay start ")

	orderid := wxutils.CompWxPayTradeId(oldOrder.OrderId, oldOrder.OrderType, orderutils.Order_Status_Wait_Customer_PrePay)

	err, rsp := wxutils.WXPayRefund(orderid, oldOrder.PrepayMoney)
	if err != nil {
		log.Println("wxpay refund wrong : ", err)
		return err
	}

	log.Println("wxpay refund rsp : ", rsp)

	err = pays.AddPayRecord(db, tx, oldOrder.CustomerId, oldOrder.CustomerId, fmt.Sprintf("order id %d, %d pay rollback with wxpay %d", oldOrder.OrderId, oldOrder.CustomerId, oldOrder.PrepayMoney))
	if err != nil {
		log.Println("wxpay rollback record : ", err)
		return err
	}
	return nil
}
