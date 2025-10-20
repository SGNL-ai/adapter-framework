package logs

import "go.uber.org/zap"

// Log field constants.
const (
	FieldDatasourceId     = "datasourceId"
	FieldDatasourceType   = "datasourceType"
	FieldEntityExternalId = "entityExternalId"
	FieldEntityID         = "entityId"
	FieldRequestCursor    = "requestCursor"
	FieldRequestPageSize  = "requestPageSize"
	FieldTenantID         = "tenantId"
)

func DatasourceID(value string) zap.Field {
	return zap.String(FieldDatasourceId, value)
}

func DatasourceType(value string) zap.Field {
	return zap.String(FieldDatasourceType, value)
}

func EntityExternalID(value string) zap.Field {
	return zap.String(FieldEntityExternalId, value)
}

func EntityID(value string) zap.Field {
	return zap.String(FieldEntityID, value)
}

func RequestCursor(value string) zap.Field {
	return zap.String(FieldRequestCursor, value)
}

func RequestPageSize(value int) zap.Field {
	return zap.Int(FieldRequestPageSize, value)
}

func TenantID(value string) zap.Field {
	return zap.String(FieldTenantID, value)
}
