package response

const (
	CodeSuccess        = 0
	CodeParameterError = 1001

	CodeProductParameterError = 1002

	CodeCreateProductFailed  = 1003
	CodeProductNotFound      = 1004
	CodeProductOnsaleFailed  = 1005
	CodeProductOffsaleFailed = 1006

	CodeInitInventoryExists      = 2001
	CodeCreateStockLogFailed     = 2002
	CodeInitInventoryFailed      = 2003
	CodeInvalidAddQuantityFailed = 2004
	CodeInventoryNotFound        = 2005
	CodeAddInventoryError        = 2006

	CodeStockLogNotFound = 3001

	CodeInsufficientStock        = 4001
	CodeCreateOrderFailed        = 4002
	CodeOrderNotFound            = 4003
	CodeOrderPayFailed           = 4004
	CodePayOrderFailed           = 4005
	CodeOrderFinishFailed        = 4006
	CodeFinishOrderUnknownFailed = 4007
	CodeOrderCancelFailed        = 4008
	CodeCancelOrderUnknownFailed = 4009
)
