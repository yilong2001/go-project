package sqlutils

import (
	"database/sql"
	"errors"
	"reflect"
)

func Sqls_Do_QueryRowAndScan(db *sql.DB, msqls string, selargs []interface{}, whereargs []interface{}) error {

	row := db.QueryRow(msqls, whereargs...)
	if row == nil {
		return errors.New("query is failed")
	}

	err := row.Scan(selargs...)
	if err != nil {
		return err
	}

	return nil
}

//orgStru will be changed per rows.scan
func Sqls_Do_QueryAndScan(db *sql.DB, msqls string, orgStru interface{}, selargs []interface{}, whereargs []interface{}) (*[]interface{}, error) {
	stmt, err := db.Prepare(msqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(whereargs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]interface{}, 0)

	// s := reflect.ValueOf(orgStru).Elem()
	// leng := s.NumField()
	// onerow := make([]interface{}, leng)
	// for i := 0; i < leng; i++ {
	// 	onerow[i] = s.Field(i).Addr().Interface()
	// }

	for rows.Next() {
		err = rows.Scan(selargs...)
		if err != nil {
			return nil, err
		}

		result = append(result, reflect.ValueOf(orgStru).Elem().Interface())
	}

	return &result, nil
}
