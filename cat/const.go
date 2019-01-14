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

const (
	DefaultAppKey   = "cat"
	DefaultHostname = "GoUnknownHost"
	DefaultEnv      = "dev"

	DefaultIp    = "127.0.0.1"
	DefaultIpHex = "7f000001"

	DefaultServer  = "cat.sankuai.com"
	DefaultXmlFile = "/data/appdatas/cat/client.xml"
	DefaultLogDir  = "/data/applogs/cat"
)

const (
	HighPriorityQueueSize   = 1000
	NormalPriorityQueueSize = 5000

	TransactionAggregatorChannelCapacity = 1000
	EventAggregatorChannelCapacity       = 1000
	MetricAggregatorChannelCapacity      = 1000

	TransactionAggregatorInterval = time.Second * 3
	EventAggregatorInterval       = time.Second * 3
	MetricAggregatorInterval      = time.Second * 3
)

const (
	System = "System"

	TransactionAggregator = "TransactionAggregator"
	EventAggregator       = "EventAggregator"
	MetricAggregator      = "MetricAggregator"
)

type Signals chan int

const (
	SignalResetConnection = iota
)
