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
	//"web/service/utils"
)

func (this *ModelInfoGetter) GetCouponByUserCouponId(db *sql.DB, usercouponid int) (*usermodel.UserCouponInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	usercouponinfo := usermodel.NewUserCouponInfo()

	_, fieldAddrIfArrs := usercouponinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["user_coupon_id"] = usercouponid
	ruleCond["user_coupon_id"] = "="

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetCouponTableName(), usercouponinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(usermodel.UserCouponInfo)
	if !ok {
		return nil, errors.New("usercouponinfo output type is wrong")
	}

	log.Println("GetCouponByUserCouponId:", out)
	return &out, nil
}
