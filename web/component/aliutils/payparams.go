package aliutils

import (
	//"github.com/ascoders/alipay"
	"log"
	"strconv"
	"time"
)

func GetAlipayH5Params(orderid, svrname string, fee int, returnUrl, notifyUrl string) *AlipayParameters {
	param := NewH5AlipayParameters()
	param.Timestatmp = time.Now().Format("2006-01-02 15:04:05")
	param.Method = "alipay.trade.wap.pay"
	param.NotifyUrl = notifyUrl
	param.ReturnUrl = returnUrl
	param.Version = "1.0"

	param.BizContent["body"] = ComposeBody(svrname, fee)
	param.BizContent["subject"] = svrname
	param.BizContent["out_trade_no"] = orderid
	param.BizContent["total_amount"] = strconv.FormatFloat(float64(fee*1.0)/100.0, 'f', 2, 32)
	param.BizContent["product_code"] = "QUICK_WAP_PAY"

	sign := param.DoSign()
	if sign == "" {
		log.Println(" alipay params sign error, ", param)
		return nil
	}

	return param
}
