package actuator

import (
	"github.com/golibs-starter/golib/config"
)

type InfoService interface {
	Info() Info
}

type DefaultInfoService struct {
	props     *config.AppProperties
	informers []Informer
}

func NewDefaultInfoService(props *config.AppProperties, informers []Informer) InfoService {
	return &DefaultInfoService{
		props:     props,
		informers: informers,
	}
}

func (d DefaultInfoService) Info() Info {
	info := Info{
		Name: d.props.Name,
		Info: make(map[string]interface{}),
	}
	for _, informer := range d.informers {
		info.Info[informer.Key()] = informer.Value()
	}
	return info
}
