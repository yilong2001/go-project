package getter

import (
	"database/sql"
	"log"

	"errors"

	"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	"web/component/sqlutils"

	"web/dal/sqldrv"
	//"web/models/tokenmodel"
	"web/models/platform/adminmodel"
)

func transferIf2AdminUserInfo(result *[]interface{}) *[]adminmodel.AdminUserInfo {
	users := []adminmodel.AdminUserInfo{}
	for _, rst := range *result {
		ui, ok := rst.(adminmodel.AdminUserInfo)
		if !ok {
			log.Println("userinfo type error ", rst)
			continue
		}

		users = append(users, ui)
	}

	return &users
}

func (this *ModelInfoGetter) GetAdminUserByLoginName(db *sql.DB, name string) (*adminmodel.AdminUserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := &adminmodel.AdminUserInfo{}

	_, fieldAddrIfArrs := userinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["login_name"] = name
	ruleCond["login_name"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetAdminUserTableName(), userinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(adminmodel.AdminUserInfo)
	if !ok {
		return nil, errors.New("userinfo output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetMultiAdminUsersByLoginNames(db *sql.DB, names []string) (*[]adminmodel.AdminUserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := &adminmodel.AdminUserInfo{}

	_, fieldAddrIfArrs := userinfo.GetWholeFields()

	ruleCond := map[string]string{}
	ruleCond["login_name"] = " = "

	allWhereConds := make([]map[string]interface{}, 0)
	for _, na := range names {
		whereCond := map[string]interface{}{}
		whereCond["login_name"] = na
		allWhereConds = append(allWhereConds, whereCond)
	}

	result, err := sqlutils.Sqls_GetMultiInfoWithMultiConds(db2, this.GetAdminUserTableName(), userinfo, fieldAddrIfArrs, allWhereConds, ruleCond)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	users := transferIf2AdminUserInfo(result)

	return users, nil
}

func (this *ModelInfoGetter) GetAdminUserByUserId(db *sql.DB, userid int, skip, specs []string) (*adminmodel.AdminUserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := &adminmodel.AdminUserInfo{}

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

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetAdminUserTableName(), userinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(adminmodel.AdminUserInfo)
	if !ok {
		return nil, errors.New("userinfo output type is wrong")
	}

	//log.Println("GetUserByUserId", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetMultiAdminUsersByUserIds(db *sql.DB, uids []int) (*[]adminmodel.AdminUserInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	userinfo := &adminmodel.AdminUserInfo{}

	_, fieldAddrIfArrs := userinfo.GetWholeFields()

	ruleCond := map[string]string{}
	ruleCond["user_id"] = " = "

	allWhereConds := make([]map[string]interface{}, 0)
	for _, uid := range uids {
		whereCond := map[string]interface{}{}
		whereCond["user_id"] = uid
		allWhereConds = append(allWhereConds, whereCond)
	}

	result, err := sqlutils.Sqls_GetMultiInfoWithMultiConds(db2, this.GetAdminUserTableName(), userinfo, fieldAddrIfArrs, allWhereConds, ruleCond)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	users := transferIf2AdminUserInfo(result)

	return users, nil
}
