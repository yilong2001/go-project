package serves

import (
	"log"
	//"github.com/go-martini/martini"
	//"os"
	"strconv"
	"web/models/reqparamodel"
)

func compWhereCondition(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	if svrid, err := strconv.ParseInt(params.RouterParams["ServiceId"], 10, 32); err == nil {
		idWhere["service_id"] = int(svrid)
		idRlue["service_id"] = "="
	}

	if userid, err := strconv.ParseInt(params.RouterParams["UserId"], 10, 32); err == nil {
		idWhere["user_id"] = int(userid)
		idRlue["user_id"] = "="
	}

	idWhere["del_status"] = int(0)
	idRlue["del_status"] = " = "

	log.Println("compWhereCondition : ")
	log.Println(idWhere)

	return idWhere, idRlue
}

func compWhereConditionForDefaultQuery(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	if svrid, err := strconv.ParseInt(params.RouterParams["ServiceId"], 10, 32); err == nil {
		idWhere["service_id"] = int(svrid)
		idRlue["service_id"] = "="
	}

	idWhere["del_status"] = int(0)
	idRlue["del_status"] = " = "

	log.Println("compWhereCondition : ")
	log.Println(idWhere)

	return idWhere, idRlue
}
