package clientmodel

import (
	"web/models/servemodel"
	"web/models/usermodel"
)

type ClientServeInfo struct {
	ServeInfo *servemodel.ServeInfo
	UserInfo  *usermodel.UserInfo
}
