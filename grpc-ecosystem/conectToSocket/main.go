package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

var (
	socketDir  = "/var/lib/kubelet/pod-resources"
	socketPath = socketDir + "/kubelet.sock"

	connectionTimeout = 60 * time.Second
)

func main() {
	fmt.Println(os.Args)
	socketPath = os.Args[len(os.Args)-1]
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, socketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	defer func() {
		if conn != nil {
			log.Println(fmt.Sprintf("connection state: %s", conn.GetState().String()))
			if er := conn.Close(); er != nil {
				log.Println(fmt.Sprintf("failure closing connection to %s: %v", socketPath, er))
			}
		} else {
			log.Println(fmt.Sprintf("connection to %s is nil", socketPath))
		}
	}()

	if err != nil {
		log.Println(fmt.Sprintf("failure connecting to %s: %v", socketPath, err))
		return
	}
}
