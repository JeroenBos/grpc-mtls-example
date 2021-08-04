package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/islishude/grpc-mtls-example/greet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func LoadKeyPair() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		panic("Load client certification failed: " + err.Error())
	}

	ca, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		panic("can't read ca file")
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(ca) {
		panic("can't add CA cert")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      capool,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(tlsConfig)
}

func main() {
	// host := "localhost" // if running from host
	host := "server" // if running from client docker
	ip := host + ":10200" 
	
	time.Sleep(2 * time.Second) // give server time to start
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(LoadKeyPair()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := greet.NewGreetingClient(conn)
	resp, err := client.SayHello(context.Background(), &greet.SayHelloRequest{Name: "world"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.GetGreet())
}
