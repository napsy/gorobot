// +build ignore

package irc

// This probably belongs to the 'bot' package

import (
	"net"
	"fmt"
	"io/ioutil"
	"io"
	"errors"
	"time"
	"path/filepath"
	"os/exec"
)

const (
	port = "8082"
	headerVersion = "gorobot 0.1"

	extensionsKeywords map[string]*Extension
)

// Extension header
type xHeader {
	Version int // protocol version
	Pid int // process PID
}

// Extension is a representation of a registered extension
type Extension struct {
	rw io.ReadWriter
	Version string
	Description string
	Keywords []string

	cmd *exec.Cmd
}

// server will run a TCP server on port 'port' and will listen
// for incoming extension connections
func runServer() error {
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

func xHandshake(r io.Reader) error {
	jsonHeader, err := ioutil.ReadAll(r)
	h := xHeader{}
	if err := json.Unmarshal([]byte(jsonHeader), &h); err != nil {
		return err
	}
	if h.Version != headerVersion {
		return errors.New("invalid protocol version")
	}
	return nil
}

func xAbout(r io.Reader) (*Extension, error) {
	jsonExtension, err := ioutil.ReadAll(r)
	extension := Extension{}
	if err := json.Unmarshal([]byte(jsonExtension), &extension); err != nil {
		return nil, err
	}
	if h.Version != headerVersion {
		return nil, errors.New("invalid protocol version")
	}
	return &extension, nil
}

func registerExtension(conn net.Conn) error {
	// Get extension information
	if err := xHandshake(conn); err != nil {
		return err
	}
	extension := &Extension{}
	if extension, err = xAbout(conn); err != nil {
		return err
	}
	for _, keyword := range extension.Keywords {
		if _, ok = extensions[keyword]; ok {
			return fmt.Errorf("keyword '%s' collision", keyword)
		}
		extensions[keyword] = extension
	}
	return nil
}

// Query sends a message to an extension and waits for a reply
// or times out.
func (x *Extension) Query(from, message string) (string, error) {
	var (
		resp chan string = make(chan string)
		line string
	)
	
	if _, err := x.rw.Write([]byte{fmt.Sprintf("!%s :%s")}); err != nil {
		return err
	}
	go func() {
		r, err := ioutil.ReadAll(x.rw)
		resp <- r
	}()
	select {
	case <- time.After(2*time.Second):
	case line = <- resp:
		return line, nil
	}
	return errors.New("timed out")
}

// runExtensions will run all executables inside the extension
// directory.
func runExtensions(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) {
		cmd := exec.Command(path)
		return cmd.Start()
	})
}
