package tokens

import (
	"database/sql"
	"web/models/tokenmodel"
)

type ApiTokenStore interface {
	SaveNewToken(db *sql.DB, token string, uuid string, uid int, expire_time string, token_type int, rctoken string) error
	GetTokenInfo(db *sql.DB, token string, uuid string) (*tokenmodel.TokenDbModel, error)
}
