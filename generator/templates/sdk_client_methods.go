package templates

import (
	"fmt"
	"strings"
)

type methodTemplater struct {
	clientName string
	typeName   string
	operation  ClientOperation
}

func (t methodTemplater) Build() (*string, error) {
	var result string
	switch strings.ToUpper(t.operation.Method) {
	case "DELETE":
		{
			if t.operation.RequestObjectName != nil {
				// TODO: maybe implement this
				return nil, fmt.Errorf("`DELETE` operations do not support Request objects at this time")
			}

			if t.operation.LongRunningOperation {
				result = t.deleteLongRunningOperation()
				break
			}

			result = t.delete()
			break
		}

	case "GET":
		{
			if t.operation.LongRunningOperation {
				return nil, fmt.Errorf("`GET` operations cannot be long-running")
			}

			if t.operation.RequestObjectName != nil {
				// TODO: implement support for this
				return nil, fmt.Errorf("`GET` operations do not support Request objects at this time")
			}

			if t.operation.ResponseObjectName == nil {
				return nil, fmt.Errorf("`GET` operations must have a Response object")
			}

			result = t.get()
			break
		}

	case "PATCH":
		{
			if t.operation.RequestObjectName == nil {
				return nil, fmt.Errorf("`PATCH` operations must have a Request Object")
			}

			if t.operation.LongRunningOperation {
				result = t.patchLongRunningOperation()
				break
			}

			result = t.patch()
			break
		}

	case "PUT":
		{
			if t.operation.RequestObjectName == nil {
				return nil, fmt.Errorf("`PUT` operations must have a Request Object")
			}

			if t.operation.LongRunningOperation {
				result = t.putLongRunningOperation()
				break
			}

			result = t.put()
			break
		}

	default:
		return nil, fmt.Errorf("unsupported method type %q..", t.operation.Method)
	}

	result = strings.TrimSpace(result)
	return &result, nil
}

func (t methodTemplater) delete() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]s) %[2]s(ctx context.Context, id %[3]sId) (*http.Response, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
%[4]s
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}
	
	return client.baseClient.Delete(ctx, req);
}
`, t.clientName, t.operation.Name, t.typeName, statusCodes)
}

func (t methodTemplater) deleteLongRunningOperation() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]s) %[2]s(ctx context.Context, id %[3]sId) (sdk.Poller, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
%[4]s,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	return client.baseClient.DeleteThenPoll(ctx, req)
}
`, t.clientName, t.operation.Name, t.typeName, statusCodes)
}

func (t methodTemplater) get() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
type %[2]s%[3]sResponse struct {
	HttpResponse *http.Response
	%[3]s    *%[4]s
}

func (client %[1]s) %[2]s(ctx context.Context, id %[3]sId) (*%[2]s%[3]sResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
%[5]s
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	var out %[4]s
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %%+v", err)
	}

	result := %[2]s%[3]sResponse{
		HttpResponse: resp,
		%[3]s:    &out,
	}
	return &result, nil
}
`, t.clientName, t.operation.Name, t.typeName, *t.operation.ResponseObjectName, statusCodes)
}

func (t methodTemplater) patch() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]s) %[2]s(ctx context.Context, id %[3]sId, input %[4]s) error {
	req := sdk.PatchHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
%[5]s
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}
	
	if _, err := client.baseClient.PatchJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %%+v", err)
	}
	return nil
}
`, t.clientName, t.operation.Name, t.typeName, *t.operation.RequestObjectName, statusCodes)
}

func (t methodTemplater) patchLongRunningOperation() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]s) %[2]s(ctx context.Context, id %[3]sId, input %[4]s) (sdk.Poller, error) {
	req := sdk.Patch%[1]sInput{
		Body: input,
		ExpectedStatusCodes: []int{
%[5]s,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	return client.baseClient.PatchJsonThenPoll(ctx, req)
}
`, t.clientName, t.operation.Name, t.typeName, *t.operation.RequestObjectName, statusCodes)
}

func (t methodTemplater) put() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]s) %[2]s(ctx context.Context, id %[3]sId, input %[4]s) error {
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
%[5]s
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}
	
	if _, err := client.baseClient.PutJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %%+v", err)
	}
	return nil
}
`, t.clientName, t.operation.Name, t.typeName, *t.operation.RequestObjectName, statusCodes)
}

func (t methodTemplater) putLongRunningOperation() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]s) %[2]s(ctx context.Context, id %[3]sId, input %[4]s) (sdk.Poller, error) {
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
%[5]s,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	return client.baseClient.PutJsonThenPoll(ctx, req)
}
`, t.clientName, t.operation.Name, t.typeName, *t.operation.RequestObjectName, statusCodes)
}

func (t methodTemplater) statusCodes(indentation string) string {
	output := make([]string, 0)

	for _, statusCode := range t.operation.ExpectedStatusCodes {
		alias := golangConstantForStatusCode(statusCode)
		description := descriptionForStatusCodeForMethod(statusCode, t.operation.Method, t.operation.LongRunningOperation)
		formatted := fmt.Sprintf("%s%s, // %s", indentation, alias, description)
		output = append(output, formatted)
	}

	return strings.Join(output, "\n")
}

func descriptionForStatusCodeForMethod(code int, method string, longRunningOperation bool) string {
	var knownStatusCodes map[int]string

	if strings.EqualFold(method, "delete") {
		if longRunningOperation {
			knownStatusCodes = map[int]string{
				200: "deletion started",
				202: "deletion accepted",
			}
		} else {
			knownStatusCodes = map[int]string{
				200: "deleted",
				204: "deleted / gone",
			}
		}
	}

	if strings.EqualFold(method, "get") {
		knownStatusCodes = map[int]string{
			200: "ok",
		}
	}

	// TODO: others

	v, ok := knownStatusCodes[code]
	if ok {
		return v
	}
	return "TODO: unknown"
}

func golangConstantForStatusCode(statusCode int) string {
	codes := map[int]string{
		200: "http.StatusOK",
		201: "http.StatusCreated",
		202: "http.StatusAccepted",
		204: "http.StatusNoContent",
		301: "http.StatusMovedPermanently",
		302: "http.StatusFound",
		307: "http.StatusTemporaryRedirect",
		308: "http.StatusPermanentRedirect",
		400: "http.StatusBadRequest",
		401: "http.StatusUnauthorized",
		403: "http.StatusForbidden",
		404: "http.StatusNotFound",
		405: "http.StatusMethodNotAllowed",
		406: "http.StatusNotAcceptable",
		407: "http.StatusProxyAuthRequired",
		408: "http.StatusRequestTimeout",
		409: "http.StatusConflict",
		410: "http.StatusGone",
		429: "http.StatusTooManyRequests",
		500: "http.StatusInternalServerError",
		501: "http.StatusNotImplemented",
		502: "http.StatusBadGateway",
		503: "http.StatusServiceUnavailable",
		504: "http.StatusGatewayTimeout",
	}
	v, ok := codes[statusCode]
	if ok {
		return v
	}

	return fmt.Sprintf("%d // TODO: document me", statusCode)
}
