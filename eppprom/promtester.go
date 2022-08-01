package eppprom

import (
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
	Zones      []Zone `yaml:"zones,omitempty"`
	Cert       string `yaml:"cert,omitempty"`
	Key        string `yaml:"key,omitempty"`
	Configfile string `yaml:"__configfile,omitempty"` // We need the config file for the certificates
}
type Zone struct {
	Name       string `yaml:"name,omitempty"`
	Promfile   string `yaml:"promfile,omitempty"`
	Host       string `yaml:"host,omitempty"`
	Port       int    `yaml:"port,omitempty"`
	Username   string `yaml:"username,omitempty"`
	Password   string `yaml:"password,omitempty"`
	Tlsversion string `yaml:"tlsversion,omitempty"`
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
	for i, zone := range config.Zones {
		if zone.Tlsversion == "" {
			config.Zones[i].Tlsversion = "any" // zone is a copy so we need to see the original.
		}
		if zone.Port == 0 {
			config.Zones[i].Port = 3121 // Default port
		}
	}
	// fmt.Printf("--- t:\n%v\n\n", *config)
}

func Testzones(config *Config) {
	for _, zone := range config.Zones {
		logrus.Info("Testing zone", zone.Name)
		for i := 1; i < retrycount; i++ {
			tstart := time.Now()
			conn, err := epp.Tlsconnectstring(config.Cert, config.Key, zone.Host, zone.Port, zone.Tlsversion)
			if err != nil {
				logrus.Fatal("Failed to connect to ", zone, ":", err)
			}
			eppsession := &epp.Session{Conn: conn}
			header, err := eppsession.Read() // The first frame is the epp greeting message.
			if err != nil {
				logrus.Fatal("Failed to read header for ", zone, ":", err)
			}
			err = eppsession.Buildgreeting(header)
			fmt.Println(eppsession.Greeting.Svid)
			duration := int32(time.Since(tstart) / time.Millisecond)
			logrus.Infoln("Run time was ", duration)
		}
	}
}

// dns_query{zone="koeln", nameserver_ip="194.0.11.114", rcode="NOERROR", hostname="monitor21.dns.net.za"} 192
