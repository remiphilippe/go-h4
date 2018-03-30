package goh4

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Role Struct for a Tetration Role
type Role struct {
	ID           string       `json:"id,omitempty"`
	Scope        *Scope       `json:"app_scope_id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Capabilities []Capability `json:"capabilities,omitempty"`
}

// MarshalJSON Converts Struct to JSON
func (r *Role) MarshalJSON() ([]byte, error) {
	type Alias Role
	return json.Marshal(&struct {
		Scope string `json:"app_scope_id"`
		*Alias
	}{
		Scope: r.Scope.ID,
		Alias: (*Alias)(r),
	})
}

// UnmarshalJSON Converts JSON to struct
func (r *Role) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Role

	aux := &struct {
		Scope string `json:"app_scope_id"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if h4 != nil {
		var scope *Scope

		if aux.Scope != "" {
			scope, err = h4.GetScope(aux.Scope)
			if err != nil {
				return err
			}
		} else {
			scope = nil
		}

		r.Scope = scope
	} else {
		return fmt.Errorf("H4 is not defined")
	}
	return nil
}

// AddCapability Add a capability to a role
func (r *Role) AddCapability(c Capability) error {
	var roleURL bytes.Buffer
	roleURL.WriteString("/roles/")
	roleURL.WriteString(r.ID)
	roleURL.WriteString("/capabilities")

	jsonStr, err := json.Marshal(&c)
	if err != nil {
		return fmt.Errorf("Error Marshalling capabilty %s", err)
	}
	_, err = h4.Post(roleURL.String(), fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error %s", err)
	}
	r.Capabilities = append(r.Capabilities, c)

	return nil
}

// Capability Struct to handle capabilities in Roles
type Capability struct {
	ID        string `json:"id,omitempty"`
	Scope     *Scope `json:"app_scope_id,omitempty"`
	Ability   string `json:"ability,omitempty"`
	Inherited bool   `json:"inherited,omitempty"`
}

// MarshalJSON Converts Struct to JSON
func (c *Capability) MarshalJSON() ([]byte, error) {
	type Alias Capability
	return json.Marshal(&struct {
		Scope string `json:"app_scope_id"`
		*Alias
	}{
		Scope: c.Scope.ID,
		Alias: (*Alias)(c),
	})
}

// UnmarshalJSON Converts JSON to struct
func (c *Capability) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Capability
	aux := &struct {
		Scope string `json:"app_scope_id"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if h4 != nil {
		var scope *Scope

		if aux.Scope != "" {
			scope, err = h4.GetScope(aux.Scope)
			if err != nil {
				return err
			}
		} else {
			scope = nil
		}

		c.Scope = scope
	} else {
		return fmt.Errorf("H4 is not defined")
	}

	return nil
}

// GetRoles Return all roles
func (h *H4) GetRoles() ([]*Role, error) {
	getResp, err := h.Get("/roles")
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}

	var jsonResp []*Role
	fmt.Printf("%s\n", getResp)
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}

// GetRole Return one role by ID
func (h *H4) GetRole(id string) (*Role, error) {
	getResp, err := h.Get(fmt.Sprintf("/roles/%s", id))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}

	var jsonResp *Role
	//fmt.Printf("%s\n", getResp)
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}

// AddRole Add a new role
func (h *H4) AddRole(r *Role) error {
	jsonStr, err := json.Marshal(&r)
	if err != nil {
		return fmt.Errorf("Error Marshalling role %s", err)
	}
	postResp, err := h.Post("/roles", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error: %s / POST: %s", err.Error(), jsonStr)
	}

	err = json.Unmarshal(postResp, &r)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), postResp)
	}

	return nil
}

// DeleteRole Delete a Role
func (h *H4) DeleteRole(roleID string) error {
	err := h.Delete(fmt.Sprintf("/roles/%s", roleID), "")
	if err != nil {
		return fmt.Errorf("Error deleting role %s: %s", roleID, err)
	}

	return nil
}

// AddRoleToUser Add a role to a user
func (h *H4) AddRoleToUser(r *Role, u *User) error {
	var userURL bytes.Buffer
	userURL.WriteString("/users/")
	userURL.WriteString(u.ID)
	userURL.WriteString("/add_role")

	rid := make(map[string]string)
	rid["role_id"] = r.ID

	jsonStr, err := json.Marshal(&rid)
	if err != nil {
		return fmt.Errorf("Error Marshalling role ID %s", err)
	}
	fmt.Printf("User ID: %s\n", u.ID)
	fmt.Printf("Request: %s\n", jsonStr)

	postResp, err := h.Put(userURL.String(), fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("PUT error %s", err)
	}
	fmt.Printf("Response: %s\n", postResp)
	return nil
}
