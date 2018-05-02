package users

import (
	"database/sql"
	"github.com/go-martini/martini"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	//"os"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"web/component/cfgutils"
	"web/component/errcode"
	"web/component/fileutils"
	"web/component/sqlutils"
	"web/models/reqparamodel"

	"web/dal/sqldrv"
	"web/service/routers"
	"web/service/uploads"
	"web/service/utils"
)

func init() {
	userPortraitRouterBuilder()
}

func userPortraitRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/user/mine/portrait", UploadPortrait)
}

var globalUserPortraitController *UsePortraitController = &UsePortraitController{
	tableName: "web_users",
}

func getGlobalUserPortraitController() *UsePortraitController {
	return globalUserPortraitController
}

type UsePortraitController struct {
	tableName string
	db        *sql.DB
}

func NewUsePortraitController() *UsePortraitController {
	ctrl := &UsePortraitController{
		tableName: "web_users",
	}
	ctrl.initDB()
	return ctrl
}

func (this *UsePortraitController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UsePortraitController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UsePortraitController) getDB() *sql.DB {
	return this.db
}

func (this *UsePortraitController) getTableName() string {
	return this.tableName
}

func (this *UsePortraitController) changePortrait(userid int, portraitPath string, r render.Render) {

	dt := time.Now().Format("2006-01-02 15:04:05")

	upCondition := map[string]interface{}{"portrait": portraitPath,
		"update_time": dt}

	whereCondition := map[string]interface{}{"user_id": userid}
	ruleCondition := map[string]string{"user_id": "="}

	msqls, args := sqlutils.Sqls_CompUpdate(this.getTableName(), upCondition, whereCondition, ruleCondition)

	log.Println(msqls, args)

	err := sqlutils.Sqls_Do_PrepareAndExec(this.getDB(), msqls, args)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Exec_Error, err.Error()))
		return
	}
}
func UploadPortrait(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUsePortraitController()
	defer ctrl.closeDB()

	log.Println("parsing form")

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
	//}

	userid, _ := strconv.ParseInt(useridStr, 10, 32)

	ctrl.changePortrait(int(userid), outFilePath, r)

	r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": outFiles})
}

func UploadPortrait1(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUsePortraitController()
	defer ctrl.closeDB()

	log.Println("parsing form")

	useridStr := headParams.RouterParams["UserId"]

	if !utils.IsFieldCorrectWithRule("user_id", useridStr) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return
	}

	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_MultiForm_Parse_Error, err.Error()))
		return
	}

	//files := req.MultipartForm.File["files"]
	//file, err := files[i].Open()
	//for i, _ := range files {
	//}

	log.Println("getting handle to file")
	file, handler, err1 := req.FormFile("upFile")
	if err1 != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_File_Upload_Handle_Open_Error, err1.Error()))
		return
	}
	defer file.Close()

	srcFileName := handler.Filename
	destDirType := "users"

	outFiles := make([]string, 0)

	log.Println("creating destination file")
	//dst, err := os.OpenFile(getUserPortraitDir(params["UserId"])+files[i].Filename, os.O_WRONLY|os.O_CREATE, 0666)
	dst, err2 := fileutils.CreateFile(fileutils.GetUpladFileDir(destDirType, useridStr), srcFileName)
	defer dst.Close()
	if err2 != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_File_Upload_Handle_Create_Error, err2.Error()))
		return
	}

	log.Println("copying the uploaded file to the destination file")
	if _, err := io.Copy(dst, file); err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_File_Upload_Handle_Copy_Error, err.Error()))
		return
	}

	outFilePath := fileutils.GetUpladFileRelativeDir(destDirType, useridStr) + srcFileName
	outFiles = append(outFiles, outFilePath)
	//}

	userid, _ := strconv.ParseInt(useridStr, 10, 32)

	ctrl.changePortrait(int(userid), outFilePath, r)

	r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": outFiles})
}
