package serves

import (
	"log"
	//"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"web/component/errcode"

	"web/models/reqparamodel"

	"web/service/utils"
)

func checkUserIdWithParamsForModify(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	log.Println("userid:" + headParams.RouterParams["UserId"])
	log.Println("jobid:" + headParams.RouterParams["JobId"])
	log.Println("serviceid:" + headParams.RouterParams["ServiceId"])

	//TODO, need to decice user id is same with session
	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	return true
}

func checkUserIdWithParamsForQuery(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	log.Println("userid:" + headParams.RouterParams["UserId"])
	log.Println("jobid:" + headParams.RouterParams["JobId"])
	log.Println("serviceid:" + headParams.RouterParams["ServiceId"])

	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	return true
}

func checkServiceIdWithParamsForModify(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	log.Println("userid:" + headParams.RouterParams["UserId"])
	log.Println("jobid:" + headParams.RouterParams["JobId"])
	log.Println("serviceid:" + headParams.RouterParams["ServiceId"])

	//TODO, need to decice user id is same with session
	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	if !utils.IsFieldCorrectWithRule("service_id", headParams.RouterParams["ServiceId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "ServiceId is not correct!"))
		return false
	}

	return true
}

func checkServeIdWithParamsForQuery(params *reqparamodel.HttpReqParams, r render.Render) bool {
	log.Println("serviceid:" + params.RouterParams["ServiceId"])

	if params.RouterParams["ServiceId"] == "" {
		return true
	}

	if !utils.IsFieldCorrectWithRule("service_id", params.RouterParams["ServiceId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "ServiceId is not correct!"))
		return false
	}

	return true
}
