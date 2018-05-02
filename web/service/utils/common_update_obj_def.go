package utils

import (
	"database/sql"
)

type UpdateObjectWithIdUtil struct {
	TableName       string
	Db              *sql.DB
	Tx              *sql.Tx
	FormUnParseFlag int
	//ObjInfo   basemodel.ObjectUtilBaseIf

	CheckParamFunc         interface{}
	ExInitFunc             interface{}
	MoreProcessFunc        interface{}
	CondCompserFunc        interface{}
	CalcedUpdateFieldsFunc interface{}
	AppendMoreResultFunc   interface{}
}

func (this *UpdateObjectWithIdUtil) getDB() *sql.DB {
	return this.Db
}

func (this *UpdateObjectWithIdUtil) getTX() *sql.Tx {
	return this.Tx
}
func (this *UpdateObjectWithIdUtil) getTableName() string {
	return this.TableName
}
