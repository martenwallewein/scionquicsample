package main

import (
	"testing"
)

func Test_keyOp_Generate(t *testing.T) {

	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				x: 5,
				y: 50,
			},
			want: "5_50",
		},
		{
			name: "success large integers",
			args: args{
				x: 50000,
				y: 999999,
			},
			want: "50000_999999",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kp := GetKeyOperator()
			if got := kp.Generate(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("keyOp.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
