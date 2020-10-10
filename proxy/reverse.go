package proxy

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
)

// ReverseConfig is the configuration to run a proxy.
type ReverseConfig struct {
	CertPath string
	KeyPath string
	ListenPort int
	OriginPort int
}

// ReverseEnabled checks the environment to see if we're configured for Proxy serving
func ReverseEnabled() (*ReverseConfig, error) {
	var exists bool
	var err error
	if os.Getenv("HABERDASHER_PROXY") != "" {
		config := ReverseConfig{}
		if config.CertPath, exists = os.LookupEnv("HABERDASHER_TLS_CERT"); !exists {
			return nil, errors.New("HABERDASHER_TLS_CERT not set")
		}
		if config.KeyPath, exists = os.LookupEnv("HABERDASHER_TLS_KEY"); !exists {
			return nil, errors.New("HABERDASHER_TLS_KEY not set")
		}
		if listenPort, exists := os.LookupEnv("HABERDASHER_LISTEN"); !exists {
			return nil, errors.New("HABERDASHER_LISTEN not set")
		} else {
			config.ListenPort, err = strconv.Atoi(listenPort)
			if err != nil {
				return nil, err
			}
		}
		if originPort, exists := os.LookupEnv("HABERDASHER_PROXY_TO"); !exists {
			return nil, errors.New("HABERDASHER_PROXY_TO not set")
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
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Forwarded-Proto", "https")
		req.Header.Add("X-Forwarded-Port", fmt.Sprintf("%d", config.ListenPort))
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("localhost:%d", config.OriginPort)
	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		proxy.ServeHTTP(writer, req)
	})

	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", config.ListenPort), config.CertPath, config.KeyPath, nil)
	log.Fatal(err)
}

