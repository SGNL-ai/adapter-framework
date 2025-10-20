package logs

import "go.uber.org/zap"

// Log field constants.
const (
	FieldClientID         = "clientId"
	FieldDatasourceID     = "datasourceId"
	FieldDatasourceType   = "datasourceType"
	FieldEntityExternalID = "entityExternalId"
	FieldEntityID         = "entityId"
	FieldRequestCursor    = "requestCursor"
	FieldRequestPageSize  = "requestPageSize"
	FieldTenantID         = "tenantId"
)

func ClientID(value string) zap.Field {
	return zap.String(FieldClientID, value)
}

func DatasourceID(value string) zap.Field {
	return zap.String(FieldDatasourceID, value)
}

func DatasourceType(value string) zap.Field {
	return zap.String(FieldDatasourceType, value)
}

func EntityExternalID(value string) zap.Field {
	return zap.String(FieldEntityExternalID, value)
}

func EntityID(value string) zap.Field {
	return zap.String(FieldEntityID, value)
}

func RequestCursor(value string) zap.Field {
	return zap.String(FieldRequestCursor, value)
}

func RequestPageSize(value int64) zap.Field {
	return zap.Int64(FieldRequestPageSize, value)
}

func TenantID(value string) zap.Field {
	return zap.String(FieldTenantID, value)
}
