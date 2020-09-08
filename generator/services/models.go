package services

type ResourceManagerApiResponse struct {
	Apis map[string]ApiDetails `json:"apis"`
}

type ApiDetails struct {
	Uri      string `json:"uri"`
	Generate bool   `json:"generate"`
}

type SupportedVersionsResponse struct {
	Versions         map[string]VersionDetails `json:"versions"`
	ResourceProvider string                    `json:"resourceProvider"`
}

type VersionDetails struct {
	Uri      string `json:"uri"`
	Generate bool   `json:"generate"`
	Preview  bool   `json:"preview"`
}

type SupportedTypesResponse struct {
	Types map[string]TypeDefinition `json:"types"`
}

type TypeDefinition struct {
	Uri        string               `json:"uri"`
	ResourceId ResourceIdDefinition `json:"resourceId"`
}

type ResourceIdDefinition struct {
	Format   string   `json:"format"`
	Segments []string `json:"segments"`
}

type OperationsResponse struct {
	ResourceManager *ResourceManagerMetaData       `json:"resourceManager,omitempty"`
	Operations      map[string]OperationDefinition `json:"operations"`
}

type ResourceManagerMetaData struct {
	ApiVersion       string               `json:"apiVersion"`
	ResourceProvider string               `json:"resourceProvider"`
	ResourceId       ResourceIdDefinition `json:"resourceId"`
}

type OperationDefinition struct {
	Method              string  `json:"method"`
	ContentType         *string `json:"contentType,omitempty"`
	ExpectedStatusCodes []int   `json:"expectedStatusCodes"`
	LongRunning         bool    `json:"longRunning"`
	RequestObject       *string `json:"requestObject,omitempty"`
	ResponseObject      *string `json:"responseObject,omitempty"`
}

type SchemaResponse struct {
	Constants map[string]ConstantDefinition `json:"constants"`
	Models    map[string]ModelDefinition    `json:"models"`
}

type ConstantDefinition struct {
	Values          []string `json:"values"`
	CaseInsensitive bool     `json:"caseInsensitive"`
	// TODO: update mappings (e.g. can go from one to another but not back again)
}

type ModelDefinition map[string]PropertyDefinition

type PropertyDefinition struct {
	Type        PropertyType          `json:"type"`
	Required    bool                  `json:"required"`
	Optional    bool                  `json:"optional"`
	DeltaUpdate bool                  `json:"deltaUpdate"`
	Default     *interface{}          `json:"default,omitempty"`
	ForceNew    bool                  `json:"forceNew,omitempty"`
	Validation  *ValidationDefinition `json:"validation,omitempty"`

	// when a constant, values can come from another reference
	ConstantReference *string `json:"constantReference,omitempty"`

	ModelReference *string `json:"modelReference,omitempty"`
}

type PropertyType string

var (
	Boolean  PropertyType = "boolean"
	Constant PropertyType = "constant"
	Integer  PropertyType = "integer"
	Location PropertyType = "location"
	Object   PropertyType = "object"
	Tags     PropertyType = "tags"
	String   PropertyType = "string"
)

type ValidationDefinition struct {
	Type   ValidationType `json:"type"`
	Values *[]interface{} `json:"values,omitempty"`
	// TODO: presumably "constant" here too in time
}

type ValidationType string

var (
	Range ValidationType = "range"
)

type OperationMetaData struct {
	Name          string `json:"name"`
	OperationsUri string `json:"operations"`
	SchemaUri     string `json:"schema"`
}
