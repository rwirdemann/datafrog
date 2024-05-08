package app

import (
	"net/http"
)

func Post(url string) error {
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	if _, err = client.Do(r); err != nil {
		return err
	}
	return nil
}
