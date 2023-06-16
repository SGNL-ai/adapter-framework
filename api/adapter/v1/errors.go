package v1

const (
	// Internal Error (ErrorCode_ERROR_CODE_INTERNAL).
	ErrorMsgFailedToConnect           = "Failed to connect to datasource"
	ErrorMsgFailedToReadResponse      = "Failed to read response from datasource"
	ErrorMsgUnhandledStatusCode       = "Unhandled status code received from datasource"
	ErrorMsgUnexpectedErrorCode       = "Unexpected error code received from datasource"
	ErrorMsgDatasourceRejectedRequest = "Datasource rejected request"

	// Invalid Page Request Config (ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG).
	ErrorMsgNoRequest       = "Request is nil"
	ErrorNoDatasourceConfig = "No datasource config provided"
	ErrorNoEntityConfig     = "No entity config provided"
	ErrorNoAuth             = "No datasource authentication credentials provided"
	ErrorMsgInvalidCursor   = "Invalid cursor provided"

	// Invalid Entity Config (ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG).
	ErrorMsgEntityDoesNotSupportChildEntities      = "Entity does not support child entities"
	ErrorMsgEntityProvidedWithNoAttributes         = "Provided entity attribute list is empty"
	ErrorMsgEntityMissingUniqueIdAttribute         = "Provided entity attribute list is missing a required unique id attribute"
	ErrorMsgEntityMissingRequiredAttributeFmt      = "Provided entity attribute list is missing the following required attributes: %s"
	ErrorMsgEntityUnsupportedAttributesProvidedFmt = "Provided entity attribute list contains the following unsupported attributes: %s"
	ErrorMsgEntityPageSizeTooSmall                 = "Invalid page size provided; value must be greater than 0"

	// Invalid Datasource Config (ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG).
	ErrorMsgDatasourceHTTPAuthEmpty      = "Provided Datasource HTTP Auth is empty"
	ErrorMsgDatasourceBasicAuthEmpty     = "Provided Datasource Basic Auth is empty"
	ErrorMsgEntityMustBeOrdered          = "Entity must be ordered"
	ErrorMsgEntityMustBeUnordered        = "Entity must be unordered"
	ErrorMsgInvalidAddressFormatProvided = "Invalid address format provided"

	// Datasource Auth Failed (ErrorCode_ERROR_CODE_DATASOURCE_AUTHENTICATION_FAILED).
	ErrorMsgFailedToAuthenticate = "Failed to authenticate with datasource; check datasource configuration details and try again"

	// Datasource Failed (ErrorCode_ERROR_CODE_DATASOURCE_FAILED).
	ErrorMsgDatasourceInternalError = "Datasource encountered an internal error"

	// Datasource Temporary Unavailable (ErrorCode_ERROR_CODE_DATASOURCE_TEMPORARILY_UNAVAILABLE).
	ErrorMsgDatasourceTemporarilyUnavailable = "Datasource temporarily unavailable; try again later"

	// Datasource Permanently Unavailable (ErrorCode_ERROR_CODE_DATASOURCE_PERMANENTLY_UNAVAILABLE).
	ErrorMsgDatasourcePermanentlyUnavailable = "Datasource permanently unavailable; check datasource configuration details or contact datasource support for assistance"
)
