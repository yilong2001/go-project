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
	//"web/models/tokenmodel"
	"web/models/coursemodel"
	//"web/service/utils"
)

func transferIf2CourseCatalogFirstInfo(result *[]interface{}) *[]coursemodel.CourseCatalogFirstInfo {
	outs := []coursemodel.CourseCatalogFirstInfo{}
	for _, rst := range *result {
		ui, ok := rst.(coursemodel.CourseCatalogFirstInfo)
		if !ok {
			log.Println("userinfo type error ", rst)
			continue
		}

		outs = append(outs, ui)
	}

	return &outs
}

func transferIf2CourseCatalogSecondInfoInfo(result *[]interface{}) *[]coursemodel.CourseCatalogSecondInfo {
	outs := []coursemodel.CourseCatalogSecondInfo{}
	for _, rst := range *result {
		ui, ok := rst.(coursemodel.CourseCatalogSecondInfo)
		if !ok {
			log.Println("userinfo type error ", rst)
			continue
		}

		outs = append(outs, ui)
	}

	return &outs
}

func (this *ModelInfoGetter) GetCourseMainByCourseId(db *sql.DB, cid int) (*coursemodel.CourseMainInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	cinfo := &coursemodel.CourseMainInfo{}

	_, fieldAddrIfArrs := cinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["course_id"] = cid
	ruleCond["course_id"] = " = "
	whereCond["del_status"] = 0
	ruleCond["del_status"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetCourseMainTableName(), cinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(coursemodel.CourseMainInfo)
	if !ok {
		return nil, errors.New("course info output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetCourseCatalogFirstByCatalogFirstId(db *sql.DB, fid int) (*coursemodel.CourseCatalogFirstInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	finfo := &coursemodel.CourseCatalogFirstInfo{}

	_, fieldAddrIfArrs := finfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["catalog_first_id"] = fid
	ruleCond["catalog_first_id"] = " = "
	whereCond["del_status"] = 0
	ruleCond["del_status"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetCourseCatalogFirstTableName(), finfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(coursemodel.CourseCatalogFirstInfo)
	if !ok {
		return nil, errors.New("course info output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetCourseCatalogFirstsByCourseId(db *sql.DB, courseid int, skips []string) (*[]coursemodel.CourseCatalogFirstInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	info := &coursemodel.CourseCatalogFirstInfo{}

	_, fieldAddrIfArrs := info.GetWholeFields()
	if skips != nil && len(skips) > 0 {
		_, fieldAddrIfArrs = info.GetFieldsWithSkip(skips)
	}

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["course_id"] = courseid
	ruleCond["course_id"] = " = "
	whereCond["del_status"] = 0
	ruleCond["del_status"] = " = "

	result, err := sqlutils.Sqls_GetMultiInfo(db2, this.GetCourseCatalogFirstTableName(), info, fieldAddrIfArrs, whereCond, ruleCond, -1)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	outs := transferIf2CourseCatalogFirstInfo(result)

	return outs, nil
}

func (this *ModelInfoGetter) GetCourseCatalogSecondByCatalogSecondId(db *sql.DB, sid int) (*coursemodel.CourseCatalogSecondInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	sinfo := &coursemodel.CourseCatalogSecondInfo{}

	_, fieldAddrIfArrs := sinfo.GetWholeFields()

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["catalog_second_id"] = sid
	ruleCond["catalog_second_id"] = " = "
	whereCond["del_status"] = 0
	ruleCond["del_status"] = " = "

	result, err := sqlutils.Sqls_GetInfo(db2, this.GetCourseCatalogSecondTableName(), sinfo, fieldAddrIfArrs, whereCond, ruleCond)
	if err != nil {
		return nil, err
	}

	out, ok := result.(coursemodel.CourseCatalogSecondInfo)
	if !ok {
		return nil, errors.New("course info output type is wrong")
	}

	//log.Println("GetUserByPhone:", out)
	return &out, nil
}

func (this *ModelInfoGetter) GetCourseCatalogSecondsByFirstId(db *sql.DB, firstid int, skips []string) (*[]coursemodel.CourseCatalogSecondInfo, error) {
	var db2 *sql.DB = db
	if db == nil {
		db2 = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db2.Close()
	}

	info := &coursemodel.CourseCatalogSecondInfo{}

	_, fieldAddrIfArrs := info.GetWholeFields()
	if skips != nil && len(skips) > 0 {
		_, fieldAddrIfArrs = info.GetFieldsWithSkip(skips)
	}

	whereCond := make(map[string]interface{})
	ruleCond := make(map[string]string)
	whereCond["catalog_first_id"] = firstid
	ruleCond["catalog_first_id"] = " = "
	whereCond["del_status"] = 0
	ruleCond["del_status"] = " = "

	result, err := sqlutils.Sqls_GetMultiInfo(db2, this.GetCourseCatalogSecondTableName(), info, fieldAddrIfArrs, whereCond, ruleCond, -1)
	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New("can not get user")
	}

	outs := transferIf2CourseCatalogSecondInfoInfo(result)

	return outs, nil
}
