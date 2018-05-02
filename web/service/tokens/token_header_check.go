package tokens

import (
	//"errors"
	"fmt"
	"log"
	//jwt "github.com/dgrijalva/jwt-go"
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/render"
	//"net/http"
	"encoding/hex"
	//"strconv"
	"strings"
	"time"
	"web/component/errcode"
	"web/component/keyutils"
	"web/component/randutils"
	//"web/component/pageutils"
	"web/models/reqparamodel"
	"web/models/tokenmodel"
	"web/service/getter"
)

func isOpenTokenUrl(url string) bool {
	if url == "/token/open" {
		return true
	}

	return false
}

func isPrivateTokenUrl(url string) bool {
	if url == "/token/private" {
		return true
	}

	return false
}

func isNormalUrl(url string) bool {
	if isOpenTokenUrl(url) {
		return false
	}

	if isPrivateTokenUrl(url) {
		return false
	}

	if isAdminTokenUrl(url) {
		return false
	}

	return true
}

func isAdminTokenUrl(url string) bool {
	if url == "/token/admin" {
		return true
	}

	return false
}

func isUserMineTokenUrl(url string) bool {
	if strings.HasPrefix(url, "/user/mine") {
		return true
	}

	return false
}

func baseCheckToken(para *reqparamodel.HttpReqParams) *errcode.ErrRsp {
	if !strings.HasSuffix(para.TokenParams["sub"], para.ShortUrl) {
		//return errcode.NewErrRsp2(errcode.Err_Token_Para_Sub_Error, "sub token is wrong")
	}

	err := checkTokenStamp(para)
	if err != nil {
		return err
	}

	err = baseCheckOpenToken(para)
	if err != nil {
		return err
	}

	err = baseCheckPrivateToken(para)
	if err != nil {
		return err
	}

	err = baseCheckNormalToken(para)
	if err != nil {
		return err
	}

	err = baseCheckAdminToken(para)
	if err != nil {
		return err
	}

	return nil
}

func checkTokenStamp(para *reqparamodel.HttpReqParams) *errcode.ErrRsp {
	return nil
}

func baseCheckOpenToken(para *reqparamodel.HttpReqParams) *errcode.ErrRsp {
	if !isOpenTokenUrl(para.ShortUrl) {
		return nil
	}

	if para.TokenParams["rsa"] != "v0.1" {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Rsa_Error, "rsa is wrong")
	}

	if para.TokenParams["jti"] != " " {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Jti_Error, "jti is wrong")
	}

	if para.TokenParams["uid"] != " " {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, "uid is wrong")
	}

	if para.TokenParams["pwt"] != "0" {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Pwt_Error, "pwt is wrong")
	}

	return nil
}

func baseCheckPrivateToken(para *reqparamodel.HttpReqParams) *errcode.ErrRsp {
	if !isPrivateTokenUrl(para.ShortUrl) {
		return nil
	}

	if para.TokenParams["rsa"] == "" {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Rsa_Error, "rsa is wrong")
	}

	if para.TokenParams["jti"] != " " {
		log.Println(para.TokenParams["jti"], len(para.TokenParams["jti"]))
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Jti_Error, "jti is wrong")
	}

	if para.TokenParams["uid"] == " " || para.TokenParams["uid"] == "" {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, "uid is wrong")
	}

	key := keyutils.GetRSAKeyByVer(para.TokenParams["rsa"])
	if key == nil {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Rsa_Error, "rsa ver is wrong")
	}

	data, err1 := hex.DecodeString(para.TokenParams["uid"])
	if err1 != nil {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, err1.Error())
	}

	dedata, err2 := randutils.RsaDecrypt(key.PrivateKey, data)
	//log.Println(dedata)
	if err2 != nil {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, err2.Error())
	}

	//update uid
	para.TokenParams["uid"] = string(dedata)

	if para.TokenParams["pwt"] != "1" && para.TokenParams["pwt"] != "2" {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Pwt_Error, "pwt is wrong")
	}

	return nil
}

func baseCheckNormalToken(para *reqparamodel.HttpReqParams) *errcode.ErrRsp {
	if !isNormalUrl(para.ShortUrl) {
		return nil
	}

	// if para.TokenParams["rsa"] != "v0.1" {
	// 	return errcode.NewErrRsp2(errcode.Err_Token_Para_Rsa_Error, "rsa is wrong")
	// }

	if len(para.TokenParams["jti"]) != 32 {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Jti_Error, "jti len is wrong")
	}

	tokenDb, err := getter.GetModelInfoGetter().GetTokenModelByTId(nil, para.TokenParams["jti"])
	if err != nil {
		return errcode.NewErrRsp2(errcode.Err_Token_GetNormalToken_Error, err.Error())
	}

	dt, err := time.Parse("2006-01-02 15:04:05", tokenDb.ExpireTime)
	if err != nil || dt.Unix() < time.Now().Unix() {
		log.Println("expire time parse err ", err)
		return errcode.NewErrRsp2(errcode.Err_Token_NormalToken_Exp_Error, "expire time is out")
	}

	destPwi := randutils.BuildRawMd5String(fmt.Sprint(tokenDb.Uid), tokenDb.Uuid)
	if destPwi != para.TokenParams["pwi"] {
		log.Println(tokenDb.Uid, tokenDb.Uuid, destPwi, para.TokenParams["pwi"])
		return errcode.NewErrRsp2(errcode.Err_Token_NormalToken_Pwi_Error, "("+para.TokenParams["pwi"]+"): pwi is wrong")
	}

	// key := keyutils.GetRSAKeyByVer(para.TokenParams["rsa"])
	// if key == nil {
	// 	return errcode.NewErrRsp2(errcode.Err_Token_Para_Rsa_Error, "rsa ver is wrong")
	// }

	// data, err1 := hex.DecodeString(para.TokenParams["uid"])
	// if err1 != nil {
	// 	return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, err1.Error())
	// }

	// dedata, err2 := randutils.RsaDecrypt(key.PrivateKey, data)
	// if err2 != nil {
	// 	return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, err2.Error())
	// }

	// if string(dedata) != tokenDb.Uid {
	// 	return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, "uid is not equal")
	// }

	//save token db uid
	if tokenDb.Uid > 10000 && tokenDb.TokenType == tokenmodel.Const_Token_Type_Private {
		para.TokenParams["UserId"] = fmt.Sprint(tokenDb.Uid)
		para.TokenParams["TokenType"] = fmt.Sprint(tokenmodel.Const_Token_Type_Private)
	}

	if tokenDb.Uid > 10000 && tokenDb.TokenType == tokenmodel.Const_Token_Type_Private && strings.HasPrefix(para.ShortUrl, "/admin/") {
		return errcode.NewErrRsp2(errcode.Err_Token_NormalToken_Pwi_Error, " router is wrong!")
	}

	if tokenDb.Uid <= 10000 && tokenDb.TokenType == tokenmodel.Const_Token_Type_Admin {
		para.TokenParams["UserId"] = fmt.Sprint(tokenDb.Uid)
		para.TokenParams["TokenType"] = fmt.Sprint(tokenmodel.Const_Token_Type_Admin)
	}

	if tokenDb.Uid <= 10000 && tokenDb.TokenType == tokenmodel.Const_Token_Type_Public {
		para.TokenParams["UserId"] = fmt.Sprint(tokenDb.Uid)
		para.TokenParams["TokenType"] = fmt.Sprint(tokenmodel.Const_Token_Type_Public)
	}

	if tokenDb.Uid <= 10000 && tokenDb.TokenType == tokenmodel.Const_Token_Type_Admin && !strings.HasPrefix(para.ShortUrl, "/admin/") {
		return errcode.NewErrRsp2(errcode.Err_Token_NormalToken_Pwi_Error, " router is wrong! ")
	}

	if para.TokenParams["pwt"] != "0" {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Pwt_Error, "pwt is wrong")
	}

	return nil
}

func baseCheckAdminToken(para *reqparamodel.HttpReqParams) *errcode.ErrRsp {
	if !isAdminTokenUrl(para.ShortUrl) {
		return nil
	}

	if para.TokenParams["rsa"] == "" {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Rsa_Error, "rsa is wrong")
	}

	if para.TokenParams["jti"] != " " {
		log.Println(para.TokenParams["jti"], len(para.TokenParams["jti"]))
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Jti_Error, "jti is wrong")
	}

	if para.TokenParams["uid"] == " " || para.TokenParams["uid"] == "" {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, "uid is wrong")
	}

	key := keyutils.GetRSAKeyByVer(para.TokenParams["rsa"])
	if key == nil {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Rsa_Error, "rsa ver is wrong")
	}

	data, err1 := hex.DecodeString(para.TokenParams["uid"])
	if err1 != nil {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, err1.Error())
	}

	dedata, err2 := randutils.RsaDecrypt(key.PrivateKey, data)
	//log.Println(dedata)
	if err2 != nil {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, err2.Error())
	}

	//update uid
	para.TokenParams["uid"] = string(dedata)

	if para.TokenParams["pwt"] != "1" && para.TokenParams["pwt"] != "2" {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Pwt_Error, "pwt is wrong")
	}

	return nil
}

func baseCheckMineToken(para *reqparamodel.HttpReqParams) *errcode.ErrRsp {
	if !isUserMineTokenUrl(para.ShortUrl) {
		return nil
	}

	if para.TokenParams["TokenType"] != fmt.Sprint(tokenmodel.Const_Token_Type_Private) {
		return errcode.NewErrRsp2(errcode.Err_Token_Para_Uid_Error, "should be private token!")
	}

	return nil
}
