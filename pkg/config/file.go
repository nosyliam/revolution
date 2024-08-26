package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

type Format int

const (
	JSON Format = iota
	YAML
)

type Savable interface {
	Save() error
}

type File[T any] struct {
	name   string
	file   *os.File
	path   string
	format Format
	obj    *T
}

func (f *File[T]) Save() error {
	var data []byte
	var err error
	switch f.format {
	case JSON:
		data, err = json.Marshal(f)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
	case YAML:
		data, err = yaml.Marshal(f)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
	default:
		panic("unknown format")
	}

	_ = f.file.Truncate(0)
	_, _ = f.file.Seek(0, 0)
	_, err = f.file.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write")
	}
	return nil
}

func (f *File[T]) load() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get working directory")
	}

	path := filepath.Join(cwd, f.path)
	f.file, err = os.OpenFile(path, os.O_RDWR, 0755)
	if os.IsNotExist(err) {
		f.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return errors.Wrap(err, "failed to create file")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	data, err := io.ReadAll(f.file)
	if err != nil {
		return errors.Wrap(err, "failed to read")
	}

	var obj T
	switch f.format {
	case JSON:
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return errors.Wrap(err, "failed to parse")
		}
	case YAML:
		err = yaml.Unmarshal(data, &obj)
		if err != nil {
			return errors.Wrap(err, "failed to parse")
		}
	default:
		panic("unknown format")
	}

	f.obj = &obj
	return nil
}

func (f *File[T]) Object() *Object[T] {
	obj := &Object[T]{obj: f.obj}
	obj.Initialize(f.name, f)
	return obj
}

func (f *File[T]) Close() {
	if f.file != nil {
		_ = f.file.Close()
	}
}
