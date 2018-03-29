package goh4

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
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
		glog.Errorf("Error Marshalling vrf %s", err)
		return err
	}
	postResp, err := h.Post("/vrfs", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		glog.Errorf("POST error %s", err)
		return err
	}

	err = json.Unmarshal(postResp, &v)
	if err != nil {
		glog.Errorf("Error unmarshalling json %s", err)
		return err
	}

	return nil
}

// GetVRF Get all VRF
func (h *H4) GetVRF() ([]VRF, error) {
	getResp, err := h.Get("/vrfs")
	if err != nil {
		glog.Errorf("GET error %s", err)
		return nil, err
	}

	var jsonResp []VRF
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		glog.Errorf("Error unmarshalling json %s", err)
		return nil, err
	}

	return jsonResp, nil
}

// DeleteVRF Delete VRF
func (h *H4) DeleteVRF(vrfID int) error {
	err := h.Delete(fmt.Sprintf("/vrfs/%d", vrfID), "")
	if err != nil {
		glog.Errorf("Error deleting vrf %d: %s", vrfID, err)
		return err
	}

	return nil
}
