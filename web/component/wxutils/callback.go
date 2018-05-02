package wxutils

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/chanxuehong/util"
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func WxpayNotifyCallback(w http.ResponseWriter, req *http.Request) (error, error, *pay.OrderQueryResponse) {

	reqinfo := map[string]string{}
	clt := core.NewClient(GetZhiEasyAppId(), GetMCHID(), GetZhiEasyAppKey(), nil)

	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	log.Println("读取http body失败，原因!", err)
	// 	http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err, nil, nil
	}
	log.Printf("[wxpay] [MCH] [API] http response body:\n%s\n", body)

	reqinfo, err = util.DecodeXMLToMap(bytes.NewReader(body))

	if err != nil {
		return err, nil, nil
	}

	fmt.Println("WxpayNotifyCallback outmap: ", reqinfo)

	returnCode, ok := reqinfo["return_code"]
	if !ok {
		return errors.New("return_code can not parse"), nil, nil
	}

	if returnCode != "SUCCESS" {
		return errors.New("return_code is fail"), nil, nil
	}

	// 安全考虑, 做下验证 appid 和 mch_id
	appId, ok := reqinfo["appid"]
	if ok && appId != clt.AppId() {
		return errors.New("appid is wrong"), nil, nil
	}
	mchId, ok := reqinfo["mch_id"]
	if ok && mchId != clt.MchId() {
		return errors.New("mch_id is wrong"), nil, nil
	}

	// 验证签名
	signature1, ok := reqinfo["sign"]
	if !ok {
		return errors.New("sign is wrong"), nil, nil
	}

	signature2 := core.Sign(reqinfo, clt.ApiKey(), nil)
	if signature1 != signature2 {
		return errors.New("sign is wrong"), nil, nil
	}

	resultCode, ok := reqinfo["result_code"]
	if !ok {
		return errors.New("result_code can not parse"), nil, nil
	}

	var resultErr error
	if resultCode != "SUCCESS" {
		resultErr = errors.New("result_code is fail")
	}

	totalfee := int64(0)
	cashfee := int64(0)
	if n, err := strconv.ParseInt(reqinfo["total_fee"], 10, 64); err != nil {
		err = fmt.Errorf("parse total_fee:%q to int64 failed: %s", reqinfo["total_fee"], err.Error())
		fmt.Println("wxpay query order total_fee : ", err)
	} else {
		totalfee = (n)
	}
	if n, err := strconv.ParseInt(reqinfo["cash_fee"], 10, 64); err != nil {
		err = fmt.Errorf("parse cash_fee:%q to int64 failed: %s", reqinfo["cash_fee"], err.Error())
		fmt.Println("wxpay query order cash_fee : ", err)
	} else {
		cashfee = (n)
	}

	notifyInfo := &pay.OrderQueryResponse{
		AppId: reqinfo["appid"],
		MchId: reqinfo["mch_id"],

		OpenId:         reqinfo["openid"],
		TradeType:      reqinfo["trade_type"],
		TradeState:     reqinfo["trade_state"],
		BankType:       reqinfo["bank_type"],
		TransactionId:  reqinfo["transaction_id"],
		OutTradeNo:     reqinfo["out_trade_no"],
		TimeEnd:        reqinfo["time_end"],
		TradeStateDesc: reqinfo["trade_state_desc"],
		TotalFee:       totalfee,
		CashFee:        cashfee,
	}

	return nil, resultErr, notifyInfo
}

func CompXmlStrForRspMsg(code, msg string) string {
	return "<xml><return_code>" + code + "</return_code><return_msg>" + msg + "</return_msg></xml>"
}

func WxpayCallback1(w http.ResponseWriter, r *http.Request) {
	// if resp, err = api.DecodeXMLHttpResponse(httpResp.Body); err != nil {
	// 	return
	// }
	// body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("读取http body失败，原因!", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("微信支付异步通知，HTTP Body:", string(body))
	var mr WXPayNotifyReq
	err = xml.Unmarshal(body, &mr)
	if err != nil {
		fmt.Println("解析HTTP Body格式到xml失败，原因!", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)

	reqMap["return_code"] = mr.Return_code
	reqMap["return_msg"] = mr.Return_msg
	reqMap["appid"] = mr.Appid
	reqMap["mch_id"] = mr.Mch_id
	reqMap["nonce_str"] = mr.Nonce
	reqMap["result_code"] = mr.Result_code
	reqMap["openid"] = mr.Openid
	reqMap["is_subscribe"] = mr.Is_subscribe
	reqMap["trade_type"] = mr.Trade_type
	reqMap["bank_type"] = mr.Bank_type
	reqMap["total_fee"] = mr.Total_fee
	reqMap["fee_type"] = mr.Fee_type
	reqMap["cash_fee"] = mr.Cash_fee
	reqMap["cash_fee_type"] = mr.Cash_fee_Type
	reqMap["transaction_id"] = mr.Transaction_id
	reqMap["out_trade_no"] = mr.Out_trade_no
	reqMap["attach"] = mr.Attach
	reqMap["time_end"] = mr.Time_end

	var resp WXPayNotifyResp
	//进行签名校验
	if WxpayVerifySign(reqMap, mr.Sign) {
		//这里就可以更新我们的后台数据库了，其他业务逻辑同理。
		resp.Return_code = "SUCCESS"
		resp.Return_msg = "OK"
	} else {
		resp.Return_code = "FAIL"
		resp.Return_msg = "failed to verify sign, please retry!"
	}

	//结果返回，微信要求如果成功需要返回return_code "SUCCESS"
	bytes, _err := xml.Marshal(resp)
	strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)
	if _err != nil {
		fmt.Println("xml编码失败，原因：", _err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.(http.ResponseWriter).WriteHeader(http.StatusOK)
	fmt.Fprint(w.(http.ResponseWriter), strResp)
}
