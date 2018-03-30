package goh4

import (
	"fmt"
	"testing"
)

func TestAllScope(t *testing.T) {
	t.Log("Starting TestAllScope Test")

	h := setupH4()

	//p, err := h.GetScope("5aa6f0fb755f023c424218f3")
	_, err := h.GetAllScope()
	if err != nil {
		t.Errorf("Error in GetAllScope: %s", err.Error())
		return
	}

	// for _, v := range s {
	// 	if v.ID == "5aa6e28c497d4f3f820310f0" {
	// 		spew.Dump(v.Query)
	// 	}
	// }
}

func TestOneScope(t *testing.T) {
	t.Log("Starting TestOneScope Test")

	h := setupH4()

	//TODO update this, UUID is not static
	_, err := h.GetScope("5aa6f0fb755f023c424218f3")
	if err != nil {
		t.Errorf("Error in GetScope: %s", err.Error())
		return
	}

	// for _, v := range s {
	// 	if v.ID == "5aa6e28c497d4f3f820310f0" {
	// 		spew.Dump(v.Query)
	// 	}
	// }
}

func testScopeSetup(h *H4) (*VRF, error) {
	id := 424242
	vrfName := "TestVRF42"
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

	return &v, nil
}

func testScopeTeardown(h *H4, id int) error {
	err := h.DeleteVRF(id)
	if err != nil {
		return fmt.Errorf("Error in DeleteVRF: %s", err.Error())
	}

	return nil
}

func TestNewScope(t *testing.T) {
	t.Log("Starting TestNewScope Test")

	scopeName := "ScopeVRF42"

	h := setupH4()

	v, err := testScopeSetup(h)
	if err != nil {
		t.Errorf("Test setup failed - %s", err.Error())
		return
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
		t.Errorf("Error in AddScope: %s", err.Error())
		return
	}

	p := s.GetParent()
	if p == nil {
		t.Errorf("Error in GetParent")
	}

	if p.ID != v.Scope.ID {
		t.Errorf("Error parent scope is invalid, expecting %s got %s", v.Scope.ID, p.ID)
	}

	err = h.DeleteScope(s.ID)
	if err != nil {
		t.Errorf("Error in DeleteScope: %s", err.Error())
		return
	}

	err = testScopeTeardown(h, v.ID)
	if err != nil {
		t.Errorf("Test teardown failed - %s", err.Error())
	}
}

func TestNewScopeDuplicate(t *testing.T) {
	t.Log("Starting TestNewScopeDuplicate Test")

	scopeName := "ScopeVRF42"

	h := setupH4()

	v, err := testScopeSetup(h)
	if err != nil {
		t.Errorf("Test setup failed - %s", err.Error())
		return
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
		t.Errorf("Error in AddScope: %s", err.Error())
		return
	}

	err = h.AddScope(&s)
	if err == nil {
		t.Errorf("Error in AddScope, this should have raised an error")
	}

	err = h.DeleteScope(s.ID)
	if err != nil {
		t.Errorf("Error in DeleteScope: %s", err.Error())
		return
	}

	err = testScopeTeardown(h, v.ID)
	if err != nil {
		t.Errorf("Test teardown failed - %s", err.Error())
	}
}

func TestNewScopeInvalid(t *testing.T) {
	t.Log("Starting TestNewScopeDuplicate Test")

	scopeName := "ScopeVRF42"

	h := setupH4()

	v, err := testScopeSetup(h)
	if err != nil {
		t.Errorf("Test setup failed - %s", err.Error())
		return
	}

	f0 := Filter{
		Type:  "eq",
		Field: "ip",
		Value: "1.2.3.4",
	}

	s1 := Scope{
		VRF:  v.ID,
		Name: scopeName,
		Query: &QueryFilter{
			Type:    "bob",
			Filters: []interface{}{&f0},
		},
		Parent: v.Scope.ID,
	}

	// s2 := Scope{
	// 	VRF:  v.ID,
	// 	Name: 1,
	// 	Query: &QueryFilter{
	// 		Type:    "eq",
	// 		Filters: []interface{}{&f0},
	// 	},
	// 	Parent: v.Scope.ID,
	// }

	err = h.AddScope(&s1)
	if err == nil {
		t.Errorf("Error in AddScope (s1), this should have raised an error")
	}

	// err = h.AddScope(&s2)
	// if err == nil {
	// 	t.Errorf("Error in AddScope (s2), this should have raised an error")
	// }

	err = testScopeTeardown(h, v.ID)
	if err != nil {
		t.Errorf("Test teardown failed - %s", err.Error())
	}
}
