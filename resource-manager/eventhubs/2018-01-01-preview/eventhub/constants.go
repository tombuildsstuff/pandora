package eventhub

type Encoding string

var (
	Avro        Encoding = "Avro"
	AvroDeflate Encoding = "AvroDeflate"
)

type EntityStatus string

var (
	Active          EntityStatus = "Active"
	Creating        EntityStatus = "Creating"
	Deleting        EntityStatus = "Deleting"
	Disabled        EntityStatus = "Disabled"
	ReceiveDisabled EntityStatus = "ReceiveDisabled"
	Renaming        EntityStatus = "Renaming"
	Restoring       EntityStatus = "Restoring"
	SendDisabled    EntityStatus = "SendDisabled"
	Unknown         EntityStatus = "Unknown"
)
