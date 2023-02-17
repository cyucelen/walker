package walker

import (
	"net/http"
)

type RequestBuilder func(start, fetchCount int) (*http.Request, error)

type httpDataSource struct {
	client         *http.Client
	requestBuilder RequestBuilder
}

func (h *httpDataSource) Fetch(start, fetchCount int) (*http.Response, error) {
	req, err := h.requestBuilder(start, fetchCount)
	if err != nil {
		return nil, err
	}

	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewApiWalker(client *http.Client, requestBuilder RequestBuilder, sink Sink[*http.Response], options ...Option) *Walker[*http.Response] {
	source := &httpDataSource{
		requestBuilder: requestBuilder,
		client:         client,
	}

	return New(source.Fetch, sink, options...)
}
