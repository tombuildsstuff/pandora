package eventhub

type Encoding string

var (
	Avro        Encoding = "Avro"
	AvroDeflate Encoding = "AvroDeflate"
)

type EntityStatus string

var (
	ReceiveDisabled EntityStatus = "ReceiveDisabled"
	Renaming        EntityStatus = "Renaming"
	Disabled        EntityStatus = "Disabled"
	Creating        EntityStatus = "Creating"
	Deleting        EntityStatus = "Deleting"
	Restoring       EntityStatus = "Restoring"
	SendDisabled    EntityStatus = "SendDisabled"
	Unknown         EntityStatus = "Unknown"
	Active          EntityStatus = "Active"
)
