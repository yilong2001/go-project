package orders

import (
	"database/sql"
	//"fmt"
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/render"
	"errors"
	//"fmt"
	"log"
	"net/http"
	//"strconv"
	"strings"
	"time"
	//"reflect"
	//"web/component/cfgutils"
	"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	"web/component/orderutils"
	//"web/component/rongcloud"
	"web/component/sqlutils"
	//"web/dal/sqldrv"
	//"web/models/basemodel"
	"web/models/ordermodel"
	"web/models/rendermodel"
	"web/models/reqparamodel"
	//"web/service/coupons"
	//"web/service/getter"
	//"web/service/pays"
	"web/service/serveups"
	//"web/service/routers"
	//"web/service/immsgs"
	"web/service/courseups"
	"web/service/orderpays"
	"web/service/userups"
	//"web/service/utils"
)

var globalSignleBusinessOrderStatusTransferMachines []*ordermodel.OrderStatusTransfer = []*ordermodel.OrderStatusTransfer{
	//provider accept order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Provider_Accept,
		NextStatus:    orderutils.Order_Status_Wait_Customer_PrePay,
		Router:        "accepter",
		ChangedFields: []string{"ArragedDateOption1", "ArragedDateOption2", "OrderStatus", "UpdateTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//provider reject order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Provider_Accept,
		NextStatus:    orderutils.Order_Status_Provider_Rejected,
		Router:        "rejecter",
		ChangedFields: []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//customer cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Provider_Accept,
		NextStatus:    orderutils.Order_Status_Customer_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//customer prepay
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Wait_Arrange_Date,
		Router:        "prepayer",
		ChangedFields: []string{"AccountLocked", "PrepayType", "PrepayStatus", "PrepayMoney", "PrepayTime", "UserCouponId01", "PayCouponId01", "PayCouponMoney01", "OrderStatus", "UpdateTime"},

		PayFunction:  orderSingle_PrePay,
		MoreFunction: order_More_NotifyMsg,
	},

	//customer cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Customer_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"AccountLocked", "OverReason", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//provider update date options
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Wait_Customer_PrePay,
		Router:        "date/option",
		ChangedFields: []string{"ArragedDateOption1", "ArragedDateOption2", "UpdateTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//provider cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Provider_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"AccountLocked", "OverReason", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//customer accept date
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Arrange_Date,
		NextStatus:    orderutils.Order_Status_Wait_Provider_Feedback,
		Router:        "date/accepter",
		ChangedFields: []string{"OrderStatus", "UpdateTime", "ArragedDateOp"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//provider update date options
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Arrange_Date,
		NextStatus:    orderutils.Order_Status_Wait_Arrange_Date,
		Router:        "date/option",
		ChangedFields: []string{"ArragedDateOption1", "ArragedDateOption2", "UpdateTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//provider cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Arrange_Date,
		NextStatus:    orderutils.Order_Status_Provider_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"PrepayStatus", "AccountLocked", "OverReason", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  orderSingle_RollbackPay,
		MoreFunction: order_More_NotifyMsg,
	},

	//provider complete session and feedback to customer
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Provider_Feedback,
		NextStatus:    orderutils.Order_Status_Wait_Comment,
		Router:        "feedback",
		ChangedFields: []string{"AccountLocked", "PrepayStatus", "PayedTotal", "Feedback", "OrderStatus", "UpdateTime"},

		PayFunction:  orderSingle_PayComplete,
		MoreFunction: order_PayComplete_More,
	},

	//customer cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Provider_Feedback,
		NextStatus:    orderutils.Order_Status_Customer_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"PrepayStatus", "AccountLocked", "OverReason", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  orderSingle_DecideRollOrPay,
		MoreFunction: order_More_NotifyMsg,
	},

	//provider cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Provider_Feedback,
		NextStatus:    orderutils.Order_Status_Provider_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"PrepayStatus", "AccountLocked", "OverReason", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  orderSingle_RollbackPay,
		MoreFunction: order_More_NotifyMsg,
	},

	//customer comment
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Comment,
		NextStatus:    orderutils.Order_Status_Complete,
		Router:        "comment",
		ChangedFields: []string{"Comment", "Star", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  nil,
		MoreFunction: orderSingle_CustomComment_More,
	},

	//provider update feedback
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Wait_Comment,
		NextStatus:    orderutils.Order_Status_Wait_Comment,
		Router:        "feedback/plus",
		ChangedFields: []string{"Feedback", "UpdateTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//customer add comment
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Complete,
		NextStatus:    orderutils.Order_Status_Complete,
		Router:        "comment/plus",
		ChangedFields: []string{"Comment", "Star", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  nil,
		MoreFunction: orderSingle_CustomComment_More,
	},

	//provider update feedback
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Provider,
		OrderType:     orderutils.Order_Type_Single_Business,
		CurrentStatus: orderutils.Order_Status_Complete,
		NextStatus:    orderutils.Order_Status_Complete,
		Router:        "feedback/plus",
		ChangedFields: []string{"Feedback", "UpdateTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},
}

func orderSingle_PrePay(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {

	log.Println("orderSingle_PrePay, newOrder.PrepayType", newOrder.PrepayType)

	fakeren := &rendermodel.FakeMrtiniRender{}

	if newOrder.UserCouponId01 > 0 && oldOrder.OrderPrice > 0 {
		err := orderpays.Order_PrepayWithCoupon(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			log.Print("orderSingle_PrePay", err)
			return err
		}
	}

	if newOrder.PrepayType == orderutils.Pay_Type_Account_Balance {
		err := orderpays.Order_PayWithAccountBalance(db, tx, newOrder.AccountLocked, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			log.Print("orderSingle_PrePay", err)
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

	//call zhifubao, weixin pay interface
	return nil
}

func orderSingle_RollbackPay(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
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

func orderSingle_PayComplete(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {
	fakeren := &rendermodel.FakeMrtiniRender{}

	if oldOrder.UserCouponId01 > 0 && oldOrder.PayCouponMoney01 > 0 {
		err := orderpays.Order_PayCompleteWithCoupon(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	}

	err := orderpays.Order_PayCompleteToProvider(db, tx, oldOrder.OrderPrice, oldOrder, newOrder, headParams, req, fakeren)
	if err != nil {
		return err
	}

	return nil
}

func orderSingle_isRollbackForCustomerCancel(oldOrder *ordermodel.OrderInfo) bool {
	if oldOrder.ArragedDateOp == orderutils.Order_DateArrange_Op_First {
		et, err := time.Parse("2006-01-02 15:04:05", oldOrder.ArragedDateOption1)
		if err == nil && et.Unix() > time.Now().Add(time.Hour*48).Unix() {
			return true
		}
	}

	if oldOrder.ArragedDateOp == orderutils.Order_DateArrange_Op_Second {
		et, err := time.Parse("2006-01-02 15:04:05", oldOrder.ArragedDateOption2)
		if err == nil && et.Unix() > time.Now().Add(time.Hour*48).Unix() {
			return true
		}
	}

	return false
}

func orderSingle_DecideRollOrPay(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {
	//fakeren := &rendermodel.FakeMrtiniRender{}

	isRollback := orderSingle_isRollbackForCustomerCancel(oldOrder)

	log.Println("orderSingle_DecideRollOrPay, isRollback = ", isRollback)

	// if oldOrder.ArragedDateOp == orderutils.Order_DateArrange_Op_First {
	// 	et, err := time.Parse("2006-01-02 15:04:05", oldOrder.ArragedDateOption1)
	// 	if err == nil && et.Unix() > time.Now().Add(time.Hour*24).Unix() {
	// 		isRollback = true
	// 	}
	// }

	// if oldOrder.ArragedDateOp == orderutils.Order_DateArrange_Op_Second {
	// 	et, err := time.Parse("2006-01-02 15:04:05", oldOrder.ArragedDateOption2)
	// 	if err == nil && et.Unix() > time.Now().Add(time.Hour*24).Unix() {
	// 		isRollback = true
	// 	}
	// }

	if isRollback {
		return orderSingle_RollbackPay(db, tx, oldOrder, newOrder, headParams, req, out)
	}

	return orderSingle_PayComplete(db, tx, oldOrder, newOrder, headParams, req, out)
}

func orderSingle_CustomComment_More(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {

	order_More_NotifyMsg(db, tx, oldOrder, newOrder, headParams, req, out)

	orderCommentCtrl := NewUserOrderCommentController(db, tx)

	where := map[string]interface{}{"order_id": oldOrder.OrderId}
	rule := map[string]string{"order_id": " = "}

	ct, err := sqlutils.Sqls_GetCounts(db, orderCommentCtrl.getTableName(), orderCommentCtrl.getCommentIdFlag(), where, rule)
	if err != nil {
		log.Println(err)
		return err
	}

	if ct >= 10 {
		log.Println("too much comments")
		return errors.New("too much comments")
	}

	commentinfo := ordermodel.NewOrderCommentInfo()
	commentinfo.OrderId = oldOrder.OrderId
	commentinfo.OrderType = oldOrder.OrderType

	if strings.Contains(headParams.ShortUrl, "customer/order") {
		commentinfo.CommentUserId = oldOrder.CustomerId
	} else {
		commentinfo.CommentUserId = oldOrder.ProviderId
	}

	//commentinfo.UserId =
	commentinfo.CustomerId = oldOrder.CustomerId
	commentinfo.ProviderId = oldOrder.ProviderId
	commentinfo.ServiceId = oldOrder.ServiceId
	commentinfo.ServiceName = oldOrder.ServiceName
	commentinfo.Comment = newOrder.Comment
	commentinfo.Star = newOrder.Star
	if commentinfo.Star < 0 {
		commentinfo.Star = 1
	}
	if commentinfo.Star > 5 {
		commentinfo.Star = 5
	}

	commentinfo.ReplyInfoNum = 0
	commentinfo.ReplyUpNum = 0
	commentinfo.ReplyDownNum = 0

	dt := time.Now().Format("2006-01-02 15:04:05")
	commentinfo.CreateTime = dt
	commentinfo.UpdateTime = dt

	err = orderCommentCtrl.addComment(commentinfo)
	if err != nil {
		return err
	}

	fakeren := &rendermodel.FakeMrtiniRender{}

	if strings.Contains(headParams.ShortUrl, "customer/order") {
		if oldOrder.OrderType == orderutils.Order_Type_Course {
			res := courseups.UpdateUserCourseCommentNum(db, tx, commentinfo.Star, oldOrder.ServiceId, fakeren)
			if !res {
				log.Println("UpdateUserCourseCommentNum failed")
				errrsp, erok := fakeren.GetVal().(*errcode.ErrRsp)
				if erok {
					return errors.New(errrsp.ErrDetail)
				}
				return errors.New(" update course comment num failed!")
			}
		} else {
			res := serveups.UpdateUserServiceCommentNum(db, tx, commentinfo.Star, oldOrder.ServiceId, fakeren)
			if !res {
				log.Println("UpdateUserServiceCommentNum failed")
				errrsp, erok := fakeren.GetVal().(*errcode.ErrRsp)
				if erok {
					return errors.New(errrsp.ErrDetail)
				}
				return errors.New(" update service comment num failed!")
			}
		}
	}

	destuser := oldOrder.CustomerId
	if strings.Contains(headParams.ShortUrl, "customer/order") {
		destuser = oldOrder.ProviderId
	}

	res := userups.UpdateUserCommentNum(db, tx, commentinfo.Star, destuser, fakeren)
	if !res {
		log.Println("UpdateUserCommentNum failed")
		errrsp, erok := fakeren.GetVal().(*errcode.ErrRsp)
		if erok {
			return errors.New(errrsp.ErrDetail)
		}
		return errors.New(" update user comment num failed!")
	}

	return nil
}
