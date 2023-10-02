package pubsub

import (
	"errors"
	"github.com/golibs-starter/golib/actuator"
)

type DefaultBusInformer struct {
	bus *DefaultEventBus
}

func NewDefaultBusInformer(bus EventBus) (actuator.Informer, error) {
	implBus, ok := bus.(*DefaultEventBus)
	if !ok {
		return nil, errors.New("EventBus is not DefaultEventBus")
	}
	return &DefaultBusInformer{bus: implBus}, nil
}

func (d DefaultBusInformer) Key() string {
	return "event_bus"
}

func (d DefaultBusInformer) Value() interface{} {
	return map[string]interface{}{
		"channel_capacity":     cap(d.bus.eventCh),
		"channel_current_size": len(d.bus.eventCh),
	}
}
