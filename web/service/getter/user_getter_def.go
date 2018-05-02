package getter

import (
//"database/sql"
//"log"
//"github.com/go-martini/martini"
//"github.com/gorilla/schema"
//"github.com/martini-contrib/binding"
//"github.com/martini-contrib/render"

//"errors"
//"net/http"
//"strconv"
//"strings"
//"time"
//"reflect"

//"web/component/cfgutils"
//"web/component/errcode"
//"web/component/idutils"
//"web/component/objutils"
//"web/component/sqlutils"

//"web/dal/sqldrv"
//"web/models/tokenmodel"
//"web/models/usermodel"
//"web/service/utils"
)

type ModelInfoGetter struct {
	userTableName                string
	jobTableName                 string
	serviceTableName             string
	tokenTableName               string
	couponTableName              string
	orderTableName               string
	platformCouponTableName      string
	adminUserTableName           string
	userFavouriteTableName       string
	orderSubTableName            string
	smsCodeTableName             string
	courseMainTableName          string
	courseCatalogFirstTableName  string
	courseCatalogSecondTableName string

	userAvAddrTableName string
}

// var modelInfoGetter *ModelInfoGetter = &ModelInfoGetter{
// 	userTableName:           "web_users",
// 	jobTableName:            "web_jobs",
// 	serviceTableName:        "web_services",
// 	tokenTableName:          "web_json_tokens",
// 	couponTableName:         "web_user_coupons",
// 	orderTableName:          "web_orders",
// 	platformCouponTableName: "platform_coupons",
// 	adminUserTableName:      "web_admin_users",
// }

func GetModelInfoGetter() *ModelInfoGetter {
	return &ModelInfoGetter{
		userTableName:                "web_users",
		jobTableName:                 "web_jobs",
		serviceTableName:             "web_services",
		tokenTableName:               "web_json_tokens",
		couponTableName:              "web_user_coupons",
		orderTableName:               "web_orders",
		platformCouponTableName:      "platform_coupons",
		adminUserTableName:           "web_admin_users",
		userFavouriteTableName:       "web_user_favourites",
		orderSubTableName:            "web_order_subs",
		smsCodeTableName:             "web_user_smscodes",
		courseMainTableName:          "web_courses",
		courseCatalogFirstTableName:  "web_course_catalog_firsts",
		courseCatalogSecondTableName: "web_course_catalog_seconds",

		userAvAddrTableName: "web_user_av_abbrs",
	}
}

func (this *ModelInfoGetter) GetUserTableName() string {
	return this.userTableName
}

func (this *ModelInfoGetter) GetAdminUserTableName() string {
	return this.adminUserTableName
}

func (this *ModelInfoGetter) GetJobTableName() string {
	return this.jobTableName
}

func (this *ModelInfoGetter) GetServiceTableName() string {
	return this.serviceTableName
}

func (this *ModelInfoGetter) GetTokenTableName() string {
	return this.tokenTableName
}

func (this *ModelInfoGetter) GetCouponTableName() string {
	return this.couponTableName
}

func (this *ModelInfoGetter) GetOrderTableName() string {
	return this.orderTableName
}

func (this *ModelInfoGetter) GetPlatformCouponTableName() string {
	return this.platformCouponTableName
}

func (this *ModelInfoGetter) GetUserFavouriteTableName() string {
	return this.userFavouriteTableName
}

func (this *ModelInfoGetter) GetOrderSubTableName() string {
	return this.orderSubTableName
}

func (this *ModelInfoGetter) GetSmsCodeTableName() string {
	return this.smsCodeTableName
}

func (this *ModelInfoGetter) GetCourseMainTableName() string {
	return this.courseMainTableName
}
func (this *ModelInfoGetter) GetCourseCatalogFirstTableName() string {
	return this.courseCatalogFirstTableName
}
func (this *ModelInfoGetter) GetCourseCatalogSecondTableName() string {
	return this.courseCatalogSecondTableName
}
func (this *ModelInfoGetter) GetUserAVAddrTableName() string {
	return this.userAvAddrTableName
}
