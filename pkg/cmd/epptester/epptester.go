package epptester

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"epptester/pkg/epp"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
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

func servercertcheck(conn *tls.Conn) {
	tlsversions := map[uint16]string{
		tls.VersionSSL30: "SSL",
		tls.VersionTLS10: "TLS 1.0",
		tls.VersionTLS11: "TLS 1.1",
		tls.VersionTLS12: "TLS 1.2",
		tls.VersionTLS13: "TLS 1.3",
	}

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		// fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println("Certificate: ", v.Subject) // , v.Verify())
		fmt.Printf("   Expiry: %s", v.NotAfter)
		if time.Now().After(v.NotAfter) {
			color.Red(" *** WARNING *** Certificate has Expired")
		} else {
			fmt.Println("")
		}
	}
	// fmt.Println("client: handshake: ", state.HandshakeComplete)
	tlsv := conn.ConnectionState().Version
	fmt.Println("TLS version: ", tlsversions[tlsv])
}

func clientcertcheck(certfile string) {
	fmt.Println("\n=============== Client side certificate check =====================")
	rootPEM, _ := ioutil.ReadFile(certfile)
	block, _ := pem.Decode([]byte(rootPEM))
	if block == nil {
		log.Fatalln("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalln("failed to parse certificate: ", err.Error())
	}

	// fmt.Println("\nCert: DNS:", cert.DNSNames)
	// fmt.Println("Issuer:", cert.Issuer)
	fmt.Println("Subject:", cert.Subject)
	fmt.Printf("    Expiry: %s", cert.NotAfter)
	if time.Now().After(cert.NotAfter) {
		color.Red(" *** WARNING *** Certificate has Expired")
	} else {
		fmt.Println("")
	}
	fingerprint := md5.Sum(cert.Raw)
	var buf bytes.Buffer
	for i, f := range fingerprint {
		if i > 0 {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%02X", f)
	}
	fmt.Println("Fingerprint:", buf.String())

}

// Use the google myaddr hack to find this servers public ip
func Findmyip() string {
	// See https://unix.stackexchange.com/questions/22615/how-can-i-get-my-external-ip-address-in-a-shell-script
	searchserver := "ns1.google.com:53"
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, "udp", searchserver)
		},
	}
	textrecords, err := r.LookupTXT(context.Background(), "o-o.myaddr.l.google.com")
	if err != nil {
		fmt.Errorf("Failed to check dns: %s", err)
		return ""
	}
	if len(textrecords) > 0 {
		return textrecords[0]
	} else {
		return ""
	}
	// print(ip[0])

}

// Find the ip address of the server we are going to connect to
func Findserverip(host string) string {
	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		return ""
	}
	if len(ips) > 0 {
		return fmt.Sprintf("%s", ips[0])
	} else {
		return ""
	}
}

// Basic checks on ip addresses
func dnschecks(host string) {
	ip := Findmyip()
	fmt.Println("=============== IP check =====================")
	fmt.Println("Your ip address is probably ", ip)
	ip = Findserverip(host)
	fmt.Println("Connecting to ", ip)
}

func logintest(session *epp.Session, username string, password string) {
	fmt.Println("\n=============== login check =====================")
	login, err := session.Login(username, password)
	if err != nil {
		log.Fatalf("Login failed with %s", err)
	}
	fmt.Println("Logged in ", login.Result.Msg)
}

func RunTest(certfile string, keyfile string, host string, port int, username string, password string, tlsversion string) {
	dnschecks(host)
	clientcertcheck(certfile)
	fmt.Println("\n=============== Server side checks =====================")
	session, err := epp.SessionStart(certfile, keyfile, host, port, tlsversion)
	if err != nil {
		fmt.Errorf("Faled to connect properly: ", err)
		log.Fatalf("Epp session failure: %s", err)
	}
	// fmt.Println("TLS session connected", session.Conn.ConnectionState().Version)
	defer session.Close()
	servercertcheck(session.Conn)
	fmt.Printf("Server: %s  -- %s\n", session.Greeting.Svid, session.Greeting.Svdate)
	fmt.Println(" ")
	logintest(session, username, password)
}
