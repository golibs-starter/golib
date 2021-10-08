package config

import (
	"encoding/json"
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
	mp := make(map[string]interface{})
	if err := y.unmarshalBytes(b, mp); err != nil {
		return nil, err
	}
	return y.expandInlineKeyInMap(mp), nil
}

func (y FileReader) unmarshalBytes(bytes []byte, c map[string]interface{}) error {
	switch strings.ToLower(y.format) {
	case "yaml", "yml":
		if err := yaml.Unmarshal(bytes, &c); err != nil {
			return err
		}
	case "json":
		if err := json.Unmarshal(bytes, &c); err != nil {
			return err
		}
	}
	return nil
}

func (y FileReader) expandInlineKeyInMap(mp map[string]interface{}) map[string]interface{} {
	var newMp = make(map[string]interface{})
	for k, v := range mp {
		var startK string
		var startV interface{}
		//hash := strings.Split(strings.ToLower(k), y.delim)
		hash := strings.Split(k, y.delim)
		if len(hash) == 1 {
			startK = hash[0]
			startV = v
		} else {
			startK = hash[0]
			startV = map[string]interface{}{
				strings.Join(hash[1:], y.delim): v,
			}
		}
		switch vOk := startV.(type) {
		case map[string]interface{}:
			newMp[startK] = y.expandInlineKeyInMap(vOk)
		default:
			newMp[startK] = startV
		}
	}
	return newMp
}
