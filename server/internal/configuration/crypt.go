package configuration

import (
	"encoding/base64"
)

type PrivateKey []byte

func ParsePrivateKey(s string) (interface{}, error) {
	bts, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return PrivateKey(bts), nil
}
