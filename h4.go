package goh4

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var h4 *H4

// NewH4 Factory to initialize H4
func NewH4(endpoint, secret, key, prefix string, verify bool) *H4 {
	h4 = new(H4)
	h4.Endpoint = endpoint
	h4.Secret = secret
	h4.Key = key
	h4.Prefix = prefix
	h4.Verify = verify

	return h4
}

func hashSHA256(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func checksum(content string, request *http.Request) *http.Request {
	hash := hashSHA256([]byte(content))
	request.Header.Set("X-Tetration-Cksum", hash)

	return request
}

func (h *H4) processRequest(req *http.Request) ([]byte, error) {
	var insecureVerify bool
	var err error

	if h.Verify == false {
		insecureVerify = true
	} else {
		insecureVerify = false
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureVerify},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request Processing Error: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GET - ReadAll Error: %s", err)
	}

	switch resp.StatusCode {
	case 200:
		err = nil
	case 400:
		contentType := resp.Header.Get("Content-Type")
		s := strings.Split(contentType, ";")
		header := strings.TrimSpace(s[0])

		if header == "application/json" {
			jsonResp := make(map[string]string)
			err = json.Unmarshal(body, &jsonResp)
			if err != nil {
				return nil, err
			}
			err = fmt.Errorf("Request Error (%d): %s", resp.StatusCode, jsonResp["error"])
		} else {
			err = fmt.Errorf("Request Error (%d): %s", resp.StatusCode, body)
		}
	case 401:
		err = fmt.Errorf("Authorization Error (%d): %s", resp.StatusCode, body)
	case 403:
		err = fmt.Errorf("Authentication Error (%d): %s", resp.StatusCode, body)
	default:
		err = fmt.Errorf("Other Error (%d): %s", resp.StatusCode, body)
	}

	if err != nil {
		return nil, err
	}

	return body, nil
}

// H4 defines a structure for the Tetration OpenAPI
type H4 struct {
	Endpoint string
	Prefix   string
	Key      string
	Secret   string
	Verify   bool
}

func (h *H4) hash(request *http.Request) *http.Request {
	var buf bytes.Buffer
	buf.WriteString(request.Method)
	buf.WriteString("\n")
	buf.WriteString(request.URL.RequestURI())
	buf.WriteString("\n")
	buf.WriteString(request.Header.Get("X-Tetration-Cksum"))
	buf.WriteString("\n")
	buf.WriteString(request.Header.Get("Content-Type"))
	buf.WriteString("\n")
	buf.WriteString(request.Header.Get("Timestamp"))
	buf.WriteString("\n")

	hash := base64.StdEncoding.EncodeToString(hmacSHA256([]byte(h.Secret), buf.String()))
	request.Header.Set("Authorization", hash)

	return request
}

func (h *H4) url(path string) string {
	var buf bytes.Buffer
	buf.WriteString(h.Endpoint)
	buf.WriteString(h.Prefix)
	buf.WriteString(path)
	url := buf.String()

	return url
}

func (h *H4) sign(request *http.Request) *http.Request {
	request.Header.Set("User-Agent", "Cisco Tetration Go client")
	request.Header.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	request.Header.Set("Id", h.Key)

	return h.hash(request)
}

// Get Perform a GET operation on Tetration OpenAPI
func (h *H4) Get(path string) ([]byte, error) {
	url := h.url(path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("GET - NewRequest Error: %s", err)
	}
	h.sign(req)

	return h.processRequest(req)
}

// Post Perform a POST operation on Tetration OpenAPI
func (h *H4) Post(path string, json string) ([]byte, error) {
	url := h.url(path)

	reader := strings.NewReader(json)

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, fmt.Errorf("POST - NewRequest Error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Hash needs to be performed here, or the reader will mark it as closed
	// Can't read twice https://groups.google.com/forum/#!topic/golang-nuts/S6qQHDtoxgo
	checksum(json, req)

	h.sign(req)

	return h.processRequest(req)
}

// Put Perform a PUT operation on Tetration OpenAPI
func (h *H4) Put(path string, json string) ([]byte, error) {
	url := h.url(path)

	reader := strings.NewReader(json)

	req, err := http.NewRequest("PUT", url, reader)
	if err != nil {
		return nil, fmt.Errorf("PUT - NewRequest Error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Hash needs to be performed here, or the reader will mark it as closed
	// Can't read twice https://groups.google.com/forum/#!topic/golang-nuts/S6qQHDtoxgo
	checksum(json, req)

	h.sign(req)

	return h.processRequest(req)
}

// Delete Perform a DELETE operation on Tetration OpenAPI
func (h *H4) Delete(path string, json string) error {
	var err error
	var req *http.Request

	url := h.url(path)

	if json != "" {
		reader := strings.NewReader(json)

		req, err = http.NewRequest("DELETE", url, reader)
	} else {
		req, err = http.NewRequest("DELETE", url, nil)
	}
	if err != nil {
		return fmt.Errorf("DELETE - NewRequest Error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Hash needs to be performed here, or the reader will mark it as closed
	// Can't read twice https://groups.google.com/forum/#!topic/golang-nuts/S6qQHDtoxgo
	// FIXME Not sure if this is needed for DELETE?
	checksum(json, req)

	h.sign(req)

	_, err = h.processRequest(req)
	if err != nil {
		return err
	}

	return nil
}
