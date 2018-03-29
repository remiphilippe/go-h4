package goh4

import (
	"testing"
)

func TestAllScope(t *testing.T) {
	t.Log("Starting TestAllScope Test")

	h := setupH4()

	//p, err := h.GetScope("5aa6f0fb755f023c424218f3")
	_, err := h.GetAllScope()
	if err != nil {
		t.Errorf("Error in GetAllScope: %s", err.Error())
		return
	}

	// for _, v := range s {
	// 	if v.ID == "5aa6e28c497d4f3f820310f0" {
	// 		spew.Dump(v.Query)
	// 	}
	// }
}

func TestOneScope(t *testing.T) {
	t.Log("Starting TestOneScope Test")

	h := setupH4()

	//TODO update this, UUID is not static
	_, err := h.GetScope("5aa6f0fb755f023c424218f3")
	if err != nil {
		t.Errorf("Error in GetScope: %s", err.Error())
		return
	}

	// for _, v := range s {
	// 	if v.ID == "5aa6e28c497d4f3f820310f0" {
	// 		spew.Dump(v.Query)
	// 	}
	// }
}
