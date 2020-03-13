package typcore

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Configuration is detail of config
type Configuration struct {
	name string
	spec interface{}
}

// ConfigField is detail field of config
type ConfigField struct {
	Name     string
	Type     string
	Default  string
	Value    interface{}
	IsZero   bool
	Required bool
}

// NewConfiguration return new instance of Configuration
func NewConfiguration(name string, spec interface{}) *Configuration {
	return &Configuration{
		name: name,
		spec: spec,
	}
}

// Name of configuration
func (c *Configuration) Name() string {
	return c.name
}

// Spec of configuration
func (c *Configuration) Spec() interface{} {
	return c.spec
}

// Fields of Config Bean
func (c *Configuration) Fields() []*ConfigField {
	return retrieveFields(c.name, c.spec)
}

func retrieveFields(name string, spec interface{}) (fields []*ConfigField) {
	val := reflect.Indirect(reflect.ValueOf(spec))
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !fieldIgnored(field) {
			name := fmt.Sprintf("%s_%s", name, fieldName(field))
			fields = append(fields, &ConfigField{
				Name:     name,
				Type:     field.Type.Name(),
				Default:  fieldDefault(field),
				Required: fieldRequired(field),
				Value:    val.Field(i).Interface(),
				IsZero:   val.Field(i).IsZero(),
			})
		}
	}
	return
}

func fieldRequired(field reflect.StructField) (required bool) {
	if v, ok := field.Tag.Lookup("required"); ok {
		required, _ = strconv.ParseBool(v)
	}
	return
}

func fieldIgnored(field reflect.StructField) (ignored bool) {
	if v, ok := field.Tag.Lookup("ignored"); ok {
		ignored, _ = strconv.ParseBool(v)
	}
	return
}

func fieldDefault(field reflect.StructField) string {
	return field.Tag.Get("default")
}

func fieldName(field reflect.StructField) (name string) {
	name = strings.ToUpper(field.Name)
	if v, ok := field.Tag.Lookup("envconfig"); ok {
		name = v
	}
	return
}
