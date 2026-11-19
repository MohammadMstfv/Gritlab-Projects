package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

var colors = map[string]string{
	"reset":   "\033[0m",
	"red":     "\033[91m",
	"green":   "\033[92m",
	"yellow":  "\033[93m",
	"blue":    "\033[94m",
	"magenta": "\033[95m",
	"cyan":    "\033[96m",
	"white":   "\033[97m",
	"gray":    "\033[90m",
}

type message struct {
	origin *client
	text   string
}

func receiveMessage(c *client) (message, error) {
	text, err := c.reader.ReadString('\n')
	if err != nil {
		return message{}, fmt.Errorf("client read message: %w", err)
	}

	return message{origin: c, text: text}, nil
}

func processMessage(msg message) error {
	if msg.origin == nil {
		return fmt.Errorf("message has no origin")
	}

	if msg.text == "" {
		return nil
	}

	if msg.text[0] == '/' {
		handleCmd(msg)
		return nil
	}

	msg.text = strings.TrimSpace(strings.TrimRight(msg.text, "\r\n"))
	if msg.text == "" {
		return nil
	}

	msg.text = fmt.Sprintf(
		colors["blue"]+"[%s]"+colors["cyan"]+"[%s/%s]"+colors["reset"]+": %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		msg.origin.g.name,
		msg.origin.name,
		msg.text,
	)

	msg.origin.s.bq <- msg
	return nil
}

func broadcastMessage(msg message) error {
	if msg.origin == nil {
		return nil
	}

	return broadcastLine(msg.origin, msg.text)
}

func broadcastSystemMessage(c *client, text string) {
	line := fmt.Sprintf(
		colors["blue"]+"[%s]"+colors["green"]+"[System]"+colors["reset"]+": %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		text,
	)

	if err := broadcastLine(c, line); err != nil {
		log.Printf("send system message: %v", err)
	}
}

func broadcastLine(origin *client, line string) error {
	s := origin.s
	s.mu.Lock()
	grp := origin.g
	if grp == nil {
		s.mu.Unlock()
		return nil
	}

	grp.history += line
	recipients := append([]*client(nil), grp.c...)
	s.mu.Unlock()

	var firstErr error
	for _, recipient := range recipients {
		if recipient == origin {
			continue
		}
		if err := recipient.write(line); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("send message to %s: %w", recipient.fd.RemoteAddr(), err)
		}
	}
	return firstErr
}
