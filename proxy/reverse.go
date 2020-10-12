package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
)

// ReverseConfig is the configuration to run a proxy.
type ReverseConfig struct {
	CaCertPath string
	CertPath string
	KeyPath string
	OriginPort int
}

// A ReverseRequestProcessor performs actions on proxied HTTP requests
type ReverseRequestProcessor interface {
	BeforeRequest(req *http.Request) (error)
	AfterRequest(req *http.Request) (error)
}

var reverseRequestProcessors = make(map[string]ReverseRequestProcessor)

// RegisterReverseProcessor will make not of new types of processors
func RegisterReverseProcessor(processorType string, processor ReverseRequestProcessor) {
	reverseRequestProcessors[processorType] = processor
}

// ReverseEnabled checks the environment to see if we're configured for Proxy serving
func ReverseEnabled() (*ReverseConfig, error) {
	var exists bool
	var err error
	if os.Getenv("HABERDASHER_REVPROXY") != "" {
		config := ReverseConfig{}
		if config.CaCertPath, exists = os.LookupEnv("HABERDASHER_TLS_CACERT"); !exists {
			return nil, errors.New("HABERDASHER_TLS_CACERT not set")
		}
		if config.CertPath, exists = os.LookupEnv("HABERDASHER_TLS_CERT"); !exists {
			return nil, errors.New("HABERDASHER_TLS_CERT not set")
		}
		if config.KeyPath, exists = os.LookupEnv("HABERDASHER_TLS_KEY"); !exists {
			return nil, errors.New("HABERDASHER_TLS_KEY not set")
		}
		if originPort, exists := os.LookupEnv("HABERDASHER_REVPROXY_TO"); !exists {
			return nil, errors.New("HABERDASHER_REVPROXY_TO not set")
		} else {
			config.OriginPort, err = strconv.Atoi(originPort)
			if err != nil {
				return nil, err
			}
		}
		return &config, nil
	}
	return nil, nil
}

// ReverseStart runs the proxy server
func ReverseStart(config *ReverseConfig) {
	// This function revises the incoming request for the original service
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Forwarded-Proto", "https")
		req.Header.Add("X-Forwarded-Port", "443")
		req.Header.Add("X-Haberdasher-Server", "1")
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("localhost:%d", config.OriginPort)
	}
	proxy := &httputil.ReverseProxy{Director: director}

	// Set up TLS with optional mTLS
	caCert, err := ioutil.ReadFile(config.CaCertPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		ClientAuth: tls.VerifyClientCertIfGiven,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		writerHeader := writer.Header()
		for _, peerCert := range req.TLS.PeerCertificates {
			writerHeader.Add("X-Haberdasher-Auth", peerCert.Subject.String())
		}
		for _, processor := range reverseRequestProcessors {
			processor.BeforeRequest(req)
		}
		proxy.ServeHTTP(writer, req)
		for _, processor := range reverseRequestProcessors {
			processor.AfterRequest(req)
		}
	})

	server := &http.Server{
		Addr: ":443",
		Handler: mux,
		TLSConfig: tlsConfig,
	}
	err = server.ListenAndServeTLS(config.CertPath, config.KeyPath)
	log.Fatal(err)
}

