package ctrlbase

import (
	"database/sql"
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/render"
	//"log"

	//"errors"
	//"net/http"
	//"strconv"
	//"strings"
	//"time"
	//"reflect"

	"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/idutils"
	//"web/component/wxutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/component/orderutils"
	//"web/component/rongcloud"
	"web/dal/sqldrv"
	//"web/models/basemodel"
	//"web/models/clientmodel"
	//"web/models/ordermodel"
	//"web/models/rendermodel"
	//"web/models/reqparamodel"
	//"web/models/usermodel"
	//"web/service/getter"
	//"web/service/immsgs"
	//"web/service/orderups"
	//"web/service/routers"
	//"web/service/utils"
)

type CtrlBaseController struct {
	TableName string
	GenIdFlag string
	Db        *sql.DB
	Tx        *sql.Tx
}

func (this *CtrlBaseController) InitDB() {
	if this.Db == nil {
		this.Db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
	}

	if this.Tx == nil {
		var err error
		this.Tx, err = this.Db.Begin()
		if err != nil {
			panic(err)
		}
	}
}

func (this *CtrlBaseController) CloseDB() {
	if this.Tx != nil {
		this.Tx.Rollback()
	}

	if this.Db != nil {
		this.Db.Close()
		this.Db = nil
	}
}

func (this *CtrlBaseController) GetDB() *sql.DB {
	return this.Db
}

func (this *CtrlBaseController) GetTX() *sql.Tx {
	return this.Tx
}

func (this *CtrlBaseController) GetTableName() string {
	return this.TableName
}

func (this *CtrlBaseController) GetGenIdFlag() string {
	return this.GenIdFlag
}

func (this *CtrlBaseController) InitNewsDB() {
	if this.Db == nil {
		this.Db = sqldrv.GetDb(cfgutils.GetNewsConfig())
	}

	if this.Tx == nil {
		var err error
		this.Tx, err = this.Db.Begin()
		if err != nil {
			panic(err)
		}
	}
}
