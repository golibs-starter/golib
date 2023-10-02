package event

import (
	"context"
	"encoding/json"
	"github.com/golibs-starter/golib/utils"
	"github.com/google/uuid"
	"time"
)

const DefaultEventSource = "not_used"

func NewApplicationEvent(ctx context.Context, name string, options ...AppEventOpt) *ApplicationEvent {
	e := ApplicationEvent{
		Ctx: ctx,
	}
	for _, opt := range options {
		opt(&e)
	}
	if e.Id == "" {
		// No error reached, ignored
		generatedId, _ := uuid.NewUUID()
		e.Id = generatedId.String()
	}
	if e.Source == "" {
		e.Source = DefaultEventSource
	}
	e.Event = name
	e.Timestamp = utils.Time2Ms(time.Now())
	return &e
}

type ApplicationEvent struct {
	Ctx            context.Context        `json:"-"`
	Id             string                 `json:"id"`
	Event          string                 `json:"event"`
	Source         string                 `json:"source"`
	ServiceCode    string                 `json:"service_code"`
	AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
	PayloadData    interface{}            `json:"payload"`
	Timestamp      int64                  `json:"timestamp"`
}

func (a *ApplicationEvent) Context() context.Context {
	return a.Ctx
}

func (a *ApplicationEvent) Identifier() string {
	return a.Id
}

func (a *ApplicationEvent) Name() string {
	return a.Event
}

func (a *ApplicationEvent) Payload() interface{} {
	return a.PayloadData
}

func (a *ApplicationEvent) AddAdditionData(key string, value interface{}) {
	if a.AdditionalData == nil {
		a.AdditionalData = make(map[string]interface{})
	}
	a.AdditionalData[key] = value
}

func (a *ApplicationEvent) DeleteAdditionData(key string) {
	delete(a.AdditionalData, key)
}

func (a *ApplicationEvent) String() string {
	return a.ToString(a)
}

func (a *ApplicationEvent) ToString(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
