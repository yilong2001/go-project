package orders

import (
	"database/sql"
	"github.com/go-martini/martini"
	"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"net/http"
	"strconv"
	//"strings"
	"errors"
	//"time"
	//"reflect"

	"web/component/cfgutils"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/component/errcode"
	"web/dal/sqldrv"
	"web/models/ordermodel"

	"web/models/clientmodel"
	"web/models/condmodel"

	"web/models/basemodel"
	"web/models/reqparamodel"
	//"web/models/servemodel"
	"web/models/usermodel"
	"web/service/getter"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	userOpenCommentRouterBuilderEx()
}

func userOpenCommentRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Get("/user/open/:DestId/comment", GetOpenUserOrderCommentInfo)
	m.Get("/user/open/:DestId/comment/:CommentId", GetOpenUserOrderCommentInfo)

	m.Get("/user/mine/comment", GetMineUserOrderCommentInfo)
	m.Get("/user/mine/comment/:CommentId", GetMineUserOrderCommentInfo)

	m.Get("/user/mine/order/:OrderId/comment", GetMineOrderUserOrderCommentInfo)
	m.Get("/user/mine/order/:OrderId/comment/:CommentId", GetMineOrderUserOrderCommentInfo)
}

func NewUserDefaultOrderCommentController() *UserDefaultOrderCommentController {
	ctrl := &UserDefaultOrderCommentController{
		tableName: "web_order_comments",
	}
	ctrl.initDB()
	return ctrl
}

func NewOpenUserDefaultCommentControllerObject(ctrl *UserDefaultOrderCommentController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInit4OpenGet,
		CheckParamFuncForGet: ctrl.check4OpenGet,
		//WhereCondFuncForGet:     ctrl.compWhereCond4OpenGet,
		WhereCondComposerForGet: ctrl.compComposer4OpenGet,
		AppendMoreResultFunc:    ctrl.appendUserInfo4Result,
	}

	return obj
}

func NewMineUserDefaultCommentControllerObject(ctrl *UserDefaultOrderCommentController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInit4MineGet,
		CheckParamFuncForGet: ctrl.check4MineGet,
		//WhereCondFuncForGet:     ctrl.compWhereCond4MineGet,
		WhereCondComposerForGet: ctrl.compComposer4MineGet,
		AppendMoreResultFunc:    ctrl.appendUserInfo4Result,
	}

	return obj
}

func NewMineOrderUserDefaultCommentControllerObject(ctrl *UserDefaultOrderCommentController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInit4MineOrderGet,
		CheckParamFuncForGet: ctrl.check4MineOrderGet,
		WhereCondFuncForGet:  ctrl.compWhereCond4MineOrderGet,
		AppendMoreResultFunc: ctrl.appendUserInfo4Result,
	}

	return obj
}

type UserDefaultOrderCommentController struct {
	db        *sql.DB
	tableName string
	userid    int
}

func (this *UserDefaultOrderCommentController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *UserDefaultOrderCommentController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *UserDefaultOrderCommentController) getDB() *sql.DB {
	return this.db
}

func (this *UserDefaultOrderCommentController) getTableName() string {
	return this.tableName
}

func (this *UserDefaultOrderCommentController) check4OpenGet(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["DestId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "DestId is not correct!"))
		return false
	}

	return true
}

func (this *UserDefaultOrderCommentController) exInit4OpenGet(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["DestId"], 10, 32)
	cmtid, _ := strconv.ParseInt(headParams.RouterParams["CommentId"], 10, 32)

	info, ok := reqInfo.(*ordermodel.OrderCommentInfo)
	if !ok {
		return errors.New("req info type is not service info")
	}
	//info := this.getServeInfo()

	info.ProviderId = int(userid)
	info.CommentId = int(cmtid)
	this.userid = int(userid)

	return nil
}

// func (this *UserDefaultOrderCommentController) compWhereCond4OpenGet(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
// 	idWhere := make(map[string]interface{})
// 	idRlue := make(map[string]string)

// 	if cmtid, err := strconv.ParseInt(params.RouterParams["CommentId"], 10, 32); err == nil {
// 		idWhere["comment_id"] = int(cmtid)
// 		idRlue["comment_id"] = "="
// 	}

// 	userid, _ := strconv.ParseInt(params.RouterParams["DestId"], 10, 32)
// 	idWhere["provider_id"] = int(userid)
// 	idRlue["provider_id"] = " = "

// 	return idWhere, idRlue
// }

func (this *UserDefaultOrderCommentController) compComposer4OpenGet(params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {

	composert1 := condmodel.NewCondComposerLinker("and")

	idWhere1 := make(map[string]interface{})
	idRlue1 := make(map[string]string)
	if svrid, err := strconv.ParseInt(params.RouterParams["CommentId"], 10, 32); err == nil {
		idWhere1["comment_id"] = int(svrid)
		idRlue1["comment_id"] = "="
	}

	compsubt1 := condmodel.NewCondComposerItem(idWhere1, idRlue1, " and ")
	composert1.SetItem(compsubt1)

	idWhere2 := make(map[string]interface{})
	idRlue2 := make(map[string]string)
	userid, _ := strconv.ParseInt(params.RouterParams["DestId"], 10, 32)
	idWhere2["customer_id"] = int(userid)
	idRlue2["customer_id"] = " = "
	idWhere2["provider_id"] = int(userid)
	idRlue2["provider_id"] = " = "

	compsubt2 := condmodel.NewCondComposerItem(idWhere2, idRlue2, " or ")
	composert2 := condmodel.NewCondComposerLinker("and")
	composert2.SetItem(compsubt2)

	cmtuseridstr := params.URLParams.Get("CommentUserId")
	if cmtuseridstr != "" {
		cmtuserid, err := strconv.ParseInt(cmtuseridstr, 10, 32)
		if err == nil {
			composert3 := condmodel.NewCondComposerLinker("and")
			idWhere3 := map[string]interface{}{"comment_user_id": int(cmtuserid)}
			idRlue3 := map[string]string{"comment_user_id": " = "}
			compsubt3 := condmodel.NewCondComposerItem(idWhere3, idRlue3, " or ")
			composert3.SetItem(compsubt3)

			composert2.SetNext(composert3)
		}
	}

	composert1.SetNext(composert2)

	//log.Println("compComposer4OpenGet", composert1)

	return composert1
}

func (this *UserDefaultOrderCommentController) appendUserInfo4Result(result *[]interface{}) *[]interface{} {
	if len(*result) < 1 {
		return result
	}

	//userids := []int{}
	clientInfos := []interface{}{}

	for _, cmtif := range *result {
		if cmtInfo, ok := cmtif.(ordermodel.OrderCommentInfo); ok {
			clientInfo := &clientmodel.ClientCommentInfo{}
			clientInfo.CommentInfo = &cmtInfo

			ui := &usermodel.UserInfo{}
			userInfo, err := getter.GetModelInfoGetter().GetUserByUserId(this.getDB(), cmtInfo.CustomerId, ui.GetSkipFieldsForOpenQuery(), nil)
			if err == nil {
				clientInfo.UserInfo = userInfo
			} else {
				log.Print("user query wrong: *** ", err)
			}

			clientInfos = append(clientInfos, clientInfo)
		}
	}

	return &clientInfos
}

func (this *UserDefaultOrderCommentController) check4MineGet(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	return true
}

func (this *UserDefaultOrderCommentController) exInit4MineGet(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	cmtid, _ := strconv.ParseInt(headParams.RouterParams["CommentId"], 10, 32)

	info, ok := reqInfo.(*ordermodel.OrderCommentInfo)
	if !ok {
		return errors.New("req info type is not service info")
	}
	//info := this.getServeInfo()

	info.CustomerId = int(userid)
	info.CommentId = int(cmtid)
	this.userid = int(userid)

	return nil
}

// func (this *UserDefaultOrderCommentController) compWhereCond4MineGet(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
// 	idWhere := make(map[string]interface{})
// 	idRlue := make(map[string]string)

// 	if svrid, err := strconv.ParseInt(params.RouterParams["CommentId"], 10, 32); err == nil {
// 		idWhere["comment_id"] = int(svrid)
// 		idRlue["comment_id"] = "="
// 	}

// 	userid, _ := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)
// 	idWhere["customer_id"] = int(userid)
// 	idRlue["customer_id"] = " = "

// 	return idWhere, idRlue
// }

func (this *UserDefaultOrderCommentController) compComposer4MineGet(params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {

	composert1 := condmodel.NewCondComposerLinker("and")

	idWhere1 := make(map[string]interface{})
	idRlue1 := make(map[string]string)
	if svrid, err := strconv.ParseInt(params.RouterParams["CommentId"], 10, 32); err == nil {
		idWhere1["comment_id"] = int(svrid)
		idRlue1["comment_id"] = "="
	}
	compsubt1 := condmodel.NewCondComposerItem(idWhere1, idRlue1, " and ")
	composert1.SetItem(compsubt1)

	idWhere2 := make(map[string]interface{})
	idRlue2 := make(map[string]string)
	userid, _ := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)
	idWhere2["customer_id"] = int(userid)
	idRlue2["customer_id"] = " = "
	idWhere2["provider_id"] = int(userid)
	idRlue2["provider_id"] = " = "

	compsubt2 := condmodel.NewCondComposerItem(idWhere2, idRlue2, " or ")
	composert2 := condmodel.NewCondComposerLinker("and")
	composert2.SetItem(compsubt2)

	cmtuseridstr := params.URLParams.Get("CommentUserId")
	if cmtuseridstr != "" {
		cmtuserid, err := strconv.ParseInt(cmtuseridstr, 10, 32)
		if err == nil {
			composert3 := condmodel.NewCondComposerLinker("and")
			idWhere3 := map[string]interface{}{"comment_user_id": int(cmtuserid)}
			idRlue3 := map[string]string{"comment_user_id": " = "}
			compsubt3 := condmodel.NewCondComposerItem(idWhere3, idRlue3, " or ")
			composert3.SetItem(compsubt3)

			composert2.SetNext(composert3)
		}
	}

	composert1.SetNext(composert2)

	return composert1
}

func (this *UserDefaultOrderCommentController) check4MineOrderGet(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "UserId is not correct!"))
		return false
	}

	if !utils.IsFieldCorrectWithRule("order_id", headParams.RouterParams["OrderId"]) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_UserId_Error, "OrderId is not correct!"))
		return false
	}

	return true
}

func (this *UserDefaultOrderCommentController) exInit4MineOrderGet(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
	cmtid, _ := strconv.ParseInt(headParams.RouterParams["CommentId"], 10, 32)
	oderid, _ := strconv.ParseInt(headParams.RouterParams["OrderId"], 10, 32)

	info, ok := reqInfo.(*ordermodel.OrderCommentInfo)
	if !ok {
		return errors.New("req info type is not service info")
	}
	//info := this.getServeInfo()

	info.CustomerId = int(userid)
	info.CommentId = int(cmtid)
	info.OrderId = int(oderid)
	this.userid = int(userid)

	return nil
}

func (this *UserDefaultOrderCommentController) compWhereCond4MineOrderGet(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	if cmtid, err := strconv.ParseInt(params.RouterParams["CommentId"], 10, 32); err == nil {
		idWhere["comment_id"] = int(cmtid)
		idRlue["comment_id"] = "="
	}

	/*userid, _ := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)
	idWhere["customer_id"] = int(userid)
	idRlue["customer_id"] = " = "*/

	cmtuseridstr := params.URLParams.Get("CommentUserId")
	if cmtuseridstr != "" {
		cmtuserid, err := strconv.ParseInt(cmtuseridstr, 10, 32)
		if err == nil {
			idWhere["comment_user_id"] = int(cmtuserid)
			idRlue["comment_user_id"] = " = "
		}
	}

	orderid, _ := strconv.ParseInt(params.RouterParams["OrderId"], 10, 32)
	idWhere["order_id"] = int(orderid)
	idRlue["order_id"] = " = "

	return idWhere, idRlue
}

func GetOpenUserOrderCommentInfo(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserDefaultOrderCommentController()
	defer ctrl.closeDB()
	obj := NewOpenUserDefaultCommentControllerObject(ctrl)

	info := &ordermodel.OrderCommentInfo{}

	obj.Util_GetObjectWithId_Composer(info, headParams, req, r, nil, nil)
}

func GetMineUserOrderCommentInfo(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserDefaultOrderCommentController()
	defer ctrl.closeDB()
	obj := NewMineUserDefaultCommentControllerObject(ctrl)

	info := &ordermodel.OrderCommentInfo{}

	obj.Util_GetObjectWithId_Composer(info, headParams, req, r, nil, nil)
}

func GetMineOrderUserOrderCommentInfo(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewUserDefaultOrderCommentController()
	defer ctrl.closeDB()
	obj := NewMineOrderUserDefaultCommentControllerObject(ctrl)

	info := &ordermodel.OrderCommentInfo{}

	obj.Util_GetObjectWithId(info, headParams, req, r, nil, nil)
}
