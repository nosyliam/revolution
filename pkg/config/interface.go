package config

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strconv"
)

type Reactive interface {
	Initialize(path string, file Savable)
	Set(chain chain, index int, value interface{}) error
	Get(chain chain, index int) (interface{}, error)
	GetConcrete(chain chain, index int) (interface{}, error)
	Append(chain chain, index int) error
	Delete(chain chain, index int) error
	Length(chain chain, index int) int
}

type reactiveList interface {
	Reactive
	list()
}

type reactiveObject interface {
	Reactive
	object() interface{}
}

type Runtime interface {
	Emit(path string, value interface{})
}

type config struct {
	path    string
	file    Savable
	runtime Runtime
}

type link struct {
	val      string
	brackets bool
}

type chain []link

func (c *config) errPath(err string) error {
	return errors.New(fmt.Sprintf("%s: %s", c.path, err))
}

func (c *config) errPanic(err string) error {
	panic(fmt.Sprintf("%s: %s", c.path, err))
}

func (c *config) setField(field reflect.Value, chain chain, index int, value interface{}) error {
	if obj, ok := field.Interface().(Reactive); ok {
		if index+1 == len(chain) {
			return errors.New("cannot set an object value")
		}
		return obj.Set(chain, index+1, value)
	}
	if len(chain)-1 != index {
		return c.errPath("cannot index a primitive value")
	}
	switch field.Kind() {
	case reflect.Int:
		field.SetInt(int64(value.(int)))
	case reflect.Bool:
		field.SetBool(value.(bool))
	case reflect.String:
		field.SetString(value.(string))
	default:
		return errors.New("unsupported value")
	}
	return nil
}

func (c *config) getField(field reflect.Value, chain chain, index int) (interface{}, error) {
	if obj, ok := field.Interface().(Reactive); ok {
		if _, ok := obj.(reactiveObject); ok {
			if index+1 == len(chain) {
				return nil, c.errPath("cannot get an object value")
			}
		}
		return obj.Get(chain, index+1)
	}
	if len(chain)-1 != index {
		return nil, c.errPath("cannot index a primitive value")
	}
	switch field.Kind() {
	case reflect.Int:
		return int(field.Int()), nil
	case reflect.Bool:
		return field.Bool(), nil
	case reflect.String:
		return field.String(), nil
	default:
		return errors.New("unsupported value"), nil
	}
}

func (c *config) getConcreteField(field reflect.Value, chain chain, index int) (interface{}, error) {
	if obj, ok := field.Interface().(Reactive); ok {
		if rObj, ok := obj.(reactiveObject); ok {
			if index+1 == len(chain) {
				return rObj.object(), nil
			}
		}
		return obj.GetConcrete(chain, index+1)
	}
	if len(chain)-1 != index {
		return nil, c.errPath("cannot index a primitive value")
	}
	return field.Interface(), nil
}

type List[T any] struct {
	config
	prim       []T
	obj        []*Object[T]
	index      map[string]*Object[T]
	key, keySz string
}

func (c *List[T]) MarshalYAML() (interface{}, error) {
	if c.prim != nil {
		return c.prim, nil
	} else {
		return c.obj, nil
	}
}

func (c *List[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var zero [0]T
	tt := reflect.TypeOf(zero).Elem()
	if tt.Kind() != reflect.Struct {
		var temp []T
		if err := unmarshal(&temp); err != nil {
			return err
		}
		c.prim = temp
		return nil
	}

	var temp []*T
	if err := unmarshal(&temp); err != nil {
		return err
	}
	var objs []*Object[T]
	for _, v := range temp {
		objs = append(objs, &Object[T]{obj: v})
	}
	c.obj = objs

	for i := 0; i < tt.NumField(); i++ {
		field := tt.Field(i)
		if key := field.Tag.Get("key"); key == "true" {
			c.key = field.Name
			c.keySz = field.Tag.Get("yaml")
			c.index = make(map[string]*Object[T])
			for _, obj := range c.obj {
				objKey := reflect.ValueOf(obj.Concrete()).Elem().FieldByName(c.key).String()
				if _, ok := c.index[objKey]; ok {
					return c.errPath(fmt.Sprintf("duplicate key: %s", objKey))
				}
				c.index[objKey] = obj
			}
		}
	}

	return nil
}

func (c *List[T]) Initialize(path string, file Savable) {
	if c.prim == nil && c.obj == nil {
		var zero [0]T
		tt := reflect.TypeOf(zero).Elem()
		if tt.Kind() != reflect.Struct {
			c.prim = make([]T, 0)
		} else {
			c.obj = make([]*Object[T], 0)
			for i := 0; i < tt.NumField(); i++ {
				field := tt.Field(i)
				if key := field.Tag.Get("key"); key == "true" {
					c.key = field.Name
					c.keySz = field.Tag.Get("yaml")
					c.index = make(map[string]*Object[T])
				}
			}
		}
	}

	c.path = path
	if c.obj != nil {
		for n, obj := range c.obj {
			if c.key != "" {
				key := reflect.ValueOf(obj).FieldByName(c.key).String()
				obj.Initialize(fmt.Sprintf("%s[%s]", path, key), file)
			} else {
				obj.Initialize(fmt.Sprintf("%s[%d]", path, n), file)
			}
		}
	}
}

func (c *List[T]) Set(chain chain, index int, value interface{}) error {
	if len(chain) == index {
		return c.errPath("a key must be provided")
	}
	if !chain[index].brackets {
		return c.errPath("list values must be indexed with brackets")
	}
	if c.index != nil {
		if len(chain) == index+1 {
			return c.errPath("cannot set an object value")
		}
		if chain[index+1].val == c.keySz {
			return c.errPath(fmt.Sprintf("cannot modify object key \"%s\"", chain[index].val))
		}
		if val, ok := c.index[chain[index].val]; ok {
			return c.setField(reflect.ValueOf(val), chain, index, value)
		} else {
			return c.errPath(fmt.Sprintf("invalid key \"%s\"", chain[index].val))
		}
	}
	idx, err := strconv.Atoi(chain[index].val)
	var count int
	var slice interface{}
	if c.prim != nil {
		count = len(c.prim)
		slice = c.prim
	} else {
		count = len(c.obj)
		slice = c.obj
	}
	if err != nil || (idx < 0 || idx > count) {
		return c.errPath("invalid index")
	}

	field := reflect.ValueOf(slice).Index(idx)
	return c.setField(field, chain, index, value)
}

func (c *List[T]) Get(chain chain, index int) (interface{}, error) {
	if len(chain) == index {
		return nil, c.errPath("a key must be provided")
	}
	if !chain[index].brackets {
		return nil, c.errPath("list values must be indexed with brackets")
	}
	if c.index != nil {
		if val, ok := c.index[chain[index].val]; ok {
			return c.getField(reflect.ValueOf(val), chain, index)
		} else {
			return nil, c.errPath(fmt.Sprintf("invalid key \"%s\"", chain[index].val))
		}
	}
	var count int
	var slice interface{}
	if c.prim != nil {
		count = len(c.prim)
		slice = c.prim
	} else {
		count = len(c.obj)
		slice = c.obj
	}
	idx, err := strconv.Atoi(chain[index].val)
	if err != nil || (idx < 0 || idx >= count) {
		return nil, c.errPath("invalid integer index")
	}

	field := reflect.ValueOf(slice).Index(idx)
	return c.getField(field, chain, index)
}

func (c *List[T]) GetConcrete(chain chain, index int) (interface{}, error) {
	if c.prim != nil {
		idx, err := strconv.Atoi(chain[index].val)
		if err != nil || (idx < 0 || idx >= len(c.prim)) {
			return nil, c.errPath("invalid integer index")
		}
		field := reflect.ValueOf(c.prim).Index(idx)
		return field.Interface(), nil
	}
	if len(chain) == index {
		return nil, c.errPath("a key must be provided")
	}
	if !chain[index].brackets {
		return nil, c.errPath("list values must be indexed with brackets")
	}
	if c.index != nil {
		if val, ok := c.index[chain[index].val]; ok {
			return c.getField(reflect.ValueOf(val), chain, index)
		} else {
			return nil, c.errPath(fmt.Sprintf("invalid key \"%s\"", chain[index].val))
		}
	}
	idx, err := strconv.Atoi(chain[index].val)
	if err != nil || (idx < 0 || idx >= len(c.obj)) {
		return nil, c.errPath("invalid integer index")
	}

	field := reflect.ValueOf(c.obj).Index(idx)
	return field.Interface().(Reactive).Get(chain, index+1)
}

func (c *List[T]) Append(chain chain, index int) error {
	if c.index != nil && !chain[index].brackets {
		return errors.New("list values must be indexed with brackets")
	}
	if len(chain)-1 != index && c.prim == nil && c.index != nil {
		if val, ok := c.index[chain[index].val]; ok {
			return val.Append(chain, index+1)
		} else {
			return c.errPath(fmt.Sprintf("invalid key \"%s\"", chain[index].val))
		}
	}
	if len(chain)-1 == index && c.index == nil {
		return c.errPath("a key must not be given when appending to a keyless list")
	}

	if c.prim != nil {
		var zero T
		c.prim = append(c.prim, zero)
		return nil
	}

	cfo := &Object[T]{}
	if c.key != "" {
		key := chain[index].val
		if _, ok := c.index[key]; ok {
			return c.errPath(fmt.Sprintf("key \"%s\" already exists", key))
		}
		cfo.Initialize(fmt.Sprintf("%s[%s]", c.path, chain[index].val), c.file)
		ref := reflect.ValueOf(cfo.obj)
		ref.Elem().FieldByName(c.key).SetString(key)
		c.index[key] = cfo
	} else {
		cfo.Initialize(fmt.Sprintf("%s[%d]", c.path, len(c.obj)), c.file)
	}
	c.obj = append(c.obj, cfo)

	return nil
}

func (c *List[T]) Delete(chain chain, index int) error {
	return nil
}

func (c *List[T]) Length(chain chain, index int) int {
	if len(chain) == index {
		if c.prim != nil {
			return len(c.prim)
		} else {
			return len(c.obj)
		}
	}
	if c.prim != nil {
		c.errPanic("cannot get the length of a primitive value")
	}
	idx, err := strconv.Atoi(chain[index].val)
	if err != nil || (idx < 0 || idx >= len(c.obj)) {
		c.errPanic("invalid integer index")
	}

	field := reflect.ValueOf(c.obj).Index(idx)
	return field.Interface().(Reactive).Length(chain, index+1)
}

func (c *List[T]) list() {}

func (c *List[T]) Concrete() []T {
	if c.prim != nil {
		return c.prim
	} else {
		var concrete []T
		for _, obj := range c.obj {
			concrete = append(concrete, *obj.obj)
		}
		return concrete // optimize?
	}
}

type Object[T any] struct {
	config
	obj *T
}

func (c *Object[T]) object() interface{} {
	return c.obj
}

func (c *Object[T]) MarshalYAML() (interface{}, error) {
	return c.obj, nil
}

func (c *Object[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var obj = new(T)
	if err := unmarshal(obj); err != nil {
		return err
	}
	c.obj = obj
	return nil
}

func (c *Object[T]) Initialize(path string, file Savable) {
	c.path = path
	if c.obj == nil {
		c.obj = new(T)
		t := reflect.TypeOf(c.obj).Elem()
		val := reflect.ValueOf(c.obj).Elem()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			meta := t.Field(i)
			def := meta.Tag.Get("default")
			if def == "" {
				continue
			}
			switch field.Kind() {
			case reflect.Int:
				if num, err := strconv.Atoi(def); err == nil {
					field.SetInt(int64(num))
				} else {
					panic("invalid number default")
				}
			case reflect.Bool:
				switch def {
				case "true":
					field.SetBool(true)
				case "false":
					field.SetBool(false)
				default:
					panic("invalid boolean default")
				}
			case reflect.String:
				field.SetString(def)
			}
		}
	}
	t := reflect.TypeOf(c.obj).Elem()
	val := reflect.ValueOf(c.obj).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		meta := t.Field(i)
		if field.Kind() == reflect.Ptr && meta.Type.Elem().Kind() == reflect.Struct && field.IsNil() {
			obj := reflect.New(meta.Type.Elem())
			ifc := obj.Interface()
			if cfg, ok := ifc.(Reactive); ok {
				cfg.Initialize(fmt.Sprintf("%s.%s", path, meta.Name), file)
			}
			field.Set(obj)
		}
	}
}

func (c *Object[T]) Set(chain chain, index int, value interface{}) error {
	t := reflect.TypeOf(c.obj).Elem()
	val := reflect.ValueOf(c.obj).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		meta := t.Field(i)
		if (meta.Tag.Get("yaml") == chain[index].val || meta.Tag.Get("state") == chain[index].val) && field.CanSet() {
			return c.setField(field, chain, index, value)
		}
	}

	return c.errPath("field path not found")
}

func (c *Object[T]) Get(chain chain, index int) (interface{}, error) {
	t := reflect.TypeOf(c.obj).Elem()
	val := reflect.ValueOf(c.obj).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		meta := t.Field(i)
		if meta.Tag.Get("yaml") == chain[index].val || meta.Tag.Get("state") == chain[index].val {
			return c.getField(field, chain, index)
		}
	}

	return nil, c.errPath(fmt.Sprintf("field %s not found", chain[index].val))
}

func (c *Object[T]) GetConcrete(chain chain, index int) (interface{}, error) {
	t := reflect.TypeOf(c.obj).Elem()
	val := reflect.ValueOf(c.obj).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		meta := t.Field(i)
		if meta.Tag.Get("yaml") == chain[index].val || meta.Tag.Get("state") == chain[index].val {
			return c.getConcreteField(field, chain, index)
		}
	}

	return nil, c.errPath(fmt.Sprintf("field %s not found", chain[index].val))
}

func (c *Object[T]) Append(chain chain, index int) error {
	if len(chain) == index {
		return c.errPath("cannot append to an object")
	}
	t := reflect.TypeOf(c.obj).Elem()
	val := reflect.ValueOf(c.obj).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		meta := t.Field(i)
		if meta.Tag.Get("yaml") == chain[index].val || meta.Tag.Get("state") == chain[index].val {
			switch obj := field.Interface().(type) {
			case reactiveObject:
				if index+1 == len(chain) {
					return errors.New("cannot set an object value")
				}
				return obj.Append(chain, index+1)
			case reactiveList:
				return obj.Append(chain, index+1)
			}
		}
	}

	return c.errPath(fmt.Sprintf("field %s not found", chain[index].val))
}

func (c *Object[T]) Delete(chain chain, index int) error {
	return nil
}

func (c *Object[T]) Length(chain chain, index int) int {
	if len(chain) == index {
		c.errPanic("cannot get the length of an object")
	}
	t := reflect.TypeOf(c.obj).Elem()
	val := reflect.ValueOf(c.obj).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		meta := t.Field(i)
		if meta.Tag.Get("yaml") == chain[index].val || meta.Tag.Get("state") == chain[index].val {
			switch obj := field.Interface().(type) {
			case reactiveObject:
				if index+1 == len(chain) {
					c.errPanic("cannot set an object value")
				}
				return obj.Length(chain, index+1)
			case reactiveList:
				return obj.Length(chain, index+1)
			}
		}
	}

	c.errPanic(fmt.Sprintf("field %s not found", chain[index].val))
	return 0
}

func (c *Object[T]) Concrete() *T {
	return c.obj
}

func (c *Object[T]) SetPath(path string, value interface{}) error {
	chain, err := compilePath(path)
	if err != nil {
		return err
	}
	return c.Set(chain, 0, value)
}

func (c *Object[T]) GetPath(path string) (interface{}, error) {
	chain, err := compilePath(path)
	if err != nil {
		return nil, err
	}
	return c.Get(chain, 0)
}

func (c *Object[T]) GetObjectPath(path string) (interface{}, error) {
	chain, err := compilePath(path)
	if err != nil {
		return nil, err
	}
	return c.GetConcrete(chain, 0)
}

func (c *Object[T]) AppendPath(path string) error {
	chain, err := compilePath(path)
	if err != nil {
		return nil
	}
	return c.Append(chain, 0)
}

func (c *Object[T]) DeletePath(path string) error {
	/*chain, err := c.compilePath(path)
	if err != nil {
		return nil
	}*/
	return nil
}

func (c *Object[T]) LengthPath(path string) int {
	chain, err := compilePath(path)
	if err != nil {
		return 0
	}
	return c.Length(chain, 0)
}

func compilePath(path string) (chain, error) {
	var chains chain
	regex := regexp.MustCompile(`(\w+)|\[(.*?)\]`)
	matches := regex.FindAllStringSubmatch(path, -1)
	if len(matches) == 0 {
		return nil, errors.New("invalid path")
	}

	for _, match := range matches {
		if match[1] != "" {
			chains = append(chains, link{match[1], false})
		}
		if match[2] != "" {
			chains = append(chains, link{match[2], true})
		}
	}

	return chains, nil
}

func Concrete[T any](object Reactive, path string, args ...interface{}) *T {
	path = fmt.Sprintf(path, args...)
	chain, err := compilePath(path)
	if err != nil {
		panic("invalid concrete path")
	}
	if obj, err := object.GetConcrete(chain, 0); err != nil {
		panic(fmt.Sprintf("invalid concrete access: %s: %v", path, err))
	} else {
		if val, ok := obj.(T); ok {
			return &val
		} else if valPtr, ok := obj.(*T); ok {
			return valPtr
		} else {
			panic(fmt.Sprintf("invalid concrete object type: %s", path))
		}
	}
}
