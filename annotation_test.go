package goh4

import (
	"testing"
)

func TestAddAnnotation(t *testing.T) {
	t.Log("Starting TestAddAnnotation Test")
	ip := "42.42.42.42"
	scope := "Default"

	h := setupH4()

	attributes := make(map[string]string)
	attributes["testTag42"] = "myTag42"

	a := Annotation{
		IP:         ip,
		Scope:      scope,
		Attributes: attributes,
	}

	err := h.AddAnnotation(&a)
	if err != nil {
		t.Errorf("Error in AddAnnotation: %s", err.Error())
	}

}

func TestGetAnnotationIP(t *testing.T) {
	t.Log("Starting TestGetAnnotation Test - IP")
	ip := "42.42.42.42"
	scope := "Default"

	h := setupH4()

	res, err := h.GetAnnotation(scope, ip)
	if err != nil {
		t.Errorf("Error in GetAnnotation: %s", err.Error())
		return
	}

	if res.IP != ip {
		t.Errorf("Invalid IP, expecting: %s got: %s\n", ip, res.IP)
	}
	if res.Scope != scope {
		t.Errorf("Invalid scope, expecting: %s got: %s\n", scope, res.Scope)
	}
	if res.Attributes["testTag42"] != "myTag42" {
		t.Errorf("Invalid tag, expecting: %s got: %s\n", "myTag42", res.Attributes["testTag42"])
	}
}

func TestGetAnnotationSubnet(t *testing.T) {
	t.Log("Starting TestGetAnnotation Test - Subnet")
	ip := "172.20.0.0/16"
	scope := "Default"

	h := setupH4()

	res, err := h.GetAnnotation(scope, ip)
	if err != nil {
		t.Errorf("Error in GetAnnotation: %s", err.Error())
		return
	}

	if res.IP != ip {
		t.Errorf("Invalid IP, expecting: %s got: %s\n", ip, res.IP)
	}
	if res.Scope != scope {
		t.Errorf("Invalid scope, expecting: %s got: %s\n", scope, res.Scope)
	}
}

func TestDeleteAnnotation(t *testing.T) {
	t.Log("Starting TestDeleteAnnotation Test")
	ip := "42.42.42.42"
	scope := "Default"

	h := setupH4()

	err := h.DeleteAnnotation(scope, ip)
	if err != nil {
		t.Errorf("Error in DeleteAnnotation: %s", err.Error())
	}
}

func TestGetAnnotationIPEmpty(t *testing.T) {
	t.Log("Starting TestGetAnnotationIPEmpty Test - After Deletetion")
	ip := "42.42.42.42"
	scope := "Default"

	h := setupH4()

	res, err := h.GetAnnotation(scope, ip)
	if err != nil {
		t.Errorf("Error in GetAnnotation: %s", err.Error())
		return
	}

	if res.IP != ip {
		t.Errorf("Invalid IP, expecting: %s got: %s\n", ip, res.IP)
	}
	if res.Scope != scope {
		t.Errorf("Invalid scope, expecting: %s got: %s\n", scope, res.Scope)
	}
	if len(res.Attributes) > 0 {
		t.Errorf("Invalid attributes count, expecting: 0 got: %d\n", len(res.Attributes))
	}
}
