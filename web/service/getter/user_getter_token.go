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
	"web/models/tokenmodel"
	//"web/service/utils"
)

func (this *ModelInfoGetter) GetTokenModelByTId(db *sql.DB, tokenid string) (*tokenmodel.TokenDbModel, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	tokendb := tokenmodel.NewTokenDbModel()

	_, fieldAddrIfArrs := tokendb.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["token"] = tokenid
	ruleCond["token"] = "="

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetTokenTableName(), tokendb, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(tokenmodel.TokenDbModel)
	if !ok {
		return nil, errors.New("token db output type is wrong")
	}

	log.Println("GetTokenModelByTId", out)
	return &out, nil
}
