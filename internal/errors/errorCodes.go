package errors

const prefix = "orders_"

const UnexpectedErrorMessage = "unexpected Error occurred, please try again later"

const (
	OrderGetInvalidParams = prefix + "get_invalid_params"
	OrderGetNotFound      = prefix + "get_not_found"
	OrdersGetServerError  = prefix + "get_server_error"

	OrderCreateInvalidInput = prefix + "create_invalid_input"
	OrderCreateServerError  = prefix + "create_server_error"

	OrderDeleteInvalidID   = prefix + "delete_invalid_order_id"
	OrderDeleteNotFound    = prefix + "delete_not_found"
	OrderDeleteServerError = prefix + "delete_server_error"
)
