package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"net/http"

	"bot"
	"irc"
)

var (
	keywords map[string]bot.KeywordActionFn
)

func RunBot(address, user, nick, real string, channels ...string) error {
	// Connect to the IRC server and register there
	server, err := bot.Connect(address, &bot.Identity{real, []string{nick}})
	if err != nil {
		return err
	}

	io := server.IO()
	if err = irc.Register(io, user, real); err != nil {
		return err
	}
	if err = irc.ChangeNick(io, nick); err != nil {
		return err
	}
	for _, channel := range channels {
		if err = irc.Join(io, channel); err != nil {
			return err
		}
	}
	// Run a new state machine
	fsm := irc.NewMachine(io)
	fsm.Keywords(keywords)
	fsm.Run()

	// Connection ended, deinitialize
	server.Close()
	return nil
}

func main() {
	keywords = map[string]bot.KeywordActionFn{",help": func(rw *bufio.ReadWriter, from, msg string) error {
		if len(msg) == 0 {
			return irc.Message(rw, from, "Oh hai! Usage: ,help <topic>")
		}
		resp := ""
		switch msg {
		case "about":
			resp = "gorobot, unofficial release"
		case "help":
			resp = "available topic: about, help"
		default:
			resp = "unknown topic"
		}
		return irc.Message(rw, from, resp)

	}, ",w": func(rw *bufio.ReadWriter, from, msg string) error {
		resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + msg)
		if err != nil {
			log.Printf("Error fetching weather: %v", err)
			return irc.Message(rw, from, "(error fetching data)")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error fetching weather: %v", err)
			return irc.Message(rw, from, "(error fetching data)")
		}

		// Get the most basic information
		data := struct {
			Weather []struct {
				Description string `json:"description"`
			} `json:"weather"`
			Main struct {
				Temp     float64 `json:"temp"`
				Humidity int     `json:"humidity"`
			} `json:"main"`
		}{}
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("Error fetching weather: %v", err)
			return irc.Message(rw, from, "(error fetching data)")
		}
		str := fmt.Sprintf("%s, temp.: %f, humidity: %d%%", data.Weather[0].Description, data.Main.Temp, data.Main.Humidity)
		return irc.Message(rw, from, str)
	}}
	address := os.Args[1]
	nick := os.Args[2]
	log.Fatal(RunBot(address, nick, nick, nick, os.Args[3:]...))
}
