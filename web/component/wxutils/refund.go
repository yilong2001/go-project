package wxutils

import (
	//"bytes"
	//"encoding/xml"
	"crypto/md5"
	"errors"
	"strconv"
	//"fmt"
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	//"io/ioutil"
	//"log"
	//"net/http"
	//"strings"
	//"time"
	"web/component/keyutils"
	"web/component/randutils"
)

func WXPayRefund(orderid string, totalFee int) (error, *pay.RefundResponse) {

	tlsclient, err0 := core.NewTLSHttpClient(keyutils.GetWxPayClientCertPath(), keyutils.GetWxPayClientKeyPath())
	if err0 != nil {
		return err0, nil
	}

	clt := core.NewClient(GetZhiEasyAppId(), GetMCHID(), GetZhiEasyAppKey(), tlsclient)

	m1 := make(map[string]string, 16)
	m1["appid"] = clt.AppId()
	m1["mch_id"] = clt.MchId()
	m1["nonce_str"] = string(randutils.KRandAll(12))
	m1["out_trade_no"] = orderid

	m1["out_refund_no"] = orderid
	m1["total_fee"] = strconv.FormatInt(int64(totalFee), 10)
	m1["refund_fee"] = strconv.FormatInt(int64(totalFee), 10)
	m1["op_user_id"] = clt.MchId()

	m1["sign"] = core.Sign(m1, clt.ApiKey(), md5.New)

	m2, err := pay.Refund(clt, m1)
	if err != nil {
		return err, nil
	}

	resultCode, ok := m2["result_code"]
	if !ok {
		return errors.New("result code is failed"), nil
	}
	if resultCode != "SUCCESS" {
		return errors.New("result code is failed"), nil
	}

	resp := &pay.RefundResponse{
		AppId: m2["appid"],
		MchId: m2["mch_id"],

		TransactionId: m2["transaction_id"],
		OutTradeNo:    m2["out_trade_no"],
		OutRefundNo:   m2["out_refund_no"],
		RefundId:      m2["refund_id"],

		DeviceInfo:    m2["device_info"],
		RefundChannel: m2["refund_channel"],
		FeeType:       m2["fee_type"],
	}

	return nil, resp
}
