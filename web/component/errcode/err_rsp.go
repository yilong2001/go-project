package errcode

import ()

type ErrRsp struct {
	ErrCode   int
	ErrMsg    string
	ErrDetail string
}

func NewErrRsp(rsp ErrRsp) *ErrRsp {
	return &ErrRsp{
		ErrCode:   rsp.ErrCode,
		ErrMsg:    rsp.ErrMsg,
		ErrDetail: rsp.ErrDetail,
	}
}

func NewErrRsp2(rsp ErrRsp, detail string) *ErrRsp {
	return &ErrRsp{
		ErrCode:   rsp.ErrCode,
		ErrMsg:    rsp.ErrMsg,
		ErrDetail: detail,
	}
}

var (
	Err_Success = ErrRsp{
		ErrCode: 0,
		ErrMsg:  "ok",
	}

	Err_Db_Query_Error = ErrRsp{
		ErrCode: 10001001,
		ErrMsg:  "Db QueryRow Error",
	}

	Err_Db_Scan_Error = ErrRsp{
		ErrCode: 10001002,
		ErrMsg:  "Db Row Scan Error",
	}

	Err_Db_Prepare_Error = ErrRsp{
		ErrCode: 10001003,
		ErrMsg:  "Db Stmt Prepare Error",
	}

	Err_Db_Exec_Error = ErrRsp{
		ErrCode: 10001004,
		ErrMsg:  "Db Stmt Exec Error",
	}

	Err_Db_Query_Where_Error = ErrRsp{
		ErrCode: 10001005,
		ErrMsg:  "Db Query Where Condition Error",
	}

	Err_Db_Get_Total_Error = ErrRsp{
		ErrCode: 10001006,
		ErrMsg:  "Db Get Total Failed",
	}

	Err_Form_Parse_Error = ErrRsp{
		ErrCode: 20001001,
		ErrMsg:  "Req Form Param Error",
	}

	Err_Form_Para_Duplicate_Error = ErrRsp{
		ErrCode: 20001002,
		ErrMsg:  "Req Form Param Is Duplicated",
	}

	Err_Form_Para_NotCorrect_Error = ErrRsp{
		ErrCode: 20001003,
		ErrMsg:  "Req Form Param Is Not Valid",
	}

	Err_Form_Para_UserId_Error = ErrRsp{
		ErrCode: 20001004,
		ErrMsg:  "UserId in Req Param Error",
	}

	Err_Form_Para_Pw_Error = ErrRsp{
		ErrCode: 20001005,
		ErrMsg:  "password in Req Param Error",
	}

	Err_Form_Para_Old_Pw_Error = ErrRsp{
		ErrCode: 20001006,
		ErrMsg:  "old pw in Req Param Error",
	}

	Err_Form_MultiForm_Parse_Error = ErrRsp{
		ErrCode: 20001008,
		ErrMsg:  "MultiForm parse Error",
	}

	Err_Form_Para_OrderId_Error = ErrRsp{
		ErrCode: 20001009,
		ErrMsg:  "OrderId Error",
	}

	Err_Form_Para_OrderStatus_Error = ErrRsp{
		ErrCode: 20001010,
		ErrMsg:  "Order Request Status Error",
	}

	Err_Form_Para_OrderOver_Error = ErrRsp{
		ErrCode: 20001011,
		ErrMsg:  "Order Over Error",
	}

	Err_Form_Para_AdminUser_Error = ErrRsp{
		ErrCode: 20001012,
		ErrMsg:  "not admin user",
	}

	Err_Type_Transfer_Error = ErrRsp{
		ErrCode: 20011001,
		ErrMsg:  "type transfer Error",
	}

	Err_File_Upload_Handle_Open_Error = ErrRsp{
		ErrCode: 30001001,
		ErrMsg:  "file handle open error",
	}

	Err_File_Upload_Handle_Create_Error = ErrRsp{
		ErrCode: 30001002,
		ErrMsg:  "file handle create error",
	}

	Err_File_Upload_Handle_Copy_Error = ErrRsp{
		ErrCode: 30001003,
		ErrMsg:  "file handle copy error",
	}

	Err_Token_UnAuthorized_Error = ErrRsp{
		ErrCode: 50001001,
		ErrMsg:  "Completed 401 Unauthorized",
	}

	Err_Token_Create_Error = ErrRsp{
		ErrCode: 50001002,
		ErrMsg:  "create token failed",
	}

	Err_Token_SaveStore_Error = ErrRsp{
		ErrCode: 50001003,
		ErrMsg:  "token saved into store error",
	}

	Err_Token_GetPrivateToken_Error = ErrRsp{
		ErrCode: 50001004,
		ErrMsg:  "get private token error",
	}

	Err_Token_GetNormalToken_Error = ErrRsp{
		ErrCode: 50001005,
		ErrMsg:  "get normal token error",
	}

	Err_Token_NormalToken_Exp_Error = ErrRsp{
		ErrCode: 50001006,
		ErrMsg:  "normal token expire time error",
	}

	Err_Token_NormalToken_Pwi_Error = ErrRsp{
		ErrCode: 50001007,
		ErrMsg:  "normal token pwi error",
	}

	Err_Token_GetAdminToken_Error = ErrRsp{
		ErrCode: 50001008,
		ErrMsg:  "get admin token error",
	}

	Err_Token_GetToken_PW_Error = ErrRsp{
		ErrCode: 50001009,
		ErrMsg:  "get token, but passwrod error",
	}

	Err_Token_Para_Sub_Error = ErrRsp{
		ErrCode: 50011001,
		ErrMsg:  "token header sub para wrong",
	}

	Err_Token_Para_Rsa_Error = ErrRsp{
		ErrCode: 50011002,
		ErrMsg:  "token header rsa para wrong",
	}

	Err_Token_Para_Uid_Error = ErrRsp{
		ErrCode: 50011003,
		ErrMsg:  "token header uid para wrong",
	}

	Err_Token_Para_Jti_Error = ErrRsp{
		ErrCode: 50011004,
		ErrMsg:  "token header jti para wrong",
	}

	Err_Token_Para_Pwt_Error = ErrRsp{
		ErrCode: 50011005,
		ErrMsg:  "token header pwt para wrong",
	}

	Err_Token_Para_Exp_Error = ErrRsp{
		ErrCode: 50011006,
		ErrMsg:  "token header exp para wrong",
	}

	Err_Token_Para_Pwi_Error = ErrRsp{
		ErrCode: 50011007,
		ErrMsg:  "token header pwi para wrong",
	}

	Err_Coupon_Free_Dispatch_OutOfLimit_Error = ErrRsp{
		ErrCode: 60001001,
		ErrMsg:  "free coupon counts that dispatched is out of limit",
	}

	Err_Coupon_In_Platform_Not_Exist_Error = ErrRsp{
		ErrCode: 60001002,
		ErrMsg:  "coupon in platform is not exist",
	}

	Err_Wx_AccessToken_Get_Error = ErrRsp{
		ErrCode: 70001001,
		ErrMsg:  "weixin access token get error",
	}

	Err_RongCloud_Msg_Sent_Error = ErrRsp{
		ErrCode: 80001001,
		ErrMsg:  "rongcloud im msg sent error",
	}
)
