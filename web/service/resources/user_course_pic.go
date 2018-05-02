package resources

import (
	//"database/sql"
	"github.com/go-martini/martini"
	//"io"
	//"log"
	"net/http"
	//"strconv"
	//"time"
	//"os"
	//"github.com/gorilla/schema"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	//"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/fileutils"
	//"web/component/sqlutils"
	"web/models/reqparamodel"
	"web/models/resourcemodel"

	//"web/dal/sqldrv"
	"web/service/routers"
	//"web/service/uploads"
	//"web/service/utils"
)

func init() {
	userCoursePicRouterBuilder()
}

func userCoursePicRouterBuilder() {
	m := routers.GetRouterHandler()

	m.Post("/user/mine/course/:CourseId/pic", AddPictureForCourse)
	m.Post("/user/mine/course/:CourseId/video", AddVideoForCourse)
	m.Post("/user/mine/course/:CourseId/article", AddArticleForCourse)

	m.Get("/user/mine/course/:CourseId/pic", GetUserResource)
	m.Get("/user/mine/course/:CourseId/video", GetUserResource)
	m.Get("/user/mine/course/:CourseId/article", GetUserResource)

	m.Post("/user/mine/service/:ServiceId/pic", AddPictureForService)
	m.Post("/user/mine/service/:ServiceId/video", AddVideoForService)
	m.Post("/user/mine/service/:ServiceId/article", AddArticleForService)

	m.Get("/user/mine/service/:ServiceId/pic", GetUserResource)
	m.Get("/user/mine/service/:ServiceId/video", GetUserResource)
	m.Get("/user/mine/service/:ServiceId/article", GetUserResource)
}

func addResource(resourceType, referType int, headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {

	headParams.MergeMartiniParams(params)

	ctrl := NewUserResourceController(resourceType, referType)

	defer ctrl.CloseDB()

	obj := NewUserResourceControllerObject(ctrl)

	info := &resourcemodel.UserResourceInfo{}

	res := obj.Util_CreateObjectWithId(info, headParams, req, r)
	if res {
		ctrl.GetTX().Commit()
	}
}

func AddPictureForCourse(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {

	addResource(resourcemodel.Const_User_Resource_Pic, resourcemodel.Const_User_Resource_Refer_Course, headParams, params, req, r)
}

func AddVideoForCourse(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {

	addResource(resourcemodel.Const_User_Resource_Video, resourcemodel.Const_User_Resource_Refer_Course, headParams, params, req, r)
}

func AddArticleForCourse(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {

	addResource(resourcemodel.Const_User_Resource_Article, resourcemodel.Const_User_Resource_Refer_Course, headParams, params, req, r)
}

func AddPictureForService(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {

	addResource(resourcemodel.Const_User_Resource_Pic, resourcemodel.Const_User_Resource_Refer_Service, headParams, params, req, r)
}

func AddVideoForService(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {

	addResource(resourcemodel.Const_User_Resource_Video, resourcemodel.Const_User_Resource_Refer_Service, headParams, params, req, r)
}

func AddArticleForService(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {

	addResource(resourcemodel.Const_User_Resource_Article, resourcemodel.Const_User_Resource_Refer_Service, headParams, params, req, r)
}

func GetUserResource(headParams *reqparamodel.HttpReqParams, params martini.Params, req *http.Request, r render.Render) {

	headParams.MergeMartiniParams(params)

	ctrl := NewUserResourceController(getResourceType(headParams), getReferType(headParams))

	defer ctrl.CloseDB()

	obj := NewUserResourceControllerQueryObject(ctrl)

	info := &resourcemodel.UserResourceInfo{}

	res := obj.Util_GetObjectWithId(info, headParams, req, r, nil, nil)
	if res {
		ctrl.GetTX().Commit()
	}
}
