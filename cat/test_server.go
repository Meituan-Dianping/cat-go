package cat

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type testServerManager struct {
	socketServers map[int]net.Listener
	servers       map[int]*testHTTPServer

	sampleRate float64
	block      bool
	ports      []int
}

type testHTTPServer struct {
	http.Server
	mux  *http.ServeMux
	port int
}

func (m *testServerManager) index(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world!")); err != nil {
		panic(err)
	}
}

func (m *testServerManager) router(w http.ResponseWriter, r *http.Request) {
	var properties = make([]routerConfigXMLProperty, 0, 3)

	properties = append(properties, routerConfigXMLProperty{
		Id:    propertySample,
		Value: fmt.Sprintf("%f", atomicLoadFloat64(&m.sampleRate)),
	})

	p := routerConfigXMLProperty{
		Id:    propertyBlock,
		Value: "false",
	}
	if m.block {
		p.Value = "true"
	}
	properties = append(properties, p)

	var servers = make(serverAddresses, 0)
	for _, port := range m.ports {
		servers.Add(localhost, port)
	}
	properties = append(properties, routerConfigXMLProperty{
		Id:    propertyRouters,
		Value: servers.Line(),
	})

	c := routerConfigXML{
		Properties: properties,
	}
	if data, err := xml.Marshal(c); err != nil {
		panic(err)
	} else {
		if _, err := w.Write(data); err != nil {
			panic(err)
		}
	}
}

func (m *testServerManager) startHTTP(port int) *testHTTPServer {
	server := newTestHTTPServer(port)

	server.mux.HandleFunc("/", m.index)
	server.mux.HandleFunc(routerPath, m.router)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %s", err)
		} else {
		}
	}()

	m.servers[port] = server
	return server
}

func (m *testServerManager) shutdownHTTP(port int) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	if err := m.servers[port].Shutdown(ctx); err != nil {
		log.Fatalf("Can't shutdown http server on port %d, %s", port, err)
	}
}

func (m *testServerManager) shutdownAll() {
	for port := range m.servers {
		m.shutdownHTTP(port)
	}
}

func newTestHTTPServer(port int) *testHTTPServer {
	mux := http.NewServeMux()

	return &testHTTPServer{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		mux:  mux,
		port: port,
	}
}

func newTestServerManager() *testServerManager {
	return &testServerManager{
		servers:    make(map[int]*testHTTPServer),
		sampleRate: 0.1,
		block:      false,
		ports:      []int{2280, 2281, 2282},
	}
}
