package sqlutils

import (
	"database/sql"
	"log"
)

func Sqls_GetCounts(db *sql.DB, table string, uniqIdName string, whereCond map[string]interface{}, ruleCond map[string]string) (int, error) {
	ct := 0

	msqls, selargs, whereargs := Sqls_CompSelectCount(table, uniqIdName, &ct, whereCond, ruleCond)

	log.Println(msqls)

	err := Sqls_Do_QueryRowAndScan(db, msqls, selargs, whereargs)
	if err != nil {
		return -1, err
	}

	return ct, nil
}
