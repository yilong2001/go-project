package admins

import (
	"database/sql"
	//"encoding/json"
	//"fmt"
	"github.com/go-martini/martini"
	//"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"errors"
	"net/http"
	"strconv"
	//"strings"
	//"time"
	//"reflect"

	"web/component/cfgutils"
	"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/reqparamodel"
	//"web/models/tokenmodel"

	"web/models/condmodel"
	//"web/models/platform/adminmodel"
	"web/models/firmmodel"
	//"web/models/rendermodel"
	//"web/models/servemodel"

	//"web/service/getter"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	openFirmTuijianWithinAdminRouterBuilder()
}

func openFirmTuijianWithinAdminRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/admin/reviewer/firms/tuijian", UpdateTuijianFirm4AdminReview)
}

func NewFirmTuijianWithinAdminControllerObject(ctrl *FirmTuijianWithinAdminController) *utils.UpdateObjectWithIdUtil {
	obj := &utils.UpdateObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExInitFunc:      ctrl.exInfoInitUpdate,
		CheckParamFunc:  ctrl.check4Update,
		CondCompserFunc: ctrl.condCompser,
	}

	return obj
}

func NewFirmTuijianWithinAdminController() *FirmTuijianWithinAdminController {
	ctrl := &FirmTuijianWithinAdminController{
		tableName: "web_firms",
	}

	ctrl.initDB()
	return ctrl
}

type FirmTuijianWithinAdminController struct {
	tableName string
	db        *sql.DB
}

func (this *FirmTuijianWithinAdminController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *FirmTuijianWithinAdminController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *FirmTuijianWithinAdminController) getDB() *sql.DB {
	return this.db
}

func (this *FirmTuijianWithinAdminController) getTableName() string {
	return this.tableName
}

func (this *FirmTuijianWithinAdminController) exInfoInitUpdate(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) error {

	info, ok := orginfo.(*firmmodel.FirmInfo)
	if !ok {
		return errors.New("orginfo type is not AdminReviewerResultInfo")
	}

	if info.FirmId == 0 {
		return errors.New("to be tuijian userid can not be 0")
	}

	// dtinfo, ok := destinfo.(*firmmodel.FirmInfo)
	// if !ok {
	// 	return errors.New("orginfo type is not AdminReviewerResultInfo")
	// }

	//dtinfo.IsTuijian = info.IsTuijian
	//dtinfo.TuijianInfo = info.TuijianInfo
	//dtinfo.TuijianImg = info.TuijianImg

	uid, err := strconv.ParseInt(params.RouterParams["UserId"], 10, 32)
	if err != nil {
		return err
	}

	info.TuijianUid = int(uid)

	return nil
}

func (this *FirmTuijianWithinAdminController) check4Update(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	if !isAdminToken(headParams) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "not admin user"))
		return false
	}

	return true
}

func (this *FirmTuijianWithinAdminController) condCompser(orginfo basemodel.ObjectUtilBaseIf, destinfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {

	compser := condmodel.NewCondComposerLinker("and")
	root := compser

	where := map[string]interface{}{}
	rule := map[string]string{}

	info, ok := orginfo.(*firmmodel.FirmInfo)
	if !ok {
		where["firm_id"] = 0
		rule["firm_id"] = " = "
	} else {
		where["firm_id"] = info.FirmId
		rule["firm_id"] = " = "
	}

	compSub := condmodel.NewCondComposerItem(where, rule, " and ")

	compser.SetItem(compSub)

	return root
}

func UpdateTuijianFirm4AdminReview(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	info := &firmmodel.FirmInfo{}
	FirmId := headParams.URLParams.Get("FirmId")
	if firmid, err := strconv.ParseInt(FirmId, 10, 32); err == nil && FirmId != "" {
		info.FirmId = int(firmid)
	} else {
		info.FirmId = 0
	}

	ctrl := NewFirmTuijianWithinAdminController()
	defer ctrl.closeDB()

	obj := NewFirmTuijianWithinAdminControllerObject(ctrl)
	obj.Update_With_MultiInObject(info, info, headParams, req, r, nil, []string{"IsTuijian", "TuijianInfo", "TuijianImg", "TuijianUid"})
}
