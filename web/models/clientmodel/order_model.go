package clientmodel

import (
	"web/models/ordermodel"
	"web/models/usermodel"
)

type ClientOrderInfo struct {
	OrderInfo     *ordermodel.OrderInfo
	CustomerInfo  *usermodel.UserInfo
	ProviderInfo  *usermodel.UserInfo
	OrderSubInfos *[]ordermodel.OrderSubInfo
	IsExpired     bool
	IsForbidden   bool
}
