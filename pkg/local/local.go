package local

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Project wide Local Client
var Client *LocalClient

type LocalClient struct {
	http   *http.Client
	dialer *websocket.Dialer

	url    url.URL
	header http.Header
}

func NewClient(port, pw, pc string) *LocalClient {
	tlsConf := &tls.Config{InsecureSkipVerify: true}

	tr := &http.Transport{
		TLSClientConfig: tlsConf,
	}

	header := http.Header{
		"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte("riot:"+pw))},
	}

	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = tlsConf

	return &LocalClient{
		http: &http.Client{Transport: tr, Timeout: 10 * time.Second},
		dialer: &dialer,

		url:    url.URL{Scheme: pc, Host: "127.0.0.1:" + port},
		header: header,
	}
}

func (l *LocalClient) InitWS() (*websocket.Conn, error) {
	wsURL := url.URL{
		Scheme: "wss",
		Host: l.url.Host,
	}

	c, resp, err := l.dialer.Dial(wsURL.String(), l.header)
	if errors.Is(err, websocket.ErrBadHandshake) {
		return c, fmt.Errorf("handshake failed with status %d", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (l *LocalClient) GET(endpoint string) (*http.Response, error) {
	l.url.Path = endpoint

	req, err := http.NewRequest(http.MethodGet, l.url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header = l.header
	return l.http.Do(req)
}
