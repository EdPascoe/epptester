package epp

import (
	"crypto/tls"
	"fmt"
)

// Login to epp server.

type Session struct {
	Conn        *tls.Conn
	Header      string
	Lastmessage string
}

func Login(session *Session, username string, password string) error {
	template := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
   <command>
      <login>
         <clID>%s</clID>
         <pw>%s</pw>
         <options>
            <version>1.0</version>
            <lang>en</lang>
         </options>
         <svcs>
            <objURI>urn:ietf:params:xml:ns:domain-1.0</objURI>
            <objURI>urn:ietf:params:xml:ns:host-1.0</objURI>
            <objURI>urn:ietf:params:xml:ns:contact-1.0</objURI>
         </svcs>
      </login>
   </command>
</epp>
`, username, password)
	msg, err := Frame(session.Conn, template)
	if err != nil {
		return err
	}
	session.Lastmessage = msg
	return err
}

func Connect(certfile string, keyfile string, host string, port int) (*Session, error) {
	conn, err := Tlsconnect(certfile, keyfile, host, port)
	if err != nil {
		return &Session{}, err
	}
	eppsession := Session{Conn: conn}
	header, err := Read(conn) // The first frame is the epp greeting message.
	if err != nil {
		return &eppsession, fmt.Errorf("Failed to get header %s", err)
	}
	eppsession.Header = header
	return &eppsession, nil
}
