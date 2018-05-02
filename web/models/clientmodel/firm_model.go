package clientmodel

import (
	"web/models/firmmodel"
	"web/models/usermodel"
)

type ClientFirmInfo struct {
	FirmInfo  *firmmodel.FirmInfo
	UserInfos *[]usermodel.UserInfo
}
