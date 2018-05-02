package usermodel

import (
	"log"
	//"reflect"
	"testing"
	"web/component/sqlutils"
	"web/component/strutils"
)

func MyPrintf(name string, args ...interface{}) {
	log.Println(name)
	log.Println("myprintf the args len is : %d", len(args))

	ct := 0
	for _, arg := range args {
		ct = ct + 1
		switch arg.(type) {
		case int:
			log.Println("%d, %v is an int value.", ct, arg)
		case string:
			log.Println("%d, %v is a string value.", ct, arg)
		case int64:
			log.Println("%d, %v is an int64 value.", ct, arg)
		default:
			log.Println("%d, %v is an unknown type.", ct, arg)
		}
	}
}

func Test_UserModelReflect(t *testing.T) {
	userInfo := NewUserInfo()
	userInfo.UserId = 1
	userInfo.LoginName = "LoginName"
	userInfo.Sex = 4
	userInfo.Industry = "Industry"

	fieldArrs := userInfo.GetFieldsWithNotDefaultValue()

	where := make(map[string]interface{})

	nm, _ := strutils.CamelToUnderLine("UserId")
	where[nm] = 1

	ma, args := sqlutils.Sqls_CompUpdate("xxxtable", fieldArrs, where)
	log.Println("user_model_test: field interface len is : %d", len(args))

	MyPrintf(ma, args...)
}
