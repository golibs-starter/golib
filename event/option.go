package event

type AppEventOpt func(event *ApplicationEvent)

func WithId(id string) AppEventOpt {
	return func(event *ApplicationEvent) {
		event.Id = id
	}
}

func WithSource(source string) AppEventOpt {
	return func(event *ApplicationEvent) {
		event.Source = source
	}
}

func WithServiceCode(serviceCode string) AppEventOpt {
	return func(event *ApplicationEvent) {
		event.ServiceCode = serviceCode
	}
}

func WithAdditionalData(additionalData map[string]interface{}) AppEventOpt {
	return func(event *ApplicationEvent) {
		event.AdditionalData = additionalData
	}
}

func WithPayload(payload interface{}) AppEventOpt {
	return func(event *ApplicationEvent) {
		event.PayloadData = payload
	}
}
