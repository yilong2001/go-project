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
	"web/models/usermodel"
	//"web/service/utils"
)

func transferIf2UserFavoutieInfo(result *[]interface{}) *[]usermodel.UserFavouriteInfo {
	users := []usermodel.UserFavouriteInfo{}
	for _, rst := range *result {
		ui, ok := rst.(usermodel.UserFavouriteInfo)
		if !ok {
			log.Println("userinfo type error ", rst)
			continue
		}

		users = append(users, ui)
	}

	return &users
}

func (this *ModelInfoGetter) GetUserFavouriteByDestId(db *sql.DB, userid, destid int) (*usermodel.UserFavouriteInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := &usermodel.UserFavouriteInfo{}

	_, fieldAddrIfArrs := userinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["user_id"] = userid
	ruleCond["user_id"] = " = "

	whereCond["dest_id"] = destid
	ruleCond["dest_id"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetUserFavouriteTableName(), userinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(usermodel.UserFavouriteInfo)
	if !ok {
		return nil, errors.New("userinfo output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetUserFavouriteByFavouriteId(db *sql.DB, fid int) (*usermodel.UserFavouriteInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := &usermodel.UserFavouriteInfo{}

	_, fieldAddrIfArrs := userinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["favourite_id"] = fid
	ruleCond["favourite_id"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetUserFavouriteTableName(), userinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(usermodel.UserFavouriteInfo)
	if !ok {
		return nil, errors.New("userinfo output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}
