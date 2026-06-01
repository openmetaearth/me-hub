package types

var (
	// ModuleName defines the module name.
	ModuleName = "denommetadata"

	// RouterKey is the message route for the denommetadata module
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName

	// EventTypeCreateDenomMetadataFailed is emitted when metadata persistence fails after a successful receive.
	EventTypeCreateDenomMetadataFailed = "create_denom_metadata_failed"

	// AttributeKeyDenom defines an event attribute for a denom.
	AttributeKeyDenom = "denom"
	// AttributeKeyError defines an event attribute for an error message.
	AttributeKeyError = "error"
)
