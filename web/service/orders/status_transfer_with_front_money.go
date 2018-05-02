package orders

import (
	"database/sql"
	// "fmt"
	// "github.com/go-martini/martini"
	// "github.com/martini-contrib/render"
	"log"

	"errors"
	"net/http"
	// "strconv"
	// "strings"
	// "time"
	// //"reflect"
	//"math"

	// "web/component/cfgutils"
	// "web/component/errcode"
	// //"web/component/idutils"
	// //"web/component/objutils"
	// "web/component/sqlutils"

	"web/component/orderutils"
	// "web/dal/sqldrv"
	// "web/models/basemodel"
	"web/models/ordermodel"
	"web/models/rendermodel"
	"web/models/reqparamodel"
	"web/service/orderpays"
	//"web/service/getter"
	// "web/service/routers"
	// "web/service/users"
	// "web/service/utils"
)

var globalWithFrontMoneyOrderStatusTransferMachines []*ordermodel.OrderStatusTransfer = []*ordermodel.OrderStatusTransfer{
	//provider accept order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Provider_Accept,
		NextStatus:    orderutils.Order_Status_Wait_Customer_PrePay,
		Router:        "accepter",
		ChangedFields: []string{"FrontMoney", "OrderStatus", "UpdateTime"},
		PayFunction:   nil,
		MoreFunction:  order_More_NotifyMsg,
	},

	//provider reject order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Provider_Accept,
		NextStatus:    orderutils.Order_Status_Provider_Rejected,
		Router:        "rejecter",
		ChangedFields: []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   nil,
		MoreFunction:  order_More_NotifyMsg,
	},

	//customer cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Provider_Accept,
		NextStatus:    orderutils.Order_Status_Customer_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   nil,
		MoreFunction:  order_More_NotifyMsg,
	},

	//customer prepay
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Wait_Arrange_Date,
		Router:        "prepayer",
		ChangedFields: []string{"AccountLocked", "PrepayType", "PrepayStatus", "PrepayMoney", "PrepayTime", "UserCouponId01", "PayCouponId01", "PayCouponMoney01", "OrderStatus", "UpdateTime"},
		PayFunction:   orderFrontMoney_PrePay,
		MoreFunction:  order_More_NotifyMsg,
	},

	//provider modify dingjin, unit price, only little, can not more
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Wait_Customer_PrePay,
		Router:        "cost/dingjin",
		ChangedFields: []string{"FrontMoney", "UpdateTime"},
		PayFunction:   nil,
		MoreFunction:  order_More_NotifyMsg,
	},

	//provider modify dingjin, unit price, only little, can not more
	// &ordermodel.OrderStatusTransfer{
	// 	Role:          orderutils.Order_Role_Provider,
	// 	OrderType:     orderutils.Order_Type_With_Front_Money_Business,
	// 	CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
	// 	NextStatus:    orderutils.Order_Status_Wait_Customer_PrePay,
	// 	Router:        "cost/unitcost",
	// 	ChangedFields: []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"},
	// 	PayFunction:   nil,
	// 	MoreFunction:  nil,
	// },

	//customer cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Customer_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   nil,
		MoreFunction:  order_More_NotifyMsg,
	},

	//provider cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Provider_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   nil,
		MoreFunction:  order_More_NotifyMsg,
	},

	//customer accept date
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Arrange_Date,
		NextStatus:    orderutils.Order_Status_Wait_Customer_Complete,
		Router:        "date/accepter",
		ChangedFields: []string{"ExpiredDate", "OrderStatus", "UpdateTime"},
		PayFunction:   nil,
		MoreFunction:  order_More_NotifyMsg,
	},

	//provider cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Arrange_Date,
		NextStatus:    orderutils.Order_Status_Provider_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"PrepayStatus", "AccountLocked", "OverReason", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   orderFrontMoney_RollbackPrepay,
		MoreFunction:  order_More_NotifyMsg,
	},

	//customer complete session
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_Complete,
		NextStatus:    orderutils.Order_Status_Wait_Comment,
		Router:        "complete",
		ChangedFields: []string{"PayStatus", "PayType", "PayMoney", "AccountLocked", "TotalPay", "OrderStatus", "UpdateTime"},
		PayFunction:   orderFrontMoney_PayComplete,
		MoreFunction:  order_PayComplete_More,
	},

	//customer complete session
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_Complete,
		NextStatus:    orderutils.Order_Status_Wait_Comment,
		Router:        "uncomplete",
		ChangedFields: []string{"AccountLocked", "Feedback", "OrderStatus", "UpdateTime"},
		PayFunction:   orderFrontMoney_RollbackPrepay,
		MoreFunction:  order_PayComplete_More,
	},

	//customer cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_Complete,
		NextStatus:    orderutils.Order_Status_Customer_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"AccountLocked", "OverReason", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   orderFrontMoney_PayFrontMoney,
		MoreFunction:  order_More_NotifyMsg,
	},

	//provider cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_Complete,
		NextStatus:    orderutils.Order_Status_Provider_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"PrepayStatus", "AccountLocked", "OverReason", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   orderFrontMoney_RollbackPrepay,
		MoreFunction:  order_More_NotifyMsg,
	},

	//customer comment, push status changed to complete
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Comment,
		NextStatus:    orderutils.Order_Status_Complete,
		Router:        "comment",
		ChangedFields: []string{"Comment", "Star", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   nil,
		MoreFunction:  orderSingle_CustomComment_More,
	},

	//provider continue add comment
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Comment,
		NextStatus:    orderutils.Order_Status_Wait_Comment,
		Router:        "comment/plus",
		ChangedFields: []string{"Comment", "Star", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   nil,
		MoreFunction:  orderSingle_CustomComment_More,
	},

	//customer continue add comment
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Complete,
		NextStatus:    orderutils.Order_Status_Complete,
		Router:        "comment/plus",
		ChangedFields: []string{"Comment", "Star", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   nil,
		MoreFunction:  orderSingle_CustomComment_More,
	},

	//provider continue add comment
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_With_Front_Money_Business,
		CurrentStatus: orderutils.Order_Status_Complete,
		NextStatus:    orderutils.Order_Status_Complete,
		Router:        "comment/plus",
		ChangedFields: []string{"Comment", "Star", "OrderStatus", "UpdateTime", "OverTime"},
		PayFunction:   nil,
		MoreFunction:  orderSingle_CustomComment_More,
	},
}

func orderFrontMoney_PrePay(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {

	fakeren := &rendermodel.FakeMrtiniRender{}

	if newOrder.UserCouponId01 > 0 && oldOrder.FrontMoney > 0 {
		err := orderpays.Order_PrepayWithCoupon(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	}

	//call zhifubao, weixin pay interface
	if newOrder.PrepayType == orderutils.Pay_Type_Account_Balance {
		err := orderpays.Order_PayWithAccountBalance(db, tx, newOrder.AccountLocked, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	} else if newOrder.PrepayType == orderutils.Pay_Type_WeiXin {
		err := orderpays.Order_PayWithWxPay(db, tx, newOrder.PrepayMoney, oldOrder, newOrder, headParams, req, fakeren, out)
		if err != nil {
			log.Print("orderSingle_PrePay, wxpay wrong : ", err)
			return err
		}
	} else if newOrder.PrepayType == orderutils.Pay_Type_ZhiFuBao {
		err := orderpays.Order_PayWithAliPay(db, tx, newOrder.PrepayMoney, oldOrder, newOrder, headParams, req, fakeren, out)
		if err != nil {
			log.Print("orderSingle_PrePay, alipay wrong : ", err)
			return err
		}
	} else {
		return errors.New("do't support 这种支付类型!")
	}

	return nil
}

func orderFrontMoney_RollbackPrepay(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {
	fakeren := &rendermodel.FakeMrtiniRender{}

	if oldOrder.UserCouponId01 > 0 {
		err := orderpays.Order_RollbackWithCoupon(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	}

	if oldOrder.PrepayType == orderutils.Pay_Type_Account_Balance {
		err := orderpays.Order_RollbackWithAccountBalance(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	} else if oldOrder.PrepayType == orderutils.Pay_Type_WeiXin && oldOrder.PrepayMoney > 0 {
		err := orderpays.Order_RollbackWithWxPay(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	} else if oldOrder.PrepayType == orderutils.Pay_Type_ZhiFuBao {
		err := orderpays.Order_RollbackWithAliPay(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	} else {
		return errors.New("do not support pay type")
	}

	return nil
}

func orderFrontMoney_PayFrontMoney(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {
	fakeren := &rendermodel.FakeMrtiniRender{}

	if oldOrder.UserCouponId01 > 0 && oldOrder.PayCouponMoney01 > 0 {
		err := orderpays.Order_PayCompleteWithCoupon(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	}

	err := orderpays.Order_PayCompleteToProvider(db, tx, oldOrder.FrontMoney, oldOrder, newOrder, headParams, req, fakeren)
	if err != nil {
		return err
	}

	return nil
}

func orderFrontMoney_CalcExtraCost(total, front, coupon int) int {
	totalCost := total
	extraCost := 0
	baseCost := front
	if baseCost < coupon {
		baseCost = coupon
	}
	if totalCost > baseCost {
		extraCost = totalCost - baseCost
	}

	return extraCost
}

func orderFrontMoney_PayComplete(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {
	fakeren := &rendermodel.FakeMrtiniRender{}

	//if use coupon on prepay state, here complete coupon pay
	if oldOrder.UserCouponId01 > 0 && oldOrder.PayCouponMoney01 > 0 {
		err := orderpays.Order_PayCompleteWithCoupon(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	}

	/*
		orderSubInfos, err := getter.GetModelInfoGetter().GetMultiOrderSubsByOrderIds(db, oldOrder.OrderId)
		if err != nil {
			log.Print(err)
			return err
		}

		totalCost := 0

		for _, subInfo := range *orderSubInfos {
			if subInfo.OverStatus == 0 {
				totalCost = totalCost + subInfo.ActualNum*subInfo.UnitCost
			}
		}
	*/
	//if front money is enough
	//for example: frontmoeny is 100, coupon is 200, totalcost is 300, so customer should pay 300-200=100
	// totalCost := newOrder.TotalPay
	// extraCost := 0
	// baseCost := oldOrder.FrontMoney
	// if baseCost < oldOrder.PayCouponMoney01 {
	// 	baseCost = oldOrder.PayCouponMoney01
	// }
	// if totalCost > baseCost {
	// 	extraCost = totalCost - baseCost
	// }

	extraCost := orderFrontMoney_CalcExtraCost(newOrder.PayedTotal, oldOrder.FrontMoney, oldOrder.PayCouponMoney01)

	lastCost := newOrder.PayedTotal
	if lastCost < oldOrder.FrontMoney {
		lastCost = oldOrder.FrontMoney
	}

	if newOrder.PayType == orderutils.Pay_Type_Account_Balance && extraCost > 0 {
		err := orderpays.Order_PayWithAccountBalance(db, tx, extraCost, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}

		err = orderpays.Order_PayCompleteToProvider(db, tx, lastCost, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	} else if newOrder.PayType == orderutils.Pay_Type_WeiXin && extraCost > 0 {
		err := orderpays.Order_PayWithWxPay(db, tx, extraCost, oldOrder, newOrder, headParams, req, fakeren, out)
		if err != nil {
			log.Print("orderFrontMoney_PayComplete, wxpay wrong : ", err)
			return err
		}
	} else if newOrder.PayType == orderutils.Pay_Type_ZhiFuBao && extraCost > 0 {
		err := orderpays.Order_PayWithAliPay(db, tx, extraCost, oldOrder, newOrder, headParams, req, fakeren, out)
		if err != nil {
			log.Print("orderFrontMoney_PayComplete, wxpay wrong : ", err)
			return err
		}
	} else if extraCost > 0 {
		return errors.New("pay type is wrong!")
	}

	//if total cost is more littler than front money

	return nil
}
