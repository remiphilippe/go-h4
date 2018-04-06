package goh4

import (
	"testing"
)

func TestGetFilters(t *testing.T) {
	t.Log("Starting TestGetFilters Test")
	h := setupH4()

	_, err := h.GetAllFilters()
	if err != nil {
		t.Errorf("Error in GetAllFilters: %s", err.Error())
		return
	}
}

func TestGetFilter(t *testing.T) {
	t.Log("Starting TestGetFilter Test")
	h := setupH4()

	_, err := h.GetFilter("5aa8df2f497d4f2ed62454c7")
	if err != nil {
		t.Errorf("Error in GetFilter: %s", err.Error())
		return
	}

	//spew.Dump(filter)

}
