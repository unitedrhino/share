package utils

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"testing"
)

func TestCopy(t *testing.T) {
	type args struct {
		toValue   interface{}
		fromValue interface{}
	}
	type (
		Str1 struct {
			Str string
		}
		Str2 struct {
			Str *wrappers.StringValue
		}
		Str3 struct {
			Str2
		}
	)

	tests := []struct {
		name string
		args args
	}{
		{name: "string-*string", args: args{
			toValue:   &Str2{},
			fromValue: &Str1{Str: "123"},
		}},
		{name: "string-*string", args: args{
			toValue:   &Str3{},
			fromValue: &Str1{Str: "123"},
		}},
	}
	for _, tt := range tests {
		err := CopyE(tt.args.toValue, tt.args.fromValue)
		fmt.Println(err)
	}
}
