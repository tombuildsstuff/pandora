package templates

type TemplateBuilder interface {
	Build() (*string, error)
}
