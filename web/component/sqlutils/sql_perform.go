package sqlutils

import (
	"database/sql"
	//"errors"
	//"reflect"
	"log"
)

func Sqls_Do_PrepareAndExec(db *sql.DB, msqls string, args []interface{}) error {
	stmt, err := db.Prepare(msqls)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err1 := stmt.Exec(args...)
	if err1 != nil {
		return err1
	}

	ct, err2 := result.RowsAffected()
	if err2 != nil {
		return err2
	}

	if ct == 0 {
		log.Println(msqls + ", updated row count is 0!")
	}

	return nil
}

func Sqls_Do_PrepareAndExec_Tx(tx *sql.Tx, msqls string, args []interface{}) error {
	stmt, err := tx.Prepare(msqls)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err1 := stmt.Exec(args...)
	if err1 != nil {
		return err1
	}

	ct, err2 := result.RowsAffected()
	if err2 != nil {
		return err2
	}

	if ct == 0 {
		log.Println(msqls + ", updated row count is 0!")
	}

	return nil
}
