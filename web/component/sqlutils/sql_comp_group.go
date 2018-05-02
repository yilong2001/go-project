package sqlutils

import (
	//"reflect"
	"log"
	//"strings"
	"database/sql"
)

func Sqls_CompGroupCount(tablename string, sqls, groupby string, whereArrs map[string]interface{}, ruleArrs map[string]string) (string, []interface{}) {
	whereFieldIfs := make([]interface{}, 0)

	sqls = sqls + " from  " + tablename + " "

	ct := 0
	if whereArrs != nil && len(whereArrs) > 0 {
		ct = 0
		sqls = sqls + " where "
		for fn, fd := range whereArrs {
			rule := ruleArrs[fn]
			ct, sqls = compWhereFields(sqls, ct, fn, rule)

			whereFieldIfs = append(whereFieldIfs, fd)
		}
	}

	sqls = sqls + " group by " + groupby + " order by count desc limit 10 "

	return sqls, whereFieldIfs
}

func Sqls_GetGroupCounts(db *sql.DB, table string, orgStru interface{}, sqls, groupby string, selArr []interface{}, whereCond map[string]interface{}, ruleCond map[string]string) (*[]interface{}, error) {

	newsqls, where := Sqls_CompGroupCount(table, sqls, groupby, whereCond, ruleCond)

	log.Println(newsqls)

	return Sqls_Do_QueryAndScan(db, newsqls, orgStru, selArr, where)
}
