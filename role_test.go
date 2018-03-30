package goh4

import (
	"fmt"
	"testing"
)

func testRoleSetup(h *H4) (*Scope, error) {
	id := 424242
	vrfName := "TestVRF42"
	scopeName := "ScopeVRF42"
	tenant := 0

	v := VRF{
		ID:       id,
		Name:     vrfName,
		TenantID: tenant,
	}

	err := h.AddVRF(&v)
	if err != nil {
		return nil, fmt.Errorf("Error in AddVRF: %s", err.Error())
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
		return nil, fmt.Errorf("Error in AddScope: %s", err.Error())
	}

	return &s, nil
}

func testRoleTeardown(h *H4, scopeID string, vrfID int) error {
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

func TestAddRole(t *testing.T) {
	t.Log("Starting TestAddRole Test")
	roleName := "TestVRF42"
	h := setupH4()

	s, err := testRoleSetup(h)
	if err != nil {
		t.Errorf("Test setup failed - %s", err.Error())
		return
	}

	r := new(Role)

	r.Scope = s.GetParent()
	r.Name = roleName

	// Add a role, make sure it succeeds before we do anything else
	err = h.AddRole(r)
	if err != nil {
		t.Errorf("Error in AddRole: %s", err.Error())
	} else {
		c1 := Capability{
			Scope:   s,
			Ability: AbilityRead,
		}

		err = r.AddCapability(c1)
		if err != nil {
			t.Errorf("Error in AddCapability c1: %s", err.Error())
		}

		testRole, err := h.GetRole(r.ID)
		if err != nil {
			t.Errorf("Error in GetRole: %s", err.Error())
		}

		if testRole == nil {
			t.Errorf("testRole is nil, something went wrong")
		} else {
			if testRole.ID != r.ID {
				t.Errorf("Invalid Role ID, expecting: %s got: %s\n", r.ID, testRole.ID)
			}
			if len(testRole.Capabilities) != 1 {
				t.Errorf("Invalid Capability c1 Length, expecting: %d got: %d\n", 1, len(testRole.Capabilities))
			} else {
				if testRole.Capabilities[0].Scope.ID != s.ID {
					t.Errorf("Invalid Capability c1 Scope, expecting: %s got: %s\n", s.ID, testRole.Capabilities[0].Scope.ID)
				}
				if testRole.Capabilities[0].Ability != AbilityRead {
					t.Errorf("Invalid Capability c1 Ability, expecting: %s got: %s\n", AbilityRead, testRole.Capabilities[0].Ability)
				}

				c2 := Capability{
					Scope:   s,
					Ability: AbilityWrite,
				}
				err = testRole.AddCapability(c2)
				// Make sure it succeeds
				if err != nil {
					t.Errorf("Error in AddCapability c2: %s", err.Error())
				}
				// Make sure the object gets updated
				if len(testRole.Capabilities) != 2 {
					t.Errorf("Invalid Capability c2 Length, expecting: %d got: %d\n", 2, len(testRole.Capabilities))
				}
			}
		}

		testRole2, err := h.GetRole(r.ID)
		if err != nil {
			t.Errorf("Error in GetRole (testRole2): %s", err.Error())
		}

		if testRole2 == nil {
			t.Errorf("testRole2 is nil, something went wrong")
		} else {
			if len(testRole2.Capabilities) != 2 {
				t.Errorf("Invalid Capability testRole2 Length, expecting: %d got: %d\n", 2, len(testRole2.Capabilities))
			}
		}

		// Clean up by deleting the role we just created
		err = h.DeleteRole(r.ID)
		if err != nil {
			t.Errorf("Error in DeleteRole: %s", err.Error())
		}

		err = h.DeleteRole(r.ID)
		if err == nil {
			t.Errorf("Trying to delete unknown role, where is my error?")
		}
	}

	err = testRoleTeardown(h, s.ID, s.VRF)
	if err != nil {
		t.Errorf("Test teardown failed - %s", err.Error())
	}

}
