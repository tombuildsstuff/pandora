package templates

import (
	"fmt"
	"strings"
)

type methodTemplater struct {
	typeName             string
	name                 string
	method               string
	longRunningOperation bool
	expectedStatusCodes  []int
}

func (t methodTemplater) Build() (*string, error) {
	var result string
	switch strings.ToUpper(t.method) {
	case "DELETE":
		{
			if t.longRunningOperation {
				result = t.deleteLongRunningOperation()
				break
			}

			result = t.delete()
			break
		}

	case "GET":
		{
			if t.longRunningOperation {
				return nil, fmt.Errorf("`GET` operations cannot be long-running")
			}

			result = t.get()
			break
		}

	case "PATCH":
		{
			if t.longRunningOperation {
				result = t.patchLongRunningOperation()
				break
			}

			result = t.patch()
			break
		}

	case "PUT":
		{
			if t.longRunningOperation {
				result = t.putLongRunningOperation()
				break
			}

			result = t.put()
			break
		}

	default:
		return nil, fmt.Errorf("unsupported method type %q..", t.method)
	}

	result = strings.TrimSpace(result)
	return &result, nil
}

func (t methodTemplater) delete() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]ssClient) %[2]s(ctx context.Context, id %[1]sID) (*http.Response, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
%[3]s
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}
	
	return client.baseClient.Delete(ctx, req);
}
`, t.typeName, t.name, statusCodes)
}

func (t methodTemplater) deleteLongRunningOperation() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]ssClient) %[2]s(ctx context.Context, id %[1]sID) (sdk.Poller, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
%[3]s,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	return client.baseClient.DeleteThenPoll(ctx, req)
}
`, t.typeName, t.name, statusCodes)
}

func (t methodTemplater) get() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]ssClient) %[2]s(ctx context.Context, id %[1]sID) (*%[2]s%[1]sResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
%[3]s,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	var out %[2]s%[1]s
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %%+v", err)
	}

	result := %[2]s%[1]sResponse{
		HttpResponse: resp,
		%[1]s:    &out,
	}
	return &result, nil
}
`, t.typeName, t.name, statusCodes)
}

func (t methodTemplater) patch() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]ssClient) %[2]s(ctx context.Context, id %[1]sID, input %[2]s%[1]sInput) error {
	req := sdk.PatchHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
%[3]s
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}
	
	if _, err := client.baseClient.PatchJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %%+v", err)
	}
	return nil
}
`, t.typeName, t.name, statusCodes)
}

func (t methodTemplater) patchLongRunningOperation() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]ssClient) %[2]s(ctx context.Context, id %[1]sID, input %[2]s%[1]sInput) (sdk.Poller, error) {
	req := sdk.Patch%[1]sInput{
		Body: input,
		ExpectedStatusCodes: []int{
%[3]s,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	return client.baseClient.PatchJsonThenPoll(ctx, req)
}
`, t.typeName, t.name, statusCodes)
}

func (t methodTemplater) put() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]ssClient) %[2]s(ctx context.Context, id %[1]sID, input %[2]s%[1]sInput) error {
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
%[3]s
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}
	
	if _, err := client.baseClient.PutJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %%+v", err)
	}
	return nil
}
`, t.typeName, t.name, statusCodes)
}

func (t methodTemplater) putLongRunningOperation() string {
	statusCodes := t.statusCodes("\t\t\t")
	return fmt.Sprintf(`
func (client %[1]ssClient) %[2]s(ctx context.Context, id %[1]sID, input %[2]s%[1]sInput) (sdk.Poller, error) {
	req := sdk.Put%[1]sInput{
		Body: input,
		ExpectedStatusCodes: []int{
%[3]s,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	return client.baseClient.PutJsonThenPoll(ctx, req)
}
`, t.typeName, t.name, statusCodes)
}

func (t methodTemplater) statusCodes(indentation string) string {
	output := make([]string, 0)

	for _, statusCode := range t.expectedStatusCodes {
		alias := golangConstantForStatusCode(statusCode)
		description := descriptionForStatusCodeForMethod(statusCode, t.method, t.longRunningOperation)
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
