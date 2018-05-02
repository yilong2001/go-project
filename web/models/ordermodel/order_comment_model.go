package ordermodel

import (
	"web/component/objutils"
)

type OrderCommentInfo struct {
	CommentId     int
	CommentUserId int
	OrderId       int
	OrderType     int
	CustomerId    int
	ProviderId    int
	ServiceId     int
	ServiceName   string
	Comment       string `schema:"Comment"`
	Star          int
	ReplyInfoNum  int
	ReplyUpNum    int
	ReplyDownNum  int
	CreateTime    string
	UpdateTime    string
}

func (this *OrderCommentInfo) GetUniqId() int {
	return this.CommentId
}

func (this *OrderCommentInfo) GetUniqIdName() string {
	return "CommentId"
}

func (this *OrderCommentInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *OrderCommentInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *OrderCommentInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

func NewOrderCommentInfo() *OrderCommentInfo {
	return &OrderCommentInfo{
		CommentId:     -1,
		CommentUserId: -1,
		OrderId:       -1,
		CustomerId:    -1,
		ProviderId:    -1,
		ServiceId:     -1,
		ServiceName:   "",
		Comment:       "",
		Star:          -1,
		ReplyInfoNum:  0,
		ReplyUpNum:    0,
		ReplyDownNum:  0,
		CreateTime:    "",
		UpdateTime:    "",
	}
}
