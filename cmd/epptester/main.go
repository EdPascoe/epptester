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

package main

// runtime "github.com/banzaicloud/logrus-runtime-formatter"
import (
	"epptester/pkg/cmd/epptester"
	_ "epptester/pkg/log"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var Version = "dev"
var GitTag = "-"

// The flag help message.
func Usage() {
	w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
	fmt.Fprintln(w, `EPPTester 
A really simple tool for testing connections to an EPP server.  
Cert and Key can be in the same file.. Just specify twice.
`)
	flag.PrintDefaults()
}

func main() {
	logrus.Debugln("Start")
	certFile := flag.String("cert", "cert.pem", "certificate PEM file")
	keyFile := flag.String("key", "key.pem", "key PEM file")
	host := flag.String("host", "127.0.0.1", "Epp server ")
	port := flag.Int("port", 3121, "Server port ")
	username := flag.String("username", "", "EPP Username")
	password := flag.String("password", "", "EPP Password")
	version := flag.Bool("version", false, "Show version")
	tlsversion := flag.String("tls", "any", "TLS version to test with: 1.0, 1.1, 1.2, 1.3, any")
	flag.Usage = Usage
	flag.Parse()
	if *version {
		fmt.Println("epptester", Version)
		fmt.Println("Git version.\t", GitTag)
		os.Exit(0)
	}
	epptester.RunTest(*certFile, *keyFile, *host, *port, *username, *password, *tlsversion)
}
