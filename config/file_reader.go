package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.id.vin/vincart/golib/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type FileReader struct {
	path   string
	format string
	delim  string
}

func NewFileReader(path string, format string, delim string) *FileReader {
	return &FileReader{
		path:   filepath.Clean(path),
		format: format,
		delim:  delim,
	}
}

func (y FileReader) ReadBytes() ([]byte, error) {
	return ioutil.ReadFile(y.path)
}

func (y FileReader) Read() (map[string]interface{}, error) {
	b, err := y.ReadBytes()
	if err != nil {
		return nil, err
	}
	mp, err := y.unmarshalBytes(b)
	if err != nil {
		return nil, err
	}
	return mp, nil
}

func (y FileReader) unmarshalBytes(bytes []byte) (map[string]interface{}, error) {
	c := make(map[string]interface{})
	switch strings.ToLower(y.format) {
	case "yaml", "yml":
		var ms yaml.MapSlice
		if err := yaml.Unmarshal(bytes, &ms); err != nil {
			return nil, err
		}
		hMap := yamlMapSliceToLinkedHMap(ms)
		expandHMap := expandInlineKeyInMap(hMap, y.delim)
		c = utils.LinkedHMapToMapStr(expandHMap)
		return c, nil
	case "json":
		if err := json.Unmarshal(bytes, &c); err != nil {
			return nil, err
		}
		return c, nil
	default:
		return nil, errors.New("format not supported")
	}
}
