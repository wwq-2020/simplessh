package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

type handler interface {
	check(line string) bool
	handle() bool
}

var (
	cReg     = regexp.MustCompile(`^c (\d+|[\s\S]+)$`)
	nReg     = regexp.MustCompile(`^n(\s[\S]+\s[\S]+\s[\S]+\s[\s\S]+)?$`)
	ncReg    = regexp.MustCompile(`^nc(\s[\S]+\s[\S]+\s[\S]+\s[\s\S]+)?$`)
	uReg     = regexp.MustCompile(`^u (\d+|[\s\S]+)(\s[\S]+\s[\S]+\s[\S]+\s[\s\S]+)?$`)
	ucReg    = regexp.MustCompile(`^uc (\d+|[\s\S]+)(\s[\S]+\s[\S]+\s[\S]+\s[\s\S]+)?$`)
	qReg     = regexp.MustCompile(`^q$`)
	iReg     = regexp.MustCompile(`^i (\d+|[\s\S]+)$`)
	hReg     = regexp.MustCompile(`^h$`)
	lReg     = regexp.MustCompile(`^l$`)
	dReg     = regexp.MustCompile(`^d (\d+|[\s\S]+|all)$`)
	emptyReg = regexp.MustCompile(`^\s?$`)
	addrReg  = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+:\d+`)
	userReg  = regexp.MustCompile(`^\w+$`)
	pwdReg   = regexp.MustCompile(`^[\s\S]+$`)
)

type cHandler struct {
	cfg *config
	err error
}

func (h *cHandler) check(line string) bool {
	results := cReg.FindStringSubmatch(line)
	if len(results) != 2 {
		return false
	}
	h.cfg, h.err = getCfg(results[1])
	return true
}

func (h *cHandler) handle() bool {
	if h.err != nil {
		io.WriteString(os.Stdout, h.err.Error())
		h.err = nil
		return true
	}
	shell(h.cfg)
	return true
}

type nHandler struct {
	addr      string
	user      string
	pwd       string
	alias     string
	err       error
	isCmdLine bool
}

func (h *nHandler) check(line string) bool {
	return h.doCheck(nReg, line)
}

func (h *nHandler) doCheck(reg *regexp.Regexp, line string) bool {
	results := reg.FindStringSubmatch(line)
	if len(results) < 1 {
		return false
	}
	if results[1] != "" {
		parts := strings.Split(results[1], " ")
		h.isCmdLine = true
		if len(parts) == 5 {
			h.addr = parts[1]
			h.user = parts[2]
			h.pwd = parts[3]
			h.alias = parts[4]
		}
	}
	return true

}

func (h *nHandler) handle() bool {
	if h.isCmdLine {
		cfg := h.genCfg()
		if _, err := getCfg(cfg.Alias); err == nil {
			io.WriteString(os.Stdout, fmt.Sprintf("dup alias:%s for cfg\n", cfg.Alias))
			return false
		}
		if tryConnect(cfg) {
			cfgs = append(cfgs, cfg)
			storeCfgs()
			return true
		}
		io.WriteString(os.Stdout, "invalid cfg\n")
		return false
	}
	cfg, goOn := h.collectInfo()
	if cfg != nil {
		cfgs = append(cfgs, cfg)
		storeCfgs()
	}
	return goOn

}

func (h *nHandler) collectInfo() (*config, bool) {
	r := bufio.NewReader(os.Stdin)
	for {
		if !h.getAddr(r) {
			return nil, false
		}
		if !h.getUser(r) {
			return nil, false
		}
		if !h.getPwd(r) {
			return nil, false
		}
		if !h.getAlias(r) {
			return nil, false
		}
		cfg := h.genCfg()
		if tryConnect(cfg) {
			return cfg, true
		}

	}
}

func (h *nHandler) getAddr(r *bufio.Reader) bool {
	for {
		io.WriteString(os.Stdout, "addr:(192.168.1.100:22) ")
		line, _, err := r.ReadLine()
		if err != nil {
			return false
		}
		if addrReg.Match(line) {
			h.addr = string(line)
			return true
		}
	}
}

func (h *nHandler) getUser(r *bufio.Reader) bool {
	for {
		io.WriteString(os.Stdout, "user:(default root) ")
		line, _, err := r.ReadLine()
		if err != nil {
			return false
		}
		if userReg.Match(line) {
			h.user = string(line)
			return true
		}

		h.user = "root"
		return true
	}
}

func (h *nHandler) getPwd(r *bufio.Reader) bool {
	for {
		io.WriteString(os.Stdout, "password: ")
		line, _, err := r.ReadLine()
		if err != nil {
			return false
		}
		if pwdReg.Match(line) {
			h.pwd = string(line)
			return true
		}
	}
}

func (h *nHandler) getAlias(r *bufio.Reader) bool {
	for {
		io.WriteString(os.Stdout, "alias: ")
		line, _, err := r.ReadLine()
		if err != nil {
			return false
		}
		if len(line) == 0 {
			h.genAlias()
			return true
		}
		lineStr := string(line)
		_, err = getCfg(lineStr)
		if err == nil {
			io.WriteString(os.Stdout, fmt.Sprintf("dup alias for %s\n", lineStr))
			continue
		}
		h.alias = lineStr
		return true

	}
}

func (h *nHandler) genAlias() {
	for {
		alias := fmt.Sprintf("%d-%d", time.Now().Unix(), rand.Int63())
		_, err := getCfg(alias)
		if err == nil {
			continue
		}
		h.alias = alias
		break
	}
}

func (h *nHandler) genCfg() *config {
	return &config{Addr: h.addr, User: h.user, Pwd: h.pwd, Alias: h.alias}
}

type ncHandler struct {
	nHandler
}

func (h *ncHandler) check(line string) bool {
	return h.doCheck(ncReg, line)
}

func (h *ncHandler) handle() bool {
	if !h.nHandler.handle() {
		return false
	}
	cfg := h.genCfg()
	shell(cfg)
	return true
}

type uHandler struct {
	cfg       *config
	err       error
	isCmdLine bool
	nHandler
}

func (h *uHandler) check(line string) bool {
	return h.doCheck(uReg, line)
}

func (h *uHandler) doCheck(reg *regexp.Regexp, line string) bool {
	results := reg.FindStringSubmatch(line)
	if len(results) < 2 {
		return false
	}
	h.cfg, h.err = getCfg(results[1])

	if results[2] != "" {
		parts := strings.Split(results[2], " ")
		h.isCmdLine = true
		if len(parts) == 5 {
			h.addr = parts[1]
			h.user = parts[2]
			h.pwd = parts[3]
			h.alias = parts[4]
		}
	}

	return true
}

func (h *uHandler) handle() bool {
	if h.err != nil {
		io.WriteString(os.Stdout, h.err.Error())
		h.err = nil
		return true
	}
	if h.isCmdLine {
		cfg := h.genCfg()
		if _, err := getCfg(cfg.Alias); err == nil {
			io.WriteString(os.Stdout, fmt.Sprintf("dup alias:%s for cfg\n", cfg.Alias))
			return false
		}
		if tryConnect(cfg) {
			*h.cfg = *cfg
			storeCfgs()
			return true
		}
		io.WriteString(os.Stdout, "invalid cfg\n")
		return false
	}

	cfg, goOn := h.nHandler.collectInfo()
	if cfg != nil {
		*h.cfg = *cfg
		storeCfgs()
	}
	return goOn

}

type ucHandler struct {
	uHandler
}

func (h *ucHandler) check(line string) bool {
	return h.doCheck(ucReg, line)
}

func (h *ucHandler) handle() bool {
	if !h.uHandler.handle() {
		return false
	}
	cfg := h.genCfg()
	shell(cfg)
	return true
}

type qHandler struct {
}

type lHandler struct {
}

func (h *lHandler) check(line string) bool {
	return lReg.MatchString(line)
}

func (h *lHandler) handle() bool {
	for idx, cfg := range cfgs {
		io.WriteString(os.Stdout, fmt.Sprintf("%d)	%s\n", idx, cfg.Addr))
	}
	return true
}

func (h *qHandler) check(line string) bool {
	return qReg.MatchString(line)
}

func (h *qHandler) handle() bool {
	return false
}

type emptyHandler struct {
}

func (h *emptyHandler) check(line string) bool {
	return emptyReg.MatchString(line)
}

func (h *emptyHandler) handle() bool {
	return true
}

type dHandler struct {
	err error
	cfg *config
	all bool
}

func (h *dHandler) check(line string) bool {
	results := dReg.FindStringSubmatch(line)
	if len(results) != 2 {
		return false
	}
	if results[1] == "all" {
		h.all = true
		return true
	}
	h.cfg, h.err = getCfg(results[1])
	return true
}

func (h *dHandler) handle() bool {
	if h.err != nil {
		io.WriteString(os.Stdout, h.err.Error())
		h.err = nil
		return true
	}
	if h.all {
		cfgs = nil
		storeCfgs()
		return true
	}
	idx := h.cfg.Idx
	copy(cfgs[idx:], cfgs[idx+1:])
	cfgs[len(cfgs)-1] = nil
	cfgs = cfgs[:len(cfgs)-1]
	storeCfgs()

	return true
}

type iHandler struct {
	cfg *config
	err error
}

func (h *iHandler) check(line string) bool {
	results := iReg.FindStringSubmatch(line)
	if len(results) != 2 {
		return false
	}
	h.cfg, h.err = getCfg(results[1])
	return true
}

func (h *iHandler) handle() bool {
	if h.err != nil {
		io.WriteString(os.Stdout, h.err.Error())
		h.err = nil
		return true
	}
	data, err := json.Marshal(h.cfg)
	if err != nil {
		io.WriteString(os.Stdout, "invalid cfg\n")
		return true
	}
	os.Stdout.Write(data)
	io.WriteString(os.Stdout, "\n")
	return true
}

type hHandler struct {
}

func (h *hHandler) check(line string) bool {
	return hReg.MatchString(line)
}

func (h *hHandler) handle() bool {
	io.WriteString(os.Stdout, `	c id/alias to connect a cfg
	n to new a cfg
	nc to new a cfg and connect it
	u id/alias to update a cfg
	uc id/alias to update a cfg and connect it
	l to list cfg
	d id/alias/all to del a cfg
	i id/alias to get detail info
	q to quit
	h to get help info
`)
	return true
}

func getHandler(line string) handler {
	for _, handler := range handlers {
		if handler.check(line) {
			return handler
		}
	}
	return &hHandler{}
}
