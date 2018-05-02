package clientmodel

import (
	"web/component/rongcloud"
	"web/models/usermodel"
)

type ClientRCImMsgInfo struct {
	RongCloudIMSentMsg *rongcloud.RongCloudIMSentMsg
	UserInfo           *usermodel.UserInfo
}
