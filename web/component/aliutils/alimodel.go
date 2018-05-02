package aliutils

import (
	"crypto"
	//"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	//"net/url"
	"strconv"
	"strings"
	"web/component/keyutils"
	"web/component/randutils"
)

const (
	Alipay_Gateway = "https://openapi.alipay.com/gateway.do"
)

type AlipayParameters struct {
	AppId      string `json:"app_id"`  //合作者身份ID
	Method     string `json:"method"`  //卖家支付宝邮箱
	Charset    string `json:"charset"` //网站编码
	Timestatmp string `json:"timestamp"`
	Version    string `json:"version"`
	Format     string `json:"format"`

	ReturnUrl string `json:"return_url"` //回调url
	NotifyUrl string `json:"notify_url"` //异步通知页面

	Sign     string `json:"sign"`      //签名，生成签名时忽略
	SignType string `json:"sign_type"` //签名类型，生成签名时忽略

	BizContent map[string]interface{} `json:"biz_content"`
}

func NewH5AlipayParameters() *AlipayParameters {
	para := &AlipayParameters{
		AppId:      GetZhiEasyAppId(),
		Charset:    "UTF-8",
		SignType:   "RSA",
		Format:     "json",
		BizContent: map[string]interface{}{},
	}

	return para
}

func (this *AlipayParameters) DoSign() string {
	m := map[string]interface{}{}
	m["version"] = this.Version
	m["app_id"] = this.AppId
	m["method"] = this.Method
	m["charset"] = this.Charset
	m["timestamp"] = this.Timestatmp
	m["format"] = this.Format
	m["sign_type"] = this.SignType

	if this.ReturnUrl != "" && len(this.ReturnUrl) > 2 {
		m["return_url"] = this.ReturnUrl
	}

	if this.NotifyUrl != "" && len(this.NotifyUrl) > 2 {
		m["notify_url"] = (this.NotifyUrl)
	}

	//m["return_url"] = this.ReturnUrl
	//m["notify_url"] = this.NotifyUrl

	bc, err := json.Marshal(this.BizContent)
	log.Println(err)

	bcjson := strings.Replace(string(bc), "BizContent", "", -1)

	m["biz_content"] = bcjson

	strPreSign, _err := genAlipaySignString(m)
	if _err != nil {
		fmt.Println("error get sign string, reason:", _err)
		return ""
	}

	log.Println("DoSign Before Str : ", strPreSign)

	rsasign, err := randutils.RsaSign(keyutils.GetAlipayPrivateKeyStr(), []byte(strPreSign), crypto.SHA1)
	if err != nil {
		log.Println(err)
		return ""
	}

	//log.Println("rsasign", rsasign)

	sign := base64.StdEncoding.EncodeToString(rsasign)

	this.Sign = sign

	log.Println("sign : ", sign)

	return sign
}

func ComposeBody(svrname string, fee int) string {
	return "for(" + svrname + ")pay:" + strconv.FormatFloat(float64(fee*1.0)/100.0, 'f', 2, 32) + "(rmb)"
}

type AlipayNotifyResult struct {
	NotifyTime  string
	TradeNo     string
	OutTradeNo  string
	TradeStatus string
	TotalAmount int
}

type AlipayTradeRefundResponseFail struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

type AlipayTradeRefundResponse struct {
	Code         string `json:"code"`
	Msg          string `json:"msg"`
	BuyerLogonId string `json:"buyer_logon_id"`
	BuyerUserId  string `json:"buyer_user_id"`
	FundChange   string `json:"fund_change"`
	GmtRefundPay string `json:"gmt_refund_pay"`
	OutTradeNo   string `json:"out_trade_no"`
	RefundFee    string `json:"refund_fee"`
}

//alipay_trade_refund_response
type AlipayRefundReponseFail struct {
	Response AlipayTradeRefundResponseFail `json:"alipay_trade_refund_response"`

	Sign string `json:"sign"`
}

type AlipayRefundReponseSucc struct {
	Response AlipayTradeRefundResponse `json:"alipay_trade_refund_response"`

	Sign string `json:"sign"`
}

type AlipayRefundReponse struct {
	Succ *AlipayRefundReponseSucc
	Fail *AlipayRefundReponseFail
}
