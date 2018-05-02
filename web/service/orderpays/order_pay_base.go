package orderpays

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/martini-contrib/render"
	//"log"
	"net/http"
	//"strings"
	//"web/component/aliutils"
	"web/component/orderutils"
	//"web/component/wxutils"
	// "web/dal/sqldrv"
	// "web/models/basemodel"
	"web/models/ordermodel"
	//"web/models/rendermodel"
	"web/models/reqparamodel"
	"web/service/coupons"
	//"web/service/getter"
	//"web/service/immsgs"
	"web/service/pays"
	//"web/service/serveups"
	"web/service/userups"
)

func Order_PrepayWithCoupon(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render) error {
	ok := coupons.UpdateUserCouponStatus(db, tx, oldOrder.CustomerId, newOrder.UserCouponId01, orderutils.User_Coupon_Status_Lock, orderutils.User_Coupon_Status_Availabe, headParams, req, fakeren)
	if !ok {
		return errors.New("UpdateUserCouponStatus failed")
	}

	err := pays.AddPayRecord(db, tx, oldOrder.CustomerId, orderutils.GetSystmePayerId(), fmt.Sprintf("order id %d, user %d to %d, prepay with coupon %d", oldOrder.OrderId, oldOrder.CustomerId, oldOrder.ProviderId, newOrder.UserCouponId01))
	if err != nil {
		return err
	}

	return nil
}

func Order_PayWithAccountBalance(db *sql.DB, tx *sql.Tx, cost int, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render) error {
	if cost == 0 {
		return nil
	}

	ok := userups.DecAccountBanlancOfUser(db, tx, oldOrder.CustomerId, cost, headParams, req, fakeren)
	if !ok {
		return errors.New("DecAccountBanlancOfUser dec failed")
	}

	err := pays.AddPayRecord(db, tx, oldOrder.CustomerId, oldOrder.ProviderId, fmt.Sprintf("order id %d, user %d to %d, pay with account balance -%d, old order status %d", oldOrder.OrderId, oldOrder.CustomerId, oldOrder.ProviderId, cost, oldOrder.OrderStatus))
	if err != nil {
		return err
	}
	return nil
}

func Order_RollbackWithCoupon(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render) error {
	ok := coupons.UpdateUserCouponStatus(db, tx, oldOrder.CustomerId, oldOrder.UserCouponId01, orderutils.User_Coupon_Status_Availabe, orderutils.User_Coupon_Status_Lock, headParams, req, fakeren)
	if !ok {
		return errors.New("UpdateUserCouponStatus failed")
	}

	err := pays.AddPayRecord(db, tx, orderutils.GetSystmePayerId(), oldOrder.CustomerId, fmt.Sprintf("order id %d, user %d pay rollback with coupon %d", oldOrder.OrderId, oldOrder.CustomerId, oldOrder.UserCouponId01))
	if err != nil {
		return err
	}

	return nil
}

func Order_RollbackWithAccountBalance(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render) error {
	if oldOrder.AccountLocked == 0 {
		return nil
	}

	ok := userups.AddAccountBanlancOfUser(db, tx, oldOrder.CustomerId, oldOrder.AccountLocked, false, headParams, req, fakeren)
	if !ok {
		return errors.New("AddAccountBanlancOfUser dec failed")
	}

	err := pays.AddPayRecord(db, tx, oldOrder.CustomerId, oldOrder.CustomerId, fmt.Sprintf("order id %d, %d pay rollback with account balance %d", oldOrder.OrderId, oldOrder.CustomerId, oldOrder.AccountLocked))
	if err != nil {
		return err
	}
	return nil
}

func Order_PayCompleteWithCoupon(db *sql.DB, tx *sql.Tx, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render) error {
	ok := coupons.UpdateUserCouponStatus(db, tx, oldOrder.CustomerId, oldOrder.UserCouponId01, orderutils.User_Coupon_Status_Used, orderutils.User_Coupon_Status_Lock, headParams, req, fakeren)
	if !ok {
		return errors.New("UpdateUserCouponStatus failed")
	}

	err := pays.AddPayRecord(db, tx, oldOrder.CustomerId, orderutils.GetSystmePayerId(), fmt.Sprintf("order id %d, user %d to %d, coupon %d has been used", oldOrder.OrderId, oldOrder.CustomerId, oldOrder.ProviderId, oldOrder.UserCouponId01))
	if err != nil {
		return err
	}

	return nil
}

func Order_PayCompleteToProvider(db *sql.DB, tx *sql.Tx, latestCost int, oldOrder, newOrder *ordermodel.OrderInfo, headParams *reqparamodel.HttpReqParams,
	req *http.Request, fakeren render.Render) error {
	ok := userups.AddAccountBanlancOfUser(db, tx, oldOrder.ProviderId, latestCost, true, headParams, req, fakeren)
	if !ok {
		return errors.New("AddAccountBanlancOfUser failed")
	}

	err := pays.AddPayRecord(db, tx, oldOrder.CustomerId, oldOrder.ProviderId, fmt.Sprintf("order id %d, user %d pay to user %d with cost +%d, within account lock %d, prepay money %d", oldOrder.OrderId, oldOrder.CustomerId, oldOrder.ProviderId, latestCost, oldOrder.AccountLocked, oldOrder.PrepayMoney))
	if err != nil {
		return err
	}

	return nil
}
