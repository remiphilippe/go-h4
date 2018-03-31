package goh4

import (
	"fmt"
	"testing"
)

func testUserSetup(h *H4) (*Scope, *Role, error) {
	id := 424242
	vrfName := "TestVRF42"
	scopeName := "ScopeVRF42"
	tenant := 0
	roleName := "TestVRF42"

	v := VRF{
		ID:       id,
		Name:     vrfName,
		TenantID: tenant,
	}

	err := h.AddVRF(&v)
	if err != nil {
		return nil, nil, fmt.Errorf("Error in AddVRF: %s", err.Error())
	}

	f0 := Filter{
		Type:  "eq",
		Field: "ip",
		Value: "1.2.3.4",
	}

	s := Scope{
		VRF:  v.ID,
		Name: scopeName,
		Query: &QueryFilter{
			Type:    "and",
			Filters: []interface{}{&f0},
		},
		Parent: v.Scope.ID,
	}

	err = h.AddScope(&s)
	if err != nil {
		return nil, nil, fmt.Errorf("Error in AddScope: %s", err.Error())
	}

	r := new(Role)

	r.Scope = s.GetParent()
	r.Name = roleName

	// Add a role, make sure it succeeds before we do anything else
	err = h.AddRole(r)
	if err != nil {
		return nil, nil, fmt.Errorf("Error in AddRole: %s", err.Error())
	}

	c1 := Capability{
		Scope:   &s,
		Ability: AbilityRead,
	}

	err = r.AddCapability(c1)
	if err != nil {
		return nil, nil, fmt.Errorf("Error in AddCapability c1: %s", err.Error())
	}

	return &s, r, nil
}

func testUserTeardown(h *H4, scopeID string, vrfID int) error {
	err := h.DeleteScope(scopeID)
	if err != nil {
		return fmt.Errorf("Error in DeleteScope: %s", err.Error())
	}

	err = h.DeleteVRF(vrfID)
	if err != nil {
		return fmt.Errorf("Error in DeleteVRF: %s", err.Error())
	}

	return nil
}

func TestGetUser(t *testing.T) {
	t.Log("Starting TestGetUser Test")
	testID := "5aa6ee0a497d4f46d9e7fe9b"

	h := setupH4()

	//TODO update this, UUID is not static
	u, err := h.GetUser(testID)
	if err != nil {
		t.Errorf("Error in GetUser: %s", err.Error())
		return
	}

	if u.ID != testID {
		t.Errorf("Error ID is invalid, expecting %s got %s", testID, u.ID)
	}

	if len(u.Roles) != 2 {
		t.Errorf("Error not enough roles, expecting %d got %d", 2, len(u.Roles))
	}

}

func TestAddUser(t *testing.T) {
	return
	// t.Log("Starting TestAddUser Test")
	// h := setupH4()
	//
	// s, r, err := testUserSetup(h)
	// if err != nil {
	// 	t.Errorf("Test setup failed - %s", err.Error())
	// 	return
	// }
	//
	// u := User{
	// 	FirstName: "Bob",
	// 	LastName:  "Sponge",
	// 	Email:     "bobby2@sponge.com",
	// 	Scope:     s.GetParent(),
	// 	Roles:     []*Role{r},
	// }
	//
	// err = h.AddUser(&u)
	// if err != nil {
	// 	t.Errorf("Error in AddUser: %s", err.Error())
	// }
	//
	// err = testUserTeardown(h, s.ID, s.VRF)
	// if err != nil {
	// 	t.Errorf("Test teardown failed - %s", err.Error())
	// }
}
