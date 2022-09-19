package utils

import "testing"

type DummyStruct struct {
}

func TestGetStructFullname(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test normal struct",
			args: args{val: DummyStruct{}},
			want: "utils.DummyStruct",
		},
		{
			name: "Test pointer struct",
			args: args{val: &DummyStruct{}},
			want: "utils.DummyStruct",
		},
		{
			name: "Test nil val",
			args: args{val: nil},
			want: "",
		},
		{
			name: "Test not struct val",
			args: args{val: "test"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStructFullname(tt.args.val); got != tt.want {
				t.Errorf("GetStructFullname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStructShortName(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test normal struct",
			args: args{val: DummyStruct{}},
			want: "DummyStruct",
		},
		{
			name: "Test pointer struct",
			args: args{val: &DummyStruct{}},
			want: "DummyStruct",
		},
		{
			name: "Test nil val",
			args: args{val: nil},
			want: "",
		},
		{
			name: "Test not struct val",
			args: args{val: "test"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStructShortName(tt.args.val); got != tt.want {
				t.Errorf("GetStructShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}
