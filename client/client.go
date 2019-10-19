package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"io"
	"log"
	"net"
)

const localAddr = "127.0.0.1:1080"

var host = flag.String("host", "host:port", "host port")

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func main() {
	flag.Parse()
	log.Fatal(client())
}

func client() error {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		fmt.Printf("listen fail, err: %v\n", err)
		return err
	}
	log.Print("socks5 bind: ", localAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}
		//create goroutine for each connect
		go process(conn)
	}
}
func process(conn net.Conn) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(*host, tlsConf, nil)
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
