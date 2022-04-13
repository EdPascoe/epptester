package epp

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
)

// Login to epp server.

type Session struct {
	Conn        *tls.Conn
	Greeting    *Eppgreeting
	Lastmessage string
}

type Eppgreetingroot struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Greeting Eppgreeting `xml:"urn:ietf:params:xml:ns:epp-1.0 greeting"`
}
type Eppgreeting struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 greeting"`
	Svid    string   `xml:"urn:ietf:params:xml:ns:epp-1.0 svID"`
	Svdate  string   `xml:"urn:ietf:params:xml:ns:epp-1.0 svDate"`
}

//<epp:epp xmlns:epp="urn:ietf:params:xml:ns:epp-1.0">
//	<epp:response>
//		<epp:result code="1000">
//			<epp:msg>Access granted</epp:msg>
//		</epp:result>
//		<epp:trID>
//			<epp:svTRID>COZA-EPP-180227F2C85-28068</epp:svTRID>
//		</epp:trID>
//	</epp:response>
//</epp:epp>

type Eppresponseroot struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Response Eppresponse `xml:"urn:ietf:params:xml:ns:epp-1.0 response"`
}
type Eppresponse struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 response"`
	Result Eppresult `xml:"urn:ietf:params:xml:ns:epp-1.0 result"`
}
type Eppresult struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 result"`
	Code    int   `xml:"code,attr"`
	Msg string `xml:"urn:ietf:params:xml:ns:epp-1.0 msg"`
}

type EppResponseError struct{
	Code int
	Msg string
}
func (m *EppResponseError) Error() string {
	return fmt.Sprintf("%d -- %s", m.Code, m.Msg)
}


// Convert the epp greeting into a struct.
func buildgreeting(session *Session, header string) error{
	root := Eppgreetingroot{}
	err := xml.Unmarshal([]byte(header), &root)
	if err != nil {
		logrus.Error("Failed to unmarshall the epp header", err, "\n",header)
		return err
	}
	session.Greeting = &root.Greeting
	return nil
}

func Login(session *Session, username string, password string) (*Eppresponse, error) {
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
		logrus.Error("Failed to login %s ", err )
		return nil, err
	}
	session.Lastmessage = msg
	root := new(Eppresponseroot)
	err = xml.Unmarshal([]byte(msg), &root)
	if err != nil {
		logrus.Error("Failed to unmarshall the epp login response", err, "\n",msg)
		return nil, err
	}
	if root.Response.Result.Code != 1000 {
		err := EppResponseError{}
		err.Code = root.Response.Result.Code
		err.Msg = root.Response.Result.Msg
		return nil, &err
	}
	return &root.Response, err
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
	err = buildgreeting(&eppsession, header)
	return &eppsession, err
}
