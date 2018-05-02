package orderpays

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	//"strings"
	"web/component/aliutils"
	"web/component/orderutils"
	"web/component/wxutils"
	// "web/dal/sqldrv"
	// "web/models/basemodel"
	"web/models/ordermodel"
	//"web/models/rendermodel"
	"web/models/reqparamodel"
	//"web/service/coupons"
	//"web/service/getter"
	//"web/service/immsgs"
	"web/service/pays"
	//"web/service/serveups"
	//"web/service/userups"
)

func Order_PayWithAliPay(db *sql.DB, tx *sql.Tx, cost int, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render, out *map[string]interface{}) error {
	if cost == 0 {
		log.Println("alipay, but pay is zero")
		return nil
	}

	// userinfo, _err := getter.GetModelInfoGetter().GetUserByUserId(db, oldOrder.CustomerId, nil, nil)
	// if _err != nil {
	//  log.Println("Order_PayWithAliPay GetUserByUserId", _err)
	//  return _err
	// }

	returnrul := headParams.URLParams.Get("ReturnUrl")

	params := aliutils.GetAlipayH5Params(wxutils.CompWxPayTradeId(oldOrder.OrderId, oldOrder.OrderType, oldOrder.OrderStatus), oldOrder.ServiceName, cost, returnrul, aliutils.CompNotifyUrl())

	log.Println("Order_PayWithAliPay.GetAlipayH5Params ok ,", params)

	(*out)["response"] = params

	return nil
}

func Order_RollbackWithAliPay(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render) error {
	log.Println("Order_RollbackWithAliPay start ")

	orderid := wxutils.CompWxPayTradeId(oldOrder.OrderId, oldOrder.OrderType, orderutils.Order_Status_Wait_Customer_PrePay)

	err, rsp := aliutils.AlipayRefund(orderid, oldOrder.PrepayMoney)
	if err != nil {
		log.Println("alipay rollback", err)
		return err
	}

	log.Println("alipay refund rsp : ", rsp)

	if rsp.Fail != nil {
		return errors.New(rsp.Fail.Response.Msg)
	}

	err = pays.AddPayRecord(db, tx, oldOrder.CustomerId, oldOrder.CustomerId, fmt.Sprintf("order id %d, %d pay rollback with alipay %d", oldOrder.OrderId, oldOrder.CustomerId, oldOrder.PrepayMoney))
	if err != nil {
		log.Println("alipay record", err)
		return err
	}
	return nil
}
