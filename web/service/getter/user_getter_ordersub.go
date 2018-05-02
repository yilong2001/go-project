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
	//"web/models/tokenmodel"
	"web/models/ordermodel"
	//"web/service/utils"
)

func transferIf2OrderSubInfo(result *[]interface{}) *[]ordermodel.OrderSubInfo {
	outs := []ordermodel.OrderSubInfo{}
	for _, rst := range *result {
		ui, ok := rst.(ordermodel.OrderSubInfo)
		if !ok {
			log.Println("userinfo type error ", rst)
			continue
		}

		outs = append(outs, ui)
	}

	return &outs
}

func (this *ModelInfoGetter) GetOrderSubByOrderSubId(db *sql.DB, subid int) (*ordermodel.OrderSubInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	ordersubinfo := &ordermodel.OrderSubInfo{}

	_, fieldAddrIfArrs := ordersubinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["order_sub_id"] = subid
	ruleCond["order_sub_id"] = " = "
	whereCond["del_status"] = 0
	ruleCond["del_status"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetOrderSubTableName(), ordersubinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(ordermodel.OrderSubInfo)
	if !ok {
		return nil, errors.New("userinfo output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetMultiOrderSubsByOrderIds(db *sql.DB, orderid int) (*[]ordermodel.OrderSubInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	ordersubinfo := &ordermodel.OrderSubInfo{}

	_, fieldAddrIfArrs := ordersubinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["order_id"] = orderid
	ruleCond["order_id"] = " = "
	whereCond["del_status"] = 0
	ruleCond["del_status"] = " = "

	result, err := sqlutils.Sqls_GetInfoEx(db2, this.GetOrderSubTableName(), ordersubinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	outs := transferIf2OrderSubInfo(result)

	return outs, nil
}
