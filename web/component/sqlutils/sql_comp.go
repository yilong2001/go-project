package sqlutils

import (
	//"reflect"
	//"log"
	"fmt"
)

func compUpdateFields(cursql string, ct int, fn string) (int, string) {
	aftersql := cursql
	if ct > 0 {
		aftersql = aftersql + ", "
	}

	aftersql = aftersql + fn + "=?"
	return (ct + 1), aftersql
}

func compWhereFields(cursql string, ct int, fn string, rule string) (int, string) {
	aftersql := cursql
	if ct > 0 {
		aftersql = aftersql + " and "
	}

	aftersql = aftersql + fn + rule + "?"
	return (ct + 1), aftersql
}

func composeWherePart(whereArrs map[string]interface{}, ruleArrs map[string]string) (string, []interface{}) {
	whereFieldIfs := []interface{}{}
	sqls := ""
	ct := 0
	if whereArrs != nil && len(whereArrs) > 0 {
		ct = 0
		sqls = sqls + " where "
		for fn, fd := range whereArrs {
			rule := ruleArrs[fn]

			if isCondValueComposed(fd) {
				if ct > 0 {
					sqls = sqls + " and "
				}
				sqls = sqls + " ( " + fn + rule + fmt.Sprint(fd) + " ) "
				ct = ct + 1
			} else if isCondSqls(fd) {
				if ct > 0 {
					sqls = sqls + " and "
				}
				sqls = sqls + " ( " + fmt.Sprint(fd) + " ) "
				ct = ct + 1
			} else {
				ct, sqls = compWhereFields(sqls, ct, fn, rule)

				whereFieldIfs = append(whereFieldIfs, fd)
			}
		}
	}

	return sqls, whereFieldIfs
}

func Sqls_CompIdWhereCond(idflag string, val int) map[string]interface{} {
	where := make(map[string]interface{})
	where[idflag] = val
	return where
}

func Sqls_CompUpdate(tablename string, fieldArrs map[string]interface{}, whereArrs map[string]interface{}, ruleArrs map[string]string) (string, []interface{}) {
	ct := 0
	fieldIfs := make([]interface{}, 0)

	sqls := "update " + tablename + " set "
	for fn, fd := range fieldArrs {
		ct, sqls = compUpdateFields(sqls, ct, fn)

		fieldIfs = append(fieldIfs, fd)
	}

	if whereArrs != nil && len(whereArrs) > 0 {
		ct = 0
		sqls = sqls + " where "
		for fn, fd := range whereArrs {
			rule := ruleArrs[fn]
			ct, sqls = compWhereFields(sqls, ct, fn, rule)

			fieldIfs = append(fieldIfs, fd)
		}
	}

	//log.Println("sql_comp: field interface len is : %d", len(fieldIfs))

	return sqls, fieldIfs
}

func compSelectFields(cursql string, ct int, fn string) (int, string) {
	aftersql := cursql
	if ct > 0 {
		aftersql = aftersql + ", "
	}

	aftersql = aftersql + fn
	return (ct + 1), aftersql
}

func Sqls_CompWhereIfs(whereArrs map[string]interface{}) []interface{} {
	whereFieldIfs := []interface{}{}
	if whereArrs != nil && len(whereArrs) > 0 {
		for _, fd := range whereArrs {
			whereFieldIfs = append(whereFieldIfs, fd)
		}
	}

	return whereFieldIfs
}

func Sqls_CompSelect(tablename string, fieldArrs map[string]interface{}, whereArrs map[string]interface{}, ruleArrs map[string]string) (string, []interface{}, []interface{}) {
	ct := 0
	selFieldIfs := make([]interface{}, 0)

	sqls := "select "
	for fn, fd := range fieldArrs {
		ct, sqls = compSelectFields(sqls, ct, fn)

		selFieldIfs = append(selFieldIfs, fd)
	}

	sqls = sqls + " from  " + tablename + " "

	exsqls, whereFieldIfs := composeWherePart(whereArrs, ruleArrs)

	sqls = sqls + exsqls

	//log.Println("sql_comp:", selFieldIfs, whereFieldIfs)

	return sqls, selFieldIfs, whereFieldIfs
}

func Sqls_CompSelectCount(tablename string, cntField string, cntVal interface{}, whereArrs map[string]interface{}, ruleArrs map[string]string) (string, []interface{}, []interface{}) {
	selFieldIfs := make([]interface{}, 0)

	sqls := "select count(" + cntField + ") as ct "
	selFieldIfs = append(selFieldIfs, cntVal)

	sqls = sqls + " from  " + tablename + " "

	exsqls, whereFieldIfs := composeWherePart(whereArrs, ruleArrs)

	sqls = sqls + exsqls

	return sqls, selFieldIfs, whereFieldIfs
}

func compInsertFields(cursql string, curval string, ct int, fn string) (int, string, string) {
	aftersql := cursql
	afterval := curval
	if ct > 0 {
		aftersql = aftersql + ", "
		afterval = curval + ", "
	}

	aftersql = aftersql + fn
	afterval = afterval + "? "
	return (ct + 1), aftersql, afterval
}

func Sqls_CompInsert(tablename string, fieldArrs map[string]interface{}) (string, []interface{}) {
	ct := 0
	fieldIfs := make([]interface{}, 0)

	sqls := "insert into " + tablename + " ( "
	vals := " ( "
	for fn, fd := range fieldArrs {
		ct, sqls, vals = compInsertFields(sqls, vals, ct, fn)

		fieldIfs = append(fieldIfs, fd)
	}

	sqls = sqls + " ) values "
	vals = vals + " ) "

	//log.Println("sql_comp: field interface len is : %d", len(fieldIfs))

	return sqls + vals, fieldIfs
}

func Sqls_CompAnyWithSelectReady(mysqls, tablename string, whereArrs map[string]interface{}, ruleArrs map[string]string) (string, []interface{}) {

	sqls := mysqls + " from  " + tablename + " "

	exsqls, whereFieldIfs := composeWherePart(whereArrs, ruleArrs)

	sqls = sqls + exsqls

	return sqls, whereFieldIfs
}
