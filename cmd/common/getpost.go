package common

import (
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"fmt"
	"gitlab.com/crankykernel/cryptotrader/core"
)

func doPostOrGet(poster core.Poster, getter core.Getter, args []string) {
	if len(args) == 0 {
		log.Fatal("error: an endpoint is required.")
	}

	endpoint := args[0]

	var params map[string]interface{}

	if len(args) > 1 {
		params = map[string]interface{}{}
		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			params[parts[0]] = parts[1]
		}
	}

	var response *http.Response
	var err error

	if poster != nil {
		response, err = poster.Post(endpoint, params)
	} else if getter != nil {
		response, err = getter.Get(endpoint, params)
	} else {
		log.Fatal("error: no poster or getter provided")
	}

	if response == nil {
		log.Fatal("error: nil response received")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("error: ", err)
	}
	fmt.Println(string(body))
}

func Post(client core.Poster, args []string) {
	doPostOrGet(client, nil, args)
}

func Get(client core.Getter, args []string) {
	doPostOrGet(nil, client, args)
}
