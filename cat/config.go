package cat

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	domain   string
	hostname string
	env      string

	ip    string
	ipHex string

	httpServerPort      int
	httpServerAddresses []serverAddress

	serverAddress []serverAddress
}

type XMLConfig struct {
	Name    xml.Name         `xml:"config"`
	Servers XMLConfigServers `xml:"servers"`
}

type XMLConfigServers struct {
	Servers []XMLConfigServer `xml:"server"`
}

type XMLConfigServer struct {
	Host string `xml:"ip,attr"`
	Port int    `xml:"port,attr"`
}

var config = Config{
	domain:   defaultAppKey,
	hostname: defaultHostname,
	env:      defaultEnv,
	ip:       defaultIp,
	ipHex:    defaultIpHex,

	httpServerPort:      8080,
	httpServerAddresses: []serverAddress{},

	serverAddress: []serverAddress{},
}

func loadConfigFromLocalFile(filename string) (data []byte, err error) {
	file, err := os.Open(filename)
	if err != nil {
		logger.Warning("Unable to open file `%s`.", filename)
		return
	}
	defer file.Close()

	data, err = ioutil.ReadAll(file)
	if err != nil {
		logger.Warning("Unable to read content from file `%s`", filename)
	}
	return
}

func loadConfigFromRemoteServer() (data []byte, err error) {
	url := fmt.Sprintf("http://%s/cat/s/launch?ip=%s", defaultServer, config.ip)
	logger.Info("Getting config from %s", url)

	res, err := http.Get(url)
	if err != nil {
		logger.Warning("Error occurred while flush http request.")
		return
	}

	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Remote server return none 200 status code: %d", res.StatusCode))
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Warning("Unable to read content from http response")
	}
	return
}

func loadConfig() (data []byte, err error) {
	if data, err = loadConfigFromLocalFile(defaultXmlFile); err == nil {
		return
	}
	logger.Warning("Failed to load local config file, trying to get from remote server.")

	if data, err = loadConfigFromRemoteServer(); err == nil {
		return
	}
	logger.Warning("Failed to load config from remote server.")
	return
}

func parseXMLConfig(data []byte) (err error) {
	c := XMLConfig{}
	err = xml.Unmarshal(data, &c)
	if err != nil {
		logger.Warning("Failed to parse xml content")
	}

	for _, x := range c.Servers.Servers {
		config.httpServerAddresses = append(config.httpServerAddresses, serverAddress{
			host: x.Host,
			port: x.Port,
		})
	}

	logger.Info("Server addresses: %s", config.httpServerAddresses)
	return
}

func (config *Config) Init(domain string) (err error) {
	config.domain = domain

	defer func() {
		if err == nil {
			logger.Info("Cat has been initialized successfully with appkey: %s", config.domain)
		} else {
			logger.Error("Failed to initialize cat.")
		}
	}()

	// TODO load env.

	if config.ip, err = getLocalhostIp(); err != nil {
		logger.Warning("Error while getting local ip, using default ip: %s", config.ip)
		return
	} else {
		// TODO ipHex
		logger.Info("Local ip has been configured to %s", config.ip)
	}

	// TODO hostname

	var data []byte
	if data, err = loadConfig(); err != nil {
		return
	}

	// Print config content to log file.
	logger.Info("\n%s", data)

	if err = parseXMLConfig(data); err != nil {
		return
	}

	return
}
