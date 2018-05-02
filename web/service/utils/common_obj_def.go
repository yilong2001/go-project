package utils

import (
	"database/sql"
	//"log"
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/render"
	//"net/http"
	//"strconv"
	//"time"
	//"web/component/cfgutils"
	//"web/component/idutils"
	//"web/component/errcode"
	//"web/component/objutils"
	//"web/component/pageutils"
	//"web/component/sqlutils"
	//"web/dal/sqldrv"
	//"web/models/basemodel"
	//"web/models/reqparamodel"
	//"web/service/routers"
)

type ObjectWithIdUtil struct {
	TableName       string
	Db              *sql.DB
	Tx              *sql.Tx
	FormUnParseFlag int
	//ObjInfo   basemodel.ObjectUtilBaseIf

	CheckParamFuncForCreate interface{}
	ExpendInitFuncForCreate interface{}
	MoreProcessForCreate    interface{}

	CheckParamFuncForGet    interface{}
	ExpendInitFuncForGet    interface{}
	WhereCondFuncForGet     interface{}
	WhereCondComposerForGet interface{}
	AppendMoreResultFunc    interface{}

	CheckParamFuncForUpdate    interface{}
	ExpendInitFuncForUpdate    interface{}
	WhereCondFuncForUpdate     interface{}
	WhereCondComposerForUpdate interface{}
	MoreProcessForUpdate       interface{}

	//WhereCondFuncForCreate interface{}

	GetDelStatusFuncUpdateDelStatus  interface{}
	CheckParamFuncForUpdateDelStatus interface{}
	ExpendInitFuncForUpdateDelStatus interface{}
	WhereCondFuncForUpdateDelStatus  interface{}
}

func (this *ObjectWithIdUtil) getDB() *sql.DB {
	return this.Db
}
func (this *ObjectWithIdUtil) getTX() *sql.Tx {
	return this.Tx
}

func (this *ObjectWithIdUtil) getTableName() string {
	return this.TableName
}

func (this *ObjectWithIdUtil) getExpendObjInitFuncForCreate() interface{} {
	return this.ExpendInitFuncForCreate
}

func (this *ObjectWithIdUtil) getCheckParamFuncForCreate() interface{} {
	return this.CheckParamFuncForCreate
}

func (this *ObjectWithIdUtil) getMoreProcessForCreate() interface{} {
	return this.MoreProcessForCreate
}

func (this *ObjectWithIdUtil) getExpendObjInitFuncForGet() interface{} {
	return this.ExpendInitFuncForGet
}

func (this *ObjectWithIdUtil) getCheckParamFuncForGet() interface{} {
	return this.CheckParamFuncForGet
}

func (this *ObjectWithIdUtil) getWhereCondFuncForGet() interface{} {
	return this.WhereCondFuncForGet
}

func (this *ObjectWithIdUtil) getExpendObjInitFuncForUpdate() interface{} {
	return this.ExpendInitFuncForUpdate
}

func (this *ObjectWithIdUtil) getCheckParamFuncForUpdate() interface{} {
	return this.CheckParamFuncForUpdate
}

func (this *ObjectWithIdUtil) getWhereCondFuncForUpdate() interface{} {
	return this.WhereCondFuncForUpdate
}

func (this *ObjectWithIdUtil) getMoreProcessForUpdate() interface{} {
	return this.MoreProcessForUpdate
}

func (this *ObjectWithIdUtil) getAppendMoreResultFunc() interface{} {
	return this.AppendMoreResultFunc
}

func (this *ObjectWithIdUtil) getExpendObjInitFuncForUpdateDelStatus() interface{} {
	return this.ExpendInitFuncForUpdateDelStatus
}

func (this *ObjectWithIdUtil) getCheckParamFuncForUpdateDelStatus() interface{} {
	return this.CheckParamFuncForUpdateDelStatus
}

func (this *ObjectWithIdUtil) getWhereCondFuncForUpdateDelStatus() interface{} {
	return this.WhereCondFuncForUpdateDelStatus
}

func (this *ObjectWithIdUtil) getFuncForDelStatus() interface{} {
	return this.GetDelStatusFuncUpdateDelStatus
}
