package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"io"
	"log"
	"net"
)

func NewClient(client, host string) error {
	listener, err := net.Listen("tcp", client)
	if err != nil {
		fmt.Printf("listen fail, err: %v\n", err)
		return err
	}
	log.Print("socks5 bind: ", client)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}
		//create goroutine for each connect
		go process(conn, host)
	}
}
func process(conn net.Conn, host string) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(host, tlsConf, nil)
	if err != nil {
		log.Print(err)
		return
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Print(err)
		return
	}
	go func() {
		_, err := io.Copy(stream, conn)
		log.Print(err)
	}()
	go func() {
		_, err := io.Copy(conn, stream)
		log.Print(err)
	}()
}
