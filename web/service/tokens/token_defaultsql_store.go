package tokens

import (
	"database/sql"
	"log"
	"time"
	"web/dal/sqldrv"

	"web/component/cfgutils"
	"web/component/sqlutils"
	"web/models/tokenmodel"
)

type DefaultMysqlTokenStore struct {
	TableName string
}

var defaultTokenStore ApiTokenStore

func GetTokenStore() ApiTokenStore {
	return defaultTokenStore
}

func init() {
	defaultTokenStore = NewDefaultMysqlTokenStore()
}

func NewDefaultMysqlTokenStore() *DefaultMysqlTokenStore {
	return &DefaultMysqlTokenStore{
		TableName: "web_json_tokens",
	}
}

func (this *DefaultMysqlTokenStore) SaveNewToken(db1 *sql.DB, token, uuid string, uid int, expire_time string, token_type int, rctoken string) error {
	var db *sql.DB = db1
	if db == nil {
		db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db.Close()
	}

	dt := time.Now().Format("2006-01-02 15:04:05")
	log.Println(dt)

	preSql := "insert into " + this.TableName + " (token, uuid, uid, expire_time, token_type, rong_cloud_token) values(?,?,?,?,?,?)"

	stmt, err := db.Prepare(preSql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(token, uuid, uid, expire_time, token_type, rctoken); err != nil {
		return err
	}

	return nil
}

func (this *DefaultMysqlTokenStore) GetTokenInfo(db1 *sql.DB, token string, uuid string) (*tokenmodel.TokenDbModel, error) {
	var db *sql.DB = db1
	if db == nil {
		db = sqldrv.GetDb(cfgutils.GetWebApiConfig())
		defer db.Close()
	}

	tokenDbModel := &tokenmodel.TokenDbModel{}

	whereCond := map[string]interface{}{"token": token, "uuid": uuid}
	ruleCondition := map[string]string{"token": "=", "uuid": "="}

	_, fieldAddrIfArrs := tokenDbModel.GetWholeFields()

	msqls, selargs, whereargs := sqlutils.Sqls_CompSelect(this.TableName, fieldAddrIfArrs, whereCond, ruleCondition)

	log.Println(msqls)

	err := sqlutils.Sqls_Do_QueryRowAndScan(db, msqls, selargs, whereargs)
	if err != nil {
		return nil, err
	}

	return tokenDbModel, nil
}
