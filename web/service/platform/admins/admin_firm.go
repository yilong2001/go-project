package admins

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"web/component/errcode"
	"web/models/reqparamodel"
	"web/service/firms"
	"web/service/routers"
)

func init() {
	adminFirmsRouterBuilderEx()
}

func adminFirmsRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Post("/admin/reviewer/firm", AdminAddFirm)
	m.Get("/admin/reviewer/firm", firms.GetOpenFirms)
}

func AdminAddFirm(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {
	if !isAdminToken(headParams) {
		r.JSON(200, errcode.NewErrRsp2(errcode.Err_Form_Para_AdminUser_Error, "not admin user"))
		return
	}

	firms.AddUserFirm(headParams, params, req, r)
}
