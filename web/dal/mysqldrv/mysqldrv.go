package mysqldrv

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"web/component/cfgutils"
)

func GetConn(webapicfg *cfgutils.WebApiConfig) *sql.DB {
	log.Println("open mysql driver ... ")
	url := webapicfg.DbUser + ":" + webapicfg.DbPw + "@" + webapicfg.DbUrl + "/" + webapicfg.DbName

	//log.Println(url)
	con, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	return con
}
