package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetEndTime(t *testing.T) {
	type args struct {
		d time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{args: args{d: time.Now()}, want: time.Now()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetEndTime(tt.args.d), "GetEndTime(%v)", tt.args.d)
		})
	}
}
