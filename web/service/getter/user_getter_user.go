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
	"web/models/firmmodel"
	"web/models/usermodel"
	//"web/service/utils"
)

func transferIf2UserInfo(result *[]interface{}) *[]usermodel.UserInfo {
	users := []usermodel.UserInfo{}
	for _, rst := range *result {
		ui, ok := rst.(usermodel.UserInfo)
		if !ok {
			log.Println("userinfo type error ", rst)
			continue
		}

		users = append(users, ui)
	}

	return &users
}

func (this *ModelInfoGetter) GetUserByLoginName(db *sql.DB, loginname string) (*usermodel.UserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := usermodel.NewUserInfo()

	_, fieldAddrIfArrs := userinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["login_name"] = loginname
	ruleCond["login_name"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetUserTableName(), userinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(usermodel.UserInfo)
	if !ok {
		return nil, errors.New("userinfo output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetUserByPhone(db *sql.DB, phone string) (*usermodel.UserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := usermodel.NewUserInfo()

	_, fieldAddrIfArrs := userinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["phone"] = phone
	ruleCond["phone"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetUserTableName(), userinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(usermodel.UserInfo)
	if !ok {
		return nil, errors.New("userinfo output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetMultiUsersByPhones(db *sql.DB, phones []string) (*[]usermodel.UserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := usermodel.NewUserInfo()

	_, fieldAddrIfArrs := userinfo.GetWholeFields()

	ruleCond := map[string]string{}
	ruleCond["phone"] = " = "

	allWhereConds := make([]map[string]interface{}, 0)
	for _, ph := range phones {
		whereCond := map[string]interface{}{}
		whereCond["phone"] = ph
		allWhereConds = append(allWhereConds, whereCond)
	}

	result, err := sqlutils.Sqls_GetMultiInfoWithMultiConds(db2, this.GetUserTableName(), userinfo, fieldAddrIfArrs, allWhereConds, ruleCond)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	users := transferIf2UserInfo(result)

	return users, nil
}

func (this *ModelInfoGetter) GetUserByUserId(db *sql.DB, userid int, skip, specs []string) (*usermodel.UserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := usermodel.NewUserInfo()

	var fieldAddrIfArrs map[string]interface{}

	if skip == nil && specs == nil {
		_, fieldAddrIfArrs = userinfo.GetWholeFields()
	} else if skip == nil {
		_, fieldAddrIfArrs = userinfo.GetFieldsWithSpecs(specs)
	} else {
		_, fieldAddrIfArrs = userinfo.GetFieldsWithSkip(skip)
	}

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["user_id"] = userid
	ruleCond["user_id"] = "="

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetUserTableName(), userinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(usermodel.UserInfo)
	if !ok {
		return nil, errors.New("userinfo output type is wrong")
	}

	//log.Println("GetUserByUserId", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetMultiUsersByUserIds(db *sql.DB, uids []int) (*[]usermodel.UserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := usermodel.NewUserInfo()

	_, fieldAddrIfArrs := userinfo.GetWholeFields()

	ruleCond := map[string]string{}
	ruleCond["user_id"] = " = "

	allWhereConds := make([]map[string]interface{}, 0)
	for _, uid := range uids {
		whereCond := map[string]interface{}{}
		whereCond["user_id"] = uid
		allWhereConds = append(allWhereConds, whereCond)
	}

	result, err := sqlutils.Sqls_GetMultiInfoWithMultiConds(db2, this.GetUserTableName(), userinfo, fieldAddrIfArrs, allWhereConds, ruleCond)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	users := transferIf2UserInfo(result)

	return users, nil
}

func (this *ModelInfoGetter) GetMultiUsersByFirm(db *sql.DB, firm string, limit int) (*[]usermodel.UserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	simplefirm := firmmodel.MakeCompareName(firm)

	userinfo := usermodel.NewUserInfo()

	_, fieldAddrIfArrs := userinfo.GetFieldsWithSkip(userinfo.GetSkipFieldsForOpenQuery())

	whereCond := map[string]interface{}{}
	ruleCond := map[string]string{}
	whereCond["simple_firm"] = simplefirm
	ruleCond["simple_firm"] = " = "

	result, err := sqlutils.Sqls_GetMultiInfo(db2, this.GetUserTableName(), userinfo, fieldAddrIfArrs, whereCond, ruleCond, limit)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	users := transferIf2UserInfo(result)

	return users, nil
}
