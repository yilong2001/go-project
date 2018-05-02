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
	//"strings"
	//"time"
	//"reflect"
	//"web/component/cfgutils"
	//"web/component/errcode"
	//"web/component/idutils"
	//"web/component/objutils"
	"web/component/orderutils"
	//"web/component/rongcloud"
	//"web/component/sqlutils"
	//"web/dal/sqldrv"
	//"web/models/basemodel"
	"web/models/ordermodel"
	"web/models/rendermodel"
	"web/models/reqparamodel"
	//"web/service/coupons"
	//"web/service/getter"
	//"web/service/pays"
	//"web/service/serveups"
	//"web/service/routers"
	//"web/service/immsgs"
	"web/service/orderpays"
	"web/service/videos"
	//"web/service/userups"
	//"web/service/utils"
)

var globalCourseOnlinePurchaseOrderStatusTransferMachines []*ordermodel.OrderStatusTransfer = []*ordermodel.OrderStatusTransfer{

	//customer prepay
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Course,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Wait_Comment,
		Router:        "prepayer",
		ChangedFields: []string{"PrepayType", "PrepayStatus", "PrepayMoney", "PrepayTime", "UserCouponId01", "PayCouponId01", "PayCouponMoney01", "OrderStatus", "UpdateTime"},

		PayFunction:  order_OnlineCourse_PayComplete,
		MoreFunction: order_OnlineCourse_Complete_More,
	},

	//customer cancel order
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Course,
		CurrentStatus: orderutils.Order_Status_Wait_Customer_PrePay,
		NextStatus:    orderutils.Order_Status_Customer_Cancel,
		Router:        "cancel",
		ChangedFields: []string{"OverReason", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  nil,
		MoreFunction: order_More_NotifyMsg,
	},

	//customer comment
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Course,
		CurrentStatus: orderutils.Order_Status_Wait_Comment,
		NextStatus:    orderutils.Order_Status_Complete,
		Router:        "comment",
		ChangedFields: []string{"Comment", "Star", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  nil,
		MoreFunction: orderSingle_CustomComment_More,
	},

	//customer add comment
	&ordermodel.OrderStatusTransfer{
		Role:          orderutils.Order_Role_Customer,
		OrderType:     orderutils.Order_Type_Course,
		CurrentStatus: orderutils.Order_Status_Complete,
		NextStatus:    orderutils.Order_Status_Complete,
		Router:        "comment/plus",
		ChangedFields: []string{"Comment", "Star", "OrderStatus", "UpdateTime", "OverTime"},

		PayFunction:  nil,
		MoreFunction: orderSingle_CustomComment_More,
	},
}

func order_OnlineCourse_PayComplete(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {

	log.Println("orderCourse_PrePay, newOrder.PrepayType", newOrder.PrepayType)

	fakeren := &rendermodel.FakeMrtiniRender{}

	if oldOrder.OrderPrice == 0 {
		return nil
	}

	if newOrder.UserCouponId01 > 0 && oldOrder.OrderPrice > 0 {
		err := orderpays.Order_PayCompleteWithCoupon(db, tx, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	}

	if newOrder.PrepayType == orderutils.Pay_Type_Account_Balance {
		err := orderpays.Order_PayWithAccountBalance(db, tx, newOrder.PrepayMoney, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			log.Print("orderCourse_PrePay", err)
			return err
		}

		err = orderpays.Order_PayCompleteToProvider(db, tx, oldOrder.OrderPrice, oldOrder, newOrder, headParams, req, fakeren)
		if err != nil {
			return err
		}
	} else if newOrder.PrepayType == orderutils.Pay_Type_WeiXin {
		err := orderpays.Order_PayWithWxPay(db, tx, newOrder.PrepayMoney, oldOrder, newOrder, headParams, req, fakeren, out)
		if err != nil {
			log.Print("orderCourse_PrePay, wxpay wrong : ", err)
			return err
		}
	} else if newOrder.PrepayType == orderutils.Pay_Type_ZhiFuBao {
		err := orderpays.Order_PayWithAliPay(db, tx, newOrder.PrepayMoney, oldOrder, newOrder, headParams, req, fakeren, out)
		if err != nil {
			log.Print("orderCourse_PrePay, alipay wrong : ", err)
			return err
		}
	} else {
		return errors.New("do't support 这种支付类型!")
	}

	//call zhifubao, weixin pay interface
	return nil
}

func order_OnlineCourse_Complete_More(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {

	fakeren := &rendermodel.FakeMrtiniRender{}
	err := videos.PurchaseCourseVideo(db, tx, oldOrder.ServiceId, oldOrder.CustomerId, headParams, fakeren)
	if err != nil {
		log.Println("add user video path is wrong : ", err)
		return err
	}

	order_PayComplete_More(db, tx, oldOrder, newOrder, headParams, req, out)

	return nil
}
