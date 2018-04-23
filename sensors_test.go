package goh4

import (
	"testing"
)

func TestGetSWAgents(t *testing.T) {
	t.Log("Starting TestGetSWAgents Test")

	h := setupH4()

	res, err := h.GetSWAgents()
	if err != nil {
		t.Errorf("Error in GetSWAgents: %s", err.Error())
		return
	}

	if len(res) == 0 {
		t.Errorf("Invalid result count, expecting more than 0 got: %d\n", len(res))
	}

	//spew.Dump(res)
}
