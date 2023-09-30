package response

import (
	coreErrors "errors"
	"github.com/pkg/errors"
	"gitlab.com/golibs-starter/golib/exception"
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		code    int
		message string
		data    interface{}
	}
	tests := []struct {
		name string
		args args
		want Response
	}{
		{
			name: "Happy case",
			args: args{
				code:    100,
				message: "Test message",
				data:    "Test data",
			},
			want: Response{
				Meta: Meta{
					Code:    100,
					Message: "Test message",
				},
				Data: "Test data",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.code, tt.args.message, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError(t *testing.T) {
	var resourceIdInvalid = exception.New(40008001, "Resource id is invalid")
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Response
	}{
		{
			name: "When error is native error should return 500",
			args: args{err: errors.New("a native error")},
			want: Response{
				Meta: Meta{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				},
				Data: nil,
			},
		},
		{
			name: "When error is a exception should return with exception code",
			args: args{err: resourceIdInvalid},
			want: Response{
				Meta: Meta{
					Code:    int(resourceIdInvalid.Code()),
					Message: resourceIdInvalid.Error(),
				},
				Data: nil,
			},
		},
		{
			name: "When error is a wrapped exception should return with exception code",
			args: args{err: errors.Wrap(resourceIdInvalid, "failed to get resource")},
			want: Response{
				Meta: Meta{
					Code:    int(resourceIdInvalid.Code()),
					Message: resourceIdInvalid.Error(),
				},
				Data: nil,
			},
		},
		{
			name: "When error is a message-wrapped exception should return with exception code",
			args: args{err: errors.WithMessage(resourceIdInvalid, "failed to get resource")},
			want: Response{
				Meta: Meta{
					Code:    int(resourceIdInvalid.Code()),
					Message: resourceIdInvalid.Error(),
				},
				Data: nil,
			},
		},
		{
			name: "When error is a joined error should return with exception code",
			args: args{err: coreErrors.Join(resourceIdInvalid, errors.New("failed to get resource"))},
			want: Response{
				Meta: Meta{
					Code:    int(resourceIdInvalid.Code()),
					Message: resourceIdInvalid.Error(),
				},
				Data: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Error(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
