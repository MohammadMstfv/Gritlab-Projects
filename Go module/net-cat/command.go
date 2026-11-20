package main

import (
	"fmt"
	"strings"
)

func nameCmd(msg message, fields []string) error {
	if len(fields) < 2 {
		if err := msg.origin.write(colors["red"] + "Usage: /name <new_name>\nOne argument required\n" + colors["reset"]); err != nil {
			return fmt.Errorf("send name command usage: %w", err)
		}
		return nil
	}
	if len(fields) > 2 {
		if err := msg.origin.write(colors["red"] + "Usage: /name <new_name>\nNames may not include spaces\n" + colors["reset"]); err != nil {
			return fmt.Errorf("send name command usage: %w", err)
		}
		return nil
	}

	newName := fields[1]
	oldName := msg.origin.name
	if err := setName(msg.origin, newName); err != nil {
		if err := msg.origin.write(colors["red"] + err.Error() + colors["reset"]); err != nil {
			return fmt.Errorf("send name command usage: %w", err)
		}
		return nil
	}

	broadcastSystemMessage(
		msg.origin,
		fmt.Sprintf(colors["cyan"]+"[%s]"+colors["reset"]+" changed their name to "+colors["yellow"]+"[%s]"+colors["reset"]+".", oldName, newName),
	)
	return nil
}

func groupCmd(msg message, fields []string) error {
	newGroup := ""
	if len(fields) > 1 {
		newGroup = fields[1]
	}
	if len(fields) > 2 {
		if err := msg.origin.write(colors["red"] + "Usage: /group <group>\nGroup names may not include spaces\n" + colors["reset"]); err != nil {
			return fmt.Errorf("send group command usage: %w", err)
		}
		return nil
	}

	oldGroup := msg.origin.g.name

	if newGroup == oldGroup {
		if err := msg.origin.write(colors["red"] + "You are already in this group\n" + colors["reset"]); err != nil {
			return fmt.Errorf("send group command usage: %w", err)
		}
		return nil
	}

	clientGroupDel(msg.origin)
	clientGroupAdd(msg.origin, newGroup)
	return nil
}

func getOnline(s *server) string {
	list := ""
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, g := range s.g {
		list += "[" + g.name + "/]:\n"
		for _, c := range g.c {
			list += "  [" + g.name + "/" + c.name + "]\n"
		}
	}
	return list
}

func onlineCmd(msg message, fields []string) error {
	if len(fields) > 1 {
		if err := msg.origin.write(colors["red"] + "Usage: /online\nCommand takes no arguments\n" + colors["reset"]); err != nil {
			return fmt.Errorf("send online command usage: %w", err)
		}
		return nil
	}
	if err := msg.origin.write(colors["yellow"] + getOnline(msg.origin.s) + colors["reset"]); err != nil {
		return fmt.Errorf("send online list: %w", err)
	}
	return nil
}

type cmdHandler struct {
	fn   func(message, []string) error
	help string
}

func helpCmd(msg message, fields []string) error {
	help := "Available Commands:\n"
	for name, cmd := range cmds {
		help += fmt.Sprintf("  %s - %s\n", name, cmd.help)
	}
	if err := msg.origin.write(colors["yellow"] + help + colors["reset"]); err != nil {
		return fmt.Errorf("send help: %w", err)
	}
	return nil
}

var cmds map[string]struct {
	fn   func(message, []string) error
	help string
}

func init() {
	cmds = map[string]struct {
		fn   func(message, []string) error
		help string
	}{
		"/name":   {nameCmd, "Change name"},
		"/group":  {groupCmd, "Change group"},
		"/online": {onlineCmd, "Get list of users online"},
		"/help":   {helpCmd, "Get help message"},
	}
}

func handleCmd(msg message) error {
	fields := strings.Fields(strings.TrimSpace(msg.text))
	if len(fields) == 0 {
		return nil
	}
	cmd, valid := cmds[fields[0]]
	if !valid {
		errstr := fmt.Sprintf(
			colors["red"]+"Command not recognized: "+colors["reset"]+"%s\n"+colors["red"]+"Use "+colors["yellow"]+"/help"+colors["red"]+" to check available commands\n"+colors["yellow"]+"Prepend a space before / to send as a message\n"+colors["reset"],
			fields[0],
		)
		if err := msg.origin.write(errstr); err != nil {
			return fmt.Errorf("send name command usage: %w", err)
		}
		return nil
	}
	return cmd.fn(msg, fields)
}
