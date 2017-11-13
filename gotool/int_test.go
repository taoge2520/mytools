package gotool

import "testing"

func TestComma(t *testing.T) {
	type args struct {
		v int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "", args: args{v: 999}, want: "999"},
		{name: "", args: args{v: 9999}, want: "9,999"},
		{name: "", args: args{v: 9999999}, want: "9,999,999"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Comma(tt.args.v); got != tt.want {
				t.Errorf("Comma() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteFormat(t *testing.T) {
	type args struct {
		bytes uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "", args: args{bytes: 1024}, want: "1.00 K"},
		{name: "", args: args{bytes: 1024 * 1024}, want: "1.00 M"},
		{name: "", args: args{bytes: 1024 * 1024 * 1024}, want: "1.00 G"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByteFormat(tt.args.bytes); got != tt.want {
				t.Errorf("ByteFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
