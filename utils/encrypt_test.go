package utils

import (
	"fmt"
	"testing"
)

func TestHmac(t *testing.T) {
	type args struct {
		sign   HmacType
		data   string
		secret []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{sign: "hmacSha256", data: fmt.Sprintf("deviceName%vproductKey%v", "c5D542F00", "k0sr6A5BMbN"), secret: []byte("hxabHJVBJJyyUXDdu5Sa7xeLfyo=")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret := Hmac(tt.args.sign, tt.args.data, tt.args.secret)
			fmt.Println(ret)
		})
	}
}
