package sqldrv

import (
	"database/sql"

	"web/component/cfgutils"
	"web/dal/mysqldrv"
	//"web/dal/sqlitedrv"
)

func GetDb(webapicfg *cfgutils.WebApiConfig) *sql.DB {
	if webapicfg.DbType == "mysql" {
		return mysqldrv.GetConn(webapicfg)
	} else if webapicfg.DbType == "sqlite" {
		return nil //return sqlitedrv.GetConn(webapicfg)
	} else {
		return nil
	}
}
