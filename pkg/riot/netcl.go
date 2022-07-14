package riot

import (
	"crypto/tls"
	"net/http"
)

type NetCL struct {
	http *http.Client

	port string
	pw   string
	pc   string
}

func NewNetCL(port, pw, pc string) *NetCL {
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

	return &NetCL{
		http: &http.Client{Transport: tr},

		port: port,
		pw: pw,
		pc: pc,
	}
}

func (n *NetCL) GET(endpoint string) (*http.Response, error) {
	url := n.pc + "://127.0.0.1:" + n.port + endpoint

	if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	} else {
		req.SetBasicAuth("riot", n.pw)
		return n.http.Do(req)
	}
}
