package cat

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type serverAddress struct {
	host string
	port int
}

type serverAddresses []serverAddress

func (s *serverAddresses) Add(host string, port int) {
	*s = append(*s, serverAddress{host, port})
}

func (s serverAddresses) String() string {
	buf := bytes.NewBufferString("")
	for i, x := range s {
		buf.WriteString("\n\t")
		buf.WriteString(fmt.Sprintf("%2d. %s:%d", i, x.host, x.port))
	}
	return buf.String()
}

func (s serverAddresses) Line() string {
	buf := bytes.NewBufferString("")
	for i, x := range s {
		if i > 0 {
			buf.WriteRune(';')
		}
		buf.WriteString(fmt.Sprintf("%s:%d", x.host, x.port))
	}
	return buf.String()
}

func (s serverAddresses) XML() string {
	var servers = make([]XMLConfigServer, len(s))
	for i, addr := range s {
		servers[i] = XMLConfigServer{
			Host: addr.host,
			Port: addr.port,
		}
	}

	config := XMLConfig{
		Servers: XMLConfigServers{
			Servers: servers,
		},
	}

	if data, err := xml.Marshal(config); err != nil {
		return ""
	} else {
		return string(data)
	}
}

func resolveServerAddresses(router string) (addresses serverAddresses) {
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
