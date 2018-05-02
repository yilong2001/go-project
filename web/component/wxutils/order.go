package wxutils

import (
	//"bytes"
	//"encoding/xml"
	//"errors"
	"fmt"
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	//"io/ioutil"
	"log"
	//"net/http"
	"strings"
	"time"
	"web/component/randutils"
)

func WXPayUnifiedOrder(shopname, clientip, orderid, notifyurl string, totalFee int, openid string) (error, *WXPayUnifiedOrderRspToClient) {
	clt := core.NewClient(GetZhiEasyAppId(), GetMCHID(), GetZhiEasyAppKey(), nil)

	req := &pay.UnifiedOrderRequest{}
	req.Body = shopname
	req.NonceStr = string(randutils.KRandAll(12))
	req.OutTradeNo = orderid
	req.TotalFee = int64(totalFee)

	cip := strings.Split(clientip, ":")
	req.SpbillCreateIP = cip[0]
	req.TimeStart = time.Now().Format("20060102150405")
	req.TimeExpire = time.Now().Add(time.Hour * 1).Format("20060102150405")
	req.NotifyURL = notifyurl
	req.TradeType = "JSAPI"
	req.OpenId = openid

	log.Println("WXPayUnifiedOrder req : ", req)

	rsp, err := pay.UnifiedOrder2(clt, req)
	if err != nil {
		log.Println("WXPayUnifiedOrder", err)
		return err, nil
	}

	outrsp := &WXPayUnifiedOrderRspToClient{}
	outrsp.Appid = GetZhiEasyAppId()
	outrsp.Timestamp = fmt.Sprint(time.Now().UnixNano() / 1000000)
	outrsp.Nonce = string(randutils.KRandAll(12))
	outrsp.Package = "prepay_id=" + rsp.PrepayId
	outrsp.Signtype = "MD5"

	outsign := core.JsapiSign(clt.AppId(), outrsp.Timestamp, outrsp.Nonce, outrsp.Package, outrsp.Signtype, clt.ApiKey())

	outrsp.Sign = outsign

	return nil, outrsp
}

