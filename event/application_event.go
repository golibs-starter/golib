package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"gitlab.id.vin/vincart/golib/log"
	"time"
)

const DefaultEventSource = "not_used"

type ApplicationEvent struct {
	Id             string                 `json:"id"`
	Event          string                 `json:"event"`
	Source         string                 `json:"source"`
	ServiceCode    string                 `json:"service_code"`
	EventPayload   interface{}            `json:"payload"`
	AdditionalData map[string]interface{} `json:"additional_data"`
	Timestamp      int64                  `json:"timestamp"`
}

func NewApplicationEvent(serviceCode string, eventName string, payload interface{}) ApplicationEvent {
	id := ""
	if genId, err := uuid.NewUUID(); err != nil {
		log.Warnf("Cannot create new event due by error [%v]", err)
	} else {
		id = genId.String()
	}
	return ApplicationEvent{
		Id:           id,
		Event:        eventName,
		Source:       DefaultEventSource,
		ServiceCode:  serviceCode,
		EventPayload: payload,
		Timestamp:    time.Now().UnixNano() / int64(time.Millisecond),
	}
}

func (a ApplicationEvent) String() string {
	data, _ := json.Marshal(a)
	return string(data)
}
