package goh4

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

type Role struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Capability struct {
	ID      string `json:"id,omitempty"`
	Scope   string `json:"app_scope_id,omitempty"`
	Ability string `json:"ability,omitempty"`
}

func (h *H4) AddRole(name, description string) (*Role, error) {
	r := new(Role)

	r.Name = name
	r.Description = description

	jsonStr, err := json.Marshal(&r)
	if err != nil {
		glog.Errorf("Error Marshalling role %s", err)
		return nil, err
	}
	postResp, err := h.Post("/roles", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		glog.Errorf("POST error %s", err)
		return nil, err
	}

	err = json.Unmarshal(postResp, &r)
	if err != nil {
		glog.Errorf("Error unmarshalling json %s", err)
		return nil, err
	}

	//TODO get the ID of the Role to update the struct before returning the Role

	return r, nil
}

func (h *H4) AddCapabilityToRole(c *Capability, r *Role) {
	var roleURL bytes.Buffer
	roleURL.WriteString("/roles/")
	roleURL.WriteString(r.ID)
	roleURL.WriteString("/capabilities")

	jsonStr, err := json.Marshal(&c)
	if err != nil {
		glog.Errorf("Error Marshalling capabilty %s", err)
		return
	}
	_, err = h.Post(roleURL.String(), fmt.Sprintf("%s", jsonStr))
	if err != nil {
		glog.Errorf("POST error %s", err)
		return
	}
}

func (h *H4) AddRoleToUser(r *Role, u *User) {
	var userURL bytes.Buffer
	userURL.WriteString("/users/")
	userURL.WriteString(u.ID)
	userURL.WriteString("/add_role")

	rid := make(map[string]string)
	rid["role_id"] = r.ID

	jsonStr, err := json.Marshal(&rid)
	if err != nil {
		glog.Errorf("Error Marshalling role ID %s", err)
	}
	fmt.Printf("User ID: %s\n", u.ID)
	fmt.Printf("Request: %s\n", jsonStr)

	postResp, err := h.Put(userURL.String(), fmt.Sprintf("%s", jsonStr))
	if err != nil {
		glog.Errorf("PUT error %s", err)
		return
	}
	fmt.Printf("Response: %s\n", postResp)
}
