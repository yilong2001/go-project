package tokenmodel

import (
	"web/component/objutils"
)

type TokenDbModel struct {
	Token          string
	Uuid           string
	Uid            int
	TokenType      int
	ExpireTime     string
	RongCloudToken string
}

func (this *TokenDbModel) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func NewTokenDbModel() *TokenDbModel {
	return &TokenDbModel{}
}
