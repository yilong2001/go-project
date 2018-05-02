package clientmodel

import (
	"web/models/jobmodel"
	"web/models/usermodel"
)

type ClientJobInfo struct {
	JobInfo  *jobmodel.JobInfo
	UserInfo *usermodel.UserInfo
}
