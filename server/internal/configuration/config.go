package configuration

import (
	"github.com/caarlos0/env/v9"
	"reflect"
)

type Config struct {
	env envModel
}

func NewConfig() (*Config, error) {
	c := &Config{}

	if err := env.ParseWithOptions(&c.env, env.Options{
		Environment:           nil,
		TagName:               "",
		RequiredIfNoDef:       false,
		OnSet:                 nil,
		Prefix:                "",
		UseFieldNameByDefault: false,
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(PrivateKey{}): ParsePrivateKey,
		},
	}); err != nil {
		return nil, err
	}

	return c, nil
}
