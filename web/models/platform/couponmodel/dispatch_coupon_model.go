package couponmodel

import ()

type DispatchCouponInfo struct {
	DispatchType string `schema:"DispatchType"`
	DispatchId   string `schema:"DispatchId"`
}

type DispatchCouponForUsersInfo struct {
	UserIds string `schema:"UserIds"`
}
