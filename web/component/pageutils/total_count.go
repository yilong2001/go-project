package pageutils

import (
	"database/sql"
	"log"
	//"strings"
	"web/component/sqlutils"
	"web/models/condmodel"
)

func GetTotals(db *sql.DB, table, countObj string, whereCond map[string]interface{}, ruleCond map[string]string, exCondVal []interface{}) (int64, error) {
	count := int64(0)

	sqls, selFieldArr, whereFieldArr := sqlutils.Sqls_CompSelectCount(table, countObj, &count, whereCond, ruleCond)

	whereFieldArr = append(whereFieldArr, exCondVal...)

	sqls = sqls
	//if strings.Contains(sqls, " group by ") {
	//	sqls = strings.Replace(sqls, "count("+countObj+")", "count(*)", 1)
	//}
	log.Println("GetTotals : ", sqls)

	//TODO, if sqls is in mem, or else get from mysql

	err := sqlutils.Sqls_Do_QueryRowAndScan(db, sqls, selFieldArr, whereFieldArr)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// func GetTotalsWithComposerLinker(db *sql.DB, table, countObj string, conds *condmodel.CondComposerLinker, orderby string) (int64, error) {
// 	count := int64(0)

// 	sqls, selFieldArr, whereFieldArr := sqlutils.Sqls_CompSelectCountWithComposerLinker(table, countObj, &count, conds)

// 	sqls = sqls + " " + orderby

// 	//if strings.Contains(sqls, " group by ") {
// 	//	sqls = strings.Replace(sqls, "count("+countObj+")", "count(*)", 1)
// 	//}
// 	//TODO, if sqls is in mem, or else get from mysql
// 	log.Println(sqls)

// 	err := sqlutils.Sqls_Do_QueryRowAndScan(db, sqls, selFieldArr, whereFieldArr)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return count, nil
// }

func GetTotalsWithComposerLinker(db *sql.DB, table, countObj string, conds *condmodel.CondComposerLinker, exWhere []interface{}) (int64, error) {
	count := int64(0)

	log.Println("total conds : ", conds)

	sqls, selFieldArr, whereFieldArr := sqlutils.Sqls_CompSelectCountWithComposerLinker(table, countObj, &count, conds)

	whereFieldArr = append(exWhere, whereFieldArr...)

	//TODO, if sqls is in mem, or else get from mysql
	log.Println("get total count : ", sqls, len(selFieldArr), len(whereFieldArr))

	err := sqlutils.Sqls_Do_QueryRowAndScan(db, sqls, selFieldArr, whereFieldArr)
	if err != nil {
		return 0, err
	}

	return count, nil
}
