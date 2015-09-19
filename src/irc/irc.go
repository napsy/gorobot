package irc

import (
	"bufio"
	"fmt"
)

func Register(rw *bufio.ReadWriter, user, real string) error {
	var (
		str string = fmt.Sprintf("USER %s 0 * :%s\r\n", user, real)
		err error
	)
	_, err = rw.Write([]byte(str))
	if err != nil {
		return err
	}
	rw.Flush()
	return nil
}

func ChangeNick(rw *bufio.ReadWriter, nick string) error {
	var (
		str string = fmt.Sprintf("NICK %s\r\n", nick)
		err error
	)
	_, err = rw.Write([]byte(str))
	if err != nil {
		return err
	}
	rw.Flush()
	return nil
}

func Join(rw *bufio.ReadWriter, channel string) error {
	var (
		str string = fmt.Sprintf("JOIN %s\r\n", channel)
		err error
	)
	_, err = rw.Write([]byte(str))
	if err != nil {
		return err
	}
	rw.Flush()
	return nil
}

func Part(rw *bufio.ReadWriter, channel string) error {
	var (
		str string = fmt.Sprintf("PART %s\r\n", channel)
		err error
	)
	_, err = rw.Write([]byte(str))
	if err != nil {
		return err
	}
	rw.Flush()
	return nil
}

func Quit(rw *bufio.ReadWriter, reason string) error {
	var (
		str string = fmt.Sprintf("QUIT %s\r\n", reason)
		err error
	)
	_, err = rw.Write([]byte(str))
	if err != nil {
		return err
	}
	rw.Flush()
	return nil
}

func Pong(rw *bufio.ReadWriter, message string) error {
	var (
		str string = fmt.Sprintf("PONG :%s\r\n", message)
		err error
	)
	_, err = rw.Write([]byte(str))
	if err != nil {
		return err
	}
	rw.Flush()
	return nil
}

func Message(rw *bufio.ReadWriter, target, message string) error {
	var (
		str string = fmt.Sprintf("PRIVMSG %s :%s\r\n", target, message)
		err error
	)
	_, err = rw.Write([]byte(str))
	if err != nil {
		return err
	}
	rw.Flush()
	return nil
}
