package util

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func ViperSetConfig(config any) {
	configValue := reflect.ValueOf(config).Elem()
	viperSetStructFields(configValue, "")
}

func viperSetStructFields(config reflect.Value, prefix string) {
	for i := 0; i < config.NumField(); i++ {
		field := config.Field(i)
		fieldType := config.Type().Field(i)

		key := prefix + fieldType.Tag.Get("yaml")
		if field.Kind() == reflect.Struct {
			viperSetStructFields(field, key+".")
		} else {
			value := field.Interface()
			viper.Set(key, value)
		}
	}
}

func ParseCfgField[T any](cfgField any) (*T, error) {
	cfgFieldMap := cfgField.(map[string]any)
	var ret T
	if err := mapstructure.Decode(cfgFieldMap, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
