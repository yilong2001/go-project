package aliutils

import (
	"fmt"
	// "strconv"
	// "strings"
	"web/component/cfgutils"
)

const (
	Prefix_ProductID = ""
	Prefix_APPID     = ""
	Prefix_APPSECRET = ""
	Partner           = ""
	Prefix_Key       = ""
	Prefix_Email     = ""
)

func GetZhiEasyProductId() string {
	return Prefix_ProductID
}

func GetZhiEasyAppId() string {
	return Prefix_APPID
}

func GetZhiEasyAppSecret() string {
	return Prefix_APPSECRET
}

func GetPartner() string {
	return Partner
}

func GetZhiEasyKey() string {
	return Prefix_Key
}

func GetZhiEasyEmail() string {
	return Prefix_Email
}

func GetAliPayNotifyUrl() string {
	return "/third/alipay/notify"
}

func GetAliPayPublicKey() string {
	return ""
}

func GetAliPayPrivateKey() string {
	return ""
}

func CompNotifyUrl() string {
	url := cfgutils.GetWebApiConfig().WxpayH5Http + "://" + cfgutils.GetWebApiConfig().WxpayH5Domain + GetAliPayNotifyUrl()
	fmt.Println(url)
	return url
}

func GetZhifubaoPublicKey() string {
	return ""
}
