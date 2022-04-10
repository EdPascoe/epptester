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
	"flag"
	"fmt"
	"github.com/EdPascoe/epptester/pkg/cmd/epptester"
	_ "github.com/EdPascoe/epptester/pkg/log"
	"github.com/sirupsen/logrus"
)

// The flag help message.
func Usage() {
	w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
	fmt.Fprintln(w, `EPPTester 
A really simp tool for testing connections to an EPP server.  Cert and Key can be in the same file.. Just speicify twice.
`)
	flag.PrintDefaults()
}

func main() {
	// log.SetFormatter(&log.JSONFormatter{})
	// loggingSetup()
	logrus.Debugln("Start")
	certFile := flag.String("cert", "cert.pem", "certificate PEM file")
	keyFile := flag.String("key", "key.pem", "key PEM file")
	host := flag.String("host", "127.0.0.1", "Epp server ")
	port := flag.Int("port", 3121, "Server port ")
	username := flag.String("username", "", "EPP Username")
	password := flag.String("password", "", "EPP Password")
	flag.Usage = Usage
	flag.Parse()

	epptester.RunTest(*certFile, *keyFile, *host, *port, *username, *password)
	//epptester.Checkcert(conn)
	//epptester.Checkepp(conn)
	//epptester.Login(conn, "dnservices", "e4a192c3")
	// fmt.Println("Login: ", login)
}
