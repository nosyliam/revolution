package config

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strconv"
)

type Defaulter interface {
	Default()
}

type Object interface {
	Initialize(path string)
	Set(chain []string, index int, value interface{}) error
	Get(chain []string, index int) (interface{}, error)
}

type KeyedObject interface {
	Object
	Key() string
}

type config struct {
	path string
}

func (c *config) setField(field reflect.Value, chain []string, index int, value interface{}) error {
	if obj, ok := field.Interface().(Object); ok {
		if len(chain) == 1 {
			return errors.New("config field path not found")
		}
		if index+1 == len(chain) {
			return errors.New("cannot get an object value")
		}
		return obj.Set(chain, index+1, value)
	}
	switch field.Kind() {
	case reflect.Int:
		field.SetInt(value.(int64))
	case reflect.Bool:
		field.SetBool(value.(bool))
	case reflect.String:
		field.SetString(value.(string))
	default:
		return errors.New("unsupported value")
	}
	return nil
}

func (c *config) getField(field reflect.Value, chain []string, index int) (interface{}, error) {
	if obj, ok := field.Interface().(Object); ok {
		if len(chain) == 1 {
			return nil, errors.New("config field path not found")
		}
		if index+1 == len(chain) {
			return nil, errors.New("cannot get an object value")
		}
		return obj.Get(chain, index+1)
	}
	switch field.Kind() {
	case reflect.Int:
		return field.Int(), nil
	case reflect.Bool:
		return field.Bool(), nil
	case reflect.String:
		return field.String(), nil
	default:
		return errors.New("unsupported value"), nil
	}
}

type configList[T Object] struct {
	config
	data []T
}

func (c *configList[T]) MarshalYAML() (interface{}, error) {
	return c.data, nil
}

func (c *configList[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var temp []T

	if err := unmarshal(&temp); err != nil {
		return err
	}

	c.data = temp
	return nil
}

func (c *configList[T]) Set(chain []string, index int, value interface{}) error {
	idx, err := strconv.Atoi(chain[index])
	if err != nil || (idx < 0 || idx > len(c.data)) {
		return errors.New("invalid index")
	}

	field := reflect.ValueOf(c.data).Field(idx)
	return c.setField(field, chain, index, value)
}

func (c *configList[T]) Get(chain []string, index int) (interface{}, error) {
	idx, err := strconv.Atoi(chain[index])
	if err != nil || (idx < 0 || idx > len(c.data)) {
		return nil, errors.New("invalid index")
	}

	field := reflect.ValueOf(c.data).Field(idx)
	return c.getField(field, chain, index), nil
}

func (c configList[T]) Append(val T) {
	c.data = append(c.data, val)
}

type indexedConfigList[T KeyedObject] struct {
	configList[T]
	index map[string]T
}

type configObject struct {
	config
}

func (c *configObject) Initialize(path string) {
	c.path = path
	t := reflect.TypeOf(c)
	val := reflect.ValueOf(c).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		meta := t.Field(i)

		if field.Kind() == reflect.Ptr && field.IsNil() {
			obj := reflect.New(field.Type().Elem())
			ifc := obj.Interface()
			if def, ok := ifc.(Defaulter); ok {
				def.Default()
			}
			if cfg, ok := ifc.(Object); ok {
				cfg.Initialize(fmt.Sprintf("%s.%s", path, meta.Name))
			}
			field.Set(obj)
		}
	}
}

func (c *configObject) Set(chain []string, index int, value interface{}) error {
	t := reflect.TypeOf(c)
	val := reflect.ValueOf(c).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		meta := t.Field(i)
		if meta.Tag.Get("yaml") == chain[index] && field.IsValid() && field.CanSet() {
			return c.setField(field, chain, index, value)
		}
	}

	return errors.New("field path not found")
}

func (c *configObject) Get(chain []string, index int) (interface{}, error) {
	t := reflect.TypeOf(c)
	val := reflect.ValueOf(c).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		meta := t.Field(i)
		if meta.Tag.Get("yaml") == chain[index] && field.IsValid() && field.CanSet() {
			return c.getField(field, chain, index)
		}
	}

	return nil, errors.New("field path not found")
}

func (c *configObject) compilePath(path string) ([]string, error) {
	re := regexp.MustCompile(`([^\.\[\]]+)`)
	matches := re.FindAllString(path, -1)
	if len(matches) == 0 {
		return nil, errors.New("invalid path")
	}
	return matches, nil
}

func (c *configObject) SetPath(path string, value interface{}) error {
	chain, err := c.compilePath(path)
	if err != nil {
		return err
	}
	return c.Set(chain, 0, value)
}

func (c *configObject) GetPath(path string, value interface{}) (interface{}, error) {
	chain, err := c.compilePath(path)
	if err != nil {
		return nil, err
	}
	return c.Get(chain, 0)
}
