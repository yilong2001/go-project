package wxutils

import ()

const (
	Prefix_APPID     = ""
	MCHID             = ""
	Prefix_KEY       = ""
	Prefix_APPSECRET = ""
)

func GetZhiEasyAppId() string {
	return Prefix_APPID
}

func GetMCHID() string {
	return MCHID
}

func GetZhiEasyAppKey() string {
	return Prefix_KEY
}

func GetZhiEasyAppSecret() string {
	return Prefix_APPSECRET
}
