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
	"web/models/videomodel"
	//"web/service/utils"
)

func transferIf2UserVideoPathInfo(result *[]interface{}) *[]videomodel.UserAVAbbrs {
	outs := []videomodel.UserAVAbbrs{}
	for _, rst := range *result {
		ui, ok := rst.(videomodel.UserAVAbbrs)
		if !ok {
			log.Println("userinfo type error ", rst)
			continue
		}

		outs = append(outs, ui)
	}

	return &outs
}

func (this *ModelInfoGetter) GetUserAVAddrByUserPathMd5(db *sql.DB, md5 string) (*videomodel.UserAVAbbrs, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	cinfo := &videomodel.UserAVAbbrs{}

	_, fieldAddrIfArrs := cinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["user_path_md5"] = md5
	ruleCond["user_path_md5"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetUserAVAddrTableName(), cinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(videomodel.UserAVAbbrs)
	if !ok {
		return nil, errors.New("user av addr info output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}
