package goh4

import (
	"encoding/json"
	"fmt"
)

// VRF Structure holding a VRF
type VRF struct {
	ID                   int      `json:"id,omitempty"`
	Name                 string   `json:"name,omitempty"`
	TenantID             int      `json:"tenant_id"`
	SwitchVRFs           []string `json:"switch_vrfs,omitempty"`
	ApplyCollectionRules bool     `json:"apply_monitoring_rules,string,omitempty"`
}

// AddVRF Create a new VRF
func (h *H4) AddVRF(v *VRF) error {
	jsonStr, err := json.Marshal(&v)
	if err != nil {
		return fmt.Errorf("Error Marshalling VRF: %s", err.Error())
	}
	postResp, err := h.Post("/vrfs", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error: %s / POST: %s", err.Error(), postResp)
	}

	err = json.Unmarshal(postResp, &v)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), postResp)
	}

	return nil
}

// GetVRF Get all VRF
func (h *H4) GetVRF() ([]VRF, error) {
	getResp, err := h.Get("/vrfs")
	if err != nil {
		return nil, fmt.Errorf("GET error: %s / GET: %s", err.Error(), getResp)
	}

	var jsonResp []VRF
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}

// DeleteVRF Delete VRF
func (h *H4) DeleteVRF(vrfID int) error {
	err := h.Delete(fmt.Sprintf("/vrfs/%d", vrfID), "")
	if err != nil {
		return fmt.Errorf("Error deleting vrf %d: %s", vrfID, err)
	}

	return nil
}
