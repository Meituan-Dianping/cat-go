package cat

import (
	"strconv"
	"strings"
)

type serverAddress struct {
	host string
	port int
}

func resolveServerAddresses(router string) (addresses []serverAddress) {
	for _, segment := range strings.Split(router, ";") {
		if len(segment) == 0 {
			continue
		}
		fragments := strings.Split(segment, ":")
		if len(fragments) != 2 {
			logger.Warning("%s isn't a valid server address.", segment)
			continue
		}

		if port, err := strconv.Atoi(fragments[1]); err != nil {
			logger.Warning("%s isn't a valid server address.", segment)
		} else {
			addresses = append(addresses, serverAddress{
				host: fragments[0],
				port: port,
			})
		}
	}
	return
}

func compareServerAddress(a, b *serverAddress) bool {
	if a == nil || b == nil {
		return false
	}
	if strings.Compare(a.host, b.host) == 0 {
		return a.port == b.port
	} else {
		return false
	}
}
