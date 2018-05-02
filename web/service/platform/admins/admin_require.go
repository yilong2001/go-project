package admins

import (
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/render"
	//"net/http"
	//"web/component/errcode"
	//"web/models/reqparamodel"
	"web/service/others"
	"web/service/routers"
)

func init() {
	adminRequiresRouterBuilderEx()
}

func adminRequiresRouterBuilderEx() {
	m := routers.GetRouterHandler()

	m.Get("/admin/reviewer/require", others.GetRequiresForAdmin)
}
