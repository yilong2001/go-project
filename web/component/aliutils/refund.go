package aliutils

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	//"web/component/fileutils"
	"web/component/keyutils"
	"web/component/randutils"
)

func AlipayRefund(orderid string, totalFee int) (error, *AlipayRefundReponse) {
	sDate := time.Now().Format("2006-01-02 15:04:05")

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)

	m["app_id"] = GetZhiEasyAppId()
	m["format"] = "json"
	m["charset"] = "UTF-8"

	m["method"] = "alipay.trade.refund"
	m["timestamp"] = sDate
	m["version"] = "1.0"
	m["sign_type"] = "RSA"

	m["biz_content"] = `{"out_trade_no":"` + orderid + `", "refund_amount":"` + strconv.FormatFloat(float64(totalFee*1.0)/100.0, 'f', 2, 32) + `", "refund_reason":"service cancel!"}`

	strPreSign, _err := genAlipaySignString(m)
	if _err != nil {
		log.Println("error get sign string, reason:", _err)
		return _err, nil
	}

	rsasign, err := randutils.RsaSign(keyutils.GetAlipayPrivateKeyStr(), []byte(strPreSign), crypto.SHA1)
	if err != nil {
		log.Println("error get sign string, reason:", err)
		return err, nil
	}

	sign := base64.StdEncoding.EncodeToString(rsasign)
	log.Println(strPreSign, sign)

	m["sign"] = sign

	urlstring := genAlipayUrlString(m)
	log.Println(urlstring)

	url := url.Values{}
	url.Set("app_id", fmt.Sprint(m["app_id"]))

	for k, v := range m {
		if k != "app_id" {
			url.Add(k, fmt.Sprint(v))
		}
	}

	body1 := ioutil.NopCloser(strings.NewReader(url.Encode()))
	log.Println("body1=", body1)

	req, err := http.NewRequest("POST", Alipay_Gateway, body1)
	if err != nil {
		log.Println(err)
		return err, nil
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	// body, err := fileutils.GetUrlFile(Alipay_Gateway + "?" + urlstring)
	// if err != nil {
	// 	log.Println("get url file ", urlstring, err)
	// 	return err, nil
	// }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err, nil
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	log.Println("alipay refund rsp body = ", string(data))

	outrsp := &AlipayRefundReponse{}
	outrsp.Fail = nil
	outrsp.Succ = nil

	if !strings.Contains(string(data), `"code":"10000"`) {
		refundrsp := &AlipayRefundReponseFail{}

		err = json.Unmarshal((data), refundrsp)
		if err != nil {
			log.Println("json refund failed ", err)
			return err, nil
		}

		outrsp.Fail = refundrsp
		return nil, outrsp
	}

	refundrsp := &AlipayRefundReponseSucc{}

	err = json.Unmarshal((data), refundrsp)
	if err != nil {
		log.Println("json refund failed ", err)
		return err, nil
	}

	log.Println("alipay refund rsp = ", refundrsp)
	outrsp.Succ = refundrsp

	return nil, outrsp
}
