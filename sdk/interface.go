package sdk

type ModelWithValidation interface {
	Validate() error
}
