package httpclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"raind/internal/utils"
)

func NewHttpClient() *HttpClient {
	return &HttpClient{
		BaseUrl: "https://localhost:7755",
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: ReadCACert(),
				},
			},
		},
	}
}

type HttpClient struct {
	BaseUrl string
	Client  *http.Client
	Request *http.Request
}

func (c *HttpClient) NewRequest(method string, path string, body []byte) error {
	var err error
	if method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodPatch ||
		method == http.MethodDelete {
		c.Request, err = http.NewRequest(
			method,
			c.BaseUrl+path,
			bytes.NewReader(body),
		)
	} else {
		c.Request, err = http.NewRequest(
			method,
			c.BaseUrl+path,
			nil,
		)
	}
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	c.Request.Header.Set("Content-Type", "application/json")
	return nil
}

func (c *HttpClient) IsStatusOk(resp *http.Response) bool {
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return false
	}
	return true
}

func ReadCACert() *x509.CertPool {
	certPool := x509.NewCertPool()

	pemBytes, err := os.ReadFile(utils.PublicCertPath)
	if err != nil {
		return nil
	}

	if ok := certPool.AppendCertsFromPEM(pemBytes); !ok {
		return nil
	}

	return certPool
}
