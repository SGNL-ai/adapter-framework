package logs

// Log field constants.
const (
	FieldClientID               = "clientId"
	FieldDatasourceAddress      = "datasourceAddress"
	FieldDatasourceID           = "datasourceId"
	FieldDatasourceType         = "datasourceType"
	FieldEntityExternalID       = "entityExternalId"
	FieldEntityID               = "entityId"
	FieldAdapterRequestCursor   = "adapterRequestCursor" // Prefix with "adapter" because there may be multiple page cursors used in a single adapter request.
	FieldAdapterRequestPageSize = "adapterRequestPageSize"
	FieldTenantID               = "tenantId"
)

// ClientID returns a log field for the client ID.
func ClientID(value string) Field {
	return Field{Key: FieldClientID, Value: value}
}

// DatasourceAddress returns a log field for the datasource address.
func DatasourceAddress(value string) Field {
	return Field{Key: FieldDatasourceAddress, Value: value}
}

// DatasourceID returns a log field for the datasource ID.
func DatasourceID(value string) Field {
	return Field{Key: FieldDatasourceID, Value: value}
}

// DatasourceType returns a log field for the datasource type.
func DatasourceType(value string) Field {
	return Field{Key: FieldDatasourceType, Value: value}
}

// EntityExternalID returns a log field for the entity external ID.
func EntityExternalID(value string) Field {
	return Field{Key: FieldEntityExternalID, Value: value}
}

// EntityID returns a log field for the entity ID.
func EntityID(value string) Field {
	return Field{Key: FieldEntityID, Value: value}
}

// RequestCursor returns a log field for the request cursor.
func RequestCursor(value any) Field {
	return Field{Key: FieldAdapterRequestCursor, Value: value}
}

// RequestPageSize returns a log field for the request page size.
func RequestPageSize(value int64) Field {
	return Field{Key: FieldAdapterRequestPageSize, Value: value}
}

// TenantID returns a log field for the tenant ID.
func TenantID(value string) Field {
	return Field{Key: FieldTenantID, Value: value}
}
