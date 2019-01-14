package cat

import (
	"net"
	"time"
)

func getLocalhostIp() (ip string, err error) {
	ip = DEFAULT_IP

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return
			}
		}
	}
	return
}

func duration2Millis(duration time.Duration) int64 {
	return duration.Nanoseconds() / time.Millisecond.Nanoseconds()
}

func duration2Micros(duration time.Duration) int64 {
	return duration.Nanoseconds() / time.Microsecond.Nanoseconds()
}
