package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInt64ToBStr(t *testing.T) {
	type args struct {
		num int64
		bit int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{num: 4323, bit: 12}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Int64ToBStr(tt.args.num, tt.args.bit), "Int64ToBStr(%v, %v)", tt.args.num, tt.args.bit)
		})
	}
}
