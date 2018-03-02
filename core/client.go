package core

import "net/http"

type NewClient func(apiKey string, apiSecret string)

type Poster interface {
	Post(endpoint string, params map[string]interface{}) (*http.Response, error)
}

type Getter interface {
	Get(endpoint string, params map[string]interface{}) (*http.Response, error)
}

