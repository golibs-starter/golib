package config

import (
	"github.com/golibs-starter/golib/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ProfileReader interface {

	// Read config in a profile
	Read(profile string) (map[string]interface{}, error)
}

type DefaultProfileReader struct {
	scanPaths        []string
	format           string
	delim            string
	formatExtMapping map[string][]string
}

func NewDefaultProfileReader(scanPaths []string, format string, delim string) (*DefaultProfileReader, error) {
	if len(scanPaths) == 0 {
		return nil, errors.New("missing scanPaths parameter")
	}
	if len(format) == 0 {
		return nil, errors.New("missing format parameter")
	}
	if len(delim) == 0 {
		return nil, errors.New("missing delim parameter")
	}
	formatExtMapping := map[string][]string{
		"yaml": {"yaml", "yml"},
		"yml":  {"yaml", "yml"},
	}
	if _, exists := formatExtMapping[format]; !exists {
		return nil, ErrFormatNotSupported
	}
	return &DefaultProfileReader{
		scanPaths:        scanPaths,
		format:           format,
		delim:            delim,
		formatExtMapping: formatExtMapping,
	}, nil
}

func (p DefaultProfileReader) Read(profile string) (map[string]interface{}, error) {
	b, err := p.readBytes(profile)
	if err != nil {
		return nil, err
	}
	mp, err := p.unmarshalBytes(b)
	if err != nil {
		return nil, err
	}
	return mp, nil
}

func (p DefaultProfileReader) readBytes(profile string) ([]byte, error) {
	file, err := p.findFile(profile)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(file)
}

func (p DefaultProfileReader) unmarshalBytes(bytes []byte) (map[string]interface{}, error) {
	switch strings.ToLower(p.format) {
	case "yaml", "yml":
		var ms yaml.MapSlice
		if err := yaml.Unmarshal(bytes, &ms); err != nil {
			return nil, err
		}
		hMap := utils.YamlMapSliceToLinkedHMap(ms)
		expandedHMap := utils.ExpandInlineKeyInLinkedHMap(hMap, p.delim)
		return utils.LinkedHMapToMapStr(expandedHMap), nil
	default:
		return nil, errors.New("format not supported")
	}
}

func (p DefaultProfileReader) findFile(profile string) (string, error) {
	extensions := p.formatExtMapping[p.format]
	for _, scanPath := range p.scanPaths {
		for _, ext := range extensions {
			file := filepath.Join(scanPath, profile+"."+ext)
			stat, err := os.Stat(file)
			if os.IsNotExist(err) {
				continue
			}
			if stat.IsDir() {
				continue
			}
			return file, nil
		}
	}
	return "", errors.New("no profile file found")
}
