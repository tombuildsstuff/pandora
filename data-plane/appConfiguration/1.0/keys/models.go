package keys

type GetKeysResponse struct {
	Keys     []KeyName `json:"keys"`
	NextLink string    `json:"@nextLink"`
}

type KeyName struct {
	Name string `json:"name"`
}
