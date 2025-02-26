package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/astaxie/beego/logs"
)

type ProxyConfig struct {
	Enable   bool   `json:"enable"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Auth     bool   `json:"auth"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type DomainConfig struct {
	Enable          bool   `json:"enable"`
	IPv6            bool   `json:"ipv6"`
	GenerateType    string `json:"generate_type"`
	ConnectivityURL string `json:"connectivity_url"`
	NetInterface    string `json:"network_interface"`
	ScriptCommand   string `json:"script_command"`
	FilterRegexp    string `json:"filter_regexp"`
	DomainName      string `json:"domain_name"`
	SubDomain       string `json:"sub_domain"`
	CustomParams    string `json:"custom_params"`
}

type ProviderConfig struct {
	Enable  bool           `json:"enable"`
	Name    string         `json:"name"`
	Key     string         `json:"key"`
	Secret  string         `json:"sercet"`
	TTL     int            `json:"ttl"`
	Domains []DomainConfig `json:"domains"`
}

type Config struct {
	Proxy     ProxyConfig      `json:"proxy"`
	Providers []ProviderConfig `json:"providers"`
}

var configCache = Config{
	Proxy: ProxyConfig{
		Enable:   false,
		Address:  "192.168.1.1",
		Port:     8080,
		Protocol: "HTTP",
		Auth:     false,
		User:     "",
		Password: "",
	},
	Providers: []ProviderConfig{},
}

var configFilePath string
var configLock sync.Mutex

func configSyncToFile() error {
	configLock.Lock()
	defer configLock.Unlock()

	value, err := json.MarshalIndent(configCache, "\t", " ")
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
		return err
	}

	return os.WriteFile(configFilePath, value, 0664)
}

func ConfigGet() *Config {
	return &configCache
}

func ProxyConfigSave(proxy ProxyConfig) error {
	configCache.Proxy = proxy
	return configSyncToFile()
}

func ProviderConfigSave(provider ProviderConfig) error {
	configCache.Providers = append(configCache.Providers, provider)
	return configSyncToFile()
}

func ConfigInit() {
	var err error
	var value []byte

	defer func() {
		if err != nil {
			err := configSyncToFile()
			if err != nil {
				logs.Error("config sync to file fail, %s", err.Error())
			}
		}
	}()

	configFilePath = filepath.Join(ConfigDirGet(), "config.json")

	_, err = os.Stat(configFilePath)
	if err != nil {
		logs.Info("config file not exist, create a new one")
	}

	value, err = os.ReadFile(configFilePath)
	if err != nil {
		logs.Error("read config file from app data dir fail, %s", err.Error())
		return
	}

	err = json.Unmarshal(value, &configCache)
	if err != nil {
		logs.Error("json unmarshal config fail, %s", err.Error())
		return
	}
}
