package getter

// import (
// 	"database/sql"
// 	"log"
// 	//"github.com/go-martini/martini"
// 	//"github.com/gorilla/schema"
// 	//"github.com/martini-contrib/binding"
// 	//"github.com/martini-contrib/render"

// 	"errors"
// 	//"net/http"
// 	//"strconv"
// 	//"strings"
// 	//"time"
// 	"reflect"

// 	//"web/component/cfgutils"
// 	//"web/component/errcode"
// 	//"web/component/idutils"
// 	//"web/component/objutils"
// 	"web/component/sqlutils"
// 	//"web/dal/sqldrv"
// 	//"web/models/tokenmodel"
// 	//"web/models/usermodel"
// 	//"web/service/utils"
// )

// func getInfo(db *sql.DB, table string, obj interface{}, fieldArrs, whereCond map[string]interface{}, ruleCond map[string]string) (interface{}, error) {
// 	msqls, selargs, whereargs := sqlutils.Sqls_CompSelect(table, fieldArrs, whereCond, ruleCond)

// 	result, err := sqlutils.Sqls_Do_QueryAndScan(db, msqls, obj, selargs, whereargs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(*result) == 0 {
// 		return nil, errors.New("item is not exist")
// 	}

// 	return (*result)[0], nil
// }

// func getMultiInfo(db *sql.DB, table string, obj interface{}, fieldArrs, whereCond map[string]interface{}, ruleCond map[string]string) (*[]interface{}, error) {
// 	msqls, selargs, whereargs := sqlutils.Sqls_CompSelect(table, fieldArrs, whereCond, ruleCond)

// 	log.Println("pre sql is : ", msqls)

// 	result, err := sqlutils.Sqls_Do_QueryAndScan(db, msqls, obj, selargs, whereargs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(*result) == 0 {
// 		return nil, errors.New("item is not exist")
// 	}

// 	return result, nil
// }

// func getMultiInfoV2(db *sql.DB, table string, obj interface{}, fieldArrs map[string]interface{}, allWhereConds [](map[string]interface{}), ruleCond map[string]string) (*[]interface{}, error) {

// 	msqls, selargs, whereargs := sqlutils.Sqls_CompSelect(table, fieldArrs, allWhereConds[0], ruleCond)

// 	stmt, err := db.Prepare(msqls)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	result := []interface{}{}

// 	var errOut error = nil

// 	for _, whereCond := range allWhereConds {
// 		whereIf := sqlutils.Sqls_CompWhereIfs(whereCond)

// 		rows, err := stmt.Query(whereIf...)
// 		if err != nil {
// 			errOut = err
// 			log.Println(whereIf, err)
// 			continue
// 		}

// 		defer rows.Close()

// 		for rows.Next() {
// 			err = rows.Scan(selargs...)
// 			if err != nil {
// 				errOut = err
// 				log.Println(whereIf, err)
// 				continue
// 			}

// 			result = append(result, reflect.ValueOf(obj).Elem().Interface())
// 		}
// 	}

// 	if len(result) < 1 {
// 		if errOut != nil {
// 			return nil, errOut
// 		}
// 		return nil, errors.New("there is no info!")
// 	}

// 	return &result, nil
// }
