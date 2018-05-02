package utils

import (
	//"database/sql"
	"log"
	//"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	//"strconv"
	"strings"
	//"time"

	//"web/component/cfgutils"
	//"web/component/idutils"
	"web/component/errcode"
	"web/component/objutils"
	"web/component/pageutils"
	"web/component/sqlutils"

	//"web/dal/sqldrv"
	"web/models/basemodel"
	"web/models/condmodel"
	"web/models/reqparamodel"
	//"web/service/routers"
)

func compCondWithGroupAndOrder(groupby, orderby, uniqname string) (string, string) {
	if groupby == "" {
		return "", ""
	}

	if orderby == "" {
		orderby = "order by " + uniqname + " desc"
	}

	condpart1 := ""
	condpart2 := ""

	gbtmp := strings.Replace(groupby, "group by", "", -1)
	obtmp := strings.Replace(orderby, "order by", "", -1)

	gbtmp2 := strings.Replace(gbtmp, " ", "", -1)
	gbtmp2s := strings.Split(gbtmp2, ",")

	for id, fd := range gbtmp2s {
		if id > 0 {
			condpart1 = condpart1 + ", "
		}
		condpart1 = condpart1 + fd
	}
	condpart2 = condpart1

	obtmp2s := strings.Split(obtmp, ",")
	for _, fdt := range obtmp2s {
		condpart1 = condpart1 + ", "
		condpart2 = condpart2 + ", "

		fd := ""
		cond := "min"
		if strings.Contains(fdt, " desc") {
			cond = "max"
		}

		fdts := strings.Split(fdt, " ")
		for _, fdtsfd := range fdts {
			if fdtsfd != "" {
				fd = fdtsfd
				break
			}
		}

		if fd != "" {
			condpart1 = condpart1 + fd
			condpart2 = condpart2 + cond + "(" + fd + ")"
		}
	}

	return " ( " + condpart1 + " ) ", condpart2
}

func (this *ObjectWithIdUtil) Util_GetObjectWithId_Composer(info basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams, req *http.Request, r render.Render, skipFields []string, specFields []string) bool {
	//log.Println(params)

	if !this.getCheckParamFuncForGet().(func(*reqparamodel.HttpReqParams, render.Render) bool)(params, r) {
		//r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, "req params is not correct"))
		return false
	}

	err := this.getExpendObjInitFuncForGet().(func(basemodel.ObjectUtilBaseIf, *reqparamodel.HttpReqParams) error)(info, params)
	if err != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Parse_Error, err.Error()))
		return false
	}

	uniqIdName, _ := objutils.CamelToUnderLine(info.GetUniqIdName())

	orderBys := compOrderByForGet(info, params)
	groupbys := compGroupByForGet(info, params)
	log.Println("order by and group by : ", orderBys, groupbys)
	//expectLatest := isExpectLatestForGet(info, params)

	//info := this.getObjInfo()
	var fieldAddrIfArrs map[string]interface{}
	if skipFields == nil && specFields == nil {
		_, fieldAddrIfArrs = info.GetWholeFields()
	} else if skipFields == nil {
		_, fieldAddrIfArrs = info.GetFieldsWithSpecs(specFields)
	} else {
		_, fieldAddrIfArrs = info.GetFieldsWithSkip(skipFields)
	}

	composer := this.WhereCondComposerForGet.(func(*reqparamodel.HttpReqParams) *condmodel.CondComposerLinker)(params)

	rootcomposer := composer
	var condExAddrIfArrs []interface{}
	if groupbys != "" {
		fds1, fds2 := compCondWithGroupAndOrder(groupbys, orderBys, uniqIdName)
		log.Println("compCondWithGroupAndOrder:", fds1, fds2)

		condsql, arrs := sqlutils.Sqls_CompAnyWithComposerLinker(" select "+fds2, this.getTableName(), composer)
		condExAddrIfArrs = arrs

		composert1 := condmodel.NewCondComposerLinker("and")

		idWhere := make(map[string]interface{})
		idRlue := make(map[string]string)

		idWhere[fds1] = " ( " + condsql + " " + groupbys + " " + " ) "
		idRlue[fds1] = " in "

		compsubt := condmodel.NewCondComposerItem(idWhere, idRlue, " and ")
		composert1.SetItem(compsubt)

		composert1.SetNext(composer)

		rootcomposer = composert1
	}

	//log.Println("rootcomposer", composer)

	total, err1 := pageutils.GetTotalsWithComposerLinker(this.getDB(), this.getTableName(), uniqIdName, rootcomposer, condExAddrIfArrs)
	if err1 != nil {
		log.Println("get total wrong : ", err1)
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Get_Total_Error, err1.Error()))
		return false
	}

	pageSize := pageutils.GetRequestPageSize(this.getTableName(), params.URLParams)
	reqPageNo := pageutils.GetRequestPageNo(this.getTableName(), params.URLParams)
	pagnitor := pageutils.NewPaginator(total, pageSize, reqPageNo)

	//newcomp := composer
	//if groupbys == "" {
	//	newcomp = pagnitor.MergePaingWithComposerLinker(this.getDB(), uniqIdName, this.getTableName(), composer, orderBys)
	//}

	newcomp := rootcomposer
	newcomp = pagnitor.MergePaingWithComposerLinker(this.getDB(), uniqIdName, this.getTableName(), rootcomposer, orderBys, condExAddrIfArrs)

	msqls, selargs, whereargs := sqlutils.Sqls_CompSelectWithComposerLinker(this.getTableName(), fieldAddrIfArrs, newcomp)

	msqls = msqls + " " + " " + orderBys

	msqls = pagnitor.MergeSQLs(msqls)

	log.Println("get record sql : ", msqls)

	whereargs = append(condExAddrIfArrs, whereargs...)

	result, err := sqlutils.Sqls_Do_QueryAndScan(this.getDB(), msqls, info, selargs, whereargs)
	if err != nil {
		log.Println(err)
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Scan_Error, err.Error()))
		return false
	}

	var newresult *[]interface{} = result

	if this.getAppendMoreResultFunc() != nil {
		newresult = this.getAppendMoreResultFunc().(func(*[]interface{}) *[]interface{})(result)
	}

	//log.Println(result)

	r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": newresult, "total": total, "pageSize": pageSize, "pageNo": reqPageNo})
	return true
}
