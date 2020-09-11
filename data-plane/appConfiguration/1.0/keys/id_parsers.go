package keys

import "fmt"

type KeysId struct {
	ConfigurationEndpoint string
	KeyName               string
}

func NewKeysId(configurationEndpoint string, keyName string) KeysId {
	return KeysId{
		ConfigurationEndpoint: configurationEndpoint,
		KeyName:               keyName,
	}
}

func (id KeysId) ID() string {
	return fmt.Sprintf("%s/keys?name=%s", id.ConfigurationEndpoint, id.KeyName)
}
