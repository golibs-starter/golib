package utils

import (
	"reflect"
	"testing"
)

func TestPrependString(t *testing.T) {
	type args struct {
		slice []string
		e     string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test prepend to not empty slice",
			args: args{
				slice: []string{"a", "b"},
				e:     "c",
			},
			want: []string{"c", "a", "b"},
		},
		{
			name: "Test prepend to empty slice",
			args: args{
				slice: []string{},
				e:     "c",
			},
			want: []string{"c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrependString(tt.args.slice, tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrependString() = %v, want %v", got, tt.want)
			}
		})
	}
}
