package goh4

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var supportedMethods = []string{"GET", "POST", "PUT"}

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

type H4 struct {
	endpoint string
	prefix   string
	key      string
	secret   string
	verify   bool
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

	hash := base64.StdEncoding.EncodeToString(hmacSHA256([]byte(h.secret), buf.String()))
	request.Header.Set("Authorization", hash)

	return request
}

func (h *H4) url(path string) string {
	var buf bytes.Buffer
	buf.WriteString(h.endpoint)
	buf.WriteString(h.prefix)
	buf.WriteString(path)
	url := buf.String()

	return url
}

func (h *H4) Sign(request *http.Request) *http.Request {

	request.Header.Set("User-Agent", "Cisco Tetration Go client")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	request.Header.Set("Id", h.key)

	return h.hash(request)
}

func (h *H4) Get(path string) []byte {
	var insecureVerify bool

	url := h.url(path)

	if h.verify == false {
		insecureVerify = true
	} else {
		insecureVerify = false
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureVerify},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest:", err)
	}
	h.Sign(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Readall:", err)
	}

	return body
}

func (h *H4) Post(path string, json string) []byte {
	var insecureVerify bool

	url := h.url(path)

	if h.verify == false {
		insecureVerify = true
	} else {
		insecureVerify = false
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureVerify},
	}
	client := &http.Client{Transport: tr}
	reader := strings.NewReader(json)

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Fatal("NewRequest:", err)
	}

	// Hash needs to be performed here, or the reader will mark it as closed
	// Can't read twice https://groups.google.com/forum/#!topic/golang-nuts/S6qQHDtoxgo
	hash := hashSHA256([]byte(json))
	req.Header.Set("X-Tetration-Cksum", hash)

	h.Sign(req)

	//spew.Dump(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Readall:", err)
	}

	return body
}
