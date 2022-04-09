package epp

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"time"
)

const TIMEOUT = 5

func Tlsconnect(certfile string, keyfile string, host string, port int) (*tls.Conn, error) {
	cert, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	hostport := fmt.Sprintf("%s:%d", host, port)
	conn, err := tls.Dial("tcp", hostport, &config)
	return conn, err
}

// Read a message from an epp server.
func Read(conn *tls.Conn) (string, error) {
	data := make([]byte, 8192)
	var buffer []byte
	var msgsize uint32 = 0 // -1 i.e we don't have a size yet
	var rbytes uint32 = 0  // Bytes received so far.
	for {
		conn.SetReadDeadline(time.Now().Add(TIMEOUT * time.Second))
		n, err := conn.Read(data)
		if err != nil {
			return "", fmt.Errorf("Failed to read from the EPP server", err)
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
func Write(conn *tls.Conn, message string) (int, error) {
	header := make([]byte, 4)
	conn.SetWriteDeadline(time.Now().Add(TIMEOUT * time.Second))
	binary.BigEndian.PutUint32(header, uint32(len(message))+4)
	msg := []byte(message)
	out := append(header, msg...)
	n, err := conn.Write(out)
	return n, err
}

// Utility. Send a Message and return the response.
func Frame(conn *tls.Conn, messsage string) (string, error) {
	_, err := Write(conn, messsage)
	if err != nil {
		return "", err
	}
	return Read(conn)
}
