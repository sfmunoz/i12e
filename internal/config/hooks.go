package config

import (
	"fmt"
	"net/netip"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

func PrefixDecodeHook() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data any) (any, error) {
		if from.Kind() != reflect.String {
			return data, nil
		}
		if to != reflect.TypeFor[*netip.Prefix]() {
			return data, nil
		}
		s := data.(string)
		p, err := netip.ParsePrefix(s)
		if err != nil {
			return nil, fmt.Errorf("invalid prefix %q: %w", s, err)
		}
		return &p, nil
	}
}
