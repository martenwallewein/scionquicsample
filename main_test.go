package main

import (
	"fmt"
	"net"
	"testing"

	"github.com/netsec-ethz/scion-apps/pkg/appnet"
)

func Test_SCIONConn_Listen(t *testing.T) {

	type args struct {
		remote string
		local  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				remote: "127.0.0.1:41000",
				local:  "127.0.0.1:40000",
			},
			want: "5_50",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serverAddr, err := net.ResolveUDPAddr("udp", tt.args.local)
			if err != nil {
				t.Error(err)
			}
			c, err := appnet.Listen(serverAddr)
			if err != nil {
				t.Error(err)
			}

			fmt.Println(c.LocalAddr())

		})
	}
}
