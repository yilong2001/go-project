package sqlitedrv

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"web/component/cfgutils"
)

func GetConn(webapicfg *cfgutils.WebApiConfig) *sql.DB {
	log.Println("open sqlite driver")
	db, err := sql.Open("sqlite3", webapicfg.DbFile)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
