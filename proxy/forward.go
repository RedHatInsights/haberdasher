package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// ForwardConfig is the configuration to run a forward proxy
type ForwardConfig struct {
	CaCertPath string
	CertPath string
	KeyPath string
	ListenPort int
	MTLS bool
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
		if os.Getenv("HABERDASHER_FWDPROXY_MTLS") != "" {
			if config.CertPath, exists = os.LookupEnv("HABERDASHER_TLS_CERT"); !exists {
				return nil, errors.New("HABERDASHER_TLS_CERT not set")
			}
			if config.KeyPath, exists = os.LookupEnv("HABERDASHER_TLS_KEY"); !exists {
				return nil, errors.New("HABERDASHER_TLS_KEY not set")
			}
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

func forwardHandler(client *http.Client, writer http.ResponseWriter, req *http.Request) {
	url := req.URL
	// TODO: We need to intelligently copy HTTP headers onto the subrequest
	proxyReq, err := http.NewRequest(req.Method, url.String(), req.Body)
	// Do an SRV lookup on the service to see if it's secure or not
	cname, addrs, err := net.LookupSRV("", "", url.Hostname())
	if (err == nil && strings.HasSuffix(cname, ".svc.cluster.local.") && url.Scheme == "http") {
		// If the service listens on 443, upgrade this connection to HTTPS
		hasHTTPSPort := false
		for _, addr := range addrs {
			if addr.Port == 443 {
				hasHTTPSPort = true
				break
			}
		}
		if hasHTTPSPort {
			url.Scheme = "https"
			url.Host = url.Hostname()
			proxyReq, err = http.NewRequest(req.Method, url.String(), req.Body)
			proxyReq.Header["Host"] = []string{url.Host}		
		}
	}

	proxyReq.Header.Add("X-Haberdasher-Client", "1")

	// Execute the proxy request
	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Proxy failure: %s", err), 502)
	} else {
		respHeaders := writer.Header()
		for key, value := range proxyResp.Header {
			respHeaders[key] = value
		}
		writer.WriteHeader(proxyResp.StatusCode)
		if proxyResp.ContentLength != 0 {
			if proxyResp.ContentLength > 0 {
				io.CopyN(writer, proxyResp.Body, proxyResp.ContentLength)
			} else {
				io.Copy(writer, proxyResp.Body)
			}
			proxyResp.Body.Close()
		}
	}
}


// ForwardStart runs the forward proxy server
func ForwardStart(config *ForwardConfig) {
	// Set up client for TLS & mTLS
	caCert, err := ioutil.ReadFile(config.CaCertPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	certs := []tls.Certificate{}
	if config.MTLS {
		cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
		if err != nil {
			log.Fatal(err)
		}
		certs = append(certs, cert)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
				Certificates: certs,
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		forwardHandler(client, writer, req)
	})
	server := &http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", config.ListenPort),
		Handler: mux,
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}
