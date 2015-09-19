// +build ignore
package irc

// This probably belongs to the 'bot' package

import (
	"net"
	"io/ioutil"
	"io"
	"time"
)

const (
	port = "8082"
)

// Extension header
type xHeader {
	Version int // protocol version
	
}

// Extension is a representation of a registered extension
type Extension struct {
	s io.ReadWriter
	Version string
	Description string
	Keywords []string
}

// server will run a TCP server on port 'port' and will listen
// for incoming extension connections
func server() error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		if err = registerExtension(conn); err != nil {
			log.Printf("Error registering extension: %v", err)
		}
	}
}

func handshake(conn net.Conn) error {
	r, err := ioutil.ReadAll(x.s)
}

func registerExtension(conn net.Conn) error {
	// Get extension information
	if err := handshake(conn); err != nil {
		return err
	}
}

// Query sends a message to an extension and waits for a reply
// or times out.
// 
func (x *Extension) Query(from, message string) (string, error) {
	var (
		resp chan string = make(chan string)
		line string
	)
	
	if _, err := x.s.Write([]byte{fmt.Sprintf("!%s :%s")}); err != nil {
		return err
	}
	go func() {
		r, err := ioutil.ReadAll(x.s)
		resp <- r
	}()
	select {
	case <- time.After(2*time.Second):
		
	case line = <- resp:
		return line, nil
	}
	return errors.New("timed out")
}

