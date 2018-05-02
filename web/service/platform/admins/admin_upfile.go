package admins

import (
	//"database/sql"
	//"encoding/json"
	//"fmt"
	"github.com/go-martini/martini"
	//"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	//"errors"
	"net/http"
	//"strconv"
	//"strings"
	//"time"
	//"reflect"

	//"web/component/cfgutils"
	"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/dal/sqldrv"
	//"web/models/basemodel"
	"web/models/reqparamodel"
	//"web/models/tokenmodel"

	//"web/models/condmodel"
	//"web/models/platform/adminmodel"
	//"web/models/rendermodel"
	//"web/models/servemodel"

	//"web/service/getter"
	"web/service/routers"
	"web/service/uploads"
	//"web/service/utils"
)

func init() {
	adminUploadRouterBuilderEx()
}

func adminUploadRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Post("/admin/reviewer/upload", UploadFile)
}

func UploadFile(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	useridStr := headParams.RouterParams["UserId"]

	if !isAdminToken(headParams) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "not admin user"))
		return
	}

	err, outFilePath := uploads.UploadFile(req, "admin", useridStr)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_File_Upload_Handle_Copy_Error, err.Error()))
		return
	}

	outFiles := make([]string, 0)
	outFiles = append(outFiles, outFilePath)

	r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": outFiles})
}
