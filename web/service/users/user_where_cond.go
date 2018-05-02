package users

import (
	"log"
	//"github.com/go-martini/martini"
	//"os"
	"strconv"
	"web/component/objutils"
	"web/models/basemodel"
	"web/models/reqparamodel"
)

func compWhereCondition(params *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string) {
	idWhere := make(map[string]interface{})
	idRlue := make(map[string]string)

	log.Println("compWhereCondition HttpReqParams", params)

	for _, idname := range basemodel.Default_All_UniqId_Names {
		log.Println(idname)

		if id, err := strconv.ParseInt(params.RouterParams[idname], 10, 32); err == nil {
			ulName, err := objutils.CamelToUnderLine(idname)
			if err == nil {
				idWhere[ulName] = int(id)
				idRlue[ulName] = "="
			} else {
				log.Println(err)
			}
		} else {
			//log.Println(err)
		}
	}

	log.Println("compWhereCondition : ", idWhere)

	return idWhere, idRlue
}
