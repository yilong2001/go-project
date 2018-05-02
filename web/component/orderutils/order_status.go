package orderutils

import ()

const (
	Order_Status_Init_Build_Info      = "新订单"
	Order_Status_Init_Invite_Info     = "用户发起预约请求，等待达人接受"
	Order_Status_Customer_Cancel_Info = "用户取消订单"
	Order_Status_Provider_Deny_Info   = "达人拒绝订单"
	Order_Status_Agenda_Ready_Info    = "达人接受预约，并确定时间"
	Order_Status_Provider_Cancel_Info = "达人取消订单"
	Order_Status_Service_Ok_Info      = "在预订时间，服务完成；达人已完成对用户的反馈，等待用户确认"
	Order_Status_Customer_Assure_Info = "用户已确认"

	Order_Status_Start                  = 1
	Order_Status_Wait_Provider_Accept   = 2
	Order_Status_Wait_Customer_PrePay   = 3
	Order_Status_Wait_Arrange_Date      = 4
	Order_Status_Wait_Provider_Feedback = 5
	Order_Status_Wait_Customer_Complete = 6
	Order_Status_Wait_Comment           = 9

	Order_Status_Over              = 10
	Order_Status_Customer_Cancel   = 11
	Order_Status_Provider_Rejected = 12
	Order_Status_Provider_Cancel   = 13
	Order_Status_Pay_Closed        = 14

	Order_Status_Complete = 0xf
)

const (
	Pay_Status_Init        = 0
	Pay_Status_Wait_Notify = 1
	Pay_Status_Fail        = 2
	Pay_Status_Ok          = 3
	Pay_Status_Refund      = 4
)

const (
	Order_DateArrange_Op_UnSelect  = 0
	Order_DateArrange_Op_First     = 1
	Order_DateArrange_Op_Second    = 2
	Order_DateArrange_Op_Negoation = 3
)

const (
	Order_Type_Single_Business           = 0
	Order_Type_With_Front_Money_Business = 1
	Order_Type_Course                    = 2
)

const (
	Order_Role_Customer = 0
	Order_Role_Provider = 1
)

const (
	Order_Sub_Refer_Type_Job = 0
)

const (
	Pay_Type_Account_Balance = 0
	Pay_Type_WeiXin          = 1
	Pay_Type_ZhiFuBao        = 2
)

func GetTableNameByOrderSubReferType(refertype int) string {
	switch refertype {
	case Order_Sub_Refer_Type_Job:
		return "web_jobs"
	}
	return "web_jobs"
}

func GetSystmePayerId() int {
	return 1
}
