package gotool

import "testing"

func TestIsDomainName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"", args{s: "abc.com"}, true},
		{"", args{s: "localhost"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDomainName(tt.args.s); got != tt.want {
				t.Errorf("IsDomainName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDomain(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"xxx.", args{s: "xxx."}, false},
		{"abc.com", args{s: "abc.com"}, true},
		{"abc.com.", args{s: "abc.com."}, true},
		{"localhost", args{s: "localhost"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDomain(tt.args.s); got != tt.want {
				t.Errorf("IsDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
