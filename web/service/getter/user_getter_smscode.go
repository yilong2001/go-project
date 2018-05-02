package getter

import (
	"database/sql"
	"log"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	//"github.com/martini-contrib/render"

	"errors"
	//"net/http"
	//"strconv"
	//"strings"
	//"time"
	//"reflect"

	"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/usermodel"
)

func (this *ModelInfoGetter) GetSmsCodeByPhone(db *sql.DB, phone string) (*usermodel.SmsCodeInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	info := &usermodel.SmsCodeInfo{}

	_, fieldAddrIfArrs := info.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["phone"] = phone
	ruleCond["phone"] = "="

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetSmsCodeTableName(), info, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(usermodel.SmsCodeInfo)
	if !ok {
		return nil, errors.New("smscode output type is wrong")
	}

	log.Println("GetSmsCodeByPhone", out)
	return &out, nil
}
