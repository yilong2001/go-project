package ordermodel

import (
	//"strings"
	"web/component/objutils"
)

type OrderSubInfo struct {
	OrderSubId int
	OrderId    int `schema:"OrderId"`
	ReferId    int `schema:"ReferId"`
	ReferType  int
	ReferName  string `schema:"ReferName"`
	ExpectNum  int    `schema:"ExpectNum"`
	ActualNum  int    `schema:"ActualNum"`
	UnitCost   int    `schema:"UnitCost"`
	OverStatus int
	OverTime   string
}

func (this *OrderSubInfo) GetUniqId() int {
	return this.OrderSubId
}

func (this *OrderSubInfo) GetUniqIdName() string {
	return "OrderSubId"
}

func (this *OrderSubInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *OrderSubInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *OrderSubInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

type OrderInfo struct {
	OrderId   int
	OrderType int

	CustomerId  int
	ProviderId  int
	ServiceId   int `schema:"ServiceId"`
	ServiceName string

	FrontMoney       int `schema:"FrontMoney"`
	ContractDuration int `schema:"ContractDuration"`

	OrderPrice  int `schema:"OrderPrice"`
	OrderStatus int

	PrepayType   int `schema:"PrepayType"`
	PrepayStatus int
	PrepayMoney  int
	PrepayTime   string

	AccountLocked int

	PayType   int `schema:"PayType"`
	PayStatus int
	PayMoney  int
	PayTime   string

	PayedTotal  int
	ExpiredDate string

	UserCouponId01   int `schema:"UserCouponId01"`
	PayCouponId01    int
	PayCouponMoney01 int

	CustomerExpect     string `schema:"CustomerExpect"`
	CustomerIntroduce  string `schema:"CustomerIntroduce"`
	CustomerDateAdvice string `schema:"CustomerDateAdvice"`

	ArragedDateOption1 string `schema:"ArragedDateOption1"`
	ArragedDateOption2 string `schema:"ArragedDateOption2"`
	ArragedDateOp      int    `schema:"ArragedDateOp"`

	Feedback string `schema:"Feedback"`
	Comment  string `schema:"Comment"`

	OverReason string `schema:"OverReason"`

	Memo string `schema:"Memo"`

	Level int `schema:"Level"`
	Star  int `schema:"Star"`

	CreateTime string
	UpdateTime string
	OverTime   string

	OrderOutId string

	OrderSubs []*OrderSubInfo
}

func (this *OrderInfo) IsUserIdValid() bool {
	if this.CustomerId < 10000 || this.ProviderId < 10000 {
		return false
	}

	return true
}

func (this *OrderInfo) GetUserId() int {
	return this.CustomerId
}

func (this *OrderInfo) GetUniqId() int {
	return this.OrderId
}

func (this *OrderInfo) GetUniqIdName() string {
	return "OrderId"
}

func (this *OrderInfo) GetSkipFieldsForCustomerQuery() []string {
	return []string{}
}

func (this *OrderInfo) GetSkipFieldsForProviderQuery() []string {
	return []string{"UserCouponId01", "PayCouponId01", "PayCouponMoney01", "PrepayType"}
}

func (this *OrderInfo) GetSpecFieldsForProviderAccept() []string {
	return []string{"ArragedDateOption1", "ArragedDateOption2", "OrderStatus", "UpdateTime", "FrontMoney"}
}

func (this *OrderInfo) GetSpecFieldsForProviderReject() []string {
	return []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"}
}

func (this *OrderInfo) GetSpecFieldsForPayer() []string {
	return []string{"PayType", "UserCouponId01", "OrderStatus", "UpdateTime"}
}

func (this *OrderInfo) GetSpecFieldsForAcceptDate() []string {
	return []string{"ArragedDate", "OrderStatus", "UpdateTime", "ArragedDateOp"}
}

func (this *OrderInfo) GetSpecFieldsForAddFeedback() []string {
	return []string{"Feedback", "OrderStatus", "UpdateTime"}
}

func (this *OrderInfo) GetSpecFieldsForAddComment() []string {
	return []string{"Comment", "Star", "OrderStatus", "UpdateTime", "OverTime"}
}

func (this *OrderInfo) GetSpecFieldsForCustomerCancel() []string {
	return []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"}
}

func (this *OrderInfo) GetSpecFieldsForProviderCancel() []string {
	return []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"}
}

func (this *OrderInfo) GetSpecFieldsForUpdateDateOption() []string {
	return []string{"ArragedDateOption1", "ArragedDateOption2", "ArragedDateOption3", "UpdateTime"}
}

func (this *OrderInfo) GetSpecFieldsForUpdateFeedback() []string {
	return []string{"Feedback", "UpdateTime"}
}

func (this *OrderInfo) GetSpecFieldsForUpdateComment() []string {
	return []string{"Comment", "Star", "UpdateTime"}
}

func (this *OrderInfo) GetWholeFields() (map[string]interface{}, map[string]interface{}) {
	return objutils.GetWholeFields(this)
}

func (this *OrderInfo) GetFieldsWithSkip(skips []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSkip(this, skips)

}
func (this *OrderInfo) GetFieldsWithSpecs(specs []string) (map[string]interface{}, map[string]interface{}) {
	return objutils.GetFieldsWithSpecs(this, specs)
}

func NewOrderInfo() *OrderInfo {
	return &OrderInfo{}
}
