package goh4

import "os"

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
	h := H4{
		Endpoint: endpoint,
		Secret:   secret,
		Key:      key,
		Prefix:   "/openapi/v1",
		Verify:   false,
	}

	return &h
}
