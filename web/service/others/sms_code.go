package others

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"

	"errors"
	"net/http"
	"strconv"
	//"crypto/md5"
	//"strings"
	"time"
	//"reflect"

	"web/component/aliutils"
	"web/component/cfgutils"
	"web/component/errcode"
	"web/component/randutils"
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
	"web/service/getter"
	//"web/service/immsgs"
	//"web/service/orderups"
	"web/service/routers"
	"web/service/utils"
)

func init() {
	SmsCodeRouterBuilder()
}

func SmsCodeRouterBuilder() {
	m := routers.GetRouterHandler()

	//sms code get ?phone=
	m.Get("/sms/code", SendSmsCode)
}

func NewSmsCodeControllerObject(ctrl *SmsCodeController) *utils.ObjectWithIdUtil {
	obj := &utils.ObjectWithIdUtil{
		TableName: ctrl.getTableName(),
		Db:        ctrl.getDB(),

		ExpendInitFuncForCreate: ctrl.exInit4Create,
		CheckParamFuncForCreate: ctrl.check,

		ExpendInitFuncForUpdate: ctrl.exInit4Up,
		CheckParamFuncForUpdate: ctrl.check,
		WhereCondFuncForUpdate:  ctrl.cond4Up,
	}

	return obj
}

func NewSmsCodeController() *SmsCodeController {
	ctrl := &SmsCodeController{
		tableName: "web_user_smscodes",
	}

	ctrl.initDB()

	return ctrl
}

type SmsCodeController struct {
	tableName string
	db        *sql.DB
	genIdFlag string
}

func (this *SmsCodeController) initDB() {
	if this.db == nil {
		this.db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}
}

func (this *SmsCodeController) closeDB() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *SmsCodeController) getDB() *sql.DB {
	return this.db
}

func (this *SmsCodeController) getTableName() string {
	return this.tableName
}

func (this *SmsCodeController) exInit4Create(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	info, ok := reqInfo.(*usermodel.SmsCodeInfo)
	if !ok {
		return errors.New("req info type is not user feedback info")
	}

	info.Phone = headParams.URLParams.Get("Phone")

	if len(info.Phone) < 11 {
		log.Println("phone is wrong : ", info.Phone)
		return errors.New("phone is wrong!")
	}

	info.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	randnum := string(randutils.KRandNum(6))

	err := aliutils.AliSnsValidateCodeSend(info.Phone, randnum)
	if err != nil {
		log.Println("ali sms code send failed, ", err)
		return err
	}

	info.Code = randutils.BuildMd5PWPhoneStringV2(randnum, "")

	return nil
}

func (this *SmsCodeController) exInit4Up(reqInfo basemodel.ObjectUtilBaseIf, headParams *reqparamodel.HttpReqParams) error {
	info, ok := reqInfo.(*usermodel.SmsCodeInfo)
	if !ok {
		return errors.New("req info type is not user feedback info")
	}

	info.Phone = headParams.URLParams.Get("Phone")

	if len(info.Phone) < 11 {
		log.Println("phone is wrong : ", info.Phone)
		return errors.New("phone is wrong!")
	}

	info.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	randnum := string(randutils.KRandNum(6))

	err := aliutils.AliSnsValidateCodeSend(info.Phone, randnum)
	if err != nil {
		log.Println("ali sms code send failed, ", err)
		return err
	}

	info.Code = randutils.BuildMd5PWPhoneStringV2(randnum, "")

	return nil
}

func (this *SmsCodeController) cond4Up(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	idWhere["Phone"] = params.URLParams.Get("Phone")
	idRlue["Phone"] = " = "

	return idWhere, idRlue
}

func (this *SmsCodeController) check(headParams *reqparamodel.HttpReqParams, r render.Render) bool {
	phone := headParams.URLParams.Get("Phone")
	ok := false
	if len(phone) >= 11 {
		_, err := strconv.ParseInt(phone, 10, 64)
		if err == nil {
			ok = true
		}
	}

	if !ok {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_NotCorrect_Error, "phone is wrong"))
	}

	return ok
}

func SendSmsCode(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, ren render.Render) {
	headParams.MergeMartiniParams(params)

	ctrl := NewSmsCodeController()
	defer ctrl.closeDB()

	obj := NewSmsCodeControllerObject(ctrl)

	info := &usermodel.SmsCodeInfo{}

	if !ctrl.check(headParams, ren) {
		return
	}

	_, err := getter.GetModelInfoGetter().GetSmsCodeByPhone(ctrl.getDB(), headParams.URLParams.Get("Phone"))
	if err != nil {
		obj.Util_CreateObjectWithId(info, headParams, req, ren)
		return
	}

	obj.Util_UpdateObjectInfoWithId(info, headParams, req, ren, nil, []string{"Code", "CreateTime"})
}
