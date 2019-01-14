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
	defaultAppKey   = "cat"
	defaultHostname = "GoUnknownHost"
	defaultEnv      = "dev"

	defaultIp    = "127.0.0.1"
	defaultIpHex = "7f000001"

	defaultServer  = "cat.sankuai.com"
	defaultXmlFile = "/data/appdatas/cat/client.xml"
	defaultLogDir  = "/data/applogs/cat"
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
	typeSystem = "System"

	nameReboot = "Reboot"

	nameTransactionAggregator = "TransactionAggregator"
	nameEventAggregator       = "EventAggregator"
	nameMetricAggregator      = "MetricAggregator"
)

const (
	signal0 = iota

	signalResetConnection

	signalShutdown

	signalSenderExit
	signalMonitorExit
	signalRouterExit
	signalTransactionAggregatorExit
	signalEventAggregatorExit
	signalMetricAggregatorExit
)
