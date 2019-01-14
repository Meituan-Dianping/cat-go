package cat

import (
	"time"
)

const (
	GoCatVersion = "2.0.0"
)

const (
	SUCCESS = "0"
	ERROR   = "-1"
	FAIL    = "fail"
)

const ( // Declared default values.
	DefaultAppKey   = "cat"
	DefaultHostname = "GoUnknownHost"
	DefaultEnv      = "dev"

	DefaultIp    = "127.0.0.1"
	DefaultIpHex = "7f000001"

	DefaultServer  = "cat.sankuai.com"
	DefaultXmlFile = "/data/appdatas/cat/client.xml"
	DefaultLogDir  = "/data/applogs/cat"
)

const ( // Declared properties given by the router server.
	propertySample  = "sample"
	propertyRouters = "routers"
	propertyBlock   = "block"
)

const (
	highPriorityQueueSize   = 1000
	normalPriorityQueueSize = 5000

	transactionAggregatorChannelCapacity = 1000
	eventAggregatorChannelCapacity       = 1000
	metricAggregatorChannelCapacity      = 1000

	transactionAggregatorInterval = time.Second * 3
	eventAggregatorInterval       = time.Second * 3
	metricAggregatorInterval      = time.Second * 3
)

const ( // Declared a series of reserved type and names.
	typeSystem = "typeSystem"

	nameTransactionAggregator = "nameTransactionAggregator"
	nameEventAggregator       = "nameEventAggregator"
	nameMetricAggregator      = "nameMetricAggregator"
)

type signals chan int

const (
	signalResetConnection = iota
)
