package aliutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/pborman/uuid"
	//"io/ioutil"
	"log"
	"net/http"
	//"strings"
	"time"
)

const (
	SmsCode_AppId     = ""
	SmsCode_AppSecret = ""
)

func AliSnsValidateCodeSend(phone string, code string) error {
	//app_key := SmsCode_AppId
	//app_secret := SmsCode_AppSecret
	//http://sms.market.alicloudapi.com/singleSendSms

	req_paras := `ParamString={"no":"` + code + `"}&RecNum=` + phone + `&SignName=xx&TemplateCode=xx`

	request_host := "http://sms.market.alicloudapi.com"
	request_uri := "/singleSendSms" + "?" + req_paras

	url := request_host + request_uri
	log.Println(url)

	request_method := "GET"

	headers := map[string]string{}

	headers["X-Ca-Key"] = SmsCode_AppId
	//headers["X-Ca-Request-Mode"] = "debug"
	headers["X-Ca-Nonce"] = uuid.NewUUID().String()
	//headers["X-Ca-Stage"] = "TEST"
	headers["X-Ca-Timestamp"] = fmt.Sprint(time.Now().Unix() * 1000)
	//headers["X-Ca-Signature-Headers"] = "X-Ca-Key,X-Ca-Nonce,X-Ca-Request-Mode,X-Ca-Stage,X-Ca-Timestamp"
	headers["X-Ca-Signature-Headers"] = "X-Ca-Key,X-Ca-Nonce,X-Ca-Timestamp"
	headers["Accept"] = "application/json;charset=utf-8"

	str_header := "X-Ca-Key:" + headers["X-Ca-Key"] + "\n"
	str_header = str_header + "X-Ca-Nonce:" + headers["X-Ca-Nonce"] + "\n"
	//str_header = str_header + "X-Ca-Request-Mode:" + headers["X-Ca-Request-Mode"] + "\n"
	//str_header = str_header + "X-Ca-Stage:" + headers["X-Ca-Stage"] + "\n"
	str_header = str_header + "X-Ca-Timestamp:" + headers["X-Ca-Timestamp"]

	str_to_sign := request_method + "\n" + headers["Accept"] + "\n\n\n\n" + str_header + "\n" + request_uri

	log.Println("str_to_sign", str_to_sign)

	mac := hmac.New(sha256.New, []byte(SmsCode_AppSecret))
	mac.Write([]byte(str_to_sign))
	out := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	log.Println("sign", out)

	headers["X-Ca-Signature"] = out

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()

	log.Println(resp)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return errors.New(resp.Status)
}

func AliSmsMsgSend(phone string, name string) error {
	//app_key := SmsCode_AppId
	//app_secret := SmsCode_AppSecret
	//http://sms.market.alicloudapi.com/singleSendSms

	req_paras := `ParamString={"name":"` + name + `"}&RecNum=` + phone + `&SignName=xx&TemplateCode=xx`

	request_host := "http://sms.market.alicloudapi.com"
	request_uri := "/singleSendSms" + "?" + req_paras

	url := request_host + request_uri
	log.Println(url)

	request_method := "GET"

	headers := map[string]string{}

	headers["X-Ca-Key"] = SmsCode_AppId
	//headers["X-Ca-Request-Mode"] = "debug"
	headers["X-Ca-Nonce"] = uuid.NewUUID().String()
	//headers["X-Ca-Stage"] = "TEST"
	headers["X-Ca-Timestamp"] = fmt.Sprint(time.Now().Unix() * 1000)
	//headers["X-Ca-Signature-Headers"] = "X-Ca-Key,X-Ca-Nonce,X-Ca-Request-Mode,X-Ca-Stage,X-Ca-Timestamp"
	headers["X-Ca-Signature-Headers"] = "X-Ca-Key,X-Ca-Nonce,X-Ca-Timestamp"
	headers["Accept"] = "application/json;charset=utf-8"

	str_header := "X-Ca-Key:" + headers["X-Ca-Key"] + "\n"
	str_header = str_header + "X-Ca-Nonce:" + headers["X-Ca-Nonce"] + "\n"
	//str_header = str_header + "X-Ca-Request-Mode:" + headers["X-Ca-Request-Mode"] + "\n"
	//str_header = str_header + "X-Ca-Stage:" + headers["X-Ca-Stage"] + "\n"
	str_header = str_header + "X-Ca-Timestamp:" + headers["X-Ca-Timestamp"]

	str_to_sign := request_method + "\n" + headers["Accept"] + "\n\n\n\n" + str_header + "\n" + request_uri

	log.Println("str_to_sign", str_to_sign)

	mac := hmac.New(sha256.New, []byte(SmsCode_AppSecret))
	mac.Write([]byte(str_to_sign))
	out := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	log.Println("sign", out)

	headers["X-Ca-Signature"] = out

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()

	log.Println(resp)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return errors.New(resp.Status)
}
