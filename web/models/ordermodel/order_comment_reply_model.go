package ordermodel

import ()

type OerderCommentReplyInfo struct {
	ReplyId       int
	ReplyUserId   int
	CommentId     int
	CommentUserId int
	ReplyType     int
	ReplyInfo     string
	CreateTime    string
	UpdateTime    string
}

func NewOerderCommentReplyInfo() *OerderCommentReplyInfo {
	return &OerderCommentReplyInfo{
		ReplyId:       -1,
		ReplyUserId:   -1,
		CommentId:     -1,
		CommentUserId: -1,
		ReplyType:     -1,
		ReplyInfo:     "",
		CreateTime:    "",
		UpdateTime:    "",
	}
}
