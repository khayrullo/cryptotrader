package binance

import (
	"net/http"
	"strings"
	"io"
	"log"
)

// BinanceApiProxy is a standard web handler function that will proxy requests
// to the Binance API.
func BinanceApiProxy(w http.ResponseWriter, r *http.Request) {
	target := "https://api.binance.com"
	request, err := http.NewRequest(r.Method, target, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error: failed to create request: %v\n", err)
		return
	}
	request.URL.Path = r.URL.Path
	request.URL.RawQuery = r.URL.RawQuery

	for key, val := range r.Header {
		switch strings.ToLower(key) {
		case "x-mbx-apikey":
			request.Header[key] = val
		default:
		}
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error: failed to send request: %v\n", err)
		return
	}

	for key, val := range response.Header {
		w.Header()[key] = val
	}

	io.Copy(w, response.Body)
}

// NewBinanceApiProxyHandler return the Binance API proxy as a http.Handler.
//
// Useful if you need to strip the prefix, for example:
//     router.PathPrefix("/proxy/binance").Handler(
//         http.StripPrefix("/proxy/binance", NewBinanceApiProxyHandler()))
func NewBinanceApiProxyHandler() http.Handler {
	return http.HandlerFunc(BinanceApiProxy)
}
