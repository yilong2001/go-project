package sqlutils

import (
	"database/sql"
	//"log"
)

// func Sqls_GetCounts(db *sql.DB, table string, uniqIdName string, fieldArrs, whereCond map[string]interface{}, ruleCond map[string]string) {

// }

func Sqls_UpdateDelStatus(db *sql.DB, table string, status int, where map[string]interface{}, rule map[string]string) error {
	upateField := make(map[string]interface{})
	upateField["del_status"] = status

	msqls, args := Sqls_CompUpdate(table, upateField, where, rule)
	err := Sqls_Do_PrepareAndExec(db, msqls, args)

	if err != nil {
		return err
	}

	return nil
}

func Sqls_UpdateDelStatus_Tx(tx *sql.Tx, table string, status int, where map[string]interface{}, rule map[string]string) error {
	upateField := make(map[string]interface{})
	upateField["del_status"] = status

	msqls, args := Sqls_CompUpdate(table, upateField, where, rule)
	err := Sqls_Do_PrepareAndExec_Tx(tx, msqls, args)

	if err != nil {
		return err
	}

	return nil
}
