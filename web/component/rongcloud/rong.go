package rongcloud

import (
	"github.com/rongcloud/server-sdk-go/RCServerSDK"
	"log"
)

func NewRCServer_json(cfg *RongCloudAppConfig) *RCServerSDK.RCServer {
	var rcServer *RCServerSDK.RCServer = nil
	if _rcServer, initError := RCServerSDK.NewRCServer(cfg.AppId, cfg.AppKey, "json"); initError != nil || _rcServer == nil {
		log.Println("初始化RCServer_json：失败! ")
	} else {
		rcServer = _rcServer
		log.Println("初始化RCServer_json：ok! ")
	}

	return rcServer
}

func UserGetToken(id, name string) *RongCloudResult {
	cfg := GetRongCloudAppConfig()
	rcServer := NewRCServer_json(cfg)

	var rcResult *RongCloudResult = nil
	if byteData, tokenError := rcServer.UserGetToken(id, name, "http://www.testPortrait.com"); tokenError != nil || len(byteData) == 0 {
		log.Println("获取 Token：测试失败！！！")
		return nil
	} else {
		log.Println("获取 Token：测试通过。returnData:", string(byteData))
		rcResult = GetRongCloudResult(string(byteData))
	}

	return rcResult
}

func Demo_GetTokens() {
	result1 := UserGetToken("13099101011", "nickname")
	log.Println(result1.Token)

	result2 := UserGetToken("13099101012", "nickname2")
	log.Println(result2.Token)
}
