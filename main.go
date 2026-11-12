package main

import (
	"dorsone/goredis/internal"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := internal.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		_ = value

		writer := internal.NewWriter(conn)
		err = writer.Write(internal.NewValue("OK"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
