package tokens

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"strings"
	"web/component/errcode"
	"web/component/keyutils"
	//"web/component/pageutils"
	"web/models/reqparamodel"
	"web/models/tokenmodel"
)

func TokenHeaderHandler(c martini.Context, ren render.Render, req *http.Request) {
	for _, thridUrl := range thirdUrlAccessTables {
		if strings.Contains(req.URL.String(), thridUrl) {
			return
		}
	}

	httpReqPara := reqparamodel.NewHttpReqParams()

	//log.Println("req.header", req.Header)
	log.Println("req.remoteaddr", req.RemoteAddr)

	err := errors.New("")
	ok := false
	//logger := &MyCustomLogger{req}
	//c.Map(logger) // mapped as *MyCustomLogger
	if req.Header.Get("X-ACCESS-TOKEN") == "" {
		ren.JSON(http.StatusUnauthorized, errcode.NewErrRsp2(errcode.Err_Token_UnAuthorized_Error, "access token wrong"))
		return
	}

	jwtToken, err := decodeJWT(req.Header.Get("X-ACCESS-TOKEN"))
	if err != nil {
		ren.JSON(http.StatusUnauthorized, errcode.NewErrRsp2(errcode.Err_Token_UnAuthorized_Error, err.Error()))
		return
	}

	log.Println("jwtToken Header", jwtToken.Header)
	log.Println("jwtToken Claims", jwtToken.Claims)

	httpReqPara.TokenParams["iss"], _ = jwtToken.Claims["iss"].(string)
	httpReqPara.TokenParams["sub"], _ = jwtToken.Claims["sub"].(string)
	httpReqPara.TokenParams["rsa"], _ = jwtToken.Claims["rsa"].(string)

	httpReqPara.TokenParams["uid"], _ = jwtToken.Claims["uid"].(string)
	httpReqPara.TokenParams["jti"], ok = jwtToken.Claims["jti"].(string)
	log.Println("jti ok:", ok)

	httpReqPara.TokenParams["exp"], _ = jwtToken.Claims["exp"].(string)
	httpReqPara.TokenParams["pwt"], _ = jwtToken.Claims["pwt"].(string)
	httpReqPara.TokenParams["pwi"], _ = jwtToken.Claims["pwi"].(string)

	httpReqPara.TokenParams["stamp"], _ = jwtToken.Claims["stamp"].(string)

	httpReqPara.TokenParams["addr"] = fmt.Sprint(req.RemoteAddr)

	//httpReqPara.URLParams = pageutils.Parse_URL_Form(req)
	req.ParseForm()

	httpReqPara.Url = req.URL.String()

	httpReqPara.URLParams = req.Form
	httpReqPara.ShortUrl = strings.Split(req.URL.String(), "?")[0]

	if req.Method == "POST" {
		for field, _ := range req.PostForm {
			httpReqPara.PostFields = append(httpReqPara.PostFields, field)
		}
	}

	//log.Println("httpreqpara", httpReqPara)

	if errrsp := baseCheckToken(httpReqPara); errrsp != nil {
		ren.JSON(http.StatusUnauthorized, errrsp)
		return
	}

	log.Println("httpreqpara", httpReqPara)

	c.Map(httpReqPara)
}

func IsAdminToken(para *reqparamodel.HttpReqParams) bool {
	if para.TokenParams["TokenType"] == fmt.Sprint(tokenmodel.Const_Token_Type_Admin) {
		return true
	}

	return false
}

func decodeJWT(myToken string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		// jwt.SigningMethodHMAC  jwt.SigningMethodRSA
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return keyutils.GetHS256Key(), nil
	})

	//log.Println("decodeJWT, ", jwtToken)

	if err == nil && jwtToken.Valid {
		return jwtToken, nil
	} else {
		if err == nil {
			return nil, fmt.Errorf("Token is not valid")
		}

		return nil, err
	}
}
