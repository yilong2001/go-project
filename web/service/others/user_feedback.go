package others

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	//"log"

	"errors"
	"net/http"
	//"strconv"
	//"strings"
	"time"
	//"reflect"

	"web/component/cfgutils"
	//"web/component/errcode"
	"web/component/idutils"
	//"web/component/wxutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/component/orderutils"
	//"web/component/rongcloud"
	"web/dal/sqldrv"
	"web/models/basemodel"
	//"web/models/clientmodel"
	//"web/models/ordermodel"
	//"web/models/rendermodel"
	"web/models/reqparamodel"
	"web/models/usermodel"
	//"web/service/getter"
	//"web/service/immsgs"
	//"web/service/orderups"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userFeedbackRouterBuilder()
}

func userFeedbackRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/feedback", CreateUserFeedback)
}

func NewUserFeedbackControllerObject(ctrl *UserFeedbackController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForCreate: ctrl.exInit4Create,
		CheckParamFuncForCreate: ctrl.check,
	}

	return obj
}

func NewUserFeedbackController() *UserFeedbackController {
	ctrl := &UserFeedbackController{
		tableName: "web_user_feedbacks",
		genIdFlag: "user_feedback_id",
	}

	ctrl.initDB()

	return ctrl
}

type UserFeedbackController struct {
	tableName string
	db        *sql.DB
	genIdFlag string
}

func (this *UserFeedbackController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UserFeedbackController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UserFeedbackController) getDB() *sql.DB {
	return this.db
}

func (this *UserFeedbackController) getTableName() string {
	return this.tableName
}

func (this *UserFeedbackController) getGenIdFlag() string {
	return this.genIdFlag
}

func (this *UserFeedbackController) exInit4Create(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	info, ok := reqInfo.(*usermodel.UserFeedbackInfo)
	if !ok {
		return errors.New("req info type is not user feedback info")
	}

	if info.Info == "" {
		info.Info = "nothing!"
	}

	info.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	info.UserFeedbackId = idutils.GetId(this.getGenIdFlag())

	return nil
}

func (this *UserFeedbackController) check(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	return true
}

func CreateUserFeedback(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserFeedbackController()
	defer ctrl.closeDB()

	obj := NewUserFeedbackControllerObject(ctrl)

	info := &usermodel.UserFeedbackInfo{}

	obj.Util_CreateObjectWithId(info, headParams, req, ren)
}
