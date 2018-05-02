package wxutils

import (
	"fmt"
	"strconv"
	"strings"
	"web/component/cfgutils"
	"web/component/orderutils"
)

func CompWxPayTradeId(orderId, orderType, orderStatus int) string {
	if orderStatus == orderutils.Order_Status_Wait_Customer_PrePay {
		return fmt.Sprintf("%dA%dA%dA%d", orderId, orderType, orderStatus, 0)
	}

	if orderStatus == orderutils.Order_Status_Wait_Customer_Complete {
		return fmt.Sprintf("%dA%dA%dA%d", orderId, orderType, orderStatus, 0)
	}

	return fmt.Sprintf("%dAfAfAf", orderId)
}

//return orderid, ordertype, orderstatus, appendcode
func ParseWxPayTradeId(tradeId string) (int, int, int, int) {
	ids := strings.Split(tradeId, "A")
	if len(ids) != 4 {
		return 0, 0, 0, 0
	}

	orderid, err := strconv.ParseInt(ids[0], 10, 32)
	if err != nil {
		return 0, 0, 0, 0
	}

	ordertype, err := strconv.ParseInt(ids[1], 10, 32)
	if err != nil {
		return int(orderid), 0, 0, 0
	}

	orderstatus, err := strconv.ParseInt(ids[2], 10, 32)
	if err != nil {
		return int(orderid), int(ordertype), 0, 0
	}

	apcode, err := strconv.ParseInt(ids[3], 10, 32)
	if err != nil {
		return int(orderid), int(ordertype), int(orderstatus), 0
	}

	return int(orderid), int(ordertype), int(orderstatus), int(apcode)
}

//var wxPayH5CallbackUrl string
func init() {
	url := cfgutils.GetWebApiConfig().WxpayH5Http + "://" + cfgutils.GetWebApiConfig().WxpayH5Domain
	fmt.Println("wxpay notify url init : ", url)
}

func GetWxPayNotifyUrl() string {
	return "/third/wxpay/notify"
}

func CompNotifyUrl() string {
	url := cfgutils.GetWebApiConfig().WxpayH5Http + "://" + cfgutils.GetWebApiConfig().WxpayH5Domain + GetWxPayNotifyUrl()
	fmt.Println(url)
	return url
}
