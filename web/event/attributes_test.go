package event

import (
	"context"
	"github.com/golibs-starter/golib/event"
	"github.com/golibs-starter/golib/web/constant"
	"reflect"
	"testing"
)

func TestMakeAttributes(t *testing.T) {
	type args struct {
		e *AbstractEvent
	}
	tests := []struct {
		name string
		args args
		want *Attributes
	}{
		{
			name: "When event is abstract event should return attributes",
			args: args{e: &AbstractEvent{
				ApplicationEvent: &event.ApplicationEvent{
					AdditionalData: map[string]interface{}{
						constant.HeaderDeviceId:        "dv1",
						constant.HeaderDeviceSessionId: "ss1",
					},
				},
				RequestId:         "req1",
				UserId:            "user1",
				TechnicalUsername: "user2",
			}},
			want: &Attributes{
				CorrelationId:     "req1",
				UserId:            "user1",
				DeviceId:          "dv1",
				DeviceSessionId:   "ss1",
				TechnicalUsername: "user2",
			},
		},
		{
			name: "When event is abstract event and addition data is not provide should not error",
			args: args{e: &AbstractEvent{
				ApplicationEvent:  &event.ApplicationEvent{},
				RequestId:         "req1",
				UserId:            "user1",
				TechnicalUsername: "user2",
			}},
			want: &Attributes{
				CorrelationId:     "req1",
				UserId:            "user1",
				DeviceId:          "",
				DeviceSessionId:   "",
				TechnicalUsername: "user2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeAttributes(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAttributes(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *Attributes
	}{
		{
			name: "When ctx is not constants attributes should return nil",
			args: args{ctx: context.Background()},
			want: nil,
		},
		{
			name: "When ctx is constants event_attributes filed but not event attributes should return nil",
			args: args{
				ctx: context.WithValue(context.Background(), constant.ContextEventAttributes, "not event attributes"),
			},
			want: nil,
		},
		{
			name: "When ctx is constants attributes should return attributes",
			args: args{
				ctx: context.WithValue(context.Background(), constant.ContextEventAttributes, &Attributes{
					CorrelationId:     "1",
					UserId:            "2",
					DeviceId:          "3",
					DeviceSessionId:   "4",
					TechnicalUsername: "5",
				}),
			},
			want: &Attributes{
				CorrelationId:     "1",
				UserId:            "2",
				DeviceId:          "3",
				DeviceSessionId:   "4",
				TechnicalUsername: "5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAttributes(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}
