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
	"web/models/ordermodel"
	//"web/models/usermodel"
	//"web/service/utils"
)

func (this *ModelInfoGetter) GetOrderByOrderId(db *sql.DB, orderid int) (*ordermodel.OrderInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	orderinfo := ordermodel.NewOrderInfo()

	_, fieldAddrIfArrs := orderinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["order_id"] = orderid
	ruleCond["order_id"] = "="

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetOrderTableName(), orderinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(ordermodel.OrderInfo)
	if !ok {
		return nil, errors.New("order info output type is wrong")
	}

	log.Println("GetServiceByServiceId", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetOrderByOrderOutId(db *sql.DB, orderoutid string, providerid int) (*ordermodel.OrderInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	orderinfo := ordermodel.NewOrderInfo()

	_, fieldAddrIfArrs := orderinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["order_out_id"] = orderoutid
	ruleCond["order_out_id"] = "="
	whereCond["provider_id"] = providerid
	ruleCond["provider_id"] = "="

	log.Println("GetOrderByOrderOutId", whereCond)

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetOrderTableName(), orderinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(ordermodel.OrderInfo)
	if !ok {
		return nil, errors.New("order info output type is wrong")
	}

	log.Println("GetServiceByServiceId", out)
	return &out, nil
}
