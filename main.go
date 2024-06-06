package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

var SERVER_ADDY string = "127.0.0.1:2409"

var (
	version byte = byte(1)
	request byte = byte(0)

	contentType byte
	content     string

	path string

	clientAddy string
	inc        bool

	send bool
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[byte]().
				Title("Request").
				Options(
					huh.NewOption("GET", byte(1)),
					huh.NewOption("POST", byte(2)),
				).
				Value(&request),

			huh.NewInput().
				Title("Endpoint").
				Value(&path),

			huh.NewSelect[byte]().
				Title("Content Type").
				Options(
					huh.NewOption("plaintext", byte(0)),
					huh.NewOption("code_", byte(1)),
					huh.NewOption("code_html", byte(2)),
					huh.NewOption("code_css", byte(3)),
					huh.NewOption("code_lua", byte(4)),
					huh.NewOption("code_json", byte(5)),
				).
				Value(&contentType),

			huh.NewInput().
				Title("Content:").
				Value(&content),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", SERVER_ADDY)

	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Fatal(err)
	}

	bytes := createBytes(byte(1), request, contentType, content, "127.0.0.1", true, path)

	_, err = conn.Write(bytes)

	conn.Close()
}

func createBytes(
	version byte,
	reqType byte,
	contentType byte,
	content string,
	client_addy string,
	inc bool,
	path string,
) []byte {
	contentLen := make([]byte, 4)

	binary.BigEndian.PutUint32(contentLen, uint32(len(content)))

	var bytes []byte

	bytes = append(bytes, version)
	bytes = append(bytes, reqType)
	spl := strings.Split(client_addy, ".")

	addr := make([]byte, 4)
	for i, s := range spl {
		u, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			fmt.Println(err)
			continue
		}

		addr[i] = byte(u)
	}

	bytes = append(bytes, addr...)

	if inc {
		bytes = append(bytes, byte(1))
	} else {
		bytes = append(bytes, byte(0))
	}

    if len(path) == 0 {
        path = "/"
    }

	pathLen := make([]byte, 2)
	binary.BigEndian.PutUint16(pathLen, uint16(len(path)))

	bytes = append(bytes, pathLen...)
	bytes = append(bytes, []byte(path)...)

	bytes = append(bytes, contentLen...)
	bytes = append(bytes, contentType)
	bytes = append(bytes, []byte(content)...)

	for i, b := range bytes {
		fmt.Printf("Byte %d = %d\n", i, b)
	}

	return bytes
}
