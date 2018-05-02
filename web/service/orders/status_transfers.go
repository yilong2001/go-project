package orders

import (
	"database/sql"
	"errors"
	//"fmt"
	//"github.com/martini-contrib/render"
	"log"
	"net/http"
	"strings"
	"web/component/orderutils"
	//"web/component/wxutils"
	// "web/dal/sqldrv"
	// "web/models/basemodel"
	"web/models/coursemodel"
	"web/models/ordermodel"
	"web/models/rendermodel"
	"web/models/reqparamodel"
	//"web/service/coupons"
	"web/service/immsgs"
	//"web/service/pays"
	"web/service/courseups"
	"web/service/getter"
	"web/service/serveups"
	//"web/service/userups"
)

func getStatusTransferMachines( /*ordertype int*/ orderinfo *ordermodel.OrderInfo) []*ordermodel.OrderStatusTransfer {

	ordertype := orderinfo.OrderType
	svrid := orderinfo.ServiceId

	switch ordertype {
	case orderutils.Order_Type_Single_Business:
		return globalSignleBusinessOrderStatusTransferMachines

	case orderutils.Order_Type_With_Front_Money_Business:
		return globalWithFrontMoneyOrderStatusTransferMachines

	case orderutils.Order_Type_Course:
		cm, err := getter.GetModelInfoGetter().GetCourseMainByCourseId(nil, svrid)
		if err != nil {
			return nil
		}

		switch cm.CourseType {
		case coursemodel.Const_Course_Type_Online:
			return globalCourseOnlinePurchaseOrderStatusTransferMachines
		}

		return globalCoursePurchaseOrderStatusTransferMachines
	}

	return globalSignleBusinessOrderStatusTransferMachines
}

func GetStatusTransferMachine(ordertype, orderstatus, role int, url string, orderinfo *ordermodel.OrderInfo) *ordermodel.OrderStatusTransfer {
	log.Println("GetStatusTransferMachine, in: ", ordertype, orderstatus, role, url)

	statusTransfers := getStatusTransferMachines(orderinfo)
	if statusTransfers == nil {
		log.Println(orderinfo.ServiceId, ", is wrong course id!")
		return nil
	}

	var destMachine *ordermodel.OrderStatusTransfer = nil
	for _, statusMachine := range statusTransfers {
		log.Println(statusMachine.CurrentStatus, statusMachine.Role, statusMachine.Router)
		if statusMachine.CurrentStatus == orderstatus &&
			statusMachine.Role == role &&
			strings.HasSuffix(url, statusMachine.Router) {
			destMachine = statusMachine
			break
		}
	}

	return destMachine
}

func order_More_NotifyMsg(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {

	id := oldOrder.CustomerId
	if strings.Contains(headParams.ShortUrl, "customer") {
		id = oldOrder.ProviderId
	}

	err := immsgs.SendPlatformNotifyMsg(id, "Order", "您的订单更新啦...", oldOrder.OrderId, oldOrder.OrderType)

	log.Println("order, more for update ", err)

	return nil
}

func order_PayComplete_More(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, out *map[string]interface{}) error {
	fakeren := &rendermodel.FakeMrtiniRender{}

	if oldOrder.OrderType == orderutils.Order_Type_Course {
		res := courseups.AddNumOfUserCourseServicedNum(db, tx, oldOrder.ServiceId, fakeren)
		if !res {
			if !res {
				return errors.New("AddNumOfUserCourseServicedNum failed")
			}
		}
	} else {
		res := serveups.AddNumOfServicedNum(db, tx, oldOrder.ServiceId, fakeren)
		if !res {
			return errors.New("AddNumOfServicedNum failed")
		}
	}

	// res = userups.AddServedNumOfUserSelf(db, tx, oldOrder.ProviderId, headParams, req, fakeren)
	// if !res {
	// 	return errors.New("AddServedNumOfUserSelf failed")
	// }

	order_More_NotifyMsg(db, tx, oldOrder, newOrder, headParams, req, nil)

	return nil
}
