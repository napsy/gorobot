package irc

import (
	"strings"
	"bufio"
	"log"

	"bot"
)

// This is a state function
type stateFn func(*machine) stateFn

type machine struct {
	rw *bufio.ReadWriter
	line string
	tokens []string
	tokenN int
	errCh chan error

	keywords map[string]bot.KeywordActionFn
}

func (m *machine) Run() {
	for state := stateStart; state != nil; {
		state = state(m)
	}
}

func NewMachine(rw *bufio.ReadWriter) *machine {
	return &machine{rw: rw, errCh: make(chan error)}
}

func (m *machine) Keywords(keywords map[string]bot.KeywordActionFn) {
	m.keywords = keywords
}

func stateStart(m *machine) stateFn {
	var err error
	m.line, err = m.rw.ReadString('\n')
	if err != nil {
		return nil
	}
	m.tokens = strings.Split(m.line, " ")
	// Server commands
	switch m.tokens[0] {
	case "PING":
		return statePing
	}
	// Server messages
	switch m.tokens[1] {
	case "PRIVMSG":
		return stateMsg
	}
	log.Printf("%v", m.line)
	return stateStart
}

func stateMsg(m *machine) stateFn {
	log.Printf("PRIVMSG: '%s'", m.line)
	msgEnd := strings.Index(m.line, "\r\n")
	msg := m.line[strings.LastIndex(m.line, ":")+1:msgEnd]
	if m.tokens[2][0] == '#' {
		// Channel messages
		idx := strings.Index(m.tokens[0], "!~")
		fromUser := m.tokens[0][1:idx]
		log.Printf("CHANNEL '%s' <%s>: %s", from, fromUser, msg)
	} else {
		// Private messages
		log.Printf("PRIVATE '%s': %s", from, msg)
	}

	if msg[0] == ',' {
		args := strings.SplitN(msg, " ", 2)
		keywordMsg := ""
		if len(args) > 1 {
			keywordMsg = args[1]
		}
		if action, _ := m.keywords[args[0]]; action != nil {
			if err := action(m.rw, from, keywordMsg); err != nil {
				log.Printf("Error handling '%s': %v", args[0], err)
			}
		}
	}
	return stateStart
}

func statePing(m *machine) stateFn {
	msgEnd := strings.Index(m.line, "\r\n")
	msg := m.line[strings.LastIndex(m.line, ":")+1:msgEnd]
	log.Printf("PING from %s", msg)
	if err := Pong(m.rw, m.tokens[1]); err != nil {
		return nil
	}
	return stateStart
}
