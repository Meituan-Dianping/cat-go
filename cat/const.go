package cat

import (
	"time"
)

const (
	GoCatVersion = "2.0.2"
)

const (
	SUCCESS = "0"
	ERROR   = "-1"
	FAIL    = "fail"
)

const (
	routerPath = "/cat/s/router"
)

const ( // Declared default values.
	localhost = "127.0.0.1"

	defaultAppKey   = "cat"
	defaultHostname = "GoUnknownHost"
	defaultEnv      = "dev"

	defaultIp    = "127.0.0.1"
	defaultIpHex = "7f000001"

	defaultXmlFile = "/data/appdatas/cat/client.xml"
	defaultLogDir  = "/data/applogs/cat"
	tmpLogDir      = "/tmp"
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

	defaultRouterInterval = time.Minute * 3
	defaultWriteDeadline  = time.Second

	senderBlockingTimeoutTime = time.Second * 3
)

var (
	routerInterval = defaultRouterInterval
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
