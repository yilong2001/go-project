package aliutils

import (
	"crypto"
	"encoding/base64"
	//"encoding/hex"
	"errors"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	//"net/url"
	//"regexp"
	"strconv"
	//"strings"
	"web/component/keyutils"
	"web/component/randutils"
)

func AlipayNotify(r *http.Request) (error, *AlipayNotifyResult) {
	result := &AlipayNotifyResult{}

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)

	r.ParseForm()

	log.Println("alipay notify : ")
	presign := ""
	for field, value := range r.PostForm {
		log.Println(field, value)
		if field != "sign" && field != "sign_type" {
			kv, err := (value[0]), error(nil) //url.QueryUnescape
			if err != nil {
				log.Println(field, value, err)
			} else {
				m[field] = kv
			}
		}

		if field == "sign" {
			presign = value[0]
		}
	}

	strPreSign, _err := genAlipaySignString(m)
	if _err != nil {
		log.Println("error get sign string, reason:", _err)
		return _err, nil
	}

	log.Println("strPreSign=", strPreSign)
	log.Println("sign=", presign)

	ppBy, err := base64.StdEncoding.DecodeString(presign)
	if err != nil {
		log.Println("base64 decode sign string, reason:", err, presign)
		return err, nil
	}

	err = randutils.RsaVerify(keyutils.GetZhifubaoPublicKeyStr(), []byte(strPreSign), ppBy, crypto.SHA1)
	if err != nil {
		log.Println(" RsaVerify string, reason:", err, presign)
		return err, nil
	}

	if m["app_id"] != GetZhiEasyAppId() {
		log.Println("app id is wrong", m["app_id"])
		return errors.New("app id is wrong"), nil
	}

	result.NotifyTime = fmt.Sprintf("%v", m["notify_time"])
	result.OutTradeNo = fmt.Sprintf("%v", m["out_trade_no"])
	result.TradeNo = fmt.Sprintf("%v", m["trade_no"])
	result.TradeStatus = fmt.Sprintf("%v", m["trade_status"])

	tm, err := strconv.ParseFloat(fmt.Sprintf("%v", m["total_amount"]), 64)
	if err != nil {
		result.TotalAmount = 0
	} else {
		result.TotalAmount = int(tm * 100)
	}

	return nil, result
}
