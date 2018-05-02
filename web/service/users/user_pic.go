package users

import (
	//"database/sql"
	"github.com/go-martini/martini"
	//"io"
	//"log"
	"net/http"
	//"strconv"
	//"time"
	//"os"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	//"web/component/cfgutils"
	"web/component/errcode"
	//"web/component/fileutils"
	//"web/component/sqlutils"
	"web/models/reqparamodel"

	//"web/dal/sqldrv"
	"web/service/routers"
	"web/service/uploads"
	"web/service/utils"
)

func init() {
	userPicUploadRouterBuilder()
}

func userPicUploadRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/user/mine/pic", UploadPicture)
}

func UploadPicture(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	useridStr := headParams.RouterParams["UserId"]

	if !utils.IsFieldCorrectWithRule("user_id", useridStr) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return
	}

	err, outFilePath := uploads.UploadFile(req, "users", useridStr)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_File_Upload_Handle_Copy_Error, err.Error()))
		return
	}

	outFiles := make([]string, 0)
	outFiles = append(outFiles, outFilePath)

	r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": outFiles})
}
