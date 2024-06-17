package web

import (
	"net/http"
)

func Post(url string) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
