package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

type lsofEntry struct {
	Command string
	PID     string
	User    string
	FD      string
	Type    string
	Device  string
	Size    string
	Node    string
	Name    string

	Host     string
	Port     string
	Analyzed bool
}

func (l *lsofEntry) analyze() (bool, error) {
	l.Analyzed = true
	opened := false
	host, port, err := net.SplitHostPort(l.Name)
	if err != nil {
		return opened, err
	}
	l.Host = host
	l.Port = port
	if host != "localhost" && host != "127.0.0.1" {
		opened = true
	}
	return opened, nil
}

func parseLSOFEntry(val string) *lsofEntry {
	l := lsofEntry{}
	for i, s := range strings.Fields(val) {
		switch i {
		case 0:
			l.Command = s
		case 1:
			l.PID = s
		case 2:
			l.User = s
		case 3:
			l.FD = s
		case 4:
			l.Type = s
		case 5:
			l.Device = s
		case 6:
			l.Size = s
		case 7:
			l.Node = s
		case 8:
			l.Name = s
		default:
		}
	}
	return &l
}

var pidPort = map[string]*lsofEntry{}

func main() {
	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			resetEntries()
			output, err := lsof()
			if err != nil {
				log.Printf("lsof error: %s\n", err)
			}
			for _, s := range output {
				trimmed := strings.TrimSpace(s)
				if len(trimmed) > 0 {
					entry := parseLSOFEntry(trimmed)
					opened, err := entry.analyze()
					if err != nil {
						log.Printf("analyze error: %s\n", err)
						continue
					}
					if opened {
						if _, exists := pidPort[entry.PID]; !exists {
							notify(fmt.Sprintf("Port %s opened (%s)", entry.Port, entry.PID), entry.Name, entry.Command)
						}
						pidPort[entry.PID] = entry
					}
				}
			}
			removeEntries()
		}
	}
}

func resetEntries() {
	for _, e := range pidPort {
		e.Analyzed = false
	}
}

func removeEntries() {
	for pid := range pidPort {
		if !pidPort[pid].Analyzed {
			notify(fmt.Sprintf("Port %s closed (%s)", pidPort[pid].Port, pid), pidPort[pid].Name, pidPort[pid].Command)
			delete(pidPort, pid)
		}
	}
}

func lsof() ([]string, error) {
	cmd := exec.Command("lsof", "-PiTCP", "-sTCP:LISTEN")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return strings.Split(out.String(), "\n")[1:], nil
}

func notify(title, subtitle, text string) error {
	notification := fmt.Sprintf("display notification \"%s\" with title \"%s\" subtitle \"%s\"", text, title, subtitle)
	return exec.Command("osascript", "-e", notification).Run()
}
