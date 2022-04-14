package epp

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"time"
)

const TIMEOUT = 5

func tlsconnect(certfile string, keyfile string, host string, port int) (*tls.Conn, error) {
	cert, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	hostport := fmt.Sprintf("%s:%d", host, port)
	conn, err := tls.Dial("tcp", hostport, &config)
	return conn, err
}

type Session struct {
	Conn        *tls.Conn
	Greeting    *Eppgreeting
	Lastmessage string
}

func Connect(certfile string, keyfile string, host string, port int) (*Session, error) {
	conn, err := tlsconnect(certfile, keyfile, host, port)
	if err != nil {
		return &Session{}, err
	}
	eppsession := &Session{Conn: conn}
	header, err := eppsession.Read() // The first frame is the epp greeting message.
	if err != nil {
		return eppsession, fmt.Errorf("Failed to get header %s", err)
	}
	err = eppsession.buildgreeting(header)
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
			return "", fmt.Errorf("Failed to read from the EPP server: %s", err)
		}
		if msgsize == 0 { // This is the first frame.
			if n < 4 {
				return "", fmt.Errorf("The initial read only returned %s bytes. Aborting. ", n)
			}
			msgsize = binary.BigEndian.Uint32(data[0:4]) - 4
			buffer = data[4:]
			rbytes = msgsize
		} else {
			buffer = append(buffer, data[0:]...)
			rbytes = rbytes + uint32(len(data))
		}
		if rbytes == msgsize {
			// In go strings a utf-8 bytes already.
			return string(buffer), nil
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
