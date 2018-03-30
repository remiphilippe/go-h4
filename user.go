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
	Scope      string  `json:"app_scope_id"`
	Role       []*Role `json:"role_ids,omitempty"`
	DisabledAt int     `json:"disabled_at,omitempty"`
}

// AddUser add a user
func (h *H4) AddUser(u *User) error {
	jsonStr, err := json.Marshal(&u)
	if err != nil {
		return fmt.Errorf("Error Marshalling scope %s", err)
	}
	postResp, err := h.Post("/users", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error %s", err)
	}

	err = json.Unmarshal(postResp, &u)
	if err != nil {
		return fmt.Errorf("Error unmarshalling json %s", err)
	}
	return nil
}
