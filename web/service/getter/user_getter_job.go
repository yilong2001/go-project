package getter

import (
	"database/sql"
	"log"
	//"github.com/go-martini/martini"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	//"github.com/martini-contrib/render"

	"errors"
	//"net/http"
	//"strconv"
	//"strings"
	//"time"
	//"reflect"

	"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	"web/component/sqlutils"

	"web/dal/sqldrv"
	"web/models/jobmodel"
	//"web/service/utils"
)

func transferIf2JobInfo(result *[]interface{}) *[]jobmodel.JobInfo {
	outs := []jobmodel.JobInfo{}
	for _, rst := range *result {
		info, ok := rst.(jobmodel.JobInfo)
		if !ok {
			log.Println("userinfo type error ", rst)
			continue
		}

		outs = append(outs, info)
	}

	return &outs
}

func (this *ModelInfoGetter) GetJobByJobId(db *sql.DB, jobid int, skip, specs []string) (*jobmodel.JobInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	jobinfo := jobmodel.NewJobInfo()

	var fieldAddrIfArrs map[string]interface{}

	if skip == nil && specs == nil {
		_, fieldAddrIfArrs = jobinfo.GetWholeFields()
	} else if skip == nil {
		_, fieldAddrIfArrs = jobinfo.GetFieldsWithSpecs(specs)
	} else {
		_, fieldAddrIfArrs = jobinfo.GetFieldsWithSkip(skip)
	}

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["job_id"] = jobid
	ruleCond["job_id"] = "="
	whereCond["del_status"] = 0
	ruleCond["del_status"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetJobTableName(), jobinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(jobmodel.JobInfo)
	if !ok {
		return nil, errors.New("JobInfo output type is wrong")
	}

	log.Println("GetJobInfo", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetJobsByJobIds(db *sql.DB, jobids []int) (*[]jobmodel.JobInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	jobinfo := jobmodel.NewJobInfo()

	_, fieldAddrIfArrs := jobinfo.GetWholeFields()

	ruleCond := make(map[string]string)
	ruleCond["job_id"] = " = "
	ruleCond["del_status"] = " = "

	allWhereConds := make([]map[string]interface{}, 0)
	for _, ji := range jobids {
		whereCond := map[string]interface{}{}
		whereCond["job_id"] = ji
		whereCond["del_status"] = 0

		allWhereConds = append(allWhereConds, whereCond)
	}

	result, err := sqlutils.Sqls_GetMultiInfoWithMultiConds(db2, this.GetJobTableName(), jobinfo, fieldAddrIfArrs, allWhereConds, ruleCond)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	outs := transferIf2JobInfo(result)

	return outs, nil
}
