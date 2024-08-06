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

type configFile struct {
	file   *os.File
	path   string
	format Format
}

func (c *configFile) Save() error {
	var data []byte
	var err error
	switch c.format {
	case JSON:
		data, err = json.Marshal(c)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
	case YAML:
		data, err = yaml.Marshal(c)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
	default:
		panic("unknown format")
	}

	_ = c.file.Truncate(0)
	_, _ = c.file.Seek(0, 0)
	_, err = c.file.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write")
	}
	return nil
}

func (c *configFile) load() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get working directory")
	}

	path := filepath.Join(cwd, c.path)
	c.file, err = os.OpenFile(path, os.O_RDWR, 0755)
	if os.IsNotExist(err) {
		c.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return errors.Wrap(err, "failed to create file")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	data, err := io.ReadAll(c.file)
	if err != nil {
		return errors.Wrap(err, "failed to read")
	}

	switch c.format {
	case JSON:
		err = json.Unmarshal(data, c)
		if err != nil {
			return errors.Wrap(err, "failed to parse")
		}
	case YAML:
		err = yaml.Unmarshal(data, c)
		if err != nil {
			return errors.Wrap(err, "failed to parse")
		}
	default:
		panic("unknown format")
	}

	return nil
}

func (c *configFile) Close() {
	if c.file != nil {
		_ = c.file.Close()
	}
}
