package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"math/big"
	"net"

	socks5 "github.com/armon/go-socks5"
	"github.com/lucas-clemente/quic-go"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

const addr = "0.0.0.0:443"

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	log.Fatal(echoServer())
}

// Start a server that echos all data on the first stream opened by the client
func echoServer() error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	log.Print("listenning:", addr)
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		go func() {
			stream, err := sess.AcceptStream(context.Background())
			if err != nil {
				log.Print(err.Error())
				return
			}
			s, err := socks5.New(&socks5.Config{})
			if err != nil {
				log.Print(err.Error())
				return
			}
			err = s.ServeConn(Stream{stream})
			if err != nil {
				log.Print(err.Error())
				return
			}
		}()
	}
	return nil
}

type Stream struct {
	quic.Stream
}

func (Stream) LocalAddr() net.Addr {
	return nil
}

func (Stream) RemoteAddr() net.Addr {
	return nil
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
