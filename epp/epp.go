package epp

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const TIMEOUT = 5

// Do a tls connection with cert and key as strings.
func Tlsconnectstring(certstr string, keystr string, host string, port int, tlsversion string) (*tls.Conn, error) {
	cert, err := tls.X509KeyPair([]byte(certstr), []byte(keystr))
	if err != nil {
		return nil, err
	}
	return Tlsconnect(cert, host, port, tlsversion)
}

// Do a tls connection with cert and key as files.
func Tlsconnectfile(certfile string, keyfile string, host string, port int, tlsversion string) (*tls.Conn, error) {
	cert, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}
	return Tlsconnect(cert, host, port, tlsversion)
}

// The actual tls certificate.
func Tlsconnect(cert tls.Certificate, host string, port int, tlsversion string) (*tls.Conn, error) {
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	switch tlsversion {
	case "1.0":
		config.MinVersion = tls.VersionTLS10
		config.MaxVersion = tls.VersionTLS10
	case "1.1":
		config.MinVersion = tls.VersionTLS11
		config.MaxVersion = tls.VersionTLS11
	case "1.2":
		config.MaxVersion = tls.VersionTLS12
		config.MaxVersion = tls.VersionTLS12
	case "1.3":
		config.MinVersion = tls.VersionTLS13
		config.MaxVersion = tls.VersionTLS13
	case "any":
		// Keep the default list of cipher suites
	default:
		panic(fmt.Sprint("Unable to handle TLS version >", tlsversion, "<"))
	}
	// logrus.Warning("Minversion: ", config.MinVersion)
	hostport := fmt.Sprintf("[%s]:%d", host, port)
	conn, err := tls.Dial("tcp", hostport, &config)
	return conn, err
}

type Session struct {
	Conn        *tls.Conn
	Greeting    *Eppgreeting
	Lastmessage string
}

func SessionStart(certfile string, keyfile string, host string, port int, tlsversion string) (*Session, error) {
	conn, err := Tlsconnectfile(certfile, keyfile, host, port, tlsversion)
	if err != nil {
		fmt.Errorf("Failed to connect; %s", err)
		return &Session{}, err
	}
	eppsession := &Session{Conn: conn}
	header, err := eppsession.Read() // The first frame is the epp greeting message.
	if err != nil {
		return eppsession, fmt.Errorf("Failed to get header %s", err)
	}
	// fmt.Println("Header ", header, "Error: ", err, "Len:", len(header))
	err = eppsession.Buildgreeting(header)
	return eppsession, err
}

// Read a message from an epp server.
func (s *Session) Read() (string, error) {
	data := make([]byte, 8192)
	var buffer []byte
	var msgsize uint32 = 0 // -1 i.e we don't have a size yet
	var rbytes uint32 = 0  // Bytes received so far.
	for {
		s.Conn.SetReadDeadline(time.Now().Add(TIMEOUT * time.Second))
		n, err := s.Conn.Read(data)
		if err != nil {
			logrus.Error("Failed eppserver ", err)
			return "", fmt.Errorf("Failed to read from the EPP server: %v", err)
		}
		// logrus.Info("Got data ", n, "bytes")
		if msgsize == 0 { // This is the first frame.
			for n < 4 {
				// Dodgy epp server. Didn't get enough bytes for a length. Trying again.
				data1 := make([]byte, 8192)
				n1, err := s.Conn.Read(data1)
				if err != nil {
					logrus.Error("Failed to read while getting message length ", err)
					return "", fmt.Errorf("Failed while getting message length: %v", err)
				}
				data = append(data[0:n], data1...)
				n = n + n1
			}
			msgsize = binary.BigEndian.Uint32(data[0:4]) - 4
			buffer = data[4:]
			rbytes = uint32(n) - 4
			// logrus.Infof("First read n %v msgsize: %v rbytes %v  ", n, msgsize, rbytes)
		} else {
			buffer = append(buffer, data[0:]...)
			rbytes = rbytes + uint32(n)
			// logrus.Infof("Appended n was %v msgsize: %v rbytes %v ", n, msgsize, rbytes)
		}
		if rbytes >= msgsize {
			if rbytes > msgsize {
				logrus.Warning("Received more data than expected. Ignoring %v", string(buffer[msgsize:]))
			}
			// In go strings a utf-8 bytes already.
			return string(buffer[:msgsize]), nil
		}
	}
}

// Send message to server
func (s *Session) Write(message string) (int, error) {
	header := make([]byte, 4)
	s.Conn.SetWriteDeadline(time.Now().Add(TIMEOUT * time.Second))
	binary.BigEndian.PutUint32(header, uint32(len(message))+4)
	msg := []byte(message)
	out := append(header, msg...)
	n, err := s.Conn.Write(out)
	return n, err
}

// Utility. Send a Message and return the response.
func (s *Session) Frame(messsage string) (string, error) {
	_, err := s.Write(messsage)
	if err != nil {
		return "", err
	}
	return s.Read()
}

// Close the socket connection
func (s *Session) Close() error {
	return s.Conn.Close()
}
