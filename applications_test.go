package goh4

import (
	"testing"
)

func TestGetApplications(t *testing.T) {
	t.Log("Starting TestGetApplications Test")
	h := setupH4()

	_, err := h.GetAllApplications(true)
	if err != nil {
		t.Errorf("Error in GetAllApplications: %s", err.Error())
		return
	}
}

func TestGetApplication(t *testing.T) {
	t.Log("Starting TestGetApplication Test")
	h := setupH4()

	_, err := h.GetApplication("5f8fb81d497d4f5453e657a3")
	if err != nil {
		t.Errorf("Error in GetApplication: %s", err.Error())
		return
	}

}
