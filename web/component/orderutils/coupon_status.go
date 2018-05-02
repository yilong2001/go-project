package orderutils

import ()

const (
	Coupon_Status_Availabe = 0
	Coupon_Status_Cancel   = 1
	Coupon_Status_Expired  = 2

	Coupon_Status_Availabe_Info = "可用"
	Coupon_Status_Cancel_Info   = "删除"
	Coupon_Status_Expired_Info  = "过期"
)

const (
	User_Coupon_Status_Ready    = 0
	User_Coupon_Status_Availabe = 1
	User_Coupon_Status_Lock     = 2
	User_Coupon_Status_Used     = 3
	User_Coupon_Status_Expired  = 4

	User_Coupon_Status_Ready_Info    = "待领取"
	User_Coupon_Status_Availabe_Info = "可用"
	User_Coupon_Status_Lock_Info     = "已锁定"
	User_Coupon_Status_Used_Info     = "已使用"
	User_Coupon_Status_Expired_Info  = "已过期"
)
