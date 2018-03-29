package goh4

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

type User struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (h *H4) AddUser(first, last, email string) (*User, error) {
	u := new(User)

	u.Email = email
	u.FirstName = first
	u.LastName = last

	jsonStr, err := json.Marshal(&u)
	if err != nil {
		glog.Errorf("Error Marshalling scope %s", err)
		return nil, err
	}
	postResp, err := h.Post("/users", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		glog.Errorf("POST error %s", err)
		return nil, err
	}

	err = json.Unmarshal(postResp, &u)
	if err != nil {
		glog.Errorf("Error unmarshalling json %s", err)
		return nil, err
	}
	return u, nil
}
