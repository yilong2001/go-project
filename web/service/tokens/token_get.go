package tokens

import (
	"database/sql"
	"fmt"
	"github.com/pborman/uuid"
	"time"
	"web/component/randutils"
	"web/models/tokenmodel"
)

func GetNewToken(userid int, exp int64, tokentype int) *tokenmodel.TokenDbModel {
	token := randutils.HashUID(userid)

	//log.Println(token, len(token))

	UUID := uuid.NewUUID().String()

	if len(fmt.Sprint(exp)) > 10 {
		exp = exp / 1000
	}

	//optoken, _ := createOpenToken()
	expire_time := time.Unix(exp, 0).Format("2006-01-02 15:04:05")

	token_type := tokentype

	tokenModel := &tokenmodel.TokenDbModel{}
	tokenModel.Token = token
	tokenModel.Uuid = UUID
	tokenModel.Uid = (userid)
	tokenModel.TokenType = token_type
	tokenModel.ExpireTime = expire_time
	tokenModel.RongCloudToken = ""

	return tokenModel
}

func GetAndSaveNewToken(db *sql.DB, userid int, exp int64, tokentype int, rctoken string) (*tokenmodel.TokenDbModel, error) {
	tokenDb := GetNewToken(userid, exp, tokentype)

	err := GetTokenStore().SaveNewToken(db, tokenDb.Token, tokenDb.Uuid, userid, tokenDb.ExpireTime, tokenDb.TokenType, rctoken)

	return tokenDb, err
}
