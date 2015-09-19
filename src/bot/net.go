package bot

import (
	"bufio"
	"net"
)

type KeywordActionFn func(rw *bufio.ReadWriter, from, message string) error
type Identity struct {
	RealName  string
	Nicknames []string
}

type Server struct {
	Name    string
	Address string // host:port

	c     net.Conn
	ident *Identity
}

func (s *Server) IO() *bufio.ReadWriter {
	return bufio.NewReadWriter(bufio.NewReader(s.c), bufio.NewWriter(s.c))
}

func Connect(address string, ident *Identity) (*Server, error) {
	var (
		server *Server = &Server{Address: address, ident: ident}
		err    error
	)
	server.c, err = net.Dial("tcp4", address)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) Close() {
	s.c.Close()
}
