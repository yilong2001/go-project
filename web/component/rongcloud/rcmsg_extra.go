package rongcloud

import (
	"encoding/json"
	//"fmt"
	//"github.com/rongcloud/server-sdk-go/RCServerSDK"
	//"log"
)

type RongCloudIMMsgExtra struct {
	ExtraType   string
	ContentType string
	ReferId     int
	ReferType   int
	Info        string
}

func NewRongCloudIMMsgExtraJson(extype, ctype, info string, id int, tp int) string {
	extra := &RongCloudIMMsgExtra{
		ExtraType:   extype,
		ContentType: ctype,
		ReferId:     id,
		ReferType:   tp,
		Info:        info,
	}

	js, _ := json.Marshal(extra)
	return string(js)
}
