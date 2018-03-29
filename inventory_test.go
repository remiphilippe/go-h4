package goh4

import "testing"

func TestInventorySearch(t *testing.T) {
	t.Log("Starting InventorySearch Test")
	hostname := "collectorDatamover-1"
	scope := "Tetration"

	h := setupH4()

	var f0 Filter
	f0.Type = "contains"
	f0.Field = "host_name"
	f0.Value = hostname

	// define a filter
	f := &QueryFilter{
		Type:    "and",
		Filters: []interface{}{&f0},
	}

	// define a search
	s := &InventoryQuery{
		Limit:  100,
		Scope:  scope,
		Filter: f,
	}

	res, err := h.InventorySearch(s)
	if err != nil {
		t.Errorf("Error in InventorySearch: %s", err.Error())
		return
	}

	if len(res) != 1 {
		t.Errorf("Invalid array length, expecting 1, got %d\n", len(res))
		return
	}

	if res[0].Hostname != hostname {
		t.Errorf("Invalid hostname, expecting %s, got %s\n", hostname, res[0].Hostname)
		return
	}
}

func TestInventorySearchEmpty(t *testing.T) {
	t.Log("Starting InventorySearch Test - Empty Results")
	hostname := "collectorDatamover-X"
	scope := "Tetration"

	h := setupH4()

	var f0 Filter
	f0.Type = "contains"
	f0.Field = "host_name"
	f0.Value = hostname

	// define a filter
	f := &QueryFilter{
		Type:    "and",
		Filters: []interface{}{&f0},
	}

	// define a search
	s := &InventoryQuery{
		Limit:  100,
		Scope:  scope,
		Filter: f,
	}

	res, err := h.InventorySearch(s)
	if err != nil {
		t.Errorf("Error in InventorySearch: %s", err.Error())
	}

	if len(res) != 0 {
		t.Errorf("Invalid array length, expecting 0, got %d\n", len(res))
		return
	}
}

func TestInventorySearchInvalid(t *testing.T) {
	t.Log("Starting InventorySearch Test - Invalid Request")
	hostname := "collectorDatamover-X"
	scope := "Tetration"

	h := setupH4()

	var f0 Filter
	f0.Type = "contains"
	f0.Field = "host_name"
	f0.Value = hostname

	// define a filter
	f := &QueryFilter{
		Type:    "bob",
		Filters: []interface{}{&f0},
	}

	// define a search
	s := &InventoryQuery{
		Limit:  100,
		Scope:  scope,
		Filter: f,
	}

	_, err := h.InventorySearch(s)
	if err == nil {
		t.Error("No Error returned for TestInventorySearchInvalid")
		return
	}
}
