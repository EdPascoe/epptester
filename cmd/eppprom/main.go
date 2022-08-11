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

import (
	"epptester/eppprom"
	_ "epptester/internal/log"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Version = "dev"
var GitTag = "-"

// The flag help message.
func Usage() {
	w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
	fmt.Fprintln(w, `epppromtester
Testing remote EPP connection and exporting result to prometheus.
`)
	flag.PrintDefaults()
}

func main() {
	logrus.Debugln("Start")
	configfile := flag.String("C", "eppprom.yaml", "Yaml config file")
	version := flag.Bool("version", false, "Show version")
	flag.Usage = Usage
	flag.Parse()
	if *version {
		fmt.Println("epptester", Version)
		fmt.Println("Git version.\t", GitTag)
		os.Exit(0)
	}
	// epptester.RunTest(*certFile, *keyFile, *host, *port, *username, *password, *tlsversion)
	config := &eppprom.Config{}
	eppprom.Loadconfig(*configfile, config)
	for true {
		eppprom.Testzones(config)
		if config.Refresh == 0 {
			break
		} else {
			time.Sleep(time.Duration(config.Refresh) * time.Second)
		}
	}
	os.Exit(0)

}
