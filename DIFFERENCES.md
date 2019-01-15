# Differences between gocat & gocat.v2

## Written in pure golang.

`gocat.v2` is written in pure golang, which means it's no longer depend on [ccat](https://github.com/dianping/cat/tree/master/lib/c).

And `CGO` is not required either.

## Cat instance is not required.

In `gocat`, you have to create an instance like this:

```go
import (
    "gocat"
)

cat := gocat.Instance()
cat.LogEvent("foo", "bar")
```

while you can achieve the same goal in `gocat.v2` like this:

```go
import (
    cat "gocat.v2"
)

cat.LogEvent("foo", "bar")
```

## Event can be nested in Transaction.

See [case1](./README.md#Example)

## API return value

The following APIs **do not** return pointer anymore.

```go
func NewTransaction(mtype, name string) *message.Transaction
func NewEvent(mtype, name string) *message.Event
func NewMetricHelper(m_name string) *MetricHelper
```

Were changed to:

```go
func NewTransaction(mtype, name string) message.Transactor
func NewEvent(mtype, name string) message.Messeger
func NewMetricHelper(m_name string) MetricHelper
```

No influences if you have used `:=` or `var` to receive our returned value.

## API params

The following APIs requires **time.Time** or **time.Duration** as parameter, which used to be `int64` (timestampInNanosecond).

```go
type Message interface {
    SetTimestamp(timestampInNano int64)
    GetTimestamp() int64
}

type Transaction interface {
    SetDuration(durationInNano int64)
    GetDuration() int64
    
    SetDurationStart(durationStartInNano int64)
}
```

Were changed to:

```
type Message interface {
    SetTime(time time.Time)
    GetTime() time.Time
}

type Transaction interface {
    SetDuration(duration time.Duration)
    GetDuration() time.Duration
    
    SetDurationStart(time time.Time)
}
```

If you have used the mentioned APIs above, migrate to `gocat.v2` will take you some time. 
