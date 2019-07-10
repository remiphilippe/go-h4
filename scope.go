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

// func (s *Scope) String() string {
// 	return s.ID
// }

// GetParent Return the parent scope as a *Scope
func (s *Scope) GetParent() *Scope {
	if s.h4 != nil {
		p, err := s.h4.GetScope(s.Parent)
		if err != nil {
			return nil
		}
		return p
	}

	return nil
}

// AddScope Add a new scope
func (h *H4) AddScope(s *Scope) error {
	jsonStr, err := json.Marshal(&s)
	if err != nil {
		return fmt.Errorf("Error Marshalling scope: %s", err.Error())
	}
	postResp, err := h.Post("/app_scopes", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error: %s / POST: %s", err.Error(), jsonStr)
	}

	err = json.Unmarshal(postResp, &s)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), postResp)
	}

	s.h4 = h

	return nil
}

// GetScope Get a single scope by ID
func (h *H4) GetScope(id string) (*Scope, error) {
	getResp, err := h.Get(fmt.Sprintf("/app_scopes/%s", id))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}

	var jsonResp *Scope
	//fmt.Printf("%s\n", getResp)
	//
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	jsonResp.h4 = h

	return jsonResp, nil
}

// GetAllScopes Get all scopes
func (h *H4) GetAllScopes() ([]*Scope, error) {
	getResp, err := h.Get(fmt.Sprintf("/app_scopes"))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}
	//fmt.Printf("%s", getResp)

	var jsonResp []*Scope
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}

// DeleteScope Delete a scope by ID
func (h *H4) DeleteScope(scopeID string) error {
	err := h.Delete(fmt.Sprintf("/app_scopes/%s", scopeID), "")
	if err != nil {
		return fmt.Errorf("Error deleting scope %s: %s", scopeID, err)
	}

	return nil
}

// GetRootScope Returns the root scope for a given VRF
func (h *H4) GetRootScope(vrf int) (*Scope, error) {
	// First get all scopes
	allScopes, err := h.GetAllScopes()
	if err != nil {
		return nil, err
	}

	for v := range allScopes {
		if allScopes[v].VRF == vrf && allScopes[v].Parent == "" {
			return allScopes[v], nil
		}
	}

	return nil, nil
}

// GetScopeByName Returns a scope based on it's name
func (h *H4) GetScopeByName(name string) (*Scope, error) {
	scopes, err := h.GetAllScopes()
	if err != nil {
		return nil, err
	}

	for _, s := range scopes {
		if s.Name == name {
			return s, nil
		}
	}

	return nil, nil
}
