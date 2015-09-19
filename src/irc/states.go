package irc

import (
	"bufio"
	"time"
	"log"
	"strings"
	"runtime/debug"
	"sync"

	"bot"
)

// This is a state function
type stateFn func(*machine) stateFn

type machine struct {
	rw     *bufio.ReadWriter
	line   string
	tokens []string
	tokenN int
	errCh  chan error
	l sync.Mutex // state machine lock

	keywords map[string]bot.KeywordActionFn
}

func (m *machine) Run() {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("State mashine crashed ...")
			debug.PrintStack()
		}
	}()
	go func() {
		for {
			time.Sleep(time.Minute)
			m.l.Lock()
			if err := Ping(m.rw, ":irc.freenode.net"); err != nil {
				log.Printf("PING error: %v", err)
			}
			m.l.Unlock()
		}
	}()
	for state := stateStart; state != nil; {
		m.l.Lock()
		state = state(m)
		m.l.Unlock()
	}
	log.Printf("Ending state machine, last line was '%s'", m.line)
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
	msgEnd := strings.Index(m.line, "\r\n")
	msg := m.line[strings.LastIndex(m.line, ":")+1 : msgEnd]
	from := m.tokens[2]
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
	msg := m.line[strings.LastIndex(m.line, ":")+1 : msgEnd]
	log.Printf("PING from %s", msg)
	if err := Pong(m.rw, m.tokens[1]); err != nil {
		return nil
	}
	return stateStart
}
