package main

import (
	"dorsone/goredis/internal"
	"fmt"
	"net"
	"strings"
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

		if value.GetType() != internal.ARRAY {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.GetArray()) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.GetArray()[0].GetBulk())
		args := value.GetArray()[1:]

		fmt.Println(command, args)

		writer := internal.NewWriter(conn)

		handler, ok := internal.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			err = writer.Write(internal.NewValue(""))
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		result := handler(args)
		fmt.Println(result)
		err = writer.Write(result)
		if err != nil {
			fmt.Println(err)
		}
	}
}
