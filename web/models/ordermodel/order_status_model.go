package ordermodel

import ()

type OrderStatusTransfer struct {
	Role          int
	OrderType     int
	CurrentStatus int
	NextStatus    int
	Router        string
	ChangedFields []string
	PayFunction   interface{}
	MoreFunction  interface{}
	// ValidChangedFields interface{}
	// IsRollbackPrepay   bool
	// IsPrepay2Provider  bool
	// DecidePrepay       interface{}
	// IsSetExpireDate    bool
}
