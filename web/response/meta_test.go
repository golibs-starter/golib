package response

import "testing"

func TestMeta_HttpStatus(t *testing.T) {
	type fields struct {
		Code    int
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "When code is 0 should return 200",
			fields: fields{
				Code: 0,
			},
			want: 200,
		},
		{
			name: "When code < 0 should return 200",
			fields: fields{
				Code: -1,
			},
			want: 200,
		},
		{
			name: "When code is xx should return 200",
			fields: fields{
				Code: 25,
			},
			want: 200,
		},
		{
			name: "When code is 200 should return 200",
			fields: fields{
				Code: 200,
			},
			want: 200,
		},
		{
			name: "When code is xxx should return xxx",
			fields: fields{
				Code: 401,
			},
			want: 401,
		},
		{
			name: "When code is xxxy should return xxx",
			fields: fields{
				Code: 4030,
			},
			want: 403,
		},
		{
			name: "When code is xxxzzzz should return xxx",
			fields: fields{
				Code: 4040812,
			},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Meta{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
			}
			if got := m.HttpStatus(); got != tt.want {
				t.Errorf("HttpStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
