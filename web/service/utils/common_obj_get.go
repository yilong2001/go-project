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
	"web/models/reqparamodel"
	//"web/service/routers"
)

func isExpectLatestForGet(info basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) bool {
	EXepectLatest := params.URLParams.Get("EXepectLatest")
	if EXepectLatest == "" || EXepectLatest == "1" {
		return true
	}

	return false
}

//groupby=SentTime,
func compGroupByForGet(info basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) string {
	groupby := params.URLParams.Get("groupby")
	if groupby == "" {
		return ""
	}

	out := ""
	fields := strings.Split(groupby, ",")
	for _, fd := range fields {
		fdn := strings.Trim(fd, " ")
		if fdn != "" {
			if out != "" {
				out = out + ", "
			}
			ctu, _ := objutils.CamelToUnderLine(fdn)
			out = out + ctu
		}
	}

	if out != "" {
		out = " group by " + out
	}

	return out
}

//order by syt: orderby=SentTime+,MessageId-
func compOrderByForGet(info basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams) string {
	orderby := params.URLParams.Get("orderby")
	if orderby == "" {
		objf, _ := objutils.CamelToUnderLine(info.GetUniqIdName())
		return " order by " + objf + " desc "
	}

	//orderby = strings.Replace(orderby, " ", "", -1)
	orderbyfs := strings.Split(orderby, ",")

	out := ""
	ct := 0
	for _, ob := range orderbyfs {
		obsub := ""
		bysub := ""
		if strings.HasSuffix(ob, "asc") {
			obsub = strings.Replace(strings.TrimRight(ob, "asc"), "-", "", -1)
			bysub = " asc "
		} else if strings.HasSuffix(ob, "desc") {
			obsub = strings.Replace(strings.TrimRight(ob, "desc"), "-", "", -1)
			bysub = " desc "
		} else {

		}

		if obsub == "" {
			continue
		}

		if ct > 0 {
			out = out + ", "
		}
		ct = ct + 1

		objf, _ := objutils.CamelToUnderLine(obsub)

		out = out + objf + bysub
	}

	if out == "" {
		objf, _ := objutils.CamelToUnderLine(info.GetUniqIdName())
		return " order by " + objf + " desc "
	}

	log.Println("order by : " + out)

	return " order by " + out
}

func (this *ObjectWithIdUtil) Util_GetObjectWithId(info basemodel.ObjectUtilBaseIf, params *reqparamodel.HttpReqParams, req *http.Request, r render.Render, skipFields []string, specFields []string) bool {
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

	orderBys := compOrderByForGet(info, params)
	groupbys := compGroupByForGet(info, params)
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

	var condExAddrIfArrs []interface{}
	uniqIdName, _ := objutils.CamelToUnderLine(info.GetUniqIdName())

	whereCond, ruleCond := this.getWhereCondFuncForGet().(func(p *reqparamodel.HttpReqParams) (map[string]interface{}, map[string]string))(params)
	if groupbys != "" {
		fds1, fds2 := compCondWithGroupAndOrder(groupbys, orderBys, uniqIdName)
		log.Println("compCondWithGroupAndOrder:", fds1, fds2)

		condsql, subwhere := sqlutils.Sqls_CompAnyWithSelectReady(" select "+fds2, this.getTableName(), whereCond, ruleCond)

		whereCond[fds1] = " ( " + condsql + " " + groupbys + " " + " ) "
		ruleCond[fds1] = " in "

		condExAddrIfArrs = subwhere
	}

	total, err1 := pageutils.GetTotals(this.getDB(), this.getTableName(), uniqIdName, whereCond, ruleCond, condExAddrIfArrs)
	if err1 != nil {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Db_Get_Total_Error, err1.Error()))
		return false
	}

	pageSize := pageutils.GetRequestPageSize(this.getTableName(), params.URLParams)
	reqPageNo := pageutils.GetRequestPageNo(this.getTableName(), params.URLParams)
	pagnitor := pageutils.NewPaginator(total, pageSize, reqPageNo)

	// if groupbys == "" {
	// 	pagnitor.MergePaingCond(this.getDB(), uniqIdName, this.getTableName(), whereCond, ruleCond, orderBys)
	// }

	pagnitor.MergePaingCond(this.getDB(), uniqIdName, this.getTableName(), whereCond, ruleCond, orderBys, condExAddrIfArrs)

	msqls, selargs, whereargs := sqlutils.Sqls_CompSelect(this.getTableName(), fieldAddrIfArrs, whereCond, ruleCond)

	whereargs = append(whereargs, condExAddrIfArrs...)

	msqls = msqls + " " + orderBys
	// if len(groupbys) > 8 && expectLatest {
	// 	sortsubmsqls := " ( select * from " + this.getTableName() + " order by " + uniqIdName + " desc ) temp "

	// 	msqls = strings.Replace(msqls, this.getTableName(), sortsubmsqls, 1)
	// }

	msqls = pagnitor.MergeSQLs(msqls)

	log.Println(msqls)

	result, err := sqlutils.Sqls_Do_QueryAndScan(this.getDB(), msqls, info, selargs, whereargs)
	if err != nil {
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
