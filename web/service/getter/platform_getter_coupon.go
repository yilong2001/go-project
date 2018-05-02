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
	"web/models/platform/couponmodel"
	//"web/service/utils"
)

func (this *ModelInfoGetter) GetPlatformCouponByCouponId(db *sql.DB, couponid int) (*couponmodel.CouponInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	couponinfo := couponmodel.NewCouponInfo()

	_, fieldAddrIfArrs := couponinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["coupon_id"] = couponid
	ruleCond["coupon_id"] = "="

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetPlatformCouponTableName(), couponinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(couponmodel.CouponInfo)
	if !ok {
		return nil, errors.New("usercouponinfo output type is wrong")
	}

	log.Println("GetCouponByUserCouponId:", out)
	return &out, nil
}
