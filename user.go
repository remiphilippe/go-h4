package goh4

import (
	"encoding/json"
	"fmt"
)

// User Struct representing a User
type User struct {
	ID         string  `json:"id,omitempty"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Email      string  `json:"email"`
	Scope      *Scope  `json:"app_scope_id"`
	Roles      []*Role `json:"role_ids,omitempty"`
	DisabledAt int     `json:"disabled_at,omitempty"`
}

// MarshalJSON Converts Struct to JSON
func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	var ids []string

	for i := range u.Roles {
		ids = append(ids, u.Roles[i].ID)
	}

	return json.Marshal(&struct {
		Scope string   `json:"app_scope_id"`
		Roles []string `json:"role_ids,omitempty"`
		*Alias
	}{
		Scope: u.Scope.ID,
		Roles: ids,
		Alias: (*Alias)(u),
	})
}

// UnmarshalJSON Converts JSON to struct
func (u *User) UnmarshalJSON(data []byte) error {
	var err error
	type Alias User

	aux := &struct {
		Scope string   `json:"app_scope_id"`
		Roles []string `json:"role_ids,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(u),
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

		u.Scope = scope

		for i := range aux.Roles {
			r, err := h4.GetRole(aux.Roles[i])
			if err != nil {
				return err
			}
			u.Roles = append(u.Roles, r)
		}
	} else {
		return fmt.Errorf("H4 is not defined")
	}
	return nil
}

// AddUser add a user
func (h *H4) AddUser(u *User) error {
	jsonStr, err := json.Marshal(&u)
	if err != nil {
		return fmt.Errorf("Error Marshalling VRF: %s", err.Error())
	}
	postResp, err := h.Post("/users", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error: %s / POST: %s", err.Error(), jsonStr)
	}

	err = json.Unmarshal(postResp, &u)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), postResp)
	}
	return nil
}

// GetUser Get a single user by ID
func (h *H4) GetUser(id string) (*User, error) {
	getResp, err := h.Get(fmt.Sprintf("/users/%s", id))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}

	var jsonResp *User
	//fmt.Printf("%s\n", getResp)

	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}
