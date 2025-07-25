package tools

import "testing"

func TestGetCityByIp(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{ip: "115.193.200.155"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCityByIp(tt.args.ip); got != tt.want {
				t.Errorf("GetCityByIp() = %v, want %v", got, tt.want)
			}
		})
	}
}
