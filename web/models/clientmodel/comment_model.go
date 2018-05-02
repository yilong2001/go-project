package clientmodel

import (
	"web/models/ordermodel"
	"web/models/usermodel"
)

type ClientCommentInfo struct {
	CommentInfo *ordermodel.OrderCommentInfo
	UserInfo    *usermodel.UserInfo
}
