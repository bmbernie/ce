package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const (
	SSH_CLIENT_VERSION = "SSH-2.0-OpenSSH_7.1\r\n"
	SSH_PORT           = "22"
	BASE10             = 10
	INT_BIT_SIZE       = 64
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: %s <host>[:port]", os.Args[0])
		os.Exit(1)
	}

	for _, target := range os.Args[1:] {
		if _, port, err := net.SplitHostPort(target); err != nil {
			target = net.JoinHostPort(target, SSH_PORT)
		} else if _, err := strconv.ParseInt(port, BASE10, INT_BIT_SIZE); err != nil {
			fmt.Println("Invalid port: %s", port)
			os.Exit(1)
		}

		conn, err := net.Dial("tcp", target)
		if err != nil {
			fmt.Println("dial error:", err)
			os.Exit(1)
		}
		defer conn.Close()
		fmt.Println("Connection Established to %s", target)

		fmt.Fprintf(conn, SSH_CLIENT_VERSION)
		fmt.Printf("->\t%s", SSH_CLIENT_VERSION)

		SSH_SERVER_VERSION, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("<-\t%s", SSH_SERVER_VERSION)

		var buf bytes.Buffer
		io.Copy(&buf, conn)
		fmt.Printf("<-\t%#v", buf)
	}
}
