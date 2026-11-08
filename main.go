package main

import (
	"dorsone/goredis/internal"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

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
			if err == io.EOF {
				break
			}

			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		fmt.Println(value)

		_, _ = conn.Write([]byte("+OK\r\n"))
	}
}
