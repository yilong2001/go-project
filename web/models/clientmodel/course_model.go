package clientmodel

import (
	"web/models/coursemodel"
	"web/models/usermodel"
)

type ClientCourseCatalogFirst struct {
	FirstInfo   coursemodel.CourseCatalogFirstInfo
	SecondInfos *[]coursemodel.CourseCatalogSecondInfo
}

type FreeVideoPath struct {
	UrlPath map[string]string
}

type ClientCourseInfo struct {
	VideoPathUrl       *FreeVideoPath
	CourseMainInfo     *coursemodel.CourseMainInfo
	CourseCatalogInfos *[]ClientCourseCatalogFirst
	UserInfo           *usermodel.UserInfo
}
