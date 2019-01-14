package cat

const (
	GOCAT_VERSION = "1.0.0"
)

const (
	FAIL = "fail"
)

const (
	DEFAULT_IP = "127.0.0.1"
)

const (
	HIGH_PRIORITY_QUEUE_SIZE   = 1000
	NORMAL_PRIORITY_QUEUE_SIZE = 5000

	TRANSACTION_AGGREGATOR_CHANNEL_CAPACITY = 1000
	EVENT_AGGREGATOR_CHANNEL_CAPACITY       = 1000
	METRIC_AGGREGATOR_CHANNEL_CAPACITY = 1000
)

const (
	CAT_SYSTEM = "System"

	TRANSACTION_AGGREGATOR = "TransactionAggregator"
	EVENT_AGGREGATOR       = "EventAggregator"
	METRIC_AGGREGATOR       = "MetricAggregator"
)

type Signals chan int

const (
	S_RESET_CONNECTION = iota
)
