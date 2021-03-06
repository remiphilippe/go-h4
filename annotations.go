package goh4

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/golang/glog"
)

// Annotation Tetration Annotation Format
type Annotation struct {
	Scope      string            `json:"-"`
	IP         string            `json:"ip"`
	Attributes map[string]string `json:"attributes"`
}

// AddAnnotation Add an annotation to an IP
func (h *H4) AddAnnotation(a *Annotation) error {
	jsonStr, err := json.Marshal(&a)
	if err != nil {
		return fmt.Errorf("Error Marshalling annotation: %s", err.Error())
	}
	h.Post(fmt.Sprintf("/inventory/tags/%s", a.Scope), fmt.Sprintf("%s", jsonStr))

	return nil
}

// GetAnnotation Get Annotation for a specific IP / Subnet
func (h *H4) GetAnnotation(scope, query string) (*Annotation, error) {
	getResp, err := h.Get(fmt.Sprintf("/inventory/tags/%s?ip=%s", scope, query))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s / GET: %s", err.Error(), getResp)
	}

	jsonResp := make(map[string]string)
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		glog.Errorf("Error unmarshalling json %s", err)
		return nil, err
	}

	return &Annotation{
			Scope:      scope,
			IP:         query,
			Attributes: jsonResp,
		},
		nil
}

// DeleteAnnotation Delete all annotation associated to an IP or subnet
func (h *H4) DeleteAnnotation(scope, ip string) error {
	q := make(map[string]string)
	q["ip"] = ip

	jsonStr, err := json.Marshal(&q)
	if err != nil {
		log.Fatal(err)
	}
	err = h.Delete(fmt.Sprintf("/inventory/tags/%s", scope), fmt.Sprintf("%s", jsonStr))
	if err != nil {
		glog.Errorf("Error deleting annotation for (ip, scope): (%s, %s) - %s", ip, scope, err)
		return err
	}

	return nil
}

// UploadAnnotation Upload an annotation file to Tetration OpenAPI. Data is a CSV format
func (h *H4) UploadAnnotation(data []byte, operation, scope string) ([]byte, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// Tim says it mimics the python client
	writer.SetBoundary("CiscoTetrationClient")

	if strings.ToLower(operation) != "add" && strings.ToLower(operation) != "delete" && strings.ToLower(operation) != "overwrite" {
		glog.Errorf("Looks like the operation is not supported, should be add, overwrite or delete but got %s", operation)
		return nil, fmt.Errorf("Looks like the operation is not supported, should be add, overwrite or delete but got %s", operation)
	}

	writer.WriteField("X-Tetration-Oper", strings.ToLower(operation))

	part, err := writer.CreateFormFile("file", "filename")
	if err != nil {
		glog.Errorf("Could not create multi-part form file header - %s\n", err)
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		glog.Errorf("Could not copy data into multi-part form part - %s\n", err)
		return nil, err
	}

	writer.Close()

	url := h.url(fmt.Sprintf("/assets/cmdb/upload/%s", scope))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		glog.Errorf("GET - NewRequest Error: %s\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	h.sign(req)

	return h.processRequest(req)
}
