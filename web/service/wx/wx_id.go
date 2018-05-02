package wx

import (
	//"database/sql"
	"github.com/go-martini/martini"
	"log"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	//"errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	//"strings"
	//"time"
	//"reflect"

	//"web/component/cfgutils"
	"web/component/errcode"
	"web/component/wxutils"
	//"web/component/idutils"
	//"web/component/objutils"
	//"web/component/sqlutils"

	//"web/dal/sqldrv"
	//"web/models/basemodel"
	//"web/models/firmmodel"
	"web/models/rendermodel"
	"web/models/reqparamodel"
	//"web/models/usermodel"
	"web/service/routers"
	"web/service/userups"
	"web/service/utils"
)

func init() {
	wxIdsRouterBuilder()
}

func wxIdsRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Get("/user/wx/openid", GetUserWxOpenId)
}

func getFile(url string) ([]byte, error) {
	oid, err := http.Get(url)
	if err != nil {
		log.Println("wx_id get file : ", err)
		return nil, err
	}
	defer oid.Body.Close()
	oids, err2 := ioutil.ReadAll(oid.Body)
	if err2 != nil {
		log.Println("wx_id get file : ", err2)
		return nil, err2
	}
	return oids, nil
}

func GetUserWxOpenId(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	headParams.MergeMartiniParams(params)

	wxcode := headParams.URLParams.Get("CODE")

	appid, secret, auth_type := wxutils.GetZhiEasyAppId(), wxutils.GetZhiEasyAppSecret(), "authorization_code"

	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + appid + "&secret=" + secret + "&code=" + wxcode + "&grant_type=" + auth_type

	ids, err0 := getFile(url)
	if err0 != nil {
		errrsp := errcode.NewErrRsp2(errcode.Err_Wx_AccessToken_Get_Error, url+":"+err0.Error())
		r.JSON(200, errrsp)
		return
	}

	dat := &wxutils.WxAccessToken{}

	err := json.Unmarshal(ids, dat)
	if err != nil {
		errrsp := errcode.NewErrRsp2(errcode.Err_Wx_AccessToken_Get_Error, err.Error())
		r.JSON(200, errrsp)
		return
	}

	if utils.IsFieldCorrectWithRule("user_id", headParams.RouterParams["UserId"]) {
		fakeren := &rendermodel.FakeMrtiniRender{}
		userid, _ := strconv.ParseInt(headParams.RouterParams["UserId"], 10, 32)
		ok := userups.UpdateUserWXId(nil, nil, dat.Openid, "", int(userid), fakeren)
		if !ok {
			errrsp := errcode.NewErrRsp2(errcode.Err_Wx_AccessToken_Get_Error, "update user openid wrong")
			r.JSON(200, errrsp)
			return
		}
	}

	r.JSON(200, map[string]interface{}{"ErrCode": 0, "ErrMsg": "ok", "Data": dat.Openid})
}
