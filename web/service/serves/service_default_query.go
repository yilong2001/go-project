package serves

import (
	"database/sql"
	"github.com/go-martini/martini"
	"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"errors"
	"net/http"
	"strconv"
	//"strings"
	//"time"
	//"reflect"
	"web/component/filterutils"

	"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/condmodel"

	"web/models/reqparamodel"
	"web/models/servemodel"
	"web/models/usermodel"
	"web/service/routers"
	"web/service/utils"

	"web/service/getter"

	"web/models/clientmodel"
)

func init() {
	serveDefaultRouterBuilder()
}

func serveDefaultRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Get("/service", GetDefaultServeInfo)
	m.Get("/service/:ServiceId", GetDefaultServeInfo)
}

func NewServeDefaultControllerObject(ctrl *ServeDefaultController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForGet: ctrl.exInfoInit4Get,

		CheckParamFuncForGet: checkServeIdWithParamsForQuery,
		//WhereCondFuncForGet:     ctrl.compCond4OpenGet,
		WhereCondComposerForGet: ctrl.compComposer4OpenGet,

		AppendMoreResultFunc: ctrl.appendUserInfo4Result,
	}

	return obj
}

func NewServeDefaultController() *ServeDefaultController {
	ctrl := &ServeDefaultController{
		tableName: "web_services",
	}
	ctrl.initDB()
	return ctrl
}

type ServeDefaultController struct {
	tableName string
	//serveInfo *servemodel.ServeInfo
	db *sql.DB
}

func (this *ServeDefaultController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *ServeDefaultController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *ServeDefaultController) getDB() *sql.DB {
	return this.db
}

func (this *ServeDefaultController) getTableName() string {
	return this.tableName
}

func (this *ServeDefaultController) exInfoInit4Get(reqInfo basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) error {
	serveid, _ := strconv.ParseInt(params.RouterParams["ServiceId"], 10, 32)

	info, ok := reqInfo.(*servemodel.ServeInfo)
	if !ok {
		return errors.New("req info type is not service info")
	}
	//info := this.getServeInfo()
	order := ""
	sorttype := params.URLParams.Get("SortType")
	if sorttype != "" {
		sortcode, err := strconv.ParseInt(sorttype, 10, 32)
		if err == nil {
			switch sortcode {
			case 0:
				break
			case 1:
				break
			case 2:
				order = "Price-asc,ServiceId-desc"
				break
			case 3:
				order = "Price-desc,ServiceId-desc"
				break
			case 4:
				order = "AvgStar-desc,ServiceId-desc"
				break
			}
		}
	}

	if order != "" {
		log.Println("ServeDefaultController exInitForGet, order: ", order)
		params.URLParams.Set("orderby", order)
	}

	info.ServiceId = int(serveid)
	return nil
}

func (this *ServeDefaultController) compComposer4OpenGet(params *reqparamodel.HttpReqParams) *condmodel.CondComposerLinker {

	return filterutils.GetComposerOfFilters(params, "ServiceId", "ServiceType")
	// composer1 := condmodel.NewCondComposerLinker("and")

	// idWhere1 := make(map[string]interface{})
	// idRlue1 := make(map[string]string)

	// if svrid, err := strconv.ParseInt(params.RouterParams["ServiceId"], 10, 32); err == nil {
	// 	idWhere1["service_id"] = int(svrid)
	// 	idRlue1["service_id"] = "="
	// }

	// idWhere1["del_status"] = int(0)
	// idRlue1["del_status"] = " = "

	// //idWhere1["audit_status"] = int(1)
	// //idRlue1["audit_status"] = " = "

	// servicetype := params.URLParams.Get("ServiceType")
	// if servicetype != "" {
	// 	servicetypeCode, err := strconv.ParseInt(servicetype, 10, 32)
	// 	if err == nil && servicetypeCode > 0 {
	// 		idWhere1["service_type"] = servicetypeCode
	// 		idRlue1["service_type"] = "="
	// 	}
	// }

	// city := params.URLParams.Get("City")
	// if city != "" {
	// 	cityCode, err := strconv.ParseInt(city, 10, 32)
	// 	if err == nil && cityCode > 0 {
	// 		idWhere1["city"] = cityCode
	// 		idRlue1["city"] = "="
	// 	}
	// }

	// cur := composer1

	// industry := params.URLParams.Get("Industry")
	// if industry != "" {
	// 	industryCode, err := strconv.ParseInt(industry, 10, 32)
	// 	if err == nil && industryCode > 0 {
	// 		if industryCode%1000 == 0 {

	// 			compsubt1 := condmodel.NewCondComposerItem(map[string]interface{}{"industry": industryCode}, map[string]string{"industry": ">="}, " and ")
	// 			composert1 := condmodel.NewCondComposerLinker("and")
	// 			composert1.SetItem(compsubt1)
	// 			cur.SetNext(composert1)
	// 			cur = composert1

	// 			compsubt2 := condmodel.NewCondComposerItem(map[string]interface{}{"industry": industryCode + 99}, map[string]string{"industry": "<="}, " and ")
	// 			composert2 := condmodel.NewCondComposerLinker("and")
	// 			composert2.SetItem(compsubt2)
	// 			cur.SetNext(composert2)
	// 			cur = composert2
	// 		} else {
	// 			idWhere1["industry"] = industryCode
	// 			idRlue1["industry"] = "="
	// 		}
	// 	}
	// }

	// compsub1 := condmodel.NewCondComposerItem(idWhere1, idRlue1, " and ")
	// composer1.SetItem(compsub1)

	// return composer1
}

// func (this *ServeDefaultController) compCond4OpenGet(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
// 	idWhere := make(map[string]interface{})
// 	idRlue := make(map[string]string)

// 	if svrid, err := strconv.ParseInt(params.RouterParams["ServiceId"], 10, 32); err == nil {
// 		idWhere["service_id"] = int(svrid)
// 		idRlue["service_id"] = "="
// 	}

// 	servicetype := params.URLParams.Get("ServiceType")
// 	if servicetype != "" {
// 		servicetypeCode, err := strconv.ParseInt(servicetype, 10, 32)
// 		if err == nil && servicetypeCode > 0 {
// 			idWhere["service_type"] = servicetypeCode
// 			idRlue["service_type"] = "="
// 		}
// 	}

// 	industry := params.URLParams.Get("Industry")
// 	if industry != "" {
// 		industryCode, err := strconv.ParseInt(industry, 10, 32)
// 		if err == nil && industryCode > 0 {
// 			if industryCode%1000 == 0 {
// 				idWhere["industry"] = industryCode
// 				idRlue["industry"] = ">="

// 				idWhere["industry"] = (industryCode + 99)
// 				idRlue["industry"] = "<="
// 			} else {
// 				idWhere["industry"] = industryCode
// 				idRlue["industry"] = "="
// 			}
// 		}
// 	}

// 	city := params.URLParams.Get("City")
// 	if city != "" {
// 		cityCode, err := strconv.ParseInt(city, 10, 32)
// 		if err == nil && cityCode > 0 {
// 			idWhere["city"] = cityCode
// 			idRlue["city"] = "="
// 		}
// 	}

// 	idWhere["del_status"] = int(0)
// 	idRlue["del_status"] = " = "

// 	//idWhere["audit_status"] = int(1)
// 	//idRlue["audit_status"] = " = "

// 	return idWhere, idRlue
// }

func (this *ServeDefaultController) appendUserInfo4Result(result *[]interface{}) *[]interface{} {
	if len(*result) < 1 {
		return result
	}

	//userids := []int{}
	clientServeInfos := []interface{}{}

	for _, ServeInfoIf := range *result {
		if serveInfo, ok := ServeInfoIf.(servemodel.ServeInfo); ok {
			//userids = append(userids, serveInfo.UserId)
			clientServeInfo := &clientmodel.ClientServeInfo{}
			clientServeInfo.ServeInfo = &serveInfo
			ui := &usermodel.UserInfo{}
			userInfo, err := getter.GetModelInfoGetter().GetUserByUserId(this.getDB(), serveInfo.UserId, ui.GetSkipFieldsForOpenQuery(), nil)
			if err == nil {
				clientServeInfo.UserInfo = userInfo
			} else {
				log.Print("user query wrong: *** ")
				log.Print(serveInfo.UserId, " *** , ")
				log.Println(err)
			}

			clientServeInfos = append(clientServeInfos, clientServeInfo)
		}
	}

	return &clientServeInfos
}

func GetDefaultServeInfo(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewServeDefaultController()
	defer ctrl.closeDB()
	obj := NewServeDefaultControllerObject(ctrl)

	info := &servemodel.ServeInfo{}

	obj.Util_GetObjectWithId_Composer(info, headParams, req, r, info.GetSkipFieldsForOpenQuery(), nil)
}
