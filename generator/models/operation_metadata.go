package models

type OperationMetaData struct {
	// TODO: make a higher level abstraction etc
	Name                 string
	Method               string
	LongRunningOperation bool
	ExpectedStatusCodes  []int
}
