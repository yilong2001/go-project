package wxutils

import (
	//"bytes"
	//"encoding/xml"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	"strconv"
	//"io/ioutil"
	//"log"
	//"net/http"
	//"strings"
	//"time"
	"web/component/orderutils"
	"web/component/randutils"
)

func WxpayParseTradeStatus(tradestatus string) int {
	if tradestatus == "SUCCESS" {
		return orderutils.Pay_Status_Ok
	}

	if tradestatus == "REFUND" {
		return orderutils.Pay_Status_Refund
	}

	if tradestatus == "NOTPAY" {
		return orderutils.Pay_Status_Fail
	}

	if tradestatus == "CLOSED" {
		return orderutils.Pay_Status_Fail
	}

	if tradestatus == "REVOKED" {
		return orderutils.Pay_Status_Refund
	}

	if tradestatus == "USERPAYING" {
		return orderutils.Pay_Status_Wait_Notify
	}

	if tradestatus == "PAYERROR" {
		return orderutils.Pay_Status_Fail
	}

	return orderutils.Pay_Status_Fail
}

func GetWXQueryRspForDebug() *pay.OrderQueryResponse {
	debug := new(pay.OrderQueryResponse)
	debug.TradeState = "SUCCESS"
	return debug
}

func WXPayQueryOrder(orderid string) (error, *pay.OrderQueryResponse) {
	clt := core.NewClient(GetZhiEasyAppId(), GetMCHID(), GetZhiEasyAppKey(), nil)

	m1 := make(map[string]string, 8)
	m1["appid"] = clt.AppId()
	m1["mch_id"] = clt.MchId()
	m1["out_trade_no"] = orderid
	m1["nonce_str"] = string(randutils.KRandAll(12))

	m1["sign"] = core.Sign(m1, clt.ApiKey(), md5.New)

	m2, err := pay.OrderQuery(clt, m1)
	if err != nil {
		return err, nil
	}

	fmt.Println("WXPayQueryOrder outmap: ", m2)

	resultCode, ok := m2["result_code"]
	if !ok {
		return errors.New("result code is not exist"), nil
	}

	if resultCode != "SUCCESS" {
		return errors.New("result code is faild"), nil
	}

	resp := &pay.OrderQueryResponse{
		AppId: m2["appid"],
		MchId: m2["mch_id"],

		OpenId:         m2["openid"],
		TradeType:      m2["trade_type"],
		TradeState:     m2["trade_state"],
		BankType:       m2["bank_type"],
		TransactionId:  m2["transaction_id"],
		OutTradeNo:     m2["out_trade_no"],
		TimeEnd:        m2["time_end"],
		TradeStateDesc: m2["trade_state_desc"],
	}

	var (
		n int64
		//id  int
		//str string
	)
	if n, err = strconv.ParseInt(m2["total_fee"], 10, 64); err != nil {
		err = fmt.Errorf("parse total_fee:%q to int64 failed: %s", m2["total_fee"], err.Error())
		fmt.Println("wxpay query order total_fee : ", err)
	} else {
		resp.TotalFee = n
	}
	if n, err = strconv.ParseInt(m2["cash_fee"], 10, 64); err != nil {
		err = fmt.Errorf("parse cash_fee:%q to int64 failed: %s", m2["cash_fee"], err.Error())
		fmt.Println("wxpay query order cash_fee : ", err)
	} else {
		resp.CashFee = n
	}

	return nil, resp
}
