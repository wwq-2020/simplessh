package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	abs = filepath.Join(os.Getenv("HOME"), filename)
	alias2Cfg = make(map[string]*config)
	handlers = append(handlers, &cHandler{})
	handlers = append(handlers, &lHandler{})
	handlers = append(handlers, &hHandler{})
	handlers = append(handlers, &nHandler{})
	handlers = append(handlers, &ncHandler{})
	handlers = append(handlers, &uHandler{})
	handlers = append(handlers, &ucHandler{})
	handlers = append(handlers, &qHandler{})
	handlers = append(handlers, &emptyHandler{})
	handlers = append(handlers, &dHandler{})
	handlers = append(handlers, &iHandler{})

}

const (
	filename = ".simplessh"
)

var (
	cfgs      []*config
	abs       string
	handlers  []handler
	alias2Cfg map[string]*config
)

func main() {
	getCfgs()
	args := os.Args[1:]
	if len(args) > 0 {
		handler := getHandler(strings.Join(args, " "))
		handler.handle()
		return
	}

	io.WriteString(os.Stdout, "type h to get help info\n")
	src := bufio.NewReader(os.Stdin)
	for {
		io.WriteString(os.Stdout, ">>>")
		line, _, err := src.ReadLine()
		if err != nil {
			break
		}
		handler := getHandler(string(line))
		if !handler.handle() {
			break
		}
	}
}
