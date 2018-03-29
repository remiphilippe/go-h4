package goh4

import (
	"encoding/json"
	"fmt"
)

// Scope Structure holding a scope
type Scope struct {
	h4          *H4
	ID          string       `json:"id,omitempty"`
	Name        string       `json:"short_name,omitempty"`
	Description string       `json:"description,omitempty"`
	Query       *QueryFilter `json:"short_query,omitempty"`
	ScopeQuery  *QueryFilter `json:"query,omitempty"`
	Parent      string       `json:"parent_app_scope_id,omitempty"`
	VRF         int          `json:"vrf_id,omitempty"`
}

// GetParent Return the parent scope as a *Scope
func (s *Scope) GetParent() (*Scope, error) {
	return s.h4.GetScope(s.Parent)
}

// AddScope Add a new scope
func (h *H4) AddScope(s *Scope) (*Scope, error) {
	jsonStr, err := json.Marshal(&s)
	if err != nil {
		return nil, fmt.Errorf("Error Marshalling scope: %s", err.Error())
	}
	postResp, err := h.Post("/app_scopes", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return nil, fmt.Errorf("POST error: %s / POST: %s", err.Error(), postResp)
	}

	err = json.Unmarshal(postResp, &s)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), postResp)
	}

	s.h4 = h

	return s, nil
}

// GetScope Get a single scope by ID
func (h *H4) GetScope(id string) (*Scope, error) {
	getResp, err := h.Get(fmt.Sprintf("/app_scopes/%s", id))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s / GET: %s", err.Error(), getResp)
	}

	var jsonResp *Scope
	//fmt.Printf("%s\n", getResp)
	//
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}

// GetAllScope Get a single scope by ID
func (h *H4) GetAllScope() ([]*Scope, error) {
	getResp, err := h.Get(fmt.Sprintf("/app_scopes"))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s / GET: %s", err.Error(), getResp)
	}
	//fmt.Printf("%s", getResp)

	var jsonResp []*Scope
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}
