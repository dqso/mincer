package configuration

import (
	"encoding/base64"
	"fmt"
	"github.com/caarlos0/env/v9"
	"log/slog"
	"reflect"
	"strings"
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
			reflect.TypeOf(PrivateKey{}):  ParsePrivateKey,
			reflect.TypeOf(slog.Level(0)): ParseSlogLevel,
		},
	}); err != nil {
		return nil, err
	}

	return c, nil
}

type PrivateKey []byte

func ParsePrivateKey(s string) (interface{}, error) {
	bts, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return PrivateKey(bts), nil
}

func ParseSlogLevel(s string) (interface{}, error) {
	switch strings.ToLower(s) {
	case "debug", "-4":
		return slog.LevelDebug, nil
	case "info", "0":
		return slog.LevelInfo, nil
	case "warn", "4":
		return slog.LevelWarn, nil
	case "error", "8":
		return slog.LevelError, nil
	default:
		return nil, fmt.Errorf("incorrect log level")
	}
}
