package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
)

type config struct {
	Addr  string `json:"addr"`
	User  string `json:"user"`
	Pwd   string `json:"pwd"`
	Alias string `json:"alias"`
	Idx   int    `json:"-"`
}

func getCfgs() {
	f, err := os.Open(abs)
	if err != nil {
		if !os.IsNotExist(err) {
			logrus.WithField("err", err).Fatal("open cfg")
		}
		f, err = os.Create(abs)
		if err != nil {
			logrus.WithField("err", err).Fatal("create cfg")
		}
		return
	}

	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		logrus.WithField("err", err).Fatal("read cfg")
	}

	if len(data) != 0 {
		if err := json.Unmarshal(data, &cfgs); err != nil {
			logrus.WithField("err", err).Fatal("unmarshal cfgs")
		}
	}

	for _, cfg := range cfgs {
		alias2Cfg[cfg.Alias] = cfg
	}
}

func storeCfgs() {
	data, err := json.Marshal(cfgs)
	if err != nil {
		logrus.WithField("err", err).Error("marshal cfgs")
		return
	}
	if err := ioutil.WriteFile(abs, data, 0644); err != nil {
		logrus.WithField("err", err).Error("write cfgs")
	}
}

func getCfg(idStr string) (*config, error) {
	if len(cfgs) == 0 {
		return nil, errors.New("has no cfg\n")
	}
	idx, err := strconv.Atoi(idStr)
	if err != nil {
		for idx, cfg := range cfgs {
			if cfg.Alias == idStr {
				cfg.Idx = idx
				return cfg, nil
			}
		}
		return nil, fmt.Errorf("found no cfg for %s\n", idStr)
	}
	if idx < len(cfgs) {
		cfg := cfgs[idx]
		cfg.Idx = idx
		return cfg, nil
	}
	return nil, fmt.Errorf("id must between 0 to %d\n", len(cfgs)-1)
}
