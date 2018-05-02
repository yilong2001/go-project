package tokenmodel

import (
	"fmt"
	//"log"
	"time"
)

const (
	CONST_Pw_Type_None string = "0"
	CONST_Pw_Type_PW   string = "1"
	CONST_Pw_Type_SMS  string = "2"
	CONST_Pw_Type_UID  string = "3"

	EasyOccupation_Sign_key string = "EasyOccupation"

	Const_User_Token_Str_Len int = 21

	Const_Token_Type_Public  = 1
	Const_Token_Type_Private = 2
	Const_Token_Type_Admin   = 3
)

type JWTTokenClaims struct {
	Iss   string //publisher
	Sub   string //request url
	Exp   string //expire time
	Rsa   string //
	Uid   string //user flag, if open token, be UUID; else RSA(phone)
	Jti   string //json token id
	Pwt   string // pw type : 0 none; 1 passowrd; 2 sms validation code; 3: uid
	Pwi   string // if pwt=0, be ''; if pwt=1/2;
	Stamp string
	// md5(md5(password)+(phone))
}

func NewDefaultOpenTokenClaims() *JWTTokenClaims {
	return &JWTTokenClaims{
		Iss:   "localhost:8090",
		Sub:   "localhost:8090/token/open",
		Exp:   fmt.Sprint(time.Now().Add(time.Hour * 72).Unix()),
		Rsa:   "v0.1",
		Uid:   " ",
		Jti:   " ",
		Pwt:   "0",
		Pwi:   " ",
		Stamp: "",
	}
}

/*
1, for token/open
//JWTTokenClaims
{
Iss : "atayun.com",
Sub : "atayun.com/token/open",
Exp : time,
Rsa : " ",
Uid : " ",
Jti : " ",
Pwt : "0",
Pwi : " ",
}

result: UID(-1), random UUID, openToken

2, for token/private
//JWTTokenClaims
{
Iss : "atayun.com",
Sub : "atayun.com/token/private",
Exp : time,
Rsa : key version
Uid : RSA(phone, public key),
Jti : " ",
Pwt : "1",
Pwi : md5(md5(password)+phone),
}

result: uid, random UUID, privateToken

3, for normal usrl
for example: atayun.com/job/xxx

3.1 use open token
{
Iss : "atayun.com",
Sub : "atayun.com/job/xxx",
Exp : time,
Rsa : rsa version
Uid : UUID, // RSA(UUID, public key)
Jti : token,
Pwt : "3"
Pwi : md5(md5(uid)+UUID),
}

3.2 use private token
{
Iss : "atayun.com",
Sub : "atayun.com/job/xxx",
Exp : time,
Rsa : rsa version
Uid : UUID, // RSA(UUID, public key)
Jti : token,
Pwt : "3"
Pwi : md5(md5(uid)+UUID),
}


*/
