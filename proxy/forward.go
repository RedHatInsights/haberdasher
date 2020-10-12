package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// ForwardConfig is the configuration to run a forward proxy
type ForwardConfig struct {
	CaCertPath string
	CertPath string
	KeyPath string
	ListenPort int
}

// ForwardEnabled checks the environmen to see if we're configured for Forward Proxy serving
func ForwardEnabled() (*ForwardConfig, error) {
	var exists bool
	var err error
	if os.Getenv("HABERDASHER_FWDPROXY") != "" {
		config := ForwardConfig{}
		if config.CaCertPath, exists = os.LookupEnv("HABERDASHER_TLS_CACERT"); !exists {
			return nil, errors.New("HABERDASHER_TLS_CACERT not set")
		}
		if config.CertPath, exists = os.LookupEnv("HABERDASHER_TLS_CERT"); !exists {
			return nil, errors.New("HABERDASHER_TLS_CERT not set")
		}
		if config.KeyPath, exists = os.LookupEnv("HABERDASHER_TLS_KEY"); !exists {
			return nil, errors.New("HABERDASHER_TLS_KEY not set")
		}
		if listenPort, exists := os.LookupEnv("HABERDASHER_FWDPROXY_PORT"); !exists {
			return nil, errors.New("HABERDASHER_FWDPROXY_PORT not set")
		} else {
			config.ListenPort, err = strconv.Atoi(listenPort)
			if err != nil {
				return nil, err
			}
		}
		return &config, nil
	}
	return nil, nil
}


// ForwardStart runs the forward proxy server
func ForwardStart(config *ForwardConfig) {
	// Set up client for mTLS
	caCert, err := ioutil.ReadFile(config.CaCertPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		// Munge the URL
		url := req.URL
		url.Scheme = "https"
		url.Host = url.Hostname()
		proxyReq, err := http.NewRequest(req.Method, url.String(), req.Body)
		proxyReq.Header["Host"] = []string{url.Host}
		proxyReq.Header.Add("X-Haberdasher-Client", "1")

		// Execute the proxy request
		proxyResp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(writer, "Proxy failure", 502)
		}
		defer proxyResp.Body.Close()
		respHeaders := writer.Header()
		for key, value := range proxyResp.Header {
			respHeaders[key] = value
		}
		writer.WriteHeader(proxyResp.StatusCode)
		if proxyResp.ContentLength > 0 {
			io.CopyN(writer, proxyResp.Body, proxyResp.ContentLength)
		}
	})
	server := &http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", config.ListenPort),
		Handler: mux,
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}
