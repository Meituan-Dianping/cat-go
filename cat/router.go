package cat

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type routerConfigXMLProperty struct {
	XMLName xml.Name `xml:"property"`
	Id      string   `xml:"id,attr"`
	Value   string   `xml:"value,attr"`
}

type routerConfigXML struct {
	XMLName    xml.Name                  `xml:"property-config"`
	Properties []routerConfigXMLProperty `xml:"property"`
}

type catRouterConfig struct {
	scheduleMixin
	sample  float64
	routers []serverAddress
	current *serverAddress
	ticker *time.Ticker
}

var router = catRouterConfig{
	scheduleMixin: makeScheduleMixedIn(signalRouterExit),
	sample:        1.0,
	routers:       make([]serverAddress, 0),
	ticker: nil,
}

func (c *catRouterConfig) GetName() string {
	return "Router"
}

func (c *catRouterConfig) updateRouterConfig() {
	var query = url.Values{}
	query.Add("env", config.env)
	query.Add("domain", config.domain)
	query.Add("ip", config.ip)
	query.Add("hostname", config.hostname)
	query.Add("op", "xml")

	u := url.URL{
		Scheme:   "http",
		Path:     "/cat/s/router",
		RawQuery: query.Encode(),
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for _, server := range config.httpServerAddresses {
		u.Host = fmt.Sprintf("%s:%d", server.host, config.httpServerPort)
		logger.Info("Getting router config from %s", u.String())

		resp, err := client.Get(u.String())
		if err != nil {
			logger.Warning("Error occurred while getting router config from url %s", u.String())
			continue
		}

		c.parse(resp.Body)
		return
	}

	logger.Error("Can't get router config from remote server.")
	return
}

func (c *catRouterConfig) handle(signal int) {
	switch signal {
	case signalResetConnection:
		logger.Warning("Connection has been reset, reconnecting.")
		c.current = nil
		c.updateRouterConfig()
	default:
		c.scheduleMixin.handle(signal)
	}
}

func (c *catRouterConfig) afterStart() {
	c.ticker = time.NewTicker(time.Minute * 3)
	c.updateRouterConfig()
}

func (c *catRouterConfig) beforeStop() {
	c.ticker.Stop()
}

func (c *catRouterConfig) process() {
	select {
	case sig := <-c.signals:
		c.handle(sig)
	case <-c.ticker.C:
		c.updateRouterConfig()
	}
}

func (c *catRouterConfig) updateSample(v string) {
	sample, err := strconv.ParseFloat(v, 32)
	if err != nil {
		logger.Warning("Sample should be a valid float, %s given", v)
	} else if math.Abs(sample - c.sample) > 1e-9 {
		c.sample = sample
		logger.Info("Sample rate has been set to %f%%", c.sample*100)
	}
}

func (c *catRouterConfig) updateBlock(v string) {
	if v == "false" {
		enable()
	} else  {
		disable()
	}
}

func (c *catRouterConfig) parse(reader io.ReadCloser) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	t := new(routerConfigXML)
	if err := xml.Unmarshal(bytes, &t); err != nil {
		logger.Warning("Error occurred while parsing router config xml content.\n%s", string(bytes))
	}

	for _, property := range t.Properties {
		switch property.Id {
		case propertySample:
			c.updateSample(property.Value)
		case propertyRouters:
			c.updateRouters(property.Value)
		case propertyBlock:
			c.updateBlock(property.Value)
		}
	}
}

func (c *catRouterConfig) updateRouters(router string) {
	newRouters := resolveServerAddresses(router)

	oldLen, newLen := len(c.routers), len(newRouters)

	if newLen == 0 {
		return
	} else if oldLen == 0 {
		logger.Info("Routers has been initialized to: %s", newRouters)
		c.routers = newRouters
	} else if oldLen != newLen {
		logger.Info("Routers has been changed to: %s", newRouters)
		c.routers = newRouters
	} else {
		for i := 0; i < oldLen; i++ {
			if !compareServerAddress(&c.routers[i], &newRouters[i]) {
				logger.Info("Routers has been changed to: %s", newRouters)
				c.routers = newRouters
				break
			}
		}
	}

	for _, server := range newRouters {
		if compareServerAddress(c.current, &server) {
			return
		}

		addr := fmt.Sprintf("%s:%d", server.host, server.port)
		if conn, err := net.DialTimeout("tcp", addr, time.Second); err != nil {
			logger.Info("Failed to connect to %s, retrying...", addr)
		} else {
			c.current = &server
			logger.Info("Connected to %s.", addr)
			sender.chConn <- conn
			return
		}
	}

	logger.Info("Cannot established a connection to cat server.")
}
