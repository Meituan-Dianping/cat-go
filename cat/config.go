package cat

import (
	"encoding/xml"
	"io/ioutil"
	"net"
	"os"
)

type Config struct {
	domain   string
	hostname string
	env      string
	ip       string
	ipHex    string

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
	defer func() {
		if err := file.Close(); err != nil {
			logger.Warning("Cannot close local client.xml file.")
		}
	}()

	data, err = ioutil.ReadAll(file)
	if err != nil {
		logger.Warning("Unable to read content from file `%s`", filename)
	}
	return
}

func loadConfig() (data []byte, err error) {
	if data, err = loadConfigFromLocalFile(defaultXmlFile); err != nil {
		logger.Error("Failed to load local config file.")
		return
	}
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

	var ip net.IP
	if ip, err = getLocalhostIp(); err != nil {
		config.ip = defaultIp
		config.ipHex = defaultIpHex
		logger.Warning("Error while getting local ip, using default ip: %s", defaultIp)
	} else {
		config.ip = ip2String(ip)
		config.ipHex = ip2HexString(ip)
		logger.Info("Local ip has been configured to %s", config.ip)
	}

	if config.hostname, err = os.Hostname(); err != nil {
		config.hostname = defaultHostname
		logger.Warning("Error while getting hostname, using default hostname: %s", defaultHostname)
	} else {
		logger.Info("Hostname has been configured to %s", config.hostname)
	}

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
