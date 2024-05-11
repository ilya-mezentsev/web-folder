// Package api is for requesting folders and files infos from corresponding API

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var errStatusIsNotOk = errors.New("response-status-is-not-ok")

type Service[T any] struct {
	request *http.Request
	client  *http.Client
}

func MustNew[T any](addr string) Service[T] {
	request, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		panic(fmt.Errorf("unable to create request: %w", err))
	}

	return Service[T]{
		request: request,
		client: &http.Client{
			Timeout: 5 * time.Second, // can be configured if needed
		},
	}
}

func (s Service[T]) Fetch() (T, error) {
	return s.makeRequest(s.request)
}

func (s Service[T]) makeRequest(request *http.Request) (T, error) {
	var result T
	response, err := s.client.Do(request)
	if err != nil {
		return result, err
	} else if response.StatusCode != http.StatusOK {
		return result, errStatusIsNotOk
	}

	var responseBody []byte
	responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(responseBody, &result)

	return result, err
}

func (s Service[T]) FetchWithQuery(key, param string) (T, error) {
	newRequest := s.request.Clone(s.request.Context())

	q := newRequest.URL.Query()
	q.Add(key, param)

	newRequest.URL.RawQuery = q.Encode()

	return s.makeRequest(newRequest)
}
