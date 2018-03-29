package goh4

import (
	"testing"
)

func TestGetVRF(t *testing.T) {
	t.Log("Starting TestGetVRF Test")

	h := setupH4()

	res, err := h.GetVRF()
	if err != nil {
		t.Errorf("Error in GetVRF: %s", err.Error())
	}

	if len(res) == 0 {
		t.Errorf("Invalid result count, expecting more than 0 got: %d\n", len(res))
	}

	found := false
	for _, v := range res {
		if v.Name == "Tetration" {
			found = true
			if v.ID != 676767 {
				t.Errorf("Invalid VRF ID, expecting: 676767 got: %d\n", v.ID)
			}
		}
	}
	if !found {
		t.Errorf("VRF Tetration not found\n")
	}
}

func TestAddVRF(t *testing.T) {
	t.Log("Starting TestAddVRF Test")

	h := setupH4()

	id := 424242
	name := "TestVRF42"
	tenant := 0

	v := VRF{
		ID:       id,
		Name:     name,
		TenantID: tenant,
	}

	err := h.AddVRF(&v)
	if err != nil {
		t.Errorf("Error in AddVRF: %s", err.Error())
	}
}

func TestDeleteVRF(t *testing.T) {
	t.Log("Starting TestDeleteVRF Test")
	h := setupH4()

	id := 424242

	err := h.DeleteVRF(id)
	if err != nil {
		t.Errorf("Error in DeleteVRF: %s", err.Error())
	}
}

func TestVRFe2e(t *testing.T) {
	t.Log("Starting TestVRFe2e Test")
	var err error

	h := setupH4()

	id := 424242
	name := "TestVRF42"
	tenant := 0

	v := VRF{
		ID:       id,
		Name:     name,
		TenantID: tenant,
	}

	err = h.AddVRF(&v)
	if err != nil {
		t.Errorf("Error in AddVRF: %s", err.Error())
	}

	t.Log("Starting TestVRFe2e Test - Check Created VRF")
	res, err := h.GetVRF()
	if err != nil {
		t.Errorf("Error in GetVRF: %s", err.Error())
	}

	if len(res) == 0 {
		t.Errorf("Invalid result count, expecting more than 0 got: %d\n", len(res))
	}

	found := false
	for _, v2 := range res {
		if v2.Name == name {
			found = true
			if v2.ID != id {
				t.Errorf("Invalid VRF ID, expecting: %d got: %d\n", id, v2.ID)
			}
		}
	}
	if !found {
		t.Errorf("VRF %s not found\n", name)
	}

	t.Log("Starting TestVRFe2e Test - Duplicate")
	err = h.AddVRF(&v)
	if err == nil {
		t.Errorf("Error - Duplicate VRF should have failed")
	}

	err = h.DeleteVRF(id)
	if err != nil {
		t.Errorf("Error in DeleteVRF: %s", err.Error())
	}
}
