package sqlutils

import (
	"fmt"
	"log"
	"strings"
	"web/models/condmodel"
)

func isCondValueComposed(fd interface{}) bool {
	str := fmt.Sprint(fd)
	if strings.Contains(str, "select ") {
		return true
	}

	return false
}

func isCondSqls(fd interface{}) bool {
	str := fmt.Sprint(fd)
	if strings.Contains(str, " > ") {
		return true
	}

	if strings.Contains(str, " < ") {
		return true
	}

	if strings.Contains(str, " = ") {
		return true
	}

	if strings.Contains(str, " >= ") {
		return true
	}

	if strings.Contains(str, " <= ") {
		return true
	}

	return false
}

func transferCondCompsoserLinker(conds *condmodel.CondComposerLinker, fieldIfs *[]interface{}) string {
	if conds == nil {
		return ""
	}

	if conds.Item == nil && conds.Next == nil {
		return ""
	}

	if conds.Item == nil && conds.Next != nil {
		return " ( " + transferCondCompsoserLinker(conds.Next, fieldIfs) + " ) "
	}

	if (len(conds.Item.Where) == 0 || len(conds.Item.Rule) == 0) && conds.Next == nil {
		return ""
	}

	if (len(conds.Item.Where) == 0 || len(conds.Item.Rule) == 0) && conds.Next != nil {
		return " ( " + transferCondCompsoserLinker(conds.Next, fieldIfs) + " ) "
	}

	cursqls := " ("
	subCt := 0
	for fn, fd := range conds.Item.Where {
		if subCt > 0 {
			cursqls = cursqls + " " + conds.Item.Comp + " "
		}

		if isCondValueComposed(fd) {
			cursqls = cursqls + fn + conds.Item.Rule[fn] + fmt.Sprint(fd)
		} else if isCondSqls(fd) {
			cursqls = cursqls + fmt.Sprint(fd)
		} else {
			cursqls = cursqls + fn + conds.Item.Rule[fn] + " ? "
			//log.Print(fd, ",")
			*fieldIfs = append(*fieldIfs, fd)
		}

		subCt = subCt + 1
	}

	cursqls = cursqls + ") "

	if conds.Next == nil {
		return cursqls
	}

	ttmp := transferCondCompsoserLinker(conds.Next, fieldIfs)
	if ttmp == " (  ) " || len(ttmp) < 8 {
		return " (" + cursqls + ") "
	} else {
		return " (" + cursqls + conds.CompNext + ttmp + ") "
	}

	//return " (" + cursqls + conds.CompNext + transferCondCompsoserLinker(conds.Next, fieldIfs) + ") "
}

func Sqls_CompUpdateWithComposerLinker(tablename string, fieldArrs map[string]interface{}, conds *condmodel.CondComposerLinker, calcedFields map[string]string) (string, []interface{}) {
	ct := 0
	fieldIfs := make([]interface{}, 0)

	sqls := "update " + tablename + " set "
	for fn, fd := range fieldArrs {
		//log.Println(fn, fd)
		if calcedFields == nil || calcedFields[fn] == "" {
			ct, sqls = compUpdateFields(sqls, ct, fn)

			fieldIfs = append(fieldIfs, fd)
		} else {
			if ct > 0 {
				sqls = sqls + ", "
			}

			if strings.Contains(calcedFields[fn], "=") {
				sqls = sqls + calcedFields[fn] + " "
			} else {
				sqls = sqls + fn + "=?" + " "
				fieldIfs = append(fieldIfs, calcedFields[fn])
			}

			ct = ct + 1
		}
	}

	ct = 0

	tmp := transferCondCompsoserLinker(conds, &fieldIfs)
	log.Println("Sqls_CompUpdateWithComposerLinker : ", tmp)
	if tmp == " (  ) " || len(tmp) < 8 {
		//
	} else {
		sqls = sqls + " where "
		sqls = sqls + tmp
	}
	//sqls = sqls + " where "
	//sqls = sqls + transferCondCompsoserLinker(conds, &fieldIfs)

	log.Println("sql_comp: field interface len is : %d", len(fieldIfs))

	return sqls, fieldIfs
}

func Sqls_CompSelectWithComposerLinker(tablename string, fieldArrs map[string]interface{}, conds *condmodel.CondComposerLinker) (string, []interface{}, []interface{}) {
	ct := 0
	selFieldIfs := make([]interface{}, 0)
	whereFieldIfs := make([]interface{}, 0)

	sqls := "select "
	for fn, fd := range fieldArrs {
		ct, sqls = compSelectFields(sqls, ct, fn)

		selFieldIfs = append(selFieldIfs, fd)
	}

	sqls = sqls + " from  " + tablename + " "

	if conds != nil && conds.Item != nil {
		ct = 0
		tmp := transferCondCompsoserLinker(conds, &whereFieldIfs)
		log.Println("Sqls_CompSelectWithComposerLinker : ", tmp)
		if tmp == " (  ) " || len(tmp) < 8 {
			//
		} else {
			sqls = sqls + " where "
			sqls = sqls + tmp
		}
		//sqls = sqls + " where "
		//sqls = sqls + transferCondCompsoserLinker(conds, &whereFieldIfs)
	}

	//log.Println("sql_comp:", selFieldIfs, whereFieldIfs)

	return sqls, selFieldIfs, whereFieldIfs
}

func Sqls_CompSelectCountWithComposerLinker(tablename string, cntField string, cntVal interface{}, conds *condmodel.CondComposerLinker) (string, []interface{}, []interface{}) {
	selFieldIfs := make([]interface{}, 0)
	whereFieldIfs := make([]interface{}, 0)

	sqls := "select count(" + cntField + ") as ct "
	selFieldIfs = append(selFieldIfs, cntVal)

	sqls = sqls + " from  " + tablename + " "

	if conds != nil && conds.Item != nil {
		tmp := transferCondCompsoserLinker(conds, &whereFieldIfs)
		log.Println("Sqls_CompSelectCountWithComposerLinker : ", tmp)
		if tmp == " (  ) " || len(tmp) < 8 {
			//
		} else {
			sqls = sqls + " where "
			sqls = sqls + tmp
		}
	}

	log.Println(sqls, len(whereFieldIfs))

	return sqls, selFieldIfs, whereFieldIfs
}

func Sqls_CompAnyWithComposerLinker(mysqls, tablename string, conds *condmodel.CondComposerLinker) (string, []interface{}) {
	whereFieldIfs := make([]interface{}, 0)

	sqls := mysqls

	sqls = sqls + " from  " + tablename + " "

	if conds != nil && conds.Item != nil {
		tmp := transferCondCompsoserLinker(conds, &whereFieldIfs)
		log.Println("Sqls_CompAnyWithComposerLinker : ", tmp)
		if tmp == " (  ) " || len(tmp) < 8 {
			//
		} else {
			sqls = sqls + " where "
			sqls = sqls + tmp
		}
	}

	log.Println(sqls, len(whereFieldIfs))

	return sqls, whereFieldIfs
}
