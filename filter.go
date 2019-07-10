package goh4

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

// InventoryFilter struct representingan InventoryFilter
type InventoryFilter struct {
	ID         string       `json:"id,omitempty"`
	Name       string       `json:"name"`
	Scope      *Scope       //`json:"app_scope_id"` //this is overriden in unmarshall
	Query      *QueryFilter `json:"query"`
	ShortQuery *QueryFilter `json:"short_query,omitempty"`
	Primary    bool         `json:"primary"`
	Type       string       `json:"filter_type,omitempty"`
}

// MarshalJSON Converts Struct to JSON
func (f *InventoryFilter) MarshalJSON() ([]byte, error) {
	if f.Scope == nil {
		return nil, fmt.Errorf("scope is not defined")
	}

	type Alias InventoryFilter
	return json.Marshal(&struct {
		Scope string `json:"app_scope_id"`
		*Alias
	}{
		Scope: f.Scope.ID,
		Alias: (*Alias)(f),
	})
}

// UnmarshalJSON Converts JSON to struct
func (f *InventoryFilter) UnmarshalJSON(data []byte) error {
	var err error
	type Alias InventoryFilter

	aux := &struct {
		// TODO: check this out, doesn't make sense to change the json struct
		Scope map[string]string `json:"parent_app_scope"`
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		glog.Errorf("err: %s\n", err)
		return err
	}
	if h4 != nil {
		var scope *Scope

		if _, ok := aux.Scope["id"]; ok {
			scope, err = h4.GetScope(aux.Scope["id"])
			if err != nil {
				return err
			}
		} else {
			scope = nil
		}

		f.Scope = scope
	} else {
		return fmt.Errorf("H4 is not defined")
	}
	return nil
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

// GetFiltersByScope Get all filters for a scope
func (h *H4) GetFiltersByScope(scope *Scope) ([]*InventoryFilter, error) {
	var filters []*InventoryFilter

	iFilters, err := h.GetAllFilters()
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}

	for _, f := range iFilters {
		if f.Scope.ID == scope.ID {
			filters = append(filters, f)
		}
	}

	return filters, nil
}

// AddFilter Add a new filter
func (h *H4) AddFilter(f *InventoryFilter) error {
	jsonStr, err := json.Marshal(&f)
	if err != nil {
		return fmt.Errorf("Error Marshalling inventory filter %s", err)
	}
	postResp, err := h.Post("/filters/inventories", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error: %s / POST: %s", err.Error(), jsonStr)
	}

	err = json.Unmarshal(postResp, &f)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), postResp)
	}

	return nil
}

// DeleteFilter Delete a Filter
func (h *H4) DeleteFilter(filterID string) error {
	err := h.Delete(fmt.Sprintf("/filters/inventories/%s", filterID), "")
	if err != nil {
		return fmt.Errorf("Error deleting role %s: %s", filterID, err)
	}

	return nil
}
