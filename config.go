package client

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type odooConfig struct {
	path   string
	config *ini.File
}

func fileExist(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func loadConfig(configPath string) (*ini.File, error) {
	if !fileExist(configPath) {
		return nil, errors.New("File does not exist")
	}
	config, err := ini.LoadSources(ini.LoadOptions{AllowPythonMultilineValues: true}, configPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (rc *odooConfig) getHttpPort() (int, error) {
	var port int
	var err error
	if rc.config.Section("options").HasKey("http_port") {
		port, err = rc.config.Section("options").Key("http_port").Int()
	} else if rc.config.Section("options").HasKey("xmlrpc_port") {
		port, err = rc.config.Section("options").Key("http_port").Int()
	} else {
		port = 8069
	}
	if err != nil {
		return 0, err
	}
	return port, nil
}

func (rc *odooConfig) getAdminPasswd() (string, error) {
	for _, key := range [2]string{"admin_passwd", "ons_admin_passwd"} {
		if rc.config.Section("options").HasKey(key) {
			return rc.config.Section("options").Key(key).String(), nil
		}
	}
	return "", errors.New("Could not find admin password")
}

func GetOdooConfig() (*odooConfig, error) {
	var configPath string
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	rcs := [3]string{fmt.Sprint(home, "/odoo-prod.conf"), fmt.Sprint(home, "/.odoorc"), fmt.Sprint(home, "/.openerp_serverrc")}
	for _, rc := range rcs {
		if fileExist(rc) {
			configPath = rc
			break
		}
	}
	if configPath == "" {
		return nil, errors.New("Could not find Odoo RC")
	}

	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return &odooConfig{config: config, path: configPath}, nil
}
