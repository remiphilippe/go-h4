package goh4

import (
	"io/ioutil"
	"os"
)

func loadJSONFromFile(file string) ([]byte, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func setupH4() *H4 {
	var endpoint, secret, key string
	if os.Getenv("OPENAPI_ENDPOINT") != "" {
		endpoint = os.Getenv("OPENAPI_ENDPOINT")
	} else {
		endpoint = "https://localhost"
	}

	if os.Getenv("OPENAPI_SECRET") != "" {
		secret = os.Getenv("OPENAPI_SECRET")
	} else {
		secret = "xxx"
	}

	if os.Getenv("OPENAPI_KEY") != "" {
		key = os.Getenv("OPENAPI_KEY")
	} else {
		key = "xxx"
	}

	// H4 object
	// h := H4{
	// 	Endpoint: endpoint,
	// 	Secret:   secret,
	// 	Key:      key,
	// 	Prefix:   "/openapi/v1",
	// 	Verify:   false,
	// }

	h := NewH4(endpoint, secret, key, "/openapi/v1", false)

	return h
}
