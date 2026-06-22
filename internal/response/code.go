package response

const (
	CodeSuccess        = 0
	CodeParameterError = 1001

	CodeCreateProductFailed    = 1103
	CodeProductNotFound        = 1104
	CodeProductOnSaleFailed    = 1105
	CodeProductAlreadyOffSale  = 1106
	CodeProductOffSaleFailed   = 1107
	CodeQueryProductFailed     = 1108
	CodeQueryProductListFailed = 1109
)

const (
	CodeInitInventoryExists      = 2001
	CodeCreateStockLogFailed     = 2002
	CodeInitInventoryFailed      = 2003
	CodeInventoryInvalidQuantity = 2004
	CodeInventoryNotFound        = 2005
	CodeAddInventoryError        = 2006
	CodeQueryInventoryFailed     = 2007
)

const (
	CodeQueryStockLogFailed = 3001
)

const (
	CodeInsufficientStock      = 4001
	CodeCreateOrderFailed      = 4002
	CodeOrderNotFound          = 4003
	CodeOrderPayFailed         = 4004
	CodeOrderFinishFailed      = 4005
	CodeOrderCancelFailed      = 4006
	CodeQueryOrderListFailed   = 4007
	CodeQueryOrderDetailFailed = 4008
	CodeOrderNotPaid           = 4009
	CodeOrderAlreadyCanceled   = 4010
	CodeOrderAlreadyFinished   = 4011
	CodeOrderAlreadyPaid       = 4012
	CodeOrderParameterError    = 4013
)

const (
	CodeInternalServerError = 5000
)
