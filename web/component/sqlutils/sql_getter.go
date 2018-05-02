package sqlutils

import (
	"database/sql"
	"log"

	"errors"
	"fmt"

	"reflect"
)

func Sqls_GetInfo(db *sql.DB, table string, obj interface{}, fieldArrs, whereCond map[string]interface{}, ruleCond map[string]string) (interface{}, error) {
	msqls, selargs, whereargs := Sqls_CompSelect(table, fieldArrs, whereCond, ruleCond)

	result, err := Sqls_Do_QueryAndScan(db, msqls, obj, selargs, whereargs)
	if err != nil {
		return nil, err
	}

	if len(*result) == 0 {
		return nil, errors.New("item is not exist")
	}

	return (*result)[0], nil
}

func Sqls_GetInfoEx(db *sql.DB, table string, obj interface{}, fieldArrs, whereCond map[string]interface{}, ruleCond map[string]string) (*[]interface{}, error) {
	msqls, selargs, whereargs := Sqls_CompSelect(table, fieldArrs, whereCond, ruleCond)

	result, err := Sqls_Do_QueryAndScan(db, msqls, obj, selargs, whereargs)
	return result, err
}

func Sqls_GetMultiInfo(db *sql.DB, table string, obj interface{}, fieldArrs, whereCond map[string]interface{}, ruleCond map[string]string, limit int) (*[]interface{}, error) {
	msqls, selargs, whereargs := Sqls_CompSelect(table, fieldArrs, whereCond, ruleCond)

	log.Println("pre sql is : ", msqls)

	if limit > 0 {
		msqls = msqls + " limit 0, " + fmt.Sprint(limit) + " "
	}

	result, err := Sqls_Do_QueryAndScan(db, msqls, obj, selargs, whereargs)
	if err != nil {
		return nil, err
	}

	if len(*result) == 0 {
		return nil, errors.New("item is not exist")
	}

	return result, nil
}

func Sqls_GetMultiInfoWithMultiConds(db *sql.DB, table string, obj interface{}, fieldArrs map[string]interface{}, allWhereConds [](map[string]interface{}), ruleCond map[string]string) (*[]interface{}, error) {

	msqls, selargs, _ := Sqls_CompSelect(table, fieldArrs, allWhereConds[0], ruleCond)

	stmt, err := db.Prepare(msqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result := []interface{}{}

	var errOut error = nil

	for _, whereCond := range allWhereConds {
		whereIf := Sqls_CompWhereIfs(whereCond)

		rows, err := stmt.Query(whereIf...)
		if err != nil {
			errOut = err
			log.Println(whereIf, err)
			continue
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(selargs...)
			if err != nil {
				errOut = err
				log.Println(whereIf, err)
				continue
			}

			result = append(result, reflect.ValueOf(obj).Elem().Interface())
		}
	}

	if len(result) < 1 {
		if errOut != nil {
			return nil, errOut
		}
		return nil, errors.New("there is no info!")
	}

	return &result, nil
}
