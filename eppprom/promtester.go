package eppprom

import (
	"epptester/internal/promexport"
	"epptester/pkg/epp"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

const retrycount = 3

type Config struct {
	Refresh  int `yaml:"refresh,omitempty"`
	Promfile string
	Cert     string
	Key      string
	Zones    []struct {
		Name       string
		Host       string
		Port       int    `yaml:"port,omitempty"`
		Tlsversion string `yaml:"tlsversion,omitempty"`
	}
	Configfile string `yaml:"__notreal__configfile,omitempty"` // We need the config file for the certificates
}

func Loadconfig(configfile string, config *Config) {
	buffer, err := os.ReadFile(configfile)
	if err != nil {
		log.Panic(err)
		log.Fatalln("Failed to read", configfile, ":", err)
	}
	err = yaml.Unmarshal(buffer, config)
	if err != nil {
		fmt.Println(config)
		log.Panic(err)
		log.Fatalln("Failed to process yaml in ", configfile, ":", err)
	}
	config.Configfile = configfile
	for i, _ := range config.Zones {
		zone := &config.Zones[i] // We need to see the original zone field. Not a copy
		if zone.Tlsversion == "" {
			zone.Tlsversion = "any"
		}
		if zone.Port == 0 {
			config.Zones[i].Port = 700 // Default IANA port for EPP
		}
	}
	// fmt.Printf("--- t:\n%v\n\n", *config)
}

func Testzones(config *Config) {
	prom, err := promexport.NewPromwriter(config.Promfile)
	defer prom.Close()
	if err != nil {
		logrus.Fatal("Failed to start prometheus logging to ", config.Promfile, ":", err)
	}

	status := "ok"
	qtime := 0
	for _, zone := range config.Zones {
		logrus.Info("Testing zone ", zone.Name)
		for i := 1; i < retrycount; i++ {
			status = "ok"
			qtime = 0
			tstart := time.Now()
			conn, err := epp.Tlsconnectstring(config.Cert, config.Key, zone.Host, zone.Port, zone.Tlsversion)
			if err != nil {
				logrus.Error("Failed to connect to ", zone, ":", err)
				status = "ERROR_CONNECT"
				qtime = 0
				continue
			}
			eppsession := &epp.Session{Conn: conn}
			header, err := eppsession.Read() // The first frame is the epp greeting message.
			if err != nil {
				logrus.Fatal("Failed to read header for ", zone, ":", err)
				status = "ERROR_NOHEADER"
				qtime = 0
				continue
			}
			err = eppsession.Buildgreeting(header)
			// fmt.Println(eppsession.Greeting.Svid)
			qtime = int(time.Since(tstart) / time.Millisecond)
			logrus.Infoln(zone.Name, " Pass ", qtime, " ms")
			break // No need to retry
		}
		prom.Writeresult(zone.Name, status, qtime)
	}
}

// dns_query{zone="koeln", nameserver_ip="194.0.11.114", rcode="NOERROR", hostname="monitor21.dns.net.za"} 192
