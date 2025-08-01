package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStepFloat(t *testing.T) {
	type args struct {
		num  float64
		step float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{args: args{num: 0.12, step: 0.001}, want: 0.12},
		{args: args{num: 235.200, step: 0.001}, want: 235.200},
		{args: args{num: 235.7000001, step: 0.001}, want: 235.7000001},
		{args: args{num: 235.8, step: 0.001}, want: 235.8},
		{args: args{num: 0.12, step: 0.002}, want: 0.12},
		{args: args{num: 235.200, step: 0.002}, want: 235.200},
		{args: args{num: 235.7000001, step: 0.002}, want: 235.700000},
		{args: args{num: 235.8, step: 0.002}, want: 235.8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, StepFloat(tt.args.num, tt.args.step), "StepFloat(%v, %v)", tt.args.num, tt.args.step)
		})
	}
}
