package auth

import (
	"io"
	"net/http"
)

type Request struct {
	url    string
	method string
	req    *http.Request
}

func NewRequest(method string, url string) (*Request, error) {
	req, err := http.NewRequest(method, url, nil)
	return &Request{
		method: method,
		url:    url,
		req:    req,
	}, err
}

func (r *Request) AddHeader(key, value string) {
	r.req.Header.Add(key, value)
}

func (r *Request) Send() ([]byte, error) {
	client := &http.Client{}
	res, err := client.Do(r.req)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}
