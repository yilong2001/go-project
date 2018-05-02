package clientmodel

import (
	"web/models/coursemodel"
	"web/models/jobmodel"
	"web/models/servemodel"
	"web/models/usermodel"
)

type ClientFavouriteInfo struct {
	JobInfo       *jobmodel.JobInfo
	ServeInfo     *servemodel.ServeInfo
	CourseInfo    *coursemodel.CourseMainInfo
	FavouriteInfo *usermodel.UserFavouriteInfo
	UserInfo      *usermodel.UserInfo
}
