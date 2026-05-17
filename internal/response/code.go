package response

const (
	CodeSuccess        = 0
	CodeParameterError = 1001

	CodeProductParameterError  = 1102
	CodeCreateProductFailed    = 1103
	CodeProductNotFound        = 1104
	CodeProductOnsaleFailed    = 1105
	CodeProductOffsaleFailed   = 1106
	CodeQueryProductFailed     = 1107
	CodeQueryProductListFailed = 1108
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
	CodeInsufficientStock        = 4001
	CodeCreateOrderFailed        = 4002
	CodeOrderNotFound            = 4003
	CodeOrderPayConflict         = 4004
	CodeOrderPayFailed           = 4005
	CodeOrderFinishConflict      = 4006
	CodeOrderFinishFailed        = 4007
	CodeOrderCancelConflict      = 4008
	CodeOrderCancelFailed        = 4009
	CodeQueryOrderListFailed     = 4010
	CodeQueryOrderDetailNotFound = 4011
)
