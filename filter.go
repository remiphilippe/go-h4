package goh4

import (
	"encoding/json"
	"fmt"
)

// InventoryFilter struct representingan InventoryFilter
type InventoryFilter struct {
	ID    string       `json:"id,omitempty"`
	Name  string       `json:"name"`
	Query *QueryFilter `json:"query"`
	Type  string       `json:"filter_type"`
}

// GetFilter Get a single application by ID
func (h *H4) GetFilter(id string) (*InventoryFilter, error) {
	getResp, err := h.Get(fmt.Sprintf("/filters/inventories/%s", id))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}

	var jsonResp *InventoryFilter
	//fmt.Printf("%s\n", getResp)
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}

// GetAllFilters Get all applications
func (h *H4) GetAllFilters() ([]*InventoryFilter, error) {
	getResp, err := h.Get(fmt.Sprintf("/filters/inventories"))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}
	//fmt.Printf("%s", getResp)

	var jsonResp []*InventoryFilter
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}
