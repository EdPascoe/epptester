package epptester

import (
	"crypto/tls"
	"fmt"
	"github.com/EdPascoe/epptester/pkg/epp"
	"github.com/sirupsen/logrus"
	"log"
)

/*
Copyright 2022 Ed Pascoe.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
const VERSION = "1.0"
const TIMEOUT = 5

// Check the epp server
func Checkepp(conn *tls.Conn) {
	//fmt.Println("=============== EPP HEADER =====================")
	//header, err := read(conn) // The first frame is the epp greeting message.
	//if err != nil {
	//	log.Fatalf("Failed to get header %s", err)
	//}
	//fmt.Println("Header message:", header)
}

func Checkcert(conn *tls.Conn) {
	fmt.Println("=============== Certificate check =====================")
	log.Println("client: connected to: ", conn.RemoteAddr())
	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		// fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println("Certificate: ", v.Subject) // , v.Verify())
	}
	//log.Println("client: handshake: ", state.HandshakeComplete)
	//log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
	return
	//
	//message := "data"
	//n, err := io.WriteString(conn, message)
	//if err != nil {
	//	log.Fatalf("client: write: %s", err)
	//}
	//log.Printf("client: wrote %q (%d bytes)", message, n)

	//reply := make([]byte, 256)
	//n, err = conn.Read(reply)
	//log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	//log.Print("client: exiting")
}

func RunTest(certfile string, keyfile string, host string, port int, username string, password string) {
	session, err := epp.Connect(certfile, keyfile, host, port)
	if err != nil {
		log.Fatalf("Epp session failure: %s", err)

	}
	logrus.Infoln("Connected")
	defer session.Conn.Close()
	err = epp.Login(session, username, password)
	if err != nil {
		log.Fatalf("Login failed with %s", err)
	}
	logrus.Infoln("Logged in")
	fmt.Println(session.Lastmessage)

	//Checkcert(conn)
	//Checkepp(conn)
	//epp.Login(conn, "dnservices", "e4a192c3")
	// fmt.Println("Login: ", login)
}
