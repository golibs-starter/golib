package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/utils"
	"time"
)

const DefaultEventSource = "not_used"

type ApplicationEvent struct {
	Id             string                 `json:"id"`
	Event          string                 `json:"event"`
	Source         string                 `json:"source"`
	ServiceCode    string                 `json:"service_code"`
	AdditionalData map[string]interface{} `json:"additional_data"`
	Timestamp      int64                  `json:"timestamp"`
}

func NewApplicationEvent(eventName string) *ApplicationEvent {
	id := ""
	if genId, err := uuid.NewUUID(); err != nil {
		log.Warnf("Cannot create new event due by error [%v]", err)
	} else {
		id = genId.String()
	}
	return &ApplicationEvent{
		Id:        id,
		Event:     eventName,
		Source:    DefaultEventSource,
		Timestamp: utils.Time2Ms(time.Now()),
	}
}

func (a ApplicationEvent) Identifier() string {
	return a.Id
}

func (a ApplicationEvent) Name() string {
	return a.Event
}

func (a ApplicationEvent) String() string {
	return a.ToString(a)
}

func (a ApplicationEvent) ToString(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
